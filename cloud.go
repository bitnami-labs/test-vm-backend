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
	c.lock.Lock()
	defer c.lock.Unlock()

	vm, err := c.vm2state(id, STARTING)
	if err != nil {
		return err
	}
	c.vms[id] = vm
	time.AfterFunc(StartDelay, func() {
		if vm, err := vm.WithState(RUNNING); err == nil {
			c.vms[id] = vm
		}
	})
	return nil
}

// Stop a VM by id
func (c *Cloud) Stop(id int) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	vm, err := c.vm2state(id, STOPPING)
	if err != nil {
		return err
	}
	c.vms[id] = vm
	time.AfterFunc(StopDelay, func() {
		if vm, err := vm.WithState(STOPPED); err == nil {
			c.vms[id] = vm
		}
	})
	return nil
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

// vm2state sets the VM identified by the given id to the given state.
// Might fail if the VM is not found or transition requested is illegal.
func (c *Cloud) vm2state(id int, state VMState) (VM, error) {
	vm := c.vms.lookup(id)
	if !vm.isValid() {
		return VM{}, fmt.Errorf("not found VM with id %d", id)
	}
	mutatedVM, err := vm.WithState(state)
	if err != nil {
		return VM{}, err
	}
	return mutatedVM, nil
}
