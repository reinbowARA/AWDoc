package generator

import (
	"github.com/reinbowARA/AWDoc/pkg/analyzer"
	"github.com/reinbowARA/AWDoc/pkg/parser"
	"fmt"
	"html"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// HTMLGenerator генерирует документацию в формате HTML
type HTMLGenerator struct {
	sourceInfo     *parser.SourceInfo
	graph          *analyzer.DependencyGraph
	templates      map[string]string
	templatesDir   string
}

// NewHTMLGenerator создает генератор HTML
func NewHTMLGenerator(sourceInfo *parser.SourceInfo, graph *analyzer.DependencyGraph) *HTMLGenerator {
	return &HTMLGenerator{
		sourceInfo:   sourceInfo,
		graph:        graph,
		templatesDir: "pkg/generator/templates",
	}
}

// loadTemplate загружает один шаблон из файла
func (hg *HTMLGenerator) loadTemplate(name string) (string, error) {
	path := filepath.Join(hg.templatesDir, name+".html")
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// loadAllTemplates загружает все шаблоны
func (hg *HTMLGenerator) loadAllTemplates() error {
	if hg.templates == nil {
		hg.templates = make(map[string]string)
	}
	
	templates := []string{"head", "nav", "footer", "statistics", "architecture", "packages", "elements", "mermaid"}
	
	for _, name := range templates {
		content, err := hg.loadTemplate(name)
		if err != nil {
			return fmt.Errorf("failed to load template %s: %v", name, err)
		}
		hg.templates[name] = content
	}
	return nil
}

// GenerateProjectDoc генерирует полную HTML документацию проекта
func (hg *HTMLGenerator) GenerateProjectDoc() string {
	// Загружаем шаблоны
	if err := hg.loadAllTemplates(); err != nil {
		return fmt.Sprintf("Error loading templates: %v", err)
	}

	var doc strings.Builder

	// HTML заголовок и стили - берем содержимое head.html (которое начинается с <!DOCTYPE html>)
	doc.WriteString(hg.templates["head"])
	// открываем body
	doc.WriteString("<body>\n")
	doc.WriteString("  <div class=\"container\">\n")

	// навигация
	doc.WriteString(hg.templates["nav"])

	// основной контент
	doc.WriteString("    <main class=\"content\">\n")

	// заголовок и обзор
	doc.WriteString("      <section id=\"overview\">\n")
	doc.WriteString("        <h1>API Documentation</h1>\n")
	doc.WriteString("        <p class=\"subtitle\">Automatically generated from source code</p>\n")
	doc.WriteString(hg.generateProjectOverview())
	doc.WriteString("      </section>\n\n")

	// пакеты
	doc.WriteString("      <section id=\"packages\">\n")
	doc.WriteString("        <h2>Packages</h2>\n")
	doc.WriteString(hg.generatePackagesSection())
	doc.WriteString("      </section>\n\n")

	// архитектурный анализ
	doc.WriteString("      <section id=\"architecture\">\n")
	doc.WriteString("        <h2>Architecture Analysis</h2>\n")
	doc.WriteString(hg.generateArchitectureSection())
	doc.WriteString("      </section>\n\n")

	// статистика
	doc.WriteString("      <section id=\"statistics\">\n")
	doc.WriteString("        <h2>Statistics</h2>\n")
	doc.WriteString(hg.generateStatistics())
	doc.WriteString("      </section>\n")

	doc.WriteString("    </main>\n")

	// footer
	doc.WriteString(hg.templates["footer"])

	doc.WriteString("  </div>\n")
	doc.WriteString("</body>\n")
	doc.WriteString("</html>\n")

	return doc.String()
}

// generateProjectOverview генерирует обзор проекта
func (hg *HTMLGenerator) generateProjectOverview() string {
	var html strings.Builder

	html.WriteString("      <div class=\"stats-grid\">\n")
	html.WriteString("        <div class=\"stat-card\">\n")
	html.WriteString(fmt.Sprintf("          <div class=\"stat-number\">%d</div>\n", len(hg.sourceInfo.Packages)))
	html.WriteString("          <div class=\"stat-label\">Total Packages</div>\n")
	html.WriteString("        </div>\n")

	elementCount := hg.countElements()
	html.WriteString("        <div class=\"stat-card\">\n")
	html.WriteString(fmt.Sprintf("          <div class=\"stat-number\">%d</div>\n", elementCount))
	html.WriteString("          <div class=\"stat-label\">Total Elements</div>\n")
	html.WriteString("        </div>\n")

	// Среднее покрытие
	avgCoverage := hg.calculateAverageCoverage()
	html.WriteString("        <div class=\"stat-card\">\n")
	html.WriteString(fmt.Sprintf("          <div class=\"stat-number\">%d%%</div>\n", int(avgCoverage)))
	html.WriteString("          <div class=\"stat-label\">Avg Test Coverage</div>\n")
	html.WriteString("        </div>\n")

	html.WriteString("        <div class=\"stat-card\">\n")
	html.WriteString(fmt.Sprintf("          <div class=\"stat-number\">%d</div>\n", len(hg.graph.Cycles)))
	html.WriteString("          <div class=\"stat-label\">Circular Dependencies</div>\n")
	html.WriteString("        </div>\n")

	html.WriteString("        <div class=\"stat-card\">\n")
	html.WriteString(fmt.Sprintf("          <div class=\"stat-number\">%d</div>\n", len(hg.graph.GodObjects)))
	html.WriteString("          <div class=\"stat-label\">Complex Packages</div>\n")
	html.WriteString("        </div>\n")

	html.WriteString("      </div>\n")

	return html.String()
}

// generatePackagesSection генерирует секцию с пакетами
func (hg *HTMLGenerator) generatePackagesSection() string {
	var html strings.Builder

	// сортируем пакеты по имени
	var pkgNames []string
	for name := range hg.sourceInfo.Packages {
		pkgNames = append(pkgNames, name)
	}
	sort.Strings(pkgNames)

	for _, pkgName := range pkgNames {
		pkg := hg.sourceInfo.Packages[pkgName]
		html.WriteString(hg.generatePackageHTML(pkg))
	}

	return html.String()
}

// generatePackageHTML генерирует HTML для одного пакета
func (hg *HTMLGenerator) generatePackageHTML(pkg *parser.Package) string {
	var html strings.Builder

	html.WriteString("      <div class=\"package\">\n")

	// Collapsible header for package info
	html.WriteString("        <div class=\"collapsible-header\">\n")
	html.WriteString("          <span class=\"toggle-icon collapsed\">▶</span>\n")
	html.WriteString(fmt.Sprintf("          <span class=\"package-name\">%s</span>\n", escapeHTML(pkg.Name)))
	html.WriteString("        </div>\n")
	html.WriteString("        <div class=\"collapsible-content\">\n")

	// описание пакета
	if pkg.Doc != "" {
		html.WriteString(fmt.Sprintf("        <p class=\"package-description\">%s</p>\n", escapeHTML(pkg.Doc)))
	}

	// информация о покрытии тестами
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

		html.WriteString("        <div style=\"margin: 1rem 0;\">\n")
		html.WriteString(fmt.Sprintf("          <span class=\"coverage-badge %s\">📊 Coverage: %d%% (%d/%d)</span>\n",
			coverageClass, coveragePercent, pkg.TestedElements, pkg.TotalElements))
		html.WriteString(fmt.Sprintf("          <div class=\"coverage-bar\">\n"))
		html.WriteString(fmt.Sprintf("            <div class=\"coverage-fill\" style=\"width: %.1f%%;\"></div>\n", pkg.Coverage))
		html.WriteString(fmt.Sprintf("          </div>\n"))
		html.WriteString("        </div>\n")
	}

	// импорты
	if len(pkg.Imports) > 0 {
		html.WriteString("        <div class=\"imports\">\n")
		html.WriteString("          <div class=\"imports-title\">Imports:</div>\n")
		html.WriteString("          <ul class=\"import-list\">\n")
		for imp := range pkg.Imports {
			html.WriteString(fmt.Sprintf("            <li class=\"import-item\">%s</li>\n", escapeHTML(imp)))
		}
		html.WriteString("          </ul>\n")
		html.WriteString("        </div>\n")
	}

	// экспортируемые элементы
	if len(pkg.ExportedAPI) > 0 {
		html.WriteString("        <h4>Exported Elements</h4>\n")
		html.WriteString(hg.generateElementsHTML(pkg.ExportedAPI))
	}

	// внутренние элементы
	if len(pkg.Elements) > len(pkg.ExportedAPI) {
		html.WriteString("        <h4>Internal Elements</h4>\n")

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

		html.WriteString(hg.generateElementsHTML(internalElems))
	}

	html.WriteString("        </div>\n")
	html.WriteString("      </div>\n")

	return html.String()
}

// generateElementsHTML генерирует HTML для элементов
func (hg *HTMLGenerator) generateElementsHTML(elements []parser.CodeElement) string {
	var html strings.Builder

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
		html.WriteString("        <div class=\"elements\">\n")
		html.WriteString(fmt.Sprintf("          <div class=\"element-type\">%s</div>\n", typeLabel))
		html.WriteString("          <ul class=\"element-list\">\n")

		for _, elem := range items {
			html.WriteString("            <li class=\"element-item\">\n")

			// Special handling for structs - use collapsible header
			if elem.Type == parser.ElementStruct && len(elem.Fields) > 0 {
				html.WriteString("              <div class=\"collapsible-header\">\n")
				html.WriteString("                <span class=\"toggle-icon collapsed\">▶</span>\n")
				html.WriteString(fmt.Sprintf("                <span class=\"element-name\">%s</span>\n", escapeHTML(elem.Name)))
				html.WriteString(fmt.Sprintf("                <span class=\"element-type-badge\">%s</span>\n", elemType))
				html.WriteString("              </div>\n")
				html.WriteString("              <div class=\"collapsible-content\">\n")

				if elem.Doc != "" {
					html.WriteString(fmt.Sprintf("                <div class=\"element-doc\">%s</div>\n", escapeHTML(elem.Doc)))
				}

				html.WriteString("                <div class=\"struct-fields\">\n")
				for _, field := range elem.Fields {
					fieldName := escapeHTML(field.Name)
					fieldType := escapeHTML(field.Type)
					if !field.Exported {
						html.WriteString(fmt.Sprintf("                  <div class=\"struct-field struct-field-private\">\n"))
					} else {
						html.WriteString(fmt.Sprintf("                  <div class=\"struct-field\">\n"))
					}
					html.WriteString(fmt.Sprintf("                    <span class=\"struct-field-name\">%s</span> <span class=\"struct-field-type\">%s</span>\n", fieldName, fieldType))
					if field.Doc != "" {
						html.WriteString(fmt.Sprintf("                    <div class=\"struct-field-doc\">%s</div>\n", escapeHTML(field.Doc)))
					}
					html.WriteString("                  </div>\n")
				}
				html.WriteString("                </div>\n")
				html.WriteString("              </div>\n")
			} else {
				// Regular elements
				html.WriteString(fmt.Sprintf("              <span class=\"element-name\">%s</span>\n", escapeHTML(elem.Name)))
				html.WriteString(fmt.Sprintf("              <span class=\"element-type-badge\">%s</span>\n", elemType))

				// Значок для методов без тестов (только для функций и методов)
				if (elem.Type == parser.ElementFunc || elem.Type == parser.ElementMethod) && !elem.HasTests {
					html.WriteString("              <span class=\"no-test-warning\">⚠️ NO TEST</span>\n")
				} else if elem.HasTests && elem.TestName != "" {
					html.WriteString(fmt.Sprintf("              <span class=\"coverage-badge coverage-high\">✅ %s</span>\n", escapeHTML(elem.TestName)))
				}

				if elem.Signature != "" {
					html.WriteString(fmt.Sprintf("              <div class=\"element-signature\">%s</div>\n", escapeHTML(elem.Signature)))
				}

				if elem.Doc != "" {
					html.WriteString(fmt.Sprintf("              <div class=\"element-doc\">%s</div>\n", escapeHTML(elem.Doc)))
				}
			}

			html.WriteString("            </li>\n")
		}

		html.WriteString("          </ul>\n")
		html.WriteString("        </div>\n")
	}

	return html.String()
}

