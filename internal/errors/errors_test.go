package errors

import (
	"errors"
	"testing"
)

func TestTextStatError_Error(t *testing.T) {
	// Test error without cause
	err := NewValidationError("test validation error")
	expected := "Validation Error: test validation error"

	if err.Error() != expected {
		t.Errorf("TextStatError.Error() = %q, expected %q", err.Error(), expected)
	}

	// Test error with cause
	underlyingErr := errors.New("underlying error")
	err = NewFileError("test file error", underlyingErr)
	expected = "File Error: test file error (caused by: underlying error)"

	if err.Error() != expected {
		t.Errorf("TextStatError.Error() with cause = %q, expected %q", err.Error(), expected)
	}
}

func TestTextStatError_Unwrap(t *testing.T) {
	// Test error without cause
	err := NewValidationError("test validation error")
	if err.Unwrap() != nil {
		t.Errorf("TextStatError.Unwrap() without cause = %v, expected nil", err.Unwrap())
	}

	// Test error with cause
	underlyingErr := errors.New("underlying error")
	err = NewFileError("test file error", underlyingErr)

	if err.Unwrap() != underlyingErr {
		t.Errorf("TextStatError.Unwrap() with cause = %v, expected %v", err.Unwrap(), underlyingErr)
	}
}

func TestErrorType_String(t *testing.T) {
	tests := []struct {
		errorType ErrorType
		expected  string
	}{
		{ValidationError, "Validation Error"},
		{FileError, "File Error"},
		{ParsingError, "Parsing Error"},
		{ExtractionError, "Extraction Error"},
		{UnsupportedFormatError, "Unsupported Format Error"},
		{ErrorType(999), "Unknown Error"}, // Test unknown type
	}

	for _, tt := range tests {
		result := tt.errorType.String()
		if result != tt.expected {
			t.Errorf("ErrorType.String(%d) = %q, expected %q", tt.errorType, result, tt.expected)
		}
	}
}

func TestNewValidationError(t *testing.T) {
	message := "validation failed"
	err := NewValidationError(message)

	if err.Type != ValidationError {
		t.Errorf("NewValidationError() Type = %v, expected %v", err.Type, ValidationError)
	}

	if err.Message != message {
		t.Errorf("NewValidationError() Message = %q, expected %q", err.Message, message)
	}

	if err.Cause != nil {
		t.Errorf("NewValidationError() Cause = %v, expected nil", err.Cause)
	}
}

func TestNewFileError(t *testing.T) {
	message := "file not found"
	cause := errors.New("os error")
	err := NewFileError(message, cause)

	if err.Type != FileError {
		t.Errorf("NewFileError() Type = %v, expected %v", err.Type, FileError)
	}

	if err.Message != message {
		t.Errorf("NewFileError() Message = %q, expected %q", err.Message, message)
	}

	if err.Cause != cause {
		t.Errorf("NewFileError() Cause = %v, expected %v", err.Cause, cause)
	}
}

func TestNewParsingError(t *testing.T) {
	message := "parse error"
	cause := errors.New("syntax error")
	err := NewParsingError(message, cause)

	if err.Type != ParsingError {
		t.Errorf("NewParsingError() Type = %v, expected %v", err.Type, ParsingError)
	}

	if err.Message != message {
		t.Errorf("NewParsingError() Message = %q, expected %q", err.Message, message)
	}

	if err.Cause != cause {
		t.Errorf("NewParsingError() Cause = %v, expected %v", err.Cause, cause)
	}
}

func TestNewExtractionError(t *testing.T) {
	message := "extraction failed"
	cause := errors.New("format error")
	err := NewExtractionError(message, cause)

	if err.Type != ExtractionError {
		t.Errorf("NewExtractionError() Type = %v, expected %v", err.Type, ExtractionError)
	}

	if err.Message != message {
		t.Errorf("NewExtractionError() Message = %q, expected %q", err.Message, message)
	}

	if err.Cause != cause {
		t.Errorf("NewExtractionError() Cause = %v, expected %v", err.Cause, cause)
	}
}

func TestNewUnsupportedFormatError(t *testing.T) {
	format := "binary"
	err := NewUnsupportedFormatError(format)

	if err.Type != UnsupportedFormatError {
		t.Errorf("NewUnsupportedFormatError() Type = %v, expected %v", err.Type, UnsupportedFormatError)
	}

	expectedMessage := "unsupported file format: binary"
	if err.Message != expectedMessage {
		t.Errorf("NewUnsupportedFormatError() Message = %q, expected %q", err.Message, expectedMessage)
	}

	if err.Cause != nil {
		t.Errorf("NewUnsupportedFormatError() Cause = %v, expected nil", err.Cause)
	}
}

func TestErrorTypeConstants(t *testing.T) {
	tests := []struct {
		constant ErrorType
		expected int
	}{
		{ValidationError, 0},
		{FileError, 1},
		{ParsingError, 2},
		{ExtractionError, 3},
		{UnsupportedFormatError, 4},
	}

	for _, tt := range tests {
		if int(tt.constant) != tt.expected {
			t.Errorf("ErrorType constant %v = %d, expected %d", tt.constant, int(tt.constant), tt.expected)
		}
	}
}

func TestErrorHandlingWithErrorsIs(t *testing.T) {
	// Test that errors.Is works correctly with our custom errors
	originalErr := errors.New("underlying error")
	fileErr := NewFileError("file operation failed", originalErr)

	if !errors.Is(fileErr, originalErr) {
		t.Error("errors.Is should return true for underlying error")
	}

	// Test with a different error
	otherErr := errors.New("different error")
	if errors.Is(fileErr, otherErr) {
		t.Error("errors.Is should return false for different error")
	}
}

func BenchmarkNewValidationError(b *testing.B) {
	message := "benchmark validation error"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewValidationError(message)
	}
}

func BenchmarkNewFileError(b *testing.B) {
	message := "benchmark file error"
	cause := errors.New("benchmark underlying error")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewFileError(message, cause)
	}
}

func BenchmarkTextStatError_Error(b *testing.B) {
	err := NewValidationError("benchmark error message")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}
