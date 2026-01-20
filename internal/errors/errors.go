package errors

import (
	"fmt"
)

// TextStatError represents different types of errors that can occur in the application
type TextStatError struct {
	Type    ErrorType
	Message string
	Cause   error
}

type ErrorType int

const (
	ValidationError ErrorType = iota
	FileError
	ParsingError
	ExtractionError
	UnsupportedFormatError
)

// Error implements the error interface
func (e *TextStatError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Type.String(), e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Type.String(), e.Message)
}

// Unwrap returns the underlying cause
func (e *TextStatError) Unwrap() error {
	return e.Cause
}

// String returns the string representation of ErrorType
func (et ErrorType) String() string {
	switch et {
	case ValidationError:
		return "Validation Error"
	case FileError:
		return "File Error"
	case ParsingError:
		return "Parsing Error"
	case ExtractionError:
		return "Extraction Error"
	case UnsupportedFormatError:
		return "Unsupported Format Error"
	default:
		return "Unknown Error"
	}
}

// NewValidationError creates a new validation error
func NewValidationError(message string) *TextStatError {
	return &TextStatError{
		Type:    ValidationError,
		Message: message,
	}
}

// NewFileError creates a new file error with an underlying cause
func NewFileError(message string, cause error) *TextStatError {
	return &TextStatError{
		Type:    FileError,
		Message: message,
		Cause:   cause,
	}
}

// NewParsingError creates a new parsing error with an underlying cause
func NewParsingError(message string, cause error) *TextStatError {
	return &TextStatError{
		Type:    ParsingError,
		Message: message,
		Cause:   cause,
	}
}

// NewExtractionError creates a new extraction error with an underlying cause
func NewExtractionError(message string, cause error) *TextStatError {
	return &TextStatError{
		Type:    ExtractionError,
		Message: message,
		Cause:   cause,
	}
}

// NewUnsupportedFormatError creates a new unsupported format error
func NewUnsupportedFormatError(format string) *TextStatError {
	return &TextStatError{
		Type:    UnsupportedFormatError,
		Message: fmt.Sprintf("unsupported file format: %s", format),
	}
}