// generateArchitectureSection генерирует секцию архитектурного анализа
func (hg *HTMLGenerator) generateArchitectureSection() string {
	var html strings.Builder

	// слои архитектуры
	if len(hg.graph.Layers) > 0 {
		html.WriteString("      <h3>Architectural Layers</h3>\n")
		html.WriteString(hg.generateLayersDiagram())
		html.WriteString("      <div class=\"layers\">\n")
		for i, layer := range hg.graph.Layers {
			html.WriteString("        <div class=\"layer\">\n")
			html.WriteString(fmt.Sprintf("          <div class=\"layer-title\">Layer %d</div>\n", i))
			html.WriteString("          <ul class=\"layer-packages\">\n")
			for _, pkg := range layer {
				html.WriteString(fmt.Sprintf("            <li class=\"layer-package\">%s</li>\n", escapeHTML(pkg)))
			}
			html.WriteString("          </ul>\n")
			html.WriteString("        </div>\n")
		}
		html.WriteString("      </div>\n")
	}

	// циклические зависимости
	if len(hg.graph.Cycles) > 0 {
		html.WriteString("      <h3>⚠️ Circular Dependencies Detected</h3>\n")
		html.WriteString("      <div class=\"warning-box\">\n")
		html.WriteString("        <div class=\"warning-title\">Found " + fmt.Sprintf("%d", len(hg.graph.Cycles)) + " circular dependencies</div>\n")
		html.WriteString("        <ul class=\"warning-list\">\n")
		for _, cycle := range hg.graph.Cycles {
			cycleStr := strings.Join(cycle, " → ")
			html.WriteString(fmt.Sprintf("          <li class=\"warning-item\">%s</li>\n", escapeHTML(cycleStr)))
		}
		html.WriteString("        </ul>\n")
		html.WriteString("      </div>\n")
	}

	// "божественные" объекты
	if len(hg.graph.GodObjects) > 0 {
		html.WriteString("      <h3>Complex Packages (God Objects)</h3>\n")
		html.WriteString("      <p>Packages with high complexity that might need refactoring:</p>\n")
		html.WriteString("      <div class=\"warning-box\" style=\"background: #f8f3cd; border-color: #ffc107;\">\n")
		html.WriteString("        <ul class=\"warning-list\">\n")
		for _, pkg := range hg.graph.GodObjects {
			node := hg.graph.Nodes[pkg]
			html.WriteString("          <li class=\"warning-item\" style=\"color: #856404;\">\n")
			html.WriteString(fmt.Sprintf("            <strong>%s</strong>\n", escapeHTML(pkg)))
			html.WriteString(fmt.Sprintf("            <span class=\"complexity-indicator\">Complexity: %d | Dependencies: %d</span>\n",
				node.Complexity, len(node.Dependencies)))
			html.WriteString("          </li>\n")
		}
		html.WriteString("        </ul>\n")
		html.WriteString("      </div>\n")
	}

	// граф зависимостей
	html.WriteString("      <h3>Dependency Graph</h3>\n")
	html.WriteString(hg.generateDependencyDiagram())
	html.WriteString("      <div class=\"dependency-graph\">\n")
	for pkg, deps := range hg.graph.Edges {
		if len(deps) > 0 {
			for _, dep := range deps {
				depStr := fmt.Sprintf("%s → %s", pkg, dep)
				html.WriteString(fmt.Sprintf("        <div class=\"dependency-line\">%s</div>\n", escapeHTML(depStr)))
			}
		} else {
			html.WriteString(fmt.Sprintf("        <div class=\"dependency-line\">%s (no dependencies)</div>\n", escapeHTML(pkg)))
		}
	}
	html.WriteString("      </div>\n")

	return html.String()
}

