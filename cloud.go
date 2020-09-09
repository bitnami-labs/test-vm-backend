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

// DoneChannel to signal completion of a delayed action
// beware, it can timeout!
type DoneChannel chan struct{}

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

// Launch a VM by id
func (c *Cloud) Launch(id int) (DoneChannel, error) {
	if err := c.setVMState(id, STARTING); err != nil {
		return nil, err
	}
	return c.delayedTransition(id, RUNNING, StartDelay), nil
}

// Stop a VM by id
func (c *Cloud) Stop(id int) (DoneChannel, error) {
	if err := c.setVMState(id, STOPPING); err != nil {
		return nil, err
	}
	return c.delayedTransition(id, STOPPED, StopDelay), nil
}

// Delete VM by id. Idempotent, will return true when VM is actually deleted
// and false if it was already not there.
func (c *Cloud) Delete(id int) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	_, found := c.vms[id]
	if !found {
		return false
	}
	c.vms[id] = VM{}
	return true
}

// delayedTransition set ups a timer in the background to move the VM
// identified by the given id to state after the given delay has passed.
// Uses setVMState internally to handle a safe concurrent delayed transition.
func (c *Cloud) delayedTransition(id int, state VMState, delay time.Duration) DoneChannel {
	done := make(DoneChannel)
	time.AfterFunc(delay, func() {
		if err := c.setVMState(id, state); err != nil {
			log.Println(err)
		}
		close(done) // signal we are done
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
