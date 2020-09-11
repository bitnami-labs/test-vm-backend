// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// VMState represents the current state of a VM
type VMState string

const (
	// STOPPED VM is stopped at this state it can be removed
	STOPPED VMState = "Stopped"

	// STARTING VM is transitioning from Stopped to Starting in about 10 minutes
	STARTING VMState = "Starting"

	// RUNNING VM is online
	RUNNING VMState = "Running"

	// STOPPING VM is transitioning from Running to Stopped
	STOPPING VMState = "Stopping"
)

const (
	// DefaultStartDelay Start VM process simulated delay
	DefaultStartDelay = 10 * time.Second

	// DefaultStopDelay Stop VM process simulated delay
	DefaultStopDelay = 5 * time.Second
)

var (
	// StartDelay for launch operations (not a constant so unit test can change it)
	StartDelay = DefaultStartDelay

	// StopDelay for stop operations (not a constant so unit test can change it)
	StopDelay = DefaultStopDelay
)

// VMsJSON filename where to store initial VMs state list
const VMsJSON = "vms.json"

func dieOnError(err error, format string, args ...interface{}) {
	if err != nil {
		log.Fatalf("%s: %v\n", fmt.Sprintf(format, args...), err)
	}
}

// VM is a Virtual Machine
type VM struct {
	VCPUS   int     `json:"vcpus,omitempty"`   // Number of processors
	Clock   float32 `json:"clock,omitempty"`   // Frequency of 1 processor, in MHz (Megahertz)
	RAM     int     `json:"ram,omitempty"`     // Amount of internal memory, in MB (Megabytes)
	Storage int     `json:"storage,omitempty"` // Amount of persistent storage, in GB (Gigabytes)
	Network int     `json:"network,omitempty"` // Network device speed in Gb/s (Gigabits per second)
	State   VMState `json:"state,omitempty"`   // Value within [Running, Stopped, Starting, Stopping]
}

// VM by default dumps itself in JSON format
func (vm VM) String() string {
	vmJSON, err := json.Marshal(vm)
	dieOnError(err, "Can't generate JSON for VM object %#v", vm)
	return string(vmJSON)
}

// AllowedTransition lists allowed state transitions
var AllowedTransition = map[VMState]VMState{
	STOPPED:  STARTING,
	STARTING: RUNNING,
	RUNNING:  STOPPING,
	STOPPING: STOPPED,
}

// WithState returns a VM on the requested end state or an error,
// if the transition was illegal
func (vm VM) WithState(state VMState) (VM, error) {
	if state == vm.State {
		return vm, nil // NOP
	}
	if AllowedTransition[vm.State] != state {
		return VM{}, fmt.Errorf("illegal transition from %q to %q", vm.State, state)
	}
	vm.State = state
	return vm, nil
}

// VMs defines a map of VMs with attached methods
type VMs map[int]VM

// clone returns a deep clone of the list, useful for snapshots
func (vms VMs) clone() VMs {
	cloneList := make(VMs, len(vms))
	for k, v := range vms {
		cloneList[k] = v
	}
	return cloneList
}

// String in VMs by default dumps itself in JSON format skipping empty entries
func (vms VMs) String() string {
	vmJSON, err := json.Marshal(vms)
	dieOnError(err, "Can't generate JSON for VM object %#v", vms)
	return string(vmJSON)
}

var defaultVMs = VMs{
	0: {
		VCPUS:   1,       // Number of processors
		Clock:   1500,    // Frequency of 1 processor, expressed in MHz (Megahertz)
		RAM:     4096,    // Amount of internal memory, expressed in MB (Megabytes)
		Storage: 128,     // Amount of internal space available for storage, expressed in GB (Gigabytes)
		Network: 1000,    // Speed of the networking device, expressed in Gb/s (Gigabits per second)
		State:   STOPPED, // Value from within the set [Running, Stopped, Starting, Stopping]
	},
	1: {
		VCPUS:   4,
		Clock:   3600,
		RAM:     32768,
		Storage: 512,
		Network: 10000,
		State:   STOPPED,
	},
	2: {
		VCPUS:   2,
		Clock:   2200,
		RAM:     8192,
		Storage: 256,
		Network: 1000,
		State:   STOPPED,
	},
}
