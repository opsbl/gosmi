package smi

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadModuleFromCopiedMib(t *testing.T) {
	inst, err := NewInstance("test-load")
	if err != nil {
		t.Fatalf("new instance: %v", err)
	}
	inst.SetPath(filepath.Join("..", "testdata", "mibs", "ietf"))
	name, err := inst.LoadModuleWithError("RFC1155-SMI")
	if err != nil {
		t.Fatalf("load module: %v", err)
	}
	if name != "RFC1155-SMI" {
		t.Fatalf("unexpected module name %q", name)
	}
	if !inst.IsLoaded("RFC1155-SMI") {
		t.Fatalf("module should be marked as loaded")
	}
}

func TestInstanceIsolation(t *testing.T) {
	instA, err := NewInstance("isolation-a")
	if err != nil {
		t.Fatalf("new instance A: %v", err)
	}
	instA.SetPath(filepath.Join("..", "testdata", "mibs", "instanceA"))

	instB, err := NewInstance("isolation-b")
	if err != nil {
		t.Fatalf("new instance B: %v", err)
	}
	instB.SetPath(filepath.Join("..", "testdata", "mibs", "instanceB"))

	if _, err := instA.LoadModuleWithError("SHARED-MIB"); err != nil {
		t.Fatalf("instance A load failed: %v", err)
	}
	if _, err := instB.LoadModuleWithError("SHARED-MIB"); err != nil {
		t.Fatalf("instance B load failed: %v", err)
	}

	nodeA := instA.GetNode(nil, "shared")
	if nodeA == nil {
		t.Fatalf("instance A node not found")
	}
	nodeB := instB.GetNode(nil, "shared")
	if nodeB == nil {
		t.Fatalf("instance B node not found")
	}

	oidA := nodeA.Oid.String()
	oidB := nodeB.Oid.String()

	if oidA == oidB {
		t.Fatalf("instances share state: got identical OID %s", oidA)
	}
	if oidA != "1.10" {
		t.Fatalf("instance A OID mismatch: %s", oidA)
	}
	if oidB != "1.20" {
		t.Fatalf("instance B OID mismatch: %s", oidB)
	}
}

func TestCircularImportProtection(t *testing.T) {
	inst, err := NewInstance("cycle")
	if err != nil {
		t.Fatalf("new instance: %v", err)
	}
	inst.SetPath(filepath.Join("..", "testdata", "mibs", "cycle"))

	name, err := inst.LoadModuleWithError("CYCLE-A-MIB")
	t.Logf("loaded name=%q err=%v loadedB=%v", name, err, inst.IsLoaded("CYCLE-B-MIB"))
	if err == nil {
		t.Fatalf("expected circular import error")
	} else if !strings.Contains(strings.ToLower(err.Error()), "circular") {
		t.Fatalf("unexpected error: %v", err)
	}
}
