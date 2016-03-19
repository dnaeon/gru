package module

import (
	"bytes"
	"reflect"
	"testing"
)

func TestModuleValidHCL(t *testing.T) {
	const hclModule = `
import {
  module = [
    "base-module",
    "some-other-module",
  ]
}

resource "pacman" {
  name = "openssh"
}

resource "service" {
  name = "sshd"
  state = "running"
  enable = true
  want = [
    "pacman[openssh]",
  ]
}
`
	m, err := Load("main", bytes.NewBufferString(hclModule))
	if err != nil {
		t.Fatal(err)
	}

	wantName := "main"
	wantNumImports := 2
	wantImportNames := []string{"base-module", "some-other-module"}
	wantNumResources := 2

	if wantName != m.Name {
		t.Errorf("want module name %q, got name %q", wantName, m.Name)
	}

	if wantNumImports != len(m.ModuleImport.Module) {
		t.Errorf("want %d imports, got %d imports", wantNumImports, len(m.ModuleImport.Module))
	}

	if !reflect.DeepEqual(wantImportNames, m.ModuleImport.Module) {
		t.Errorf("want %q import names, got %q names", wantImportNames, m.ModuleImport.Module)
	}

	if wantNumResources != len(m.Resources) {
		t.Errorf("want %d resources, got %d resources", wantNumResources, len(m.Resources))
	}
}
