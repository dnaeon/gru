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
		{Name: "base-module"},
		{Name: "some-other-module"},
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
