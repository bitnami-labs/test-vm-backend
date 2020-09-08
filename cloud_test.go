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

func cloneDefaultVMList() VMList {
	defaultVMListClone := make(VMList, len(defaultVMList))
	copy(defaultVMListClone, defaultVMList)
	return defaultVMListClone
}

func NewDefaultCloud() Cloud {
	return Cloud{vms: cloneDefaultVMList()}
}

func TestList(t *testing.T) {
	c := NewDefaultCloud()
	want := cloneDefaultVMList().String()
	got := c.List().String()
	if got != want {
		t.Fatalf("got: %q, wanted: %q", got, want)
	}
}

func TestInspect(t *testing.T) {
	c := NewDefaultCloud()
	want := cloneDefaultVMList()[GoodID].String()
	got := c.Inspect(GoodID).String()
	if got != want {
		t.Fatalf("got: %q, want: %q", got, want)
	}
}

func TestBadInspect(t *testing.T) {
	c := NewDefaultCloud()
	want := false
	got := c.Inspect(BadID).isValid()
	if got != want {
		t.Fatalf("got: %v, want: %v", got, want)
	}
}

func VMByIndexInState(t *testing.T, c *Cloud, id int, state VMState) VM {
	vm, err := c.vms[id].WithState(state)
	if err != nil {
		t.Fatalf("Failed to set VM state %v: %v", vm, err)
	}
	return vm
}

func shrinkTime() {
	StartDelay = 10 * time.Millisecond
	StopDelay = 5 * time.Millisecond
}

func TestLaunch(t *testing.T) {
	shrinkTime()
	c := NewDefaultCloud()
	want := VMByIndexInState(t, &c, GoodID, STARTING).String()
	if err := c.Launch(GoodID); err != nil {
		t.Fatalf("Failed to Launch VM %q: %q", GoodID, err)
	}
	if got := c.vms[GoodID].String(); got != want {
		t.Fatalf("got: %q, want: %q", got, want)
	}
	// TODO: make this transition test reliable, it fails always on Mac tests
	// time.Sleep(2 * StartDelay)
	// want2 := VMByIndexInState(t, &c, GoodID, RUNNING).String()
	// if got2 := c.vms[GoodID].String(); got2 != wanted2 {
	// 	t.Fatalf("got %q, want: %q", got2, want2)
	// }
}

func TestBadVMLaunch(t *testing.T) {
	c := NewDefaultCloud()
	want := fmt.Sprintf(expectedNotFoundMsgFmt, BadID)
	if got := c.Launch(BadID).Error(); got != want {
		t.Fatalf("got: %q, want: %q", got, want)
	}
}

func TestBadStateLaunch(t *testing.T) {
	c := NewDefaultCloud()
	var badState VMState = RUNNING
	c.vms[GoodID].State = badState
	want := fmt.Sprintf("illegal transition from %q to %q", badState, STARTING)
	if got := c.Launch(GoodID).Error(); got != want {
		t.Fatalf("got: %q, want: %q", got, want)
	}
}

func TestStop(t *testing.T) {
	shrinkTime()
	c := NewDefaultCloud()
	c.vms[GoodID].State = RUNNING
	want := VMByIndexInState(t, &c, GoodID, STOPPING).String()
	if err := c.Stop(GoodID); err != nil {
		t.Fatalf("Failed to Stop VM %q: %q", GoodID, err)
	}
	if got := c.vms[GoodID].String(); got != want {
		t.Fatalf("got: %q, want: %q", got, want)
	}
	// TODO: make this transition test reliable, it fails always on Mac tests
	// time.Sleep(2 * StopDelay)
	// want2 := VMByIndexInState(t, &c, GoodID, STOPPED).String()
	// if got2 := c.vms[GoodID].String(); got2 != want2 {
	// 	t.Fatalf("got: %q, want: %q", got2, want2)
	// }
}
func TestBadVMStop(t *testing.T) {
	c := NewDefaultCloud()
	c.vms[GoodID].State = RUNNING
	want := fmt.Sprintf(expectedNotFoundMsgFmt, BadID)
	if got := c.Stop(BadID).Error(); got != want {
		t.Fatalf("got: %q, want: %q", got, want)
	}
}

func TestBadStateStop(t *testing.T) {
	c := NewDefaultCloud()
	// No extra setup needed: initial state Stopped is already bad for stopping
	want := fmt.Sprintf("illegal transition from %q to %q", STOPPED, STOPPING)
	if got := c.Stop(GoodID); got == nil || got.Error() != want {
		t.Fatalf("got: %q, want: %q", got, want)
	}
}

func TestDelete(t *testing.T) {
	c := NewDefaultCloud()
	want := true
	if got := c.Delete(GoodID); got != want {
		t.Fatalf("got: %v, want: %v", got, want)
	}
}

func TestBadDelete(t *testing.T) {
	c := NewDefaultCloud()
	want := false
	if got := c.Delete(BadID); got != want {
		t.Fatalf("got: %v, want: %v", got, want)
	}
}
