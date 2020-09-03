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
	m := NewDefaultCloud()
	wanted := cloneDefaultVMList().String()
	got := m.List().String()
	if got != wanted {
		t.Fatalf("Expected %q but got %q", wanted, got)
	}
}

func TestInspect(t *testing.T) {
	m := NewDefaultCloud()
	wanted := cloneDefaultVMList()[GoodID].String()
	got := m.Inspect(GoodID).String()
	if got != wanted {
		t.Fatalf("Expected %q but got %q", wanted, got)
	}
}

func TestBadInspect(t *testing.T) {
	m := NewDefaultCloud()
	var wanted VM = emptyVM
	got := m.Inspect(BadID)
	if got != wanted {
		t.Fatalf("Expected %v but got %v", wanted, got)
	}
}

func VMByIndexInState(t *testing.T, m *Cloud, id int, state VMState) VM {
	vm, err := m.vms[id].WithState(state)
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
	m := NewDefaultCloud()
	wanted := VMByIndexInState(t, &m, GoodID, STARTING).String()
	if err := m.Launch(GoodID); err != nil {
		t.Fatalf("Failed to Launch VM %q: %q", GoodID, err)
	}
	if got := m.vms[GoodID].String(); got != wanted {
		t.Fatalf("Expected %q but got %q", wanted, got)
	}
	time.Sleep(2 * StartDelay)
	wanted2 := VMByIndexInState(t, &m, GoodID, RUNNING).String()
	if got2 := m.vms[GoodID].String(); got2 != wanted2 {
		t.Fatalf("Expected %q but got %q", wanted2, got2)
	}
}

func TestBadVMLaunch(t *testing.T) {
	m := NewDefaultCloud()
	wanted := fmt.Sprintf("VM with id %d not found", BadID)
	if got := m.Launch(BadID); got.Error() != wanted {
		t.Fatalf("Unexpected Launch VM error: %q", got)
	}
}

func TestBadStateLaunch(t *testing.T) {
	m := NewDefaultCloud()
	var badState VMState = RUNNING
	m.vms[GoodID].State = badState
	wanted := fmt.Sprintf("Illegal transition from %q to %q", badState, STARTING)
	if got := m.Launch(GoodID); got.Error() != wanted {
		t.Fatalf("Unexpected Launch VM error: %q", got)
	}
}

func TestStop(t *testing.T) {
	shrinkTime()
	m := NewDefaultCloud()
	m.vms[GoodID].State = RUNNING
	wanted := VMByIndexInState(t, &m, GoodID, STOPPING).String()
	if err := m.Stop(GoodID); err != nil {
		t.Fatalf("Failed to Stop VM %q: %q", GoodID, err)
	}
	if got := m.vms[GoodID].String(); got != wanted {
		t.Fatalf("Expected %q but got %q", wanted, got)
	}
	time.Sleep(2 * StopDelay)
	wanted2 := VMByIndexInState(t, &m, GoodID, STOPPED).String()
	if got2 := m.vms[GoodID].String(); got2 != wanted2 {
		t.Fatalf("Expected %q but got %q", wanted2, got2)
	}
}
func TestBadVMStop(t *testing.T) {
	m := NewDefaultCloud()
	m.vms[GoodID].State = RUNNING
	wanted := fmt.Sprintf("VM with id %d not found", BadID)
	if got := m.Stop(BadID); got.Error() != wanted {
		t.Fatalf("Unexpected Stop VM error: %q", got)
	}
}

func TestBadStateStop(t *testing.T) {
	m := NewDefaultCloud()
	// No extra setup needed: initial state Stopped is already bad for stopping
	wanted := fmt.Sprintf("Illegal transition from %q to %q", STOPPED, STOPPING)
	if got := m.Stop(GoodID); got == nil || got.Error() != wanted {
		t.Fatalf("Unexpected Stop VM error: %q", got)
	}
}

func TestDelete(t *testing.T) {
	m := NewDefaultCloud()
	wanted := true
	if got := m.Delete(GoodID); got != wanted {
		t.Fatalf("Expected %v but got %v", wanted, got)
	}
}

func TestBadDelete(t *testing.T) {
	m := NewDefaultCloud()
	wanted := false
	if got := m.Delete(BadID); got != wanted {
		t.Fatalf("Expected deleted=%v but got deleted=%v", wanted, got)
	}
}
