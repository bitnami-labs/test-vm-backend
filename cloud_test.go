// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package main

import (
	"fmt"
	"testing"
	"time"
)

const (
	// GoodID for happy tests
	GoodID = 1

	// BadID for bad search tests
	BadID = 10000
)

const (
	expectedNotFoundMsgFmt = "not found VM with id %d"
)

func NewDefaultCloud() Cloud {
	return Cloud{vms: defaultVMs.clone()}
}

// copyInState gets a copy of the VM identified by id from cloud,
// but set to an arbitrary given state.
func copyInState(cloud *Cloud, id int, state VMState) (VM, error) {
	vm, found := cloud.Inspect(id)
	if !found {
		return VM{}, fmt.Errorf("Could not find id=%d at %#v", id, cloud)
	}
	vm.State = state
	return vm, nil
}

// forceState even if it is not a legal transition
func forceState(cloud *Cloud, id int, state VMState) error {
	vm, err := copyInState(cloud, id, state)
	if err != nil {
		return err
	}
	cloud.vms[id] = vm
	return nil
}

// shrinkTime sets up shorter delays so that test can go faster
func shrinkTime() {
	StartDelay = 10 * time.Millisecond
	StopDelay = 5 * time.Millisecond
}

// waitDone waits for a done channel to finish or a timeout to occur
func waitDone(done chan struct{}, timeout time.Duration) error {
	timeoutChannel := time.After(timeout)
	select {
	case <-done:
		return nil
	case <-timeoutChannel:
		return fmt.Errorf("Timeout expired (%v) waiting on done channel", timeout)
	}
}

func TestList(t *testing.T) {
	c := NewDefaultCloud()
	want := defaultVMs.String()
	got := c.List().String()
	if got != want {
		t.Fatalf("got: %s, wanted: %s", got, want)
	}
}

func TestInspect(t *testing.T) {
	c := NewDefaultCloud()
	want := defaultVMs[GoodID]
	got, _ := c.Inspect(GoodID)
	if got != want {
		t.Fatalf("got: %s, want: %s", got, want)
	}
}

func TestBadInspect(t *testing.T) {
	c := NewDefaultCloud()
	if _, found := c.Inspect(BadID); found == true {
		t.Fatalf("found: %v, want: false", found)
	}
}

func TestLaunch(t *testing.T) {
	shrinkTime()
	c := NewDefaultCloud()
	// Test 1st transition
	want, err := copyInState(&c, GoodID, STARTING)
	if err != nil {
		t.Fatal(err)
	}
	done, err := c.Launch(GoodID)
	if err != nil {
		t.Fatalf("Failed to Launch VM %d: %v", GoodID, err)
	}
	if got, _ := c.Inspect(GoodID); got != want {
		t.Fatalf("got: %s, want: %s", got, want)
	}
	// Wait and test 2nd transition
	if err := waitDone(done, 10*StartDelay); err != nil {
		t.Fatal(err)
	}
	want2, err := copyInState(&c, GoodID, RUNNING)
	if err != nil {
		t.Fatal(err)
	}
	if got2, _ := c.Inspect(GoodID); got2 != want2 {
		t.Fatalf("got %q, want: %q", got2, want2)
	}
}

func TestBadVMLaunch(t *testing.T) {
	c := NewDefaultCloud()
	want := fmt.Sprintf(expectedNotFoundMsgFmt, BadID)
	if _, got := c.Launch(BadID); got == nil || got.Error() != want {
		t.Fatalf("got: %q, want: %q", got, want)
	}
}

func TestBadStateLaunch(t *testing.T) {
	c := NewDefaultCloud()
	var badState VMState = RUNNING
	if err := forceState(&c, GoodID, badState); err != nil {
		t.Fatal(err)
	}
	want := fmt.Sprintf("illegal transition from %q to %q", badState, STARTING)
	if _, got := c.Launch(GoodID); got == nil || got.Error() != want {
		t.Fatalf("got: %q, want: %q", got, want)
	}
}

func TestStop(t *testing.T) {
	shrinkTime()
	c := NewDefaultCloud()
	forceState(&c, GoodID, RUNNING)
	// Test 1st transition
	want, err := copyInState(&c, GoodID, STOPPING)
	if err != nil {
		t.Fatal(err)
	}
	done, err := c.Stop(GoodID)
	if err != nil {
		t.Fatalf("Failed to Stop VM %d: %v", GoodID, err)
	}
	if got, _ := c.Inspect(GoodID); got != want {
		t.Fatalf("got: %v, want: %v", got, want)
	}
	// Wait and test 2nd transition
	if err := waitDone(done, 10*StopDelay); err != nil {
		t.Fatal(err)
	}
	want2, err := copyInState(&c, GoodID, STOPPED)
	if err != nil {
		t.Fatal(err)
	}
	if got2, _ := c.Inspect(GoodID); got2 != want2 {
		t.Fatalf("got: %v, want: %v", got2, want2)
	}
}

func TestBadVMStop(t *testing.T) {
	c := NewDefaultCloud()
	c.setVMState(GoodID, RUNNING)
	want := fmt.Sprintf(expectedNotFoundMsgFmt, BadID)
	if _, got := c.Stop(BadID); got == nil || got.Error() != want {
		t.Fatalf("got: %v, want: %v", got, want)
	}
}

func TestBadStateStop(t *testing.T) {
	c := NewDefaultCloud()
	// No extra setup needed: initial state Stopped is already bad for stopping
	want := fmt.Sprintf("illegal transition from %q to %q", STOPPED, STOPPING)
	if _, got := c.Stop(GoodID); got == nil || got.Error() != want {
		t.Fatalf("got: %v, want: %v", got, want)
	}
}

func TestDelete(t *testing.T) {
	c := NewDefaultCloud()
	if err := c.Delete(GoodID); err != nil {
		t.Fatal(err)
	}
}

func TestBadDelete(t *testing.T) {
	c := NewDefaultCloud()
	want := fmt.Sprintf("delete error: not found VM %d", BadID)
	if got := c.Delete(BadID); got == nil || got.Error() != want {
		t.Fatalf("got: %q, want: %q", got, want)
	}
}

func TestReDelete(t *testing.T) {
	c := NewDefaultCloud()
	if err := c.Delete(GoodID); err != nil {
		t.Fatalf("Unexpected deletion error: %v", err)
	}
	want := fmt.Sprintf("delete error: not found VM %d", GoodID)
	if got := c.Delete(GoodID); got == nil || got.Error() != want {
		t.Fatalf("got: %q, want: %q", got, want)
	}
}

func TestBadStateDelete(t *testing.T) {
	c := NewDefaultCloud()
	badState := RUNNING // not allowed to delete in this state
	forceState(&c, GoodID, badState)
	want := fmt.Sprintf("delete error: VM %d must be in state %v for deletion but it is %v", GoodID, STOPPED, badState)
	if got := c.Delete(GoodID); got == nil || got.Error() != want {
		t.Fatalf("got: %q, want: %q", got, want)
	}
}
