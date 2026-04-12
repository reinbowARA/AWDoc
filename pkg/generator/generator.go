package generator

import (
	"awdoc/pkg/analyzer"
	"awdoc/pkg/parser"
)

// Generator - основной интерфейс для генератора документации
type Generator interface {
	GenerateProjectDoc() string
}

// DocumentationBuilder собирает документацию
type DocumentationBuilder struct {
	sourceInfo *parser.SourceInfo
	graph      *analyzer.DependencyGraph
}

// NewDocumentationBuilder создает конструктор документации
func NewDocumentationBuilder(sourceInfo *parser.SourceInfo, graph *analyzer.DependencyGraph) *DocumentationBuilder {
	return &DocumentationBuilder{
		sourceInfo: sourceInfo,
		graph:      graph,
	}
}

// BuildMarkdown строит документацию в Markdown
func (db *DocumentationBuilder) BuildMarkdown() string {
	gen := NewMarkdownGenerator(db.sourceInfo, db.graph)
	return gen.GenerateProjectDoc()
}

// BuildHTML строит документацию в HTML
func (db *DocumentationBuilder) BuildHTML() string {
	gen := NewHTMLGenerator(db.sourceInfo, db.graph)
	return gen.GenerateProjectDoc()
}