// generateStatistics генерирует секцию со статистикой
func (hg *HTMLGenerator) generateStatistics() string {
	var html strings.Builder

	html.WriteString("      <div class=\"stats-grid\">\n")

	// общая статистика
	var totalFuncs, totalMethods, totalTypes, totalStructs, totalInterfaces int
	for _, pkg := range hg.sourceInfo.Packages {
		for _, elem := range pkg.Elements {
			switch elem.Type {
			case parser.ElementFunc:
				totalFuncs++
			case parser.ElementMethod:
				totalMethods++
			case parser.ElementType_:
				totalTypes++
			case parser.ElementStruct:
				totalStructs++
			case parser.ElementInterface:
				totalInterfaces++
			}
		}
	}

	stats := map[string]int{
		"Functions":  totalFuncs,
		"Methods":    totalMethods,
		"Types":      totalTypes,
		"Structs":    totalStructs,
		"Interfaces": totalInterfaces,
		"Layers":     len(hg.graph.Layers),
	}

	for label, count := range stats {
		html.WriteString("        <div class=\"stat-card\">\n")
		html.WriteString(fmt.Sprintf("          <div class=\"stat-number\">%d</div>\n", count))
		html.WriteString(fmt.Sprintf("          <div class=\"stat-label\">%s</div>\n", label))
		html.WriteString("        </div>\n")
	}

	html.WriteString("      </div>\n")

	return html.String()
}

