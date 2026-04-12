package generator

import (
	"awdoc/pkg/analyzer"
	"awdoc/pkg/parser"
	"fmt"
	"html"
	"sort"
	"strings"
)

// HTMLGenerator генерирует документацию в формате HTML
type HTMLGenerator struct {
	sourceInfo *parser.SourceInfo
	graph      *analyzer.DependencyGraph
}

// NewHTMLGenerator создает генератор HTML
func NewHTMLGenerator(sourceInfo *parser.SourceInfo, graph *analyzer.DependencyGraph) *HTMLGenerator {
	return &HTMLGenerator{
		sourceInfo: sourceInfo,
		graph:      graph,
	}
}

// GenerateProjectDoc генерирует полную HTML документацию проекта
func (hg *HTMLGenerator) GenerateProjectDoc() string {
	var doc strings.Builder

	// HTML заголовок и стили
	doc.WriteString(hg.generateHTMLHead())

	// открываем body
	doc.WriteString("<body>\n")
	doc.WriteString("  <div class=\"container\">\n")

	// навигация
	doc.WriteString(hg.generateNavigation())

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
	doc.WriteString(hg.generateFooter())

	doc.WriteString("  </div>\n")
	doc.WriteString("</body>\n")
	doc.WriteString("</html>\n")

	return doc.String()
}

// generateHTMLHead генерирует head раздел HTML
func (hg *HTMLGenerator) generateHTMLHead() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>API Documentation</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            line-height: 1.6;
            color: #333;
            background: #f5f5f5;
        }

        .container {
            max-width: 1200px;
            margin: 0 auto;
            background: white;
            box-shadow: 0 0 10px rgba(0,0,0,0.1);
        }

        header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 2rem;
            text-align: center;
        }

        nav {
            background: #f8f9fa;
            border-bottom: 1px solid #e9ecef;
            padding: 1rem 2rem;
            position: sticky;
            top: 0;
            z-index: 100;
        }

        nav ul {
            list-style: none;
            display: flex;
            gap: 2rem;
            flex-wrap: wrap;
        }

        nav a {
            color: #667eea;
            text-decoration: none;
            font-weight: 500;
            transition: color 0.3s;
        }

        nav a:hover {
            color: #764ba2;
        }

        .content {
            padding: 2rem;
        }

        section {
            margin-bottom: 3rem;
        }

        h1 {
            font-size: 2.5rem;
            margin-bottom: 0.5rem;
            color: white;
        }

        h2 {
            font-size: 2rem;
            margin-bottom: 1.5rem;
            padding-bottom: 0.5rem;
            border-bottom: 2px solid #667eea;
            color: #333;
        }

        h3 {
            font-size: 1.5rem;
            margin-top: 2rem;
            margin-bottom: 1rem;
            color: #555;
        }

        h4 {
            font-size: 1.2rem;
            margin-top: 1.5rem;
            margin-bottom: 0.8rem;
            color: #666;
        }

        .subtitle {
            font-size: 1.2rem;
            opacity: 0.95;
        }

        .package {
            background: #f8f9fa;
            border-left: 4px solid #667eea;
            padding: 1.5rem;
            margin-bottom: 2rem;
            border-radius: 4px;
        }

        .package-name {
            font-size: 1.3rem;
            font-weight: bold;
            color: #667eea;
            margin-bottom: 0.5rem;
            font-family: 'Courier New', monospace;
        }

        .package-description {
            color: #666;
            margin-bottom: 1rem;
        }

        .imports {
            margin-bottom: 1.5rem;
        }

        .imports-title {
            font-weight: bold;
            margin-bottom: 0.5rem;
            color: #555;
        }

        .import-list {
            list-style: none;
            padding-left: 1rem;
        }

        .import-item {
            color: #666;
            margin: 0.3rem 0;
            font-family: 'Courier New', monospace;
            font-size: 0.9rem;
        }

        .elements {
            margin-top: 1.5rem;
        }

        .element-type {
            font-weight: bold;
            color: #667eea;
            margin-top: 1rem;
            margin-bottom: 0.5rem;
        }

        .element-list {
            list-style: none;
            padding-left: 1rem;
        }

        .element-item {
            margin: 0.8rem 0;
            padding: 0.8rem;
            background: white;
            border: 1px solid #e9ecef;
            border-radius: 4px;
        }

        .element-name {
            font-family: 'Courier New', monospace;
            font-weight: bold;
            color: #333;
        }

        .element-type-badge {
            display: inline-block;
            background: #667eea;
            color: white;
            padding: 0.2rem 0.6rem;
            border-radius: 3px;
            font-size: 0.85rem;
            margin-left: 0.5rem;
        }

        .element-signature {
            background: #f8f9fa;
            border: 1px solid #e9ecef;
            padding: 0.8rem;
            margin: 0.5rem 0;
            border-radius: 4px;
            font-family: 'Courier New', monospace;
            font-size: 0.9rem;
            overflow-x: auto;
        }

        .element-doc {
            color: #666;
            margin-top: 0.5rem;
            font-size: 0.95rem;
        }

        .layers {
            background: white;
            border: 1px solid #e9ecef;
            border-radius: 4px;
            padding: 1.5rem;
        }

        .layer {
            margin-bottom: 1.5rem;
            padding: 1rem;
            background: #f8f9fa;
            border-left: 3px solid #667eea;
            border-radius: 4px;
        }

        .layer-title {
            font-weight: bold;
            color: #667eea;
            margin-bottom: 0.5rem;
        }

        .layer-packages {
            list-style: none;
            padding-left: 1rem;
        }

        .layer-package {
            color: #666;
            margin: 0.3rem 0;
            font-family: 'Courier New', monospace;
        }

        .warning-box {
            background: #fff3cd;
            border-left: 4px solid #ffc107;
            padding: 1rem;
            margin: 1rem 0;
            border-radius: 4px;
        }

        .warning-title {
            font-weight: bold;
            color: #856404;
            margin-bottom: 0.5rem;
        }

        .warning-list {
            list-style: none;
            padding-left: 1rem;
        }

        .warning-item {
            color: #856404;
            margin: 0.5rem 0;
        }

        .complexity-indicator {
            display: inline-block;
            padding: 0.3rem 0.8rem;
            background: #f8d7da;
            color: #721c24;
            border-radius: 3px;
            font-size: 0.85rem;
        }

        .dependency-graph {
            background: #f8f9fa;
            border: 1px solid #e9ecef;
            padding: 1.5rem;
            border-radius: 4px;
            overflow-x: auto;
        }

        .dependency-line {
            font-family: 'Courier New', monospace;
            font-size: 0.9rem;
            color: #666;
            margin: 0.3rem 0;
        }

        .stats-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 1.5rem;
            margin: 1rem 0;
        }

        .stat-card {
            background: white;
            border: 1px solid #e9ecef;
            border-radius: 4px;
            padding: 1.5rem;
            text-align: center;
        }

        .stat-number {
            font-size: 2rem;
            font-weight: bold;
            color: #667eea;
            margin-bottom: 0.5rem;
        }

        .stat-label {
            color: #666;
            font-size: 0.9rem;
        }

        footer {
            background: #f8f9fa;
            border-top: 1px solid #e9ecef;
            padding: 2rem;
            text-align: center;
            color: #666;
            font-size: 0.9rem;
        }

        code {
            background: #f4f4f4;
            padding: 0.2rem 0.4rem;
            border-radius: 3px;
            font-family: 'Courier New', monospace;
            font-size: 0.9em;
        }

        @media (max-width: 768px) {
            .container {
                box-shadow: none;
            }

            h1 {
                font-size: 1.8rem;
            }

            h2 {
                font-size: 1.5rem;
            }

            nav ul {
                gap: 1rem;
            }

            .stats-grid {
                grid-template-columns: 1fr;
            }

            .content {
                padding: 1rem;
            }
        }
    </style>
