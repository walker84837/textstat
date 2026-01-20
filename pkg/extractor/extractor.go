package extractor

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ledongthuc/pdf"
	"textstat/internal/errors"
)

// FileType represents supported document types
type FileType string

const (
	TextType FileType = "text"
	PDFType  FileType = "pdf"
	DocxType FileType = "docx"
)

// DocumentExtractor handles text extraction from various document types
type DocumentExtractor struct {
	filePath string
	fileType FileType
}

// NewDocumentExtractor creates a new document extractor instance
func NewDocumentExtractor(filePath string, fileType FileType) (*DocumentExtractor, error) {
	if filePath == "" {
		return nil, errors.NewValidationError("file path cannot be empty")
	}

	// Validate file type
	switch fileType {
	case TextType, PDFType, DocxType:
		// Valid types
	default:
		return nil, errors.NewUnsupportedFormatError(string(fileType))
	}

	return &DocumentExtractor{
		filePath: filePath,
		fileType: fileType,
	}, nil
}

// validateFilePath performs security validation on file paths
func (e *DocumentExtractor) validateFilePath() error {
	// Clean the path to resolve any directory traversal
	cleanPath := filepath.Clean(e.filePath)

	// Check for path traversal attempts
	if strings.Contains(cleanPath, "..") {
		return errors.NewValidationError("path traversal detected")
	}

	// Check if file exists and is accessible
	info, err := os.Stat(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.NewFileError("file does not exist", err)
		}
		if os.IsPermission(err) {
			return errors.NewFileError("permission denied", err)
		}
		return errors.NewFileError("failed to access file", err)
	}

	// Check if it's a regular file
	if !info.Mode().IsRegular() {
		return errors.NewValidationError("path is not a regular file")
	}

	// Check file size (prevent processing extremely large files)
	const maxFileSize = 100 * 1024 * 1024 // 100MB
	if info.Size() > maxFileSize {
		return errors.NewValidationError("file size exceeds maximum allowed size")
	}

	return nil
}

// ExtractText extracts text from the configured document
func (e *DocumentExtractor) ExtractText() (string, error) {
	if err := e.validateFilePath(); err != nil {
		return "", err
	}

	switch e.fileType {
	case TextType:
		return e.extractPlainText()
	case PDFType:
		return e.extractPDFText()
	case DocxType:
		return e.extractDocxText()
	default:
		return "", errors.NewUnsupportedFormatError(string(e.fileType))
	}
}

// extractPlainText extracts text from plain text files
func (e *DocumentExtractor) extractPlainText() (string, error) {
	file, err := os.Open(e.filePath)
	if err != nil {
		return "", errors.NewFileError("failed to open text file", err)
	}
	defer file.Close()

	// Use strings.Builder for efficient string building
	var builder strings.Builder
	buffer := make([]byte, 4096)

	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return "", errors.NewFileError("failed to read text file", err)
		}
		if n == 0 {
			break
		}

		builder.Write(buffer[:n])
	}

	return builder.String(), nil
}

// extractPDFText extracts text from PDF files
func (e *DocumentExtractor) extractPDFText() (string, error) {
	f, r, err := pdf.Open(e.filePath)
	if err != nil {
		return "", errors.NewExtractionError("failed to open PDF", err)
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			// Log error but don't override the main error
			_ = closeErr
		}
	}()

	if r.NumPage() == 0 {
		return "", nil
	}

	var builder strings.Builder

	for i := 1; i <= r.NumPage(); i++ {
		p := r.Page(i)
		if p.V.IsNull() {
			continue
		}

		pageText, err := p.GetPlainText(nil)
		if err != nil {
			return "", errors.NewExtractionError(fmt.Sprintf("failed to extract text from page %d", i), err)
		}

		builder.WriteString(pageText)
		if i < r.NumPage() {
			builder.WriteString("\n")
		}
	}

	return builder.String(), nil
}

// extractDocxText extracts text from DOCX files
func (e *DocumentExtractor) extractDocxText() (string, error) {
	doc, err := e.extractDocxTextFromZip()
	if err != nil {
		return "", err
	}

	text := e.extractTextFromXML(doc)
	return text, nil
}

// extractDocxTextFromZip extracts document.xml from DOCX file
func (e *DocumentExtractor) extractDocxTextFromZip() (string, error) {
	zipReader, err := zip.OpenReader(e.filePath)
	if err != nil {
		return "", errors.NewExtractionError("failed to open DOCX file as zip", err)
	}
	defer zipReader.Close()

	var docXml *zip.File
	for _, f := range zipReader.File {
		if f.Name == "word/document.xml" {
			docXml = f
			break
		}
	}

	if docXml == nil {
		return "", errors.NewExtractionError("document.xml not found in DOCX file", nil)
	}

	docXmlReader, err := docXml.Open()
	if err != nil {
		return "", errors.NewExtractionError("failed to open document.xml", err)
	}
	defer docXmlReader.Close()

	xmlContent, err := io.ReadAll(docXmlReader)
	if err != nil {
		return "", errors.NewExtractionError("failed to read document.xml", err)
	}

	return string(xmlContent), nil
}

// extractTextFromXML extracts text content from DOCX XML
func (e *DocumentExtractor) extractTextFromXML(xml string) string {
	// Process the XML in chunks to extract text
	replacements := map[string]string{
		"<w:p>":    "",
		"</w:p>":   "\n",
		"<w:t>":    "",
		"</w:t>":   "",
		"<w:tab/>": "    ",
		"<w:br/>":  "\n",
		"<w:cr/>":  "\n",
	}

	text := xml
	for old, new := range replacements {
		text = strings.ReplaceAll(text, old, new)
	}

	// Remove any remaining XML tags
	xmlTagRegex := strings.NewReplacer(
		"<w:r>", "",
		"</w:r>", "",
		"<w:rPr>", "",
		"</w:rPr>", "",
	)
	text = xmlTagRegex.Replace(text)

	// Clean up extra whitespace
	lines := strings.Split(text, "\n")
	var cleanedLines []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			cleanedLines = append(cleanedLines, line)
		}
	}

	return strings.Join(cleanedLines, "\n")
}
