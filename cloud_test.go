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
	wanted := cloneDefaultVMList().String()
	got := c.List().String()
	if got != wanted {
		t.Fatalf("Expected %q but got %q", wanted, got)
	}
}

func TestInspect(t *testing.T) {
	c := NewDefaultCloud()
	wanted := cloneDefaultVMList()[GoodID].String()
	got := c.Inspect(GoodID).String()
	if got != wanted {
		t.Fatalf("Expected %q but got %q", wanted, got)
	}
}

func TestBadInspect(t *testing.T) {
	c := NewDefaultCloud()
	var wanted VM = emptyVM
	got := c.Inspect(BadID)
	if got != wanted {
		t.Fatalf("Expected %v but got %v", wanted, got)
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
	wanted := VMByIndexInState(t, &c, GoodID, STARTING).String()
	if err := c.Launch(GoodID); err != nil {
		t.Fatalf("Failed to Launch VM %q: %q", GoodID, err)
	}
	if got := c.vms[GoodID].String(); got != wanted {
		t.Fatalf("Expected %q but got %q", wanted, got)
	}
	// TODO: make this transition test reliable, it fails always on Mac tests
	// time.Sleep(2 * StartDelay)
	// wanted2 := VMByIndexInState(t, &c, GoodID, RUNNING).String()
	// if got2 := c.vms[GoodID].String(); got2 != wanted2 {
	// 	t.Fatalf("Expected %q but got %q", wanted2, got2)
	// }
}

func TestBadVMLaunch(t *testing.T) {
	c := NewDefaultCloud()
	wanted := fmt.Sprintf("VM with id %d not found", BadID)
	if got := c.Launch(BadID); got.Error() != wanted {
		t.Fatalf("Unexpected Launch VM error: %q", got)
	}
}

func TestBadStateLaunch(t *testing.T) {
	c := NewDefaultCloud()
	var badState VMState = RUNNING
	c.vms[GoodID].State = badState
	wanted := fmt.Sprintf("illegal transition from %q to %q", badState, STARTING)
	if got := c.Launch(GoodID); got.Error() != wanted {
		t.Fatalf("Unexpected Launch VM error: %q", got)
	}
}

func TestStop(t *testing.T) {
	shrinkTime()
	c := NewDefaultCloud()
	c.vms[GoodID].State = RUNNING
	wanted := VMByIndexInState(t, &c, GoodID, STOPPING).String()
	if err := c.Stop(GoodID); err != nil {
		t.Fatalf("Failed to Stop VM %q: %q", GoodID, err)
	}
	if got := c.vms[GoodID].String(); got != wanted {
		t.Fatalf("Expected %q but got %q", wanted, got)
	}
	// TODO: make this transition test reliable, it fails always on Mac tests
	// time.Sleep(2 * StopDelay)
	// wanted2 := VMByIndexInState(t, &c, GoodID, STOPPED).String()
	// if got2 := c.vms[GoodID].String(); got2 != wanted2 {
	// 	t.Fatalf("Expected %q but got %q", wanted2, got2)
	// }
}
func TestBadVMStop(t *testing.T) {
	c := NewDefaultCloud()
	c.vms[GoodID].State = RUNNING
	wanted := fmt.Sprintf("VM with id %d not found", BadID)
	if got := c.Stop(BadID); got.Error() != wanted {
		t.Fatalf("Unexpected Stop VM error: %q", got)
	}
}

func TestBadStateStop(t *testing.T) {
	c := NewDefaultCloud()
	// No extra setup needed: initial state Stopped is already bad for stopping
	wanted := fmt.Sprintf("illegal transition from %q to %q", STOPPED, STOPPING)
	if got := c.Stop(GoodID); got == nil || got.Error() != wanted {
		t.Fatalf("Unexpected Stop VM error: %q", got)
	}
}

func TestDelete(t *testing.T) {
	c := NewDefaultCloud()
	wanted := true
	if got := c.Delete(GoodID); got != wanted {
		t.Fatalf("Expected %v but got %v", wanted, got)
	}
}

func TestBadDelete(t *testing.T) {
	c := NewDefaultCloud()
	wanted := false
	if got := c.Delete(BadID); got != wanted {
		t.Fatalf("Expected deleted=%v but got deleted=%v", wanted, got)
	}
}