</head>
`
}

// generateNavigation генерирует навигационное меню
func (hg *HTMLGenerator) generateNavigation() string {
	return `    <nav>
      <ul>
        <li><a href="#overview">Overview</a></li>
        <li><a href="#packages">Packages</a></li>
        <li><a href="#architecture">Architecture</a></li>
        <li><a href="#statistics">Statistics</a></li>
      </ul>
    </nav>
`
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
	html.WriteString(fmt.Sprintf("        <div class=\"package-name\">%s</div>\n", escapeHTML(pkg.Name)))

	if pkg.Doc != "" {
		html.WriteString(fmt.Sprintf("        <p class=\"package-description\">%s</p>\n", escapeHTML(pkg.Doc)))
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
			html.WriteString(fmt.Sprintf("              <span class=\"element-name\">%s</span>\n", escapeHTML(elem.Name)))
			html.WriteString(fmt.Sprintf("              <span class=\"element-type-badge\">%s</span>\n", elemType))

			if elem.Signature != "" {
				html.WriteString(fmt.Sprintf("              <div class=\"element-signature\">%s</div>\n", escapeHTML(elem.Signature)))
			}

			if elem.Doc != "" {
				html.WriteString(fmt.Sprintf("              <div class=\"element-doc\">%s</div>\n", escapeHTML(elem.Doc)))
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

// generateFooter генерирует footer
func (hg *HTMLGenerator) generateFooter() string {
	return `    <footer>
      <p>Generated by AWDoc - Automated Web Documentation Generator</p>
      <p>This documentation was automatically generated from source code analysis.</p>
    </footer>
`
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
