package module

import (
	"bytes"
	"reflect"
	"testing"
)

func TestModuleHCL(t *testing.T) {
	hclInput := `
import {
  module = [
    "base-module",
    "some-other-module",
  ]
}

resource "pacman" {
  name = "openssh"
}

resource "pacman" {
  name = "tmux"
  state = "present"
}
`
	hclModule, err := Load("main", bytes.NewBufferString(hclInput))
	if err != nil {
		t.Fatal(err)
	}

	wantName := "main"
	wantNumImports := 2
	wantImportNames := []string{"base-module", "some-other-module"}
	wantNumResources := 2

	if wantName != hclModule.Name {
		t.Errorf("want module name %q, got name %q", wantName, hclModule.Name)
	}

	if wantNumImports != len(hclModule.ModuleImport.Module) {
		t.Errorf("want %d imports, got %d imports", wantNumImports, len(hclModule.ModuleImport.Module))
	}

	if !reflect.DeepEqual(wantImportNames, hclModule.ModuleImport.Module) {
		t.Errorf("want %q import names, got %q names", wantImportNames, hclModule.ModuleImport.Module)
	}

	if wantNumResources != len(hclModule.Resources) {
		t.Errorf("want %d resources, got %d resources", wantNumResources, len(hclModule.Resources))
	}
}

func TestModuleJSON(t *testing.T) {
	jsonInput := `
{
  "import": {
    "module": [
      "base-module",
      "some-other-module"
    ],
  },
  "resource": [
    {
      "pacman": {
        "name": "openssh",
        "state": "present"
      }
    },
    {
      "pacman": {
        "name": "tmux",
        "state": "present",
      }
    }
  ]
}
`
	jsonModule, err := Load("main", bytes.NewBufferString(jsonInput))
	if err != nil {
		t.Fatal(err)
	}

	wantName := "main"
	wantNumImports := 2
	wantImportNames := []string{"base-module", "some-other-module"}
	wantNumResources := 2

	if wantName != jsonModule.Name {
		t.Errorf("want module name %q, got name %q", wantName, jsonModule.Name)
	}

	if wantNumImports != len(jsonModule.ModuleImport.Module) {
		t.Errorf("want %d imports, got %d imports", wantNumImports, len(jsonModule.ModuleImport.Module))
	}

	if !reflect.DeepEqual(wantImportNames, jsonModule.ModuleImport.Module) {
		t.Errorf("want %q import names, got %q names", wantImportNames, jsonModule.ModuleImport.Module)
	}

	if wantNumResources != len(jsonModule.Resources) {
		t.Errorf("want %d resources, got %d resources", wantNumResources, len(jsonModule.Resources))
	}
}
