// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

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

var withStateHappyCases = []struct {
	vm    *VM
	state VMState
	want  *VM
}{
	{vm: VMInState(STOPPED), state: STARTING, want: VMInState(STARTING)},
	{vm: VMInState(STARTING), state: RUNNING, want: VMInState(RUNNING)},
	{vm: VMInState(RUNNING), state: STOPPING, want: VMInState(STOPPING)},
	{vm: VMInState(STOPPING), state: STOPPED, want: VMInState(STOPPED)},
	{vm: VMInState(STOPPED), state: STOPPED, want: VMInState(STOPPED)},
	{vm: VMInState(STARTING), state: STARTING, want: VMInState(STARTING)},
	{vm: VMInState(RUNNING), state: RUNNING, want: VMInState(RUNNING)},
	{vm: VMInState(STOPPING), state: STOPPING, want: VMInState(STOPPING)},
}

func TestWithStateHappyCases(t *testing.T) {
	for _, tc := range withStateHappyCases {
		got, err := tc.vm.WithState(tc.state)
		if err != nil {
			t.Fatalf("Unexpected error in happy case %v: %v", tc, err)
		}
		if got != *tc.want {
			t.Fatalf("got: %v, want %v", got, *tc.want)
		}
	}
}

var withStateErrors = []struct {
	vm    *VM
	state VMState
	want  string
}{
	{vm: VMInState(STOPPED), state: RUNNING,
		want: `illegal transition from "Stopped" to "Running"`},
	{vm: VMInState(STOPPED), state: STOPPING,
		want: `illegal transition from "Stopped" to "Stopping"`},
	{vm: VMInState(RUNNING), state: STOPPED,
		want: `illegal transition from "Running" to "Stopped"`},
	{vm: VMInState(RUNNING), state: STARTING,
		want: `illegal transition from "Running" to "Starting"`},
	{vm: VMInState(STARTING), state: STOPPED,
		want: `illegal transition from "Starting" to "Stopped"`},
	{vm: VMInState(STARTING), state: STOPPING,
		want: `illegal transition from "Starting" to "Stopping"`},
}

func TestWithStateErrors(t *testing.T) {
	for _, tc := range withStateErrors {
		vm, got := tc.vm.WithState(tc.state)
		if (vm != VM{}) {
			t.Fatalf("Unexpected VM valid value in error case %v: %v", tc, vm)
		}
		if got.Error() != tc.want {
			t.Fatalf("got: %q, want %q", got, tc.want)
		}
	}
}
