package status

type Error string

func (e Error) Error() string {
	return string(e)
}

// ErrInvalidVersion is returned when an invalid version is provided.
var ErrInvalidVersion = Error("invalid version")

// ErrInvalidRange is returned when an invalid range is provided.
var ErrInvalidRange = Error("invalid range")

// ErrOutOfRange is returned when a value is out of range.
var ErrOutOfRange = Error("out of range")

// ErrAlreadyExists is returned when a value already exists.
var ErrAlreadyExists = Error("already exists")

// ErrNotFound is returned when a value is not found.
var ErrNotFound = Error("not found")

// ErrInvalidArgument is returned when an invalid argument is provided.
var ErrInvalidArgument = Error("invalid argument")

// ErrAlreadyStarted is returned when a service is already started.
var ErrAlreadyStarted = Error("already started")

// ErrNotStarted is returned when a service is not started.
var ErrNotStarted = Error("not started")
