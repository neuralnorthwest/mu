package config

// Config holds configuration variables.
type Config interface {
	// NewInt creates a new int variable.
	NewInt(name string, defaultValue int, description string, options ...IntOption) error
	// Int returns the value of the int variable with the given name. If the variable does not exist, it calls bug.Bug.
	// If bug.Bug does not panic, Int returns 0.
	Int(name string) int
	// DescribeInt returns the description of the int variable with the given name. If the variable does not exist, it calls bug.Bug.
	// If bug.Bug does not panic, DescribeInt returns "".
	DescribeInt(name string) string
	// NewString creates a new string variable.
	NewString(name string, defaultValue string, description string, options ...StringOption) error
	// String returns the value of the string variable with the given name. If the variable does not exist, it calls bug.Bug.
	// If bug.Bug does not panic, String returns "".
	String(name string) string
	// DescribeString returns the description of the string variable with the given name. If the variable does not exist, it calls bug.Bug.
	// If bug.Bug does not panic, DescribeString returns "".
	DescribeString(name string) string
	// NewBool creates a new bool variable.
	NewBool(name string, defaultValue bool, description string) error
	// Bool returns the value of the bool variable with the given name. If the variable does not exist, it calls bug.Bug.
	// If bug.Bug does not panic, Bool returns false.
	Bool(name string) bool
	// DescribeBool returns the description of the bool variable with the given name. If the variable does not exist, it calls bug.Bug.
	// If bug.Bug does not panic, DescribeBool returns "".
	DescribeBool(name string) string
}

// Source is an interface that loads the values of the configuration variables.
type Source interface {
	// SetPrefix sets the prefix of the names of the configuration variables.
	SetPrefix(prefix string)
	// LoadInt loads the value of the int variable with the given name. If the
	// variable does not exist, it must return 0 and status.ErrNotFound.
	LoadInt(name string) (int, error)
	// LoadString loads the value of the string variable with the given name. If
	// the variable does not exist, it must return "" and status.ErrNotFound.
	LoadString(name string) (string, error)
	// LoadBool loads the value of the bool variable with the given name. If the
	// variable does not exist, it must return false and status.ErrNotFound.
	LoadBool(name string) (bool, error)
}

// configImpl holds the configuration variables.
type configImpl struct {
	ints       map[string]*Int
	strings    map[string]*String
	bools      map[string]*Bool
	loadPrefix string
	source     Source
}

// Option is an option for Config.
type Option func(c *configImpl)

// WithLoadPrefix returns an Option that prepends a prefix to the names of the
// configuration variables when loading them.
func WithLoadPrefix(prefix string) Option {
	return func(c *configImpl) {
		c.loadPrefix = prefix
	}
}

// WithSource returns an Option that sets the source of the configuration
// variables.
func WithSource(source Source) Option {
	return func(c *configImpl) {
		c.source = source
	}
}

// New creates a new Config.
func New(opts ...Option) Config {
	c := &configImpl{
		ints:    make(map[string]*Int),
		strings: make(map[string]*String),
		bools:   make(map[string]*Bool),
	}
	for _, opt := range opts {
		opt(c)
	}
	if c.source == nil {
		c.source = &envSource{}
	}
	c.source.SetPrefix(c.loadPrefix)
	return c
}