// escapeHTML экранирует HTML специальные символы
func escapeHTML(s string) string {
	return html.EscapeString(s)
}

// countElements подсчитывает количество элементов
func (hg *HTMLGenerator) countElements() int {
	count := 0
	for _, pkg := range hg.sourceInfo.Packages {
		count += len(pkg.Elements)
	}
	return count
}

// calculateAverageCoverage вычисляет среднее покрытие по всем пакетам
func (hg *HTMLGenerator) calculateAverageCoverage() float64 {
	if len(hg.sourceInfo.Packages) == 0 {
		return 0
	}

	totalCoverage := 0.0
	validPackages := 0

	for _, pkg := range hg.sourceInfo.Packages {
		if pkg.TotalElements > 0 {
			totalCoverage += pkg.Coverage
			validPackages++
		}
	}

	if validPackages == 0 {
		return 0
	}

	return totalCoverage / float64(validPackages)
}

// generateDependencyDiagram генерирует диаграмму зависимостей Mermaid
func (hg *HTMLGenerator) generateDependencyDiagram() string {
	if len(hg.graph.Edges) == 0 {
		return "<p>No dependencies found</p>\n"
	}

	var diagram strings.Builder
	diagram.WriteString("      <div class=\"mermaid\">\n")
	diagram.WriteString("        graph LR\n")

	// Определяем цвет для каждого пакета на основе сложности
	nodes := hg.graph.Nodes
	for pkg := range hg.graph.Edges {
		var color string
		if node, exists := nodes[pkg]; exists {
			if node.Complexity > 20 {
				color = "#ff6b6b"
			} else if node.Complexity > 10 {
				color = "#ffd700"
			} else {
				color = "#90ee90"
			}
		} else {
			color = "#90ee90"
		}

		pkgLabel := pkg
		if node, exists := nodes[pkg]; exists {
			pkgLabel = fmt.Sprintf("%s<br/>(%d)", pkg, len(node.Dependencies))
		}

		diagram.WriteString(fmt.Sprintf("        %s[\"%s\"]:::node%d\n",
			strings.ReplaceAll(pkg, "/", "_"), pkgLabel, hashColor(color)))
	}

	// Добавляем связи
	for pkg, deps := range hg.graph.Edges {
		for _, dep := range deps {
			diagram.WriteString(fmt.Sprintf("        %s --> %s\n",
				strings.ReplaceAll(pkg, "/", "_"),
				strings.ReplaceAll(dep, "/", "_")))
		}
	}

	// Добавляем стили
	diagram.WriteString("        classDef node1 fill:#90ee90,stroke:#333,stroke-width:2px,color:#000\n")
	diagram.WriteString("        classDef node2 fill:#ffd700,stroke:#333,stroke-width:2px,color:#000\n")
	diagram.WriteString("        classDef node3 fill:#ff6b6b,stroke:#333,stroke-width:2px,color:#fff\n")

	diagram.WriteString("      </div>\n")
	return diagram.String()
}

