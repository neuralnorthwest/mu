package config

import (
	"strconv"

	"github.com/neuralnorthwest/mu/status"
)

// testSource is a Source that reads from a map.
type testSource struct {
	values map[string]string
	prefix string
}

var _ Source = (*testSource)(nil)

// newTestSource creates a new testSource.
func newTestSource() *testSource {
	return &testSource{
		values: make(map[string]string),
	}
}

// SetValue sets the value of the variable with the given name.
func (s *testSource) SetValue(name, value string) {
	s.values[name] = value
}

// SetPrefix sets the prefix for the names of the configuration variables.
func (s *testSource) SetPrefix(prefix string) {
	s.prefix = prefix
}

// LoadInt loads the value of the int variable with the given name.
func (s *testSource) LoadInt(name string) (int, error) {
	str, ok := s.values[s.prefix+name]
	if !ok {
		return 0, status.ErrNotFound
	}
	return strconv.Atoi(str)
}

// LoadString loads the value of the string variable with the given name.
func (s *testSource) LoadString(name string) (string, error) {
	str, ok := s.values[s.prefix+name]
	if !ok {
		return "", status.ErrNotFound
	}
	return str, nil
}

// LoadBool loads the value of the bool variable with the given name.
func (s *testSource) LoadBool(name string) (bool, error) {
	str, ok := s.values[s.prefix+name]
	if !ok {
		return false, status.ErrNotFound
	}
	return strconv.ParseBool(str)
}

// nullSource is a source that returns zero values.
type nullSource struct{}

var _ Source = (*nullSource)(nil)

// newNullSource creates a new nullSource.
func newNullSource() *nullSource {
	return &nullSource{}
}

// SetPrefix sets the prefix for the variables.
func (s *nullSource) SetPrefix(prefix string) {}

// LoadInt loads the value of the int variable with the given name.
func (s *nullSource) LoadInt(name string) (int, error) {
	return 0, status.ErrNotFound
}

// LoadString loads the value of the string variable with the given name.
func (s *nullSource) LoadString(name string) (string, error) {
	return "", status.ErrNotFound
}

// LoadBool loads the value of the bool variable with the given name.
func (s *nullSource) LoadBool(name string) (bool, error) {
	return false, status.ErrNotFound
}

// errorSource is a source that returns an error.
type errorSource struct {
	err error
}

var _ Source = (*errorSource)(nil)

// newErrorSource creates a new errorSource.
func newErrorSource(err error) *errorSource {
	return &errorSource{err: err}
}

// SetPrefix sets the prefix for the variables.
func (s *errorSource) SetPrefix(prefix string) {}

// LoadInt loads the value of the int variable with the given name.
func (s *errorSource) LoadInt(name string) (int, error) {
	return 0, s.err
}

// LoadString loads the value of the string variable with the given name.
func (s *errorSource) LoadString(name string) (string, error) {
	return "", s.err
}

// LoadBool loads the value of the bool variable with the given name.
func (s *errorSource) LoadBool(name string) (bool, error) {
	return false, s.err
}
