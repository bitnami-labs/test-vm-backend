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

func (vm VM) isValid() bool {
	return vm != VM{}
}

// AllowedTransition lists allowed state transitions
var AllowedTransition = map[VMState]VMState{
	STOPPED:  STARTING,
	STARTING: RUNNING,
	RUNNING:  STOPPING,
	STOPPING: STOPPED,
}

// VMList defines a list of VMs with attached methods
type VMList []VM

// clone returns a deep clone of the list, useful for snapshots
func (vms VMList) clone() VMList {
	cloneList := make(VMList, len(vms))
	copy(cloneList, vms)
	return cloneList
}

// hashize turns VMList into a map skipping empty entries (if any)
func (vms VMList) hashize() map[int]VM {
	hash := make(map[int]VM)
	for index, vm := range vms {
		if vm.isValid() {
			hash[index] = vm
		}
	}
	return hash
}

// lookup a VM returning the empty VM if out of bounds or not found
func (vms VMList) lookup(id int) VM {
	if id < 0 || id > (len(vms)-1) {
		return VM{}
	}
	return vms[id]
}

// String in VMList by default dumps itself in JSON format skipping empty entries
func (vms VMList) String() string {
	vmJSON, err := json.Marshal(vms.hashize())
	dieOnError(err, "Can't generate JSON for VM object %#v", vms)
	return string(vmJSON)
}

var defaultVMList = VMList{
	VM{
		VCPUS:   1,       // Number of processors
		Clock:   1500,    // Frequency of 1 processor, expressed in MHz (Megahertz)
		RAM:     4096,    // Amount of internal memory, expressed in MB (Megabytes)
		Storage: 128,     // Amount of internal space available for storage, expressed in GB (Gigabytes)
		Network: 1000,    // Speed of the networking device, expressed in Gb/s (Gigabits per second)
		State:   STOPPED, // Value from within the set [Running, Stopped, Starting, Stopping]
	},
	VM{
		VCPUS:   4,
		Clock:   3600,
		RAM:     32768,
		Storage: 512,
		Network: 10000,
		State:   STOPPED,
	},
	VM{
		VCPUS:   2,
		Clock:   2200,
		RAM:     8192,
		Storage: 256,
		Network: 1000,
		State:   STOPPED,
	},
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
