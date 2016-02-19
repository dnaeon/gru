package classifier

import "errors"

var errNotFound = errors.New("Classifier key not found")

// Classifier type contains a key/value pair repsenting a classifier
type Classifier struct {
	// Classifier key
	Key string `json:"key"`

	// Classifier value
	Value string `json:"value"`
}

// Type of classifier providers
// A classifier provider is what does the
// actual evaluation of a classifier
type provider func() (string, error)

// Registry provides a global registry for all classifiers
var Registry = make(map[string]provider)

// Register registers a new classifier to the registry
func Register(key string, p provider) error {
	Registry[key] = p

	return nil
}

// Get retrieves a classifier from the registry by looking up its key
func Get(key string) (*Classifier, error) {
	c := new(Classifier)

	p, ok := Registry[key]
	if ok {
		// Evaluate the classifier provider
		value, err := p()
		c := &Classifier{
			Key:   key,
			Value: value,
		}
		return c, err
	}

	return c, errNotFound
}
