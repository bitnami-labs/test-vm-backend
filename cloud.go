package main

import (
	"fmt"
	"sync"
	"time"
)

// Cloud can list VMs, and perform start/stop and remove operations on them
type Cloud struct {
	lock sync.RWMutex
	vms  VMList
}

// List the VMs handled under this Cloud
func (c *Cloud) List() VMList {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.vms.clone()
}

// Inspect a VM data by id (might not find it and return nil)
func (c *Cloud) Inspect(id int) VM {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.vms.lookup(id)
}

// Launch a VM by id
func (c *Cloud) Launch(id int) error {
	return c.delayedTransition(id, STARTING, RUNNING, StartDelay)
}

// Stop a VM by id
func (c *Cloud) Stop(id int) error {
	return c.delayedTransition(id, STOPPING, STOPPED, StopDelay)
}

// Delete VM by id. Idempotent, will return true when VM is actually deleted
// and false if it was already not there.
func (c *Cloud) Delete(id int) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	if !c.vms.lookup(id).isValid() {
		return false
	}
	c.vms[id] = VM{}
	return true
}

// delayedTransition moves the VM identified by the given id to the final state
// after setting it in the ongoing state the given delay has passed.
// Uses setVMState internally to handle safe concurrent transitions.
func (c *Cloud) delayedTransition(id int, ongoing, final VMState, delay time.Duration) error {
	vm := c.vms.lookup(id)
	if !vm.isValid() {
		return fmt.Errorf("not found VM with id %d", id)
	}

	vm, err := c.setVMState(vm, ongoing)
	if err != nil {
		return err
	}
	c.vms[id] = vm
	time.AfterFunc(delay, func() {
		c.setVMState(vm, final)
	})
	return nil
}

// setVMState sets the VM identified by the given id to the given state.
// Might fail if the VM transition requested is illegal.
// Do it in a locked transaction
func (c *Cloud) setVMState(vm VM, state VMState) (VM, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	mutatedVM, err := vm.WithState(state)
	if err != nil {
		return VM{}, err
	}
	return mutatedVM, nil
}
