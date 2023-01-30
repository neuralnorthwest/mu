package mu

import (
	"testing"

	"golang.org/x/mod/semver"
)

// Test_Version tests that the Version function returns a semver string.
func Test_Version(t *testing.T) {
	v := Version()
	if !semver.IsValid(v) {
		t.Errorf("Version() = %q; want semver string", v)
	}
}