// generateLayersDiagram генерирует диаграмму архитектурных слоев
func (hg *HTMLGenerator) generateLayersDiagram() string {
	if len(hg.graph.Layers) == 0 {
		return "<p>No architectural layers found</p>\n"
	}

	var diagram strings.Builder
	diagram.WriteString("      <div class=\"mermaid\">\n")
	diagram.WriteString("        graph TD\n")

	colors := []string{"#90ee90", "#87ceeb", "#ffd700", "#ff6b6b", "#dda0dd"}

	for i, layer := range hg.graph.Layers {
		layerName := fmt.Sprintf("Layer%d", i)
		layerLabel := fmt.Sprintf("Layer %d", i)

		diagram.WriteString(fmt.Sprintf("        %s[\"%s", layerName, layerLabel))

		for _, pkg := range layer {
			diagram.WriteString(fmt.Sprintf("<br/>%s", pkg))
		}

		diagram.WriteString("\"]\n")

		// Добавляем стиль для слоя
		colorIdx := i % len(colors)
		diagram.WriteString(fmt.Sprintf("        style %s fill:%s,stroke:#333,stroke-width:2px,color:#000\n",
			layerName, colors[colorIdx]))
	}

	// Добавляем связи между слоями
	for i := 0; i < len(hg.graph.Layers)-1; i++ {
		diagram.WriteString(fmt.Sprintf("        Layer%d --> Layer%d\n", i, i+1))
	}

	diagram.WriteString("      </div>\n")
	return diagram.String()
}

// hashColor возвращает индекс стиля на основе цвета
func hashColor(color string) int {
	if strings.Contains(color, "90ee90") {
		return 1
	} else if strings.Contains(color, "ffd700") {
		return 2
	}
	return 3
}