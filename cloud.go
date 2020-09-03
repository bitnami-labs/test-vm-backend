package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// Cloud can list VMs, and perform start/stop and remove operations on them
type Cloud struct {
	lock sync.RWMutex
	vms  VMList
}

// Cloud by default dumps itself in JSON format
func (m *Cloud) String() string {
	mJSON, err := json.Marshal(m)
	dieOnError(err, "Can't generate JSON for Cloud object: %#v", m)
	return string(mJSON)
}

// List the VMs handled under this Cloud
func (m *Cloud) List() VMList {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.vms.clone()
}

// Inspect a VM data by id (might not find it and return nil)
func (m *Cloud) Inspect(id int) VM {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.vms.lookup(id)
}

// Launch a VM by id
func (m *Cloud) Launch(id int) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	vm, err := m.vm2state(id, STARTING)
	if err != nil {
		return err
	}
	m.vms[id] = vm
	time.AfterFunc(StartDelay, func() {
		if vm, err := vm.WithState(RUNNING); err == nil {
			m.vms[id] = vm
		}
	})
	return nil
}

// Stop a VM by id
func (m *Cloud) Stop(id int) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	vm, err := m.vm2state(id, STOPPING)
	if err != nil {
		return err
	}
	m.vms[id] = vm
	time.AfterFunc(StopDelay, func() {
		if vm, err := vm.WithState(STOPPED); err == nil {
			m.vms[id] = vm
		}
	})
	return nil
}

// Delete VM by id. Idempotent, will return true when VM is actually deleted
// and false if it was already not there.
func (m *Cloud) Delete(id int) bool {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.vms.lookup(id) == emptyVM {
		return false
	}
	m.vms[id] = emptyVM
	return true
}

func (m *Cloud) vm2state(id int, state VMState) (VM, error) {
	vm := m.vms.lookup(id)
	if vm == emptyVM {
		return emptyVM, fmt.Errorf("VM with id %d not found", id)
	}
	mutatedVM, err := vm.WithState(state)
	if err != nil {
		return emptyVM, err
	}
	return mutatedVM, nil
}
