package module

import (
	"bytes"
	"reflect"
	"testing"
)

func TestModuleHCL(t *testing.T) {
	hclInput := `
import {
  name = "base-module"
}

import {
  name = "some-other-module"
}

pacman "openssh" {
  state = "present"
}

pacman "tmux" {
  state = "present"
}
`
	hclModule, err := Load("main", bytes.NewBufferString(hclInput))
	if err != nil {
		t.Fatal(err)
	}

	wantName := "main"
	wantImports := []Import{
		Import{Name: "base-module"},
		Import{Name: "some-other-module"},
	}
	wantNumResources := 2

	if wantName != hclModule.Name {
		t.Errorf("want module name %q, got name %q", wantName, hclModule.Name)
	}

	if !reflect.DeepEqual(wantImports, hclModule.Imports) {
		t.Errorf("want %q imports, got %q imports", wantImports, hclModule.Imports)
	}

	if wantNumResources != len(hclModule.Resources) {
		t.Errorf("want %d resources, got %d resources", wantNumResources, len(hclModule.Resources))
	}
}

func TestModuleJSON(t *testing.T) {
	jsonInput := `
{
  "import": [
    {
      "name": "base-module"
    },
    {
      "name": "some-other-module"
    }
  ],
  "pacman": [
    {
      "openssh": {
        "name": "openssh",
        "state": "present"
      }
    },
    {
      "valgrind": {
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
	wantImports := []Import{
		Import{Name: "base-module"},
		Import{Name: "some-other-module"},
	}
	wantNumResources := 2

	if wantName != jsonModule.Name {
		t.Errorf("want module name %q, got name %q", wantName, jsonModule.Name)
	}

	if !reflect.DeepEqual(wantImports, jsonModule.Imports) {
		t.Errorf("want %q imports, got %q imports", wantImports, jsonModule.Imports)
	}

	if wantNumResources != len(jsonModule.Resources) {
		t.Errorf("want %d resources, got %d resources", wantNumResources, len(jsonModule.Resources))
	}
}
