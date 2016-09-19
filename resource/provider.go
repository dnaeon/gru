package resource

// providerRegistry contains the registered providers
var providerRegistry = make([]ProviderItem, 0)

// Provider type is the type which creates new resources
type Provider func(name string) (Resource, error)

// ProviderItem type represents a single item from the
// provider registry.
type ProviderItem struct {
	// Type name of the provider
	Type string

	// Provider is the actual resource provider
	Provider Provider

	// Namespace represents the Lua table that the
	// provider will be registered in
	Namespace string
}

// RegisterProvider registers a provider to the registry.
func RegisterProvider(items ...ProviderItem) {
	providerRegistry = append(providerRegistry, items...)
}
