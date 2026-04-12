package generator

import (
	"awdoc/pkg/analyzer"
	"awdoc/pkg/parser"
	"fmt"
	"sort"
	"strings"
)

// MarkdownGenerator генерирует документацию в формате Markdown
type MarkdownGenerator struct {
	sourceInfo *parser.SourceInfo
	graph      *analyzer.DependencyGraph
}

// NewMarkdownGenerator создает генератор Markdown
func NewMarkdownGenerator(sourceInfo *parser.SourceInfo, graph *analyzer.DependencyGraph) *MarkdownGenerator {
	return &MarkdownGenerator{
		sourceInfo: sourceInfo,
		graph:      graph,
	}
}

// GenerateProjectDoc генерирует документацию всего проекта
func (mg *MarkdownGenerator) GenerateProjectDoc() string {
	var doc strings.Builder

	doc.WriteString("# API Documentation\n\n")
	doc.WriteString("This documentation was automatically generated from source code.\n\n")

	// оглавление
	doc.WriteString("## Table of Contents\n\n")
	doc.WriteString("- [Project Overview](#project-overview)\n")
	doc.WriteString("- [Packages](#packages)\n")
	doc.WriteString("- [Architecture Analysis](#architecture-analysis)\n\n")

	// обзор проекта
	doc.WriteString("## Project Overview\n\n")
	doc.WriteString(fmt.Sprintf("**Total Packages:** %d\n\n", len(mg.sourceInfo.Packages)))
	doc.WriteString(fmt.Sprintf("**Total Elements:** %d\n\n", mg.countElements()))

	// информация о пакетах
	doc.WriteString("## Packages\n\n")

	// сортируем пакеты по имени
	var pkgNames []string
	for name := range mg.sourceInfo.Packages {
		pkgNames = append(pkgNames, name)
	}
	sort.Strings(pkgNames)

	for _, pkgName := range pkgNames {
		pkg := mg.sourceInfo.Packages[pkgName]
		doc.WriteString(mg.generatePackageDoc(pkg))
	}

	// анализ архитектуры
	doc.WriteString("## Architecture Analysis\n\n")
	doc.WriteString(mg.generateArchitectureAnalysis())

	return doc.String()
}

// generatePackageDoc генерирует документацию для пакета
func (mg *MarkdownGenerator) generatePackageDoc(pkg *parser.Package) string {
	var doc strings.Builder

	doc.WriteString(fmt.Sprintf("### Package: `%s`\n\n", pkg.Name))

	if pkg.Doc != "" {
		doc.WriteString(fmt.Sprintf("**Description:** %s\n\n", pkg.Doc))
	}

	// coverage информация
	if pkg.TotalElements > 0 {
		coveragePercent := int(pkg.Coverage)
		coverageClass := "coverage-none"
		if pkg.Coverage >= 80 {
			coverageClass = "coverage-high"
		} else if pkg.Coverage >= 50 {
			coverageClass = "coverage-medium"
		} else if pkg.Coverage > 0 {
			coverageClass = "coverage-low"
		}

		var coverageEmoji string
		switch coverageClass {
		case "coverage-high":
			coverageEmoji = "✅"
		case "coverage-medium":
			coverageEmoji = "🟡"
		case "coverage-low":
			coverageEmoji = "🔴"
		default:
			coverageEmoji = "⚪"
		}

		doc.WriteString(fmt.Sprintf("**Coverage:** %s %d%% (%d/%d)\n\n", coverageEmoji, coveragePercent, pkg.TestedElements, pkg.TotalElements))
	}

	// импорты
	if len(pkg.Imports) > 0 {
		doc.WriteString("**Imports:**\n")
		for imp := range pkg.Imports {
			doc.WriteString(fmt.Sprintf("- `%s`\n", imp))
		}
		doc.WriteString("\n")
	}

	// экспортируемые элементы
	if len(pkg.ExportedAPI) > 0 {
		doc.WriteString("#### Exported Elements\n\n")
		doc.WriteString(mg.generateElementsDoc(pkg.ExportedAPI))
	}

	// все элементы
	if len(pkg.Elements) > len(pkg.ExportedAPI) {
		doc.WriteString("#### Internal Elements\n\n")

		internalElems := make([]parser.CodeElement, 0)
		exportedSet := make(map[string]bool)
		for _, elem := range pkg.ExportedAPI {
			exportedSet[elem.Name] = true
		}
		for _, elem := range pkg.Elements {
			if !exportedSet[elem.Name] {
				internalElems = append(internalElems, elem)
			}
		}

		doc.WriteString(mg.generateElementsDoc(internalElems))
	}

	doc.WriteString("\n---\n\n")

	return doc.String()
}

