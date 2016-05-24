// +build !windows

package resource

import (
	"errors"
	"fmt"

	"github.com/dnaeon/gru/utils"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
)

// Name and description of the resource
const packageResourceType = "package"
const packageResourceDesc = "meta resource for package management"

// ErrNoPackageProviderFound is returned when no suitable provider is found
var ErrNoPackageProviderFound = errors.New("No suitable package provider found")

// BasePackageResource is the base resource type for package management
// It's purpose is to be embeded into other package resource providers.
type BasePackageResource struct {
	// Name of the package
	Name string `hcl:"name"`

	// Version of the package
	Version string `hcl:"version"`

	// Provider to use
	Provider string `hcl:"provider"`
}

// NewPackageResource creates a new resource for managing packages
func NewPackageResource(title string, obj *ast.ObjectItem, config *Config) (Resource, error) {
	// The package providers that we know of
	providers := map[string]Provider{
		pacmanResourceType: NewPacmanResource,
	}

	// Releases files used by the various GNU/Linux distros
	releases := map[string]Provider{
		"/etc/arch-release": NewPacmanResource,
	}

	// Decode the object from HCL
	var pr BasePackageResource
	err := hcl.DecodeObject(&pr, obj)
	if err != nil {
		return nil, err
	}

	// If we have a provider for this resource, use it
	if pr.Provider != "" {
		provider, ok := providers[pr.Provider]
		if !ok {
			return nil, fmt.Errorf("Unknown package provider '%s'", pr.Provider)
		}

		r, err := provider(title, obj, config)
		if err != nil {
			return nil, err
		}

		// Replace the resource type with our meta type
		r.SetType(packageResourceType)
		return r, nil
	}

	// Do our best to determine the proper provider for this resource
	for release, provider := range releases {
		dst := utils.NewFileUtil(release)
		if dst.Exists() {
			r, err := provider(title, obj, config)
			if err != nil {
				return nil, err
			}

			// Replace the resource type with our meta type
			r.SetType(packageResourceType)
			return r, nil
		}
	}

	return nil, ErrNoPackageProviderFound
}

func init() {
	item := RegistryItem{
		Name:        packageResourceType,
		Description: packageResourceDesc,
		Provider:    NewPackageResource,
	}

	Register(item)
}
