// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// Cloud can perform concurrent-safe operations on a bunch of VMs:
// List all VMs, inspect a VM, start/stop a VM or remove it from the list
type Cloud struct {
	lock sync.RWMutex
	vms  VMs
}

// List the VMs handled under this Cloud
func (c *Cloud) List() VMs {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.vms.clone()
}

// Inspect a VM data by id (might not find it and return nil)
func (c *Cloud) Inspect(id int) (VM, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	vm, found := c.vms[id]
	return vm, found
}

// Launch a VM by id.
// The return includes a channel to optionally check completion of the launch
// process, apart from a possible error.
func (c *Cloud) Launch(id int) (chan struct{}, error) {
	if err := c.setVMState(id, STARTING); err != nil {
		return nil, err
	}
	return c.delayedTransition(id, RUNNING, StartDelay), nil
}

// Stop a VM by id.
// The return includes a channel to optionally check completion of the stop
// process, apart from a possible error.
func (c *Cloud) Stop(id int) (chan struct{}, error) {
	if err := c.setVMState(id, STOPPING); err != nil {
		return nil, err
	}
	return c.delayedTransition(id, STOPPED, StopDelay), nil
}

// Delete VM by id.
// An error is returned if the VM is missing or not in the Stopped state.
func (c *Cloud) Delete(id int) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	vm, found := c.vms[id]
	if !found {
		return fmt.Errorf("delete error: not found VM %d", id)
	}
	if vm.State != STOPPED {
		return fmt.Errorf("delete error: VM %d must be in state %v for deletion but it is %v", id, STOPPED, vm.State)
	}
	delete(c.vms, id)
	return nil
}

// delayedTransition set ups a timer in the background to move the VM
// identified by the given id to state after the given delay has passed.
// Uses setVMState internally to handle a safe concurrent delayed transition.
func (c *Cloud) delayedTransition(id int, state VMState, delay time.Duration) chan struct{} {
	done := make(chan struct{})
	time.AfterFunc(delay, func() {
		if err := c.setVMState(id, state); err != nil {
			log.Println(err)
		}
		close(done) // signal delayed transition completion
	})
	return done
}

// setVMState sets the VM identified by the given id to the given state.
// Might fail if the VM transition requested is illegal.
// Do it in a locked transaction
func (c *Cloud) setVMState(id int, state VMState) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	vm, found := c.vms[id]
	if !found {
		return fmt.Errorf("not found VM with id %d", id)
	}
	mutatedVM, err := vm.WithState(state)
	if err != nil {
		return err
	}
	c.vms[id] = mutatedVM
	return nil
}