// generateElementsDoc генерирует документацию для элементов
func (mg *MarkdownGenerator) generateElementsDoc(elements []parser.CodeElement) string {
	var doc strings.Builder

	// группируем элементы по типам
	byType := make(map[parser.ElementType][]parser.CodeElement)
	for _, elem := range elements {
		byType[elem.Type] = append(byType[elem.Type], elem)
	}

	typeOrder := []parser.ElementType{
		parser.ElementFunc,
		parser.ElementMethod,
		parser.ElementType_,
		parser.ElementStruct,
		parser.ElementInterface,
		parser.ElementConst,
		parser.ElementVar,
	}

	for _, elemType := range typeOrder {
		items, exists := byType[elemType]
		if !exists || len(items) == 0 {
			continue
		}

		typeLabel := elementTypeLabel(elemType)
		doc.WriteString(fmt.Sprintf("**%s:**\n\n", typeLabel))

		for _, elem := range items {
			doc.WriteString(fmt.Sprintf("- **`%s`** (%s)\n", elem.Name, elemType))
			if elem.Signature != "" {
				doc.WriteString(fmt.Sprintf("  ```go\n  %s\n  ```\n", elem.Signature))
			}
			if elem.Doc != "" {
				doc.WriteString(fmt.Sprintf("  %s\n", elem.Doc))
			}
			doc.WriteString("\n")
		}
	}

	return doc.String()
}

// generateArchitectureAnalysis генерирует анализ архитектуры
func (mg *MarkdownGenerator) generateArchitectureAnalysis() string {
	var doc strings.Builder

	// слои архитектуры
	if len(mg.graph.Layers) > 0 {
		doc.WriteString("### Architectural Layers\n\n")
		for i, layer := range mg.graph.Layers {
			doc.WriteString(fmt.Sprintf("**Layer %d:**\n", i))
			for _, pkg := range layer {
				doc.WriteString(fmt.Sprintf("- %s\n", pkg))
			}
			doc.WriteString("\n")
		}
	}

	// циклические зависимости
	if len(mg.graph.Cycles) > 0 {
		doc.WriteString("### ⚠️  Circular Dependencies Detected\n\n")
		for _, cycle := range mg.graph.Cycles {
			doc.WriteString(fmt.Sprintf("- %s\n", strings.Join(cycle, " → ")))
		}
		doc.WriteString("\n")
	}

	// "божественные" объекты
	if len(mg.graph.GodObjects) > 0 {
		doc.WriteString("### Complex Packages (God Objects)\n\n")
		doc.WriteString("Packages with high complexity that might need refactoring:\n\n")
		for _, pkg := range mg.graph.GodObjects {
			node := mg.graph.Nodes[pkg]
			doc.WriteString(fmt.Sprintf("- **%s** (Complexity: %d, Dependencies: %d)\n",
				pkg, node.Complexity, len(node.Dependencies)))
		}
		doc.WriteString("\n")
	}

	// граф зависимостей в текстовом виде
	doc.WriteString("### Dependency Graph\n\n")
	doc.WriteString("```\n")
	for pkg, deps := range mg.graph.Edges {
		if len(deps) > 0 {
			for _, dep := range deps {
				doc.WriteString(fmt.Sprintf("%s → %s\n", pkg, dep))
			}
		} else {
			doc.WriteString(fmt.Sprintf("%s (no dependencies)\n", pkg))
		}
	}
	doc.WriteString("```\n\n")

	return doc.String()
}

// countElements подсчитывает общее количество элементов
func (mg *MarkdownGenerator) countElements() int {
	count := 0
	for _, pkg := range mg.sourceInfo.Packages {
		count += len(pkg.Elements)
	}
	return count
}

// elementTypeLabel возвращает читаемое имя типа элемента
func elementTypeLabel(elemType parser.ElementType) string {
	switch elemType {
	case parser.ElementFunc:
		return "Functions"
	case parser.ElementMethod:
		return "Methods"
	case parser.ElementType_:
		return "Types"
	case parser.ElementStruct:
		return "Structs"
	case parser.ElementInterface:
		return "Interfaces"
	case parser.ElementConst:
		return "Constants"
	case parser.ElementVar:
		return "Variables"
	default:
		return "Other"
	}
}
