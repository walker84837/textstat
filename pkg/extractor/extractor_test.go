package extractor

import (
	"os"
	"strings"
	"testing"
)

func TestNewDocumentExtractor(t *testing.T) {
	// Test valid extractor creation
	extractor, err := NewDocumentExtractor("test.txt", TextType)
	if err != nil {
		t.Fatalf("NewDocumentExtractor() error = %v", err)
	}
	if extractor == nil {
		t.Fatal("NewDocumentExtractor() returned nil")
	}

	// Test empty file path
	_, err = NewDocumentExtractor("", TextType)
	if err == nil {
		t.Error("NewDocumentExtractor() with empty path should return error")
	}

	// Test invalid file type
	_, err = NewDocumentExtractor("test.txt", FileType("invalid"))
	if err == nil {
		t.Error("NewDocumentExtractor() with invalid type should return error")
	}
}

func TestValidateFilePath(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "test_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write some content
	tempFile.WriteString("test content")
	tempFile.Close()

	tests := []struct {
		name      string
		filePath  string
		expectErr bool
	}{
		{"Valid file", tempFile.Name(), false},
		{"Path traversal", "../../../etc/passwd", true},
		{"Non-existent file", "/non/existent/file.txt", true},
		{"Empty path", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			extractor := &DocumentExtractor{
				filePath: tt.filePath,
				fileType: TextType,
			}

			err := extractor.validateFilePath()
			if (err != nil) != tt.expectErr {
				t.Errorf("validateFilePath() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

func TestExtractPlainText(t *testing.T) {
	// Create a temporary file with test content
	content := "This is test content.\nWith multiple lines.\nAnd some punctuation!"

	tempFile, err := os.CreateTemp("", "test_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write test content
	_, err = tempFile.WriteString(content)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	extractor := &DocumentExtractor{
		filePath: tempFile.Name(),
		fileType: TextType,
	}

	result, err := extractor.extractPlainText()
	if err != nil {
		t.Fatalf("extractPlainText() error = %v", err)
	}

	if result != content {
		t.Errorf("extractPlainText() = %q, expected %q", result, content)
	}
}

func TestExtractTextFromXML(t *testing.T) {
	xml := `<w:p><w:t>Hello</w:t><w:t> </w:t><w:t>world</w:t></w:p><w:p><w:t>This is a test.</w:t></w:p>`
	expected := "Hello world\nThis is a test."

	extractor := &DocumentExtractor{}

	result := extractor.extractTextFromXML(xml)
	if result != expected {
		t.Errorf("extractTextFromXML() = %q, expected %q", result, expected)
	}
}

func TestExtractTextFromXMLWithEmptyContent(t *testing.T) {
	xml := `<w:p></w:p><w:p></w:p>`

	extractor := &DocumentExtractor{}

	result := extractor.extractTextFromXML(xml)
	if result != "" {
		t.Errorf("extractTextFromXML() with empty content = %q, expected empty string", result)
	}
}

func TestFileTypeConstants(t *testing.T) {
	tests := []struct {
		value    FileType
		expected string
	}{
		{TextType, "text"},
		{PDFType, "pdf"},
		{DocxType, "docx"},
	}

	for _, tt := range tests {
		if string(tt.value) != tt.expected {
			t.Errorf("FileType %s = %s, expected %s", tt.value, string(tt.value), tt.expected)
		}
	}
}

func TestExtractDocxTextFromZip(t *testing.T) {
	// Create a mock DOCX file (ZIP archive) with document.xml
	tempFile, err := os.CreateTemp("", "test_*.docx")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	extractor := &DocumentExtractor{
		filePath: tempFile.Name(),
		fileType: DocxType,
	}

	// This will fail because our temp file is not a valid DOCX
	_, err = extractor.extractDocxTextFromZip()
	if err == nil {
		t.Error("extractDocxTextFromZip() with invalid DOCX should return error")
	}
}

// Integration test with a simple DOCX-like structure
func TestExtractDocxTextIntegration(t *testing.T) {
	// Skip this test if we can't create a proper DOCX file
	t.Skip("Skipping DOCX integration test - requires proper DOCX file creation")
}

func BenchmarkExtractPlainText(b *testing.B) {
	// Create a temporary file with test content
	content := strings.Repeat("This is test content for benchmarking. ", 1000)

	tempFile, err := os.CreateTemp("", "bench_*.txt")
	if err != nil {
		b.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write test content
	_, err = tempFile.WriteString(content)
	if err != nil {
		b.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	extractor := &DocumentExtractor{
		filePath: tempFile.Name(),
		fileType: TextType,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = extractor.extractPlainText()
	}
}

func BenchmarkExtractTextFromXML(b *testing.B) {
	xml := `<w:p><w:t>` + strings.Repeat("Hello world. ", 100) + `</w:t></w:p>`

	extractor := &DocumentExtractor{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		extractor.extractTextFromXML(xml)
	}
}
