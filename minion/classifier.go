package minion

// Global classifier registry
// Each classifier that is registered in the
// registry is called by the minions in order to
// classify themselves
var ClassifierRegistry []MinionClassifier

// Interface for classifying minion
type MinionClassifier interface {
	// Gets the key of a classifier
	GetKey() (string, error)

	// Gets the description of a classifier
	GetDescription() (string, error)

	// Classifies a minion returning the classifier value
	GetValue(m Minion) (string, error)
}

// Register new classifiers
func RegisterClassifier(c ...MinionClassifier) error {
	ClassifierRegistry = append(ClassifierRegistry, c...)

	return nil
}

// Simple classifier
type SimpleClassifier struct {
	Key, Value, Description string
}

// Creates a new simple classifier
func NewSimpleClassifier(key, value, description string) MinionClassifier {
	c := &SimpleClassifier{
		Key: key,
		Value: value,
		Description: description,
	}

	return c
}

func (c *SimpleClassifier) GetKey() (string, error) {
	return c.Key, nil
}

func (c *SimpleClassifier) GetDescription() (string, error) {
	return c.Description, nil
}

func (c *SimpleClassifier) GetValue(m Minion) (string, error) {
	return c.Value, nil
}

// Classifier that uses callbacks for classifying minions
type cbClassifier func(Minion) (string, error)
type CallbackClassifier struct {
	Key, Description string

	// Callback used to classify the minions
	Callback cbClassifier
}

func NewCallbackClassifier(key, description string, callback cbClassifier) MinionClassifier {
	c := &CallbackClassifier{
		Key:         key,
		Description: description,
		Callback:    callback,
	}

	return c
}

func (c *CallbackClassifier) GetKey() (string, error) {
	return c.Key, nil
}

func (c *CallbackClassifier) GetDescription() (string, error) {
	return c.Description, nil
}

func (c *CallbackClassifier) GetValue(m Minion) (string, error) {
	value, err := c.Callback(m)

	return value, err
}
