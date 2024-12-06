package main

import (
	"archive/zip"
	"fmt"
	"io"
	"strings"
)

type DocumentExtractor struct {
	filePath string
	fileType string
}

func NewDocumentExtractor(filePath string, fileType string) *DocumentExtractor {
	return &DocumentExtractor{
		filePath: filePath,
		fileType: fileType,
	}
}

func extractPDFText( /*filePath string*/ ) (string, error) {
	fmt.Println("unimplemented")
	return "", nil
}

func extractDocxText(filePath string) (string, error) {
	doc, err := extractDocxTextFromZip(filePath)
	if err != nil {
		return "", err
	}

	text := extractTextFromXml(doc)
	return text, nil
}

func extractDocxTextFromZip(filePath string) (string, error) {
	// Open the DOCX file as a zip file
	zipReader, err := zip.OpenReader(filePath)
	if err != nil {
		return "", err
	}
	defer zipReader.Close()

	// Find the document.xml file inside the DOCX
	var docXml *zip.File
	for _, f := range zipReader.File {
		if f.Name == "word/document.xml" {
			docXml = f
			break
		}
	}

	if docXml == nil {
		return "", fmt.Errorf("document.xml not found in the DOCX file")
	}

	// Extract and read the content of document.xml
	docXmlReader, err := docXml.Open()
	if err != nil {
		return "", err
	}
	defer docXmlReader.Close()

	xmlContent, err := io.ReadAll(docXmlReader)
	if err != nil {
		return "", err
	}

	// Remove XML tags and extract text
	text := extractTextFromXml(string(xmlContent))

	return text, nil
}

func extractTextFromXml(xml string) string {
	text := xml
	text = strings.ReplaceAll(text, "<w:p>", "")
	text = strings.ReplaceAll(text, "</w:p>", "\n")
	text = strings.ReplaceAll(text, "<w:t>", "")
	text = strings.ReplaceAll(text, "</w:t>", "")
	text = strings.ReplaceAll(text, "<w:tab/>", "    ")
	return strings.TrimSpace(text)
}

func (d *DocumentExtractor) ExtractText() (string, error) {
	switch d.fileType {
	case "pdf":
		return extractPDFText( /*d.filePath*/ )
	case "docx":
		return extractDocxText(d.filePath)
	default:
		return "", fmt.Errorf("unsupported file type: %s", d.fileType)
	}

}
