package main

import (
	"testing"
)

func VMInState(state VMState) *VM {
	return &VM{
		VCPUS:   1,
		Clock:   1500,
		RAM:     4096,
		Storage: 128,
		Network: 1000,
		State:   state,
	}
}

var withStateHappyTestCases = []struct {
	vm     *VM
	state  VMState
	wanted *VM
}{
	{vm: VMInState(STOPPED), state: STARTING, wanted: VMInState(STARTING)},
	{vm: VMInState(STARTING), state: RUNNING, wanted: VMInState(RUNNING)},
	{vm: VMInState(RUNNING), state: STOPPING, wanted: VMInState(STOPPING)},
	{vm: VMInState(STOPPING), state: STOPPED, wanted: VMInState(STOPPED)},
	{vm: VMInState(STOPPED), state: STOPPED, wanted: VMInState(STOPPED)},
	{vm: VMInState(STARTING), state: STARTING, wanted: VMInState(STARTING)},
	{vm: VMInState(RUNNING), state: RUNNING, wanted: VMInState(RUNNING)},
	{vm: VMInState(STOPPING), state: STOPPING, wanted: VMInState(STOPPING)},
}

func TestWithStateHappyCases(t *testing.T) {
	for _, testcase := range withStateHappyTestCases {
		got, err := testcase.vm.WithState(testcase.state)
		if err != nil {
			t.Fatalf("Unexpected error in happy case %v: %v", testcase, err)
		}
		if got != *testcase.wanted {
			t.Fatalf("Expected %v but got %v", *testcase.wanted, got)
		}
	}
}

var withStateErrors = []struct {
	vm     *VM
	state  VMState
	wanted string
}{
	{vm: VMInState(STOPPED), state: RUNNING,
		wanted: `Illegal transition from "Stopped" to "Running"`},
	{vm: VMInState(STOPPED), state: STOPPING,
		wanted: `Illegal transition from "Stopped" to "Stopping"`},
	{vm: VMInState(RUNNING), state: STOPPED,
		wanted: `Illegal transition from "Running" to "Stopped"`},
	{vm: VMInState(RUNNING), state: STARTING,
		wanted: `Illegal transition from "Running" to "Starting"`},
	{vm: VMInState(STARTING), state: STOPPED,
		wanted: `Illegal transition from "Starting" to "Stopped"`},
	{vm: VMInState(STARTING), state: STOPPING,
		wanted: `Illegal transition from "Starting" to "Stopping"`},
}

func TestWithStateErrors(t *testing.T) {
	emptyVM := VM{}
	for _, testcase := range withStateErrors {
		vm, got := testcase.vm.WithState(testcase.state)
		if vm != emptyVM {
			t.Fatalf("Unexpected VM non empty value in error case %v: %v", testcase, vm)
		}
		if got.Error() != testcase.wanted {
			t.Fatalf("Expected %q but got %q", testcase.wanted, got)
		}
	}
}
