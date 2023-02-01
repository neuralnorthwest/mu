package logging

import "testing"

// Test_New tests that the New function returns a logger, and that it is a
// zapLogger.
func Test_New(t *testing.T) {
	t.Parallel()
	logger, err := New()
	if err != nil {
		t.Fatalf("New returned an error: %v", err)
	}
	if logger == nil {
		t.Fatal("New returned a nil logger")
	}
	if _, ok := logger.(*zapLogger); !ok {
		t.Fatal("New returned a logger that is not a zapLogger")
	}
	z := logger.(*zapLogger).GetZap()
	if z == nil {
		t.Fatal("zapLogger.GetZap returned a nil logger")
	}
}

// Test_With tests the With function.
func Test_With(t *testing.T) {
	t.Parallel()
	logger, err := New()
	if err != nil {
		t.Fatalf("New returned an error: %v", err)
	}
	logger = logger.With("key", "value")
	if logger == nil {
		t.Fatal("With returned a nil logger")
	}
	if _, ok := logger.(*zapLogger); !ok {
		t.Fatal("With returned a logger that is not a zapLogger")
	}
	z := logger.(*zapLogger).GetZap()
	if z == nil {
		t.Fatal("zapLogger.GetZap returned a nil logger")
	}
}
