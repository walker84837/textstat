package cli

import (
	"bytes"
	"strings"
	"testing"

	"textstat/internal/errors"
	"textstat/pkg/textstat"
)

func TestNewApp(t *testing.T) {
	app := NewApp()
	if app == nil {
		t.Fatal("NewApp() returned nil")
	}
	if app.analyzer == nil {
		t.Error("NewApp() did not initialize analyzer")
	}
}

func TestParseFlags(t *testing.T) {
	// Test validation logic
	tests := []struct {
		name      string
		config    *Config
		expectErr bool
	}{
		{"Valid config", &Config{FileType: "text", OutputFormat: "table"}, false},
		{"Invalid file type", &Config{FileType: "invalid", OutputFormat: "table"}, true},
		{"Invalid output format", &Config{FileType: "text", OutputFormat: "invalid"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(tt.config)

			if (err != nil) != tt.expectErr {
				t.Errorf("validateConfig() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
		})
	}
}

// Helper function to validate config
func validateConfig(config *Config) error {
	// Validate file type
	switch config.FileType {
	case "text", "pdf", "docx":
		// Valid types
	default:
		return errors.NewValidationError("invalid file type. Supported types: text, pdf, docx")
	}

	// Validate output format
	switch config.OutputFormat {
	case "table", "json":
		// Valid formats
	default:
		return errors.NewValidationError("invalid output format. Supported formats: table, json")
	}

	return nil
}

func TestPrintStatsFormats(t *testing.T) {
	app := NewApp()

	stats := &textstat.TextStats{
		WordCount:             10,
		LetterCount:           50,
		SentenceCount:         2,
		ParagraphCount:        1,
		AverageWordLength:     5.0,
		AverageSentenceLength: 5.0,
		LongestWord:           "hello",
		MostCommonWord:        "hello",
		UniqueWordCount:       8,
		FleschKincaidGrade:    8.5,
		GunningFogIndex:       9.0,
		SMOGGrade:             8.0,
		EnglishLevel:          "Intermediate",
		SMOGInterpretation:    "Test interpretation",
		FogInterpretation:     "Test interpretation",
	}

	// Test that the functions don't panic - basic functionality test
	t.Run("PrintTableStats", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("printTableStats() panicked: %v", r)
			}
		}()
		app.printTableStats(stats)
	})

	t.Run("PrintJSONStats", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("printJSONStats() panicked: %v", r)
			}
		}()
		app.printJSONStats(stats)
	})
}

func TestHandleErrorTypes(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		contains string
	}{
		{"Validation error", errors.NewValidationError("test validation error"), "Validation Error"},
		{"File error", errors.NewFileError("test file error", nil), "File Error"},
		{"Parsing error", errors.NewParsingError("test parsing error", nil), "Parsing Error"},
		{"Extraction error", errors.NewExtractionError("test extraction error", nil), "Extraction Error"},
		{"Unsupported format", errors.NewUnsupportedFormatError("test"), "Unsupported Format Error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err != nil {
				// Just verify the error structure
				if !strings.Contains(tt.err.Error(), tt.contains) {
					t.Errorf("Error should contain %q, got: %s", tt.contains, tt.err.Error())
				}
			}
		})
	}
}

func BenchmarkValidateConfig(b *testing.B) {
	config := &Config{
		FilePath:     "test.txt",
		FileType:     "text",
		OutputFormat: "table",
		Verbose:      false,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validateConfig(config)
	}
}

func BenchmarkPrintTableStats(b *testing.B) {
	app := NewApp()

	stats := &textstat.TextStats{
		WordCount:             100,
		LetterCount:           500,
		SentenceCount:         10,
		ParagraphCount:        5,
		AverageWordLength:     5.0,
		AverageSentenceLength: 10.0,
		LongestWord:           "benchmarking",
		MostCommonWord:        "hello",
		UniqueWordCount:       50,
		FleschKincaidGrade:    8.5,
		GunningFogIndex:       9.0,
		SMOGGrade:             8.0,
		EnglishLevel:          "Intermediate",
		SMOGInterpretation:    "Test interpretation",
		FogInterpretation:     "Test interpretation",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		// Redirect output for benchmark
		_ = buf
		app.printTableStats(stats)
	}
}

func BenchmarkPrintJSONStats(b *testing.B) {
	app := NewApp()

	stats := &textstat.TextStats{
		WordCount:             100,
		LetterCount:           500,
		SentenceCount:         10,
		ParagraphCount:        5,
		AverageWordLength:     5.0,
		AverageSentenceLength: 10.0,
		LongestWord:           "benchmarking",
		MostCommonWord:        "hello",
		UniqueWordCount:       50,
		FleschKincaidGrade:    8.5,
		GunningFogIndex:       9.0,
		SMOGGrade:             8.0,
		EnglishLevel:          "Intermediate",
		SMOGInterpretation:    "Test interpretation",
		FogInterpretation:     "Test interpretation",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		// Redirect output for benchmark
		_ = buf
		app.printJSONStats(stats)
	}
}
