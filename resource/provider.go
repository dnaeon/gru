package resource

import "fmt"

// Provider type is the type which creates new resources
type Provider func(name string) (Resource, error)

// providerRegistry contains the registered providers
var providerRegistry = make(map[string]Provider)

// RegisterProvider registers a provider to the registry
func RegisterProvider(typ string, p Provider) error {
	_, ok := providerRegistry[typ]
	if ok {
		return fmt.Errorf("Provider for '%s' is already registered", typ)
	}

	providerRegistry[typ] = p

	return nil
}
