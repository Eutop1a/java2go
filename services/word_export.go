package services

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/unidoc/unioffice/document"
)

//go:embed templates/*
var templateFS embed.FS

type WordExporter struct {
	data map[string]string
}

func NewWordExporter(data map[string]string) *WordExporter {
	return &WordExporter{data: data}
}

func (we *WordExporter) ExportTestPaper(templateType int) (*os.File, error) {
	templatePath, err := getTemplatePath(templateType)
	if err != nil {
		return nil, err
	}

	doc, err := loadTemplate(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load template: %w", err)
	}

	replacePlaceholders(doc, we.data)

	tmpFile, err := createTempFile()
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmpFile.Close()

	if err := saveDocument(doc, tmpFile); err != nil {
		return nil, fmt.Errorf("failed to save document: %w", err)
	}

	return reopenTempFile(tmpFile.Name())
}

func getTemplatePath(templateType int) (string, error) {
	switch templateType {
	case 1:
		return "templates/test-paper-template.docx", nil
	case 2:
		return "templates/answer-template.docx", nil
	default:
		return "", fmt.Errorf("invalid template type: %d", templateType)
	}
}

func loadTemplate(path string) (*document.Document, error) {
	file, err := templateFS.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return document.Read(bytes.NewReader(data), int64(len(data)))
}

func replacePlaceholders(doc *document.Document, data map[string]string) {
	for _, p := range doc.Paragraphs() {
		for _, r := range p.Runs() {
			original := r.Text()
			replaced := replaceText(original, data)
			if replaced != original {
				r.Clear()
				r.AddText(replaced)
			}
		}
	}
}

func replaceText(s string, data map[string]string) string {
	for k, v := range data {
		placeholder := "${" + k + "}"
		s = strings.ReplaceAll(s, placeholder, v)
	}
	return s
}

func createTempFile() (*os.File, error) {
	return os.CreateTemp("", "testpaper_*.docx")
}

func saveDocument(doc *document.Document, tmpFile *os.File) error {
	tmpPath := tmpFile.Name()
	tmpFile.Close() // Close first to allow SaveToFile to write
	return doc.SaveToFile(tmpPath)
}

func reopenTempFile(path string) (*os.File, error) {
	return os.Open(path)
}

func main() {
	// 示例用法
	data := map[string]string{
		"title":    "Go语言测试试卷",
		"question": "1. 什么是Go语言？",
	}

	exporter := NewWordExporter(data)
	file, err := exporter.ExportTestPaper(1)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	fmt.Printf("Generated file: %s\n", file.Name())
}
