package qmp

import "fmt"

// ErrorClass corresponds to the QMP protocol error class.
type ErrorClass string

const (
	ErrClassGenericError    ErrorClass = "GenericError"
	ErrClassCommandNotFound ErrorClass = "CommandNotFound"
	ErrClassDeviceNotFound  ErrorClass = "DeviceNotFound"
	ErrClassKVMMissing      ErrorClass = "KVMMissing"
)

// CommandError represents a QMP command execution failure returned by the server.
type CommandError struct {
	Class   ErrorClass
	Desc    string
	Command string
}

func (e *CommandError) Error() string {
	return fmt.Sprintf("QMP command %q failed [%s]: %s", e.Command, e.Class, e.Desc)
}

// TransportError represents transport-level errors (disconnections, timeouts, etc.).
type TransportError struct {
	Op  string
	Err error
}

func (e *TransportError) Error() string {
	return fmt.Sprintf("QMP transport error during %s: %v", e.Op, e.Err)
}

func (e *TransportError) Unwrap() error { return e.Err }
