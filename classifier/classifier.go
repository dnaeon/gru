package classifier

import "errors"

var errNotFound = errors.New("Classifier key not found")

type Classifier struct {
	// Classifier key
	Key string

	// Classifier value
	Value string
}

// Type of classifier providers
// A classifier provider is what does the
// actual evaluation of a classifier
type provider func() (string, error)

// Global classifier registry
var Registry = make(map[string]provider)

// Adds a classifier to the registry
func Register(key string, p provider) error {
	Registry[key] = p

	return nil
}

// Evaluates a classifier provider and returns the
// classifier value
func Get(key string) (*Classifier, error) {
	c := new(Classifier)

	p, ok := Registry[key]
	if ok {
		// Evaluate the classifier provider
		value, err := p()
		c := &Classifier{
			Key: key,
			Value: value,
		}
		return c, err
	}

	return c, errNotFound
}
