package generator

import (
	"awdoc/pkg/analyzer"
	"awdoc/pkg/parser"
	"strings"
	"testing"
)

// TestMarkdownGenerator тестирует генератор Markdown
func TestMarkdownGenerator(t *testing.T) {
	// Создаем тестовую SourceInfo
	sourceInfo := &parser.SourceInfo{
		Files: []string{},
		Packages: map[string]*parser.Package{
			"test": {
				Name: "test",
				Path: "test",
				Doc:  "Test package",
				Elements: []parser.CodeElement{
					{
						Name:       "TestFunc",
						Type:       parser.ElementFunc,
						Exported:   true,
						Doc:        "Test function",
						Signature:  "func TestFunc() error",
						SourceFile: "test.go",
						StartLine:  1,
						EndLine:    5,
					},
				},
				Imports: map[string]bool{},
				ExportedAPI: []parser.CodeElement{
					{
						Name:       "TestFunc",
						Type:       parser.ElementFunc,
						Exported:   true,
						Doc:        "Test function",
						Signature:  "func TestFunc() error",
						SourceFile: "test.go",
						StartLine:  1,
						EndLine:    5,
					},
				},
			},
		},
	}

	// Создаем пустой граф
	graph := &analyzer.DependencyGraph{
		Nodes:      map[string]*analyzer.PackageNode{},
		Edges:      map[string][]string{},
		Cycles:     [][]string{},
		Layers:     [][]string{{"test"}},
		GodObjects: []string{},
	}

	// Создаем генератор
	gen := NewMarkdownGenerator(sourceInfo, graph)
	doc := gen.GenerateProjectDoc()

	if doc == "" {
		t.Fatal("Generated document is empty")
	}

	// Проверяем ключевые элементы
	if !strings.Contains(doc, "# API Documentation") {
		t.Error("Missing main header")
	}

	if !strings.Contains(doc, "test") {
		t.Error("Package name not found in documentation")
	}

	if !strings.Contains(doc, "TestFunc") {
		t.Error("Function name not found in documentation")
	}

	if !strings.Contains(doc, "Test package") {
		t.Error("Package documentation not found")
	}
}

// TestDocumentationBuilder тестирует конструктор документации
func TestDocumentationBuilder(t *testing.T) {
	sourceInfo := &parser.SourceInfo{
		Files:    []string{},
		Packages: map[string]*parser.Package{},
	}

	graph := &analyzer.DependencyGraph{
		Nodes:  map[string]*analyzer.PackageNode{},
		Edges:  map[string][]string{},
		Cycles: [][]string{},
		Layers: [][]string{},
	}

	builder := NewDocumentationBuilder(sourceInfo, graph)

	// Тест Markdown генерации
	markdown := builder.BuildMarkdown()
	if markdown == "" {
		t.Error("Markdown document is empty")
	}

	if !strings.Contains(markdown, "# API Documentation") {
		t.Error("Markdown header missing")
	}

	// Тест HTML генерации (placeholder)
	html := builder.BuildHTML()
	if html == "" {
		t.Error("HTML document is empty")
	}

	if !strings.Contains(html, "html") {
		t.Error("HTML tag missing")
	}
}
