package analyzer

import (
	"awdoc/pkg/parser"
	"fmt"
	"sort"
)

// DependencyGraph представляет граф зависимостей пакетов
type DependencyGraph struct {
	Nodes      map[string]*PackageNode // узлы графа (пакеты)
	Edges      map[string][]string     // рёбра (зависимости)
	Cycles     [][]string              // циклические зависимости
	Layers     [][]string              // архитектурные слои
	GodObjects []string                // пакеты с большой сложностью (много зависимостей)
}

// PackageNode представляет узел графа (пакет)
type PackageNode struct {
	Package       *parser.Package
	Dependencies  []string // прямые зависимости
	Dependents    []string // пакеты, зависящие от этого
	Complexity    int      // сложность пакета
	ExportedCount int      // кол-во экспортируемых элементов
	InternalCount int      // кол-во внутренних элементов
	CircularDeps  []string // циклические зависимости
}

// Analyzer анализирует исходный код и строит граф зависимостей
type Analyzer struct {
	sourceInfo *parser.SourceInfo
	graph      *DependencyGraph
}

// NewAnalyzer создает новый анализатор
func NewAnalyzer(sourceInfo *parser.SourceInfo) *Analyzer {
	return &Analyzer{
		sourceInfo: sourceInfo,
		graph: &DependencyGraph{
			Nodes:      make(map[string]*PackageNode),
			Edges:      make(map[string][]string),
			Cycles:     [][]string{},
			Layers:     [][]string{},
			GodObjects: []string{},
		},
	}
}

// Analyze анализирует исходный код и строит граф
func (a *Analyzer) Analyze() (*DependencyGraph, error) {
	// создаем узлы графа
	for pkgName, pkg := range a.sourceInfo.Packages {
		node := &PackageNode{
			Package:       pkg,
			Dependencies:  []string{},
			Dependents:    []string{},
			ExportedCount: len(pkg.ExportedAPI),
			InternalCount: len(pkg.Elements) - len(pkg.ExportedAPI),
		}
		a.graph.Nodes[pkgName] = node
	}

	// строим рёбра (зависимости)
	for pkgName, pkg := range a.sourceInfo.Packages {
		node := a.graph.Nodes[pkgName]
		for importPath := range pkg.Imports {
			// проверяем, есть ли этот импорт в нашем проекте
			if importNode, exists := a.graph.Nodes[importPath]; exists {
				node.Dependencies = append(node.Dependencies, importPath)
				importNode.Dependents = append(importNode.Dependents, pkgName)
				a.graph.Edges[pkgName] = append(a.graph.Edges[pkgName], importPath)
			}
		}
	}

	// находим циклические зависимости
	a.detectCycles()

	// вычисляем сложность пакетов
	a.computeComplexity()

	// выявляем "божественные" объекты
	a.identifyGodObjects()

	// выявляем архитектурные слои
	a.identifyLayers()

	return a.graph, nil
}

// detectCycles находит циклические зависимости
func (a *Analyzer) detectCycles() {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	for pkg := range a.graph.Nodes {
		if !visited[pkg] {
			a.detectCyclesDFS(pkg, visited, recStack, []string{})
		}
	}
}

// detectCyclesDFS вспомогательная функция для поиска циклов
func (a *Analyzer) detectCyclesDFS(pkg string, visited map[string]bool, recStack map[string]bool, path []string) {
	visited[pkg] = true
	recStack[pkg] = true
	path = append(path, pkg)

	for _, dep := range a.graph.Edges[pkg] {
		if !visited[dep] {
			a.detectCyclesDFS(dep, visited, recStack, path)
		} else if recStack[dep] {
			// найден цикл
			cycleStart := -1
			for i, p := range path {
				if p == dep {
					cycleStart = i
					break
				}
			}
			if cycleStart != -1 {
				cycle := append(path[cycleStart:], dep)
				a.graph.Cycles = append(a.graph.Cycles, cycle)

				// отмечаем циклические зависимости в узлах
				for _, p := range cycle {
					if node, ok := a.graph.Nodes[p]; ok {
						node.CircularDeps = append(node.CircularDeps, dep)
					}
				}
			}
		}
	}

	recStack[pkg] = false
}

// computeComplexity вычисляет сложность каждого пакета
func (a *Analyzer) computeComplexity() {
	for _, node := range a.graph.Nodes {
		complexity := 0

		// фактор 1: количество элементов
		complexity += len(node.Package.Elements) * 1

		// фактор 2: количество внешних зависимостей
		complexity += len(node.Dependencies) * 3

		// фактор 3: количество пакетов, зависящих от этого (вес)
		complexity += len(node.Dependents) * 2

		// фактор 4: циклические зависимости (штраф)
		complexity += len(node.CircularDeps) * 5

		node.Complexity = complexity
	}
}

// identifyGodObjects выявляет пакеты со слишком большой сложностью
func (a *Analyzer) identifyGodObjects() {
	var complexities []int
	complexityMap := make(map[int]string)

	for pkgName, node := range a.graph.Nodes {
		complexities = append(complexities, node.Complexity)
		complexityMap[node.Complexity] = pkgName
	}

	sort.Ints(complexities)

	// находим пакеты в верхних 20% по сложности
	threshold := len(complexities) / 5
	if threshold < 1 {
		threshold = 1
	}

	seen := make(map[string]bool)
	for i := len(complexities) - 1; i >= len(complexities)-threshold && i >= 0; i-- {
		if pkgName := complexityMap[complexities[i]]; pkgName != "" && !seen[pkgName] {
			a.graph.GodObjects = append(a.graph.GodObjects, pkgName)
			seen[pkgName] = true
		}
	}
}

// identifyLayers выявляет архитектурные слои
func (a *Analyzer) identifyLayers() {
	// простой алгоритм: слои зависимостей
	// слой 0: пакеты без зависимостей
	// слой 1: пакеты, зависящие только от слоя 0
	// и т.д.

	layerAssignment := make(map[string]int)
	var currentLayer []string
	var nextLayer []string

	// слой 0: пакеты без зависимостей
	for pkgName, node := range a.graph.Nodes {
		if len(node.Dependencies) == 0 {
			currentLayer = append(currentLayer, pkgName)
			layerAssignment[pkgName] = 0
		}
	}

	layerNum := 0
	a.graph.Layers = append(a.graph.Layers, currentLayer)

	// остальные слои
	for len(currentLayer) > 0 && len(currentLayer) < len(a.graph.Nodes) {
		nextLayer = []string{}

		for pkgName, node := range a.graph.Nodes {
			if _, assigned := layerAssignment[pkgName]; assigned {
				continue
			}

			// проверяем, все ли зависимости уже в назначенных слоях
			allDepsAssigned := true
			for _, dep := range node.Dependencies {
				if _, assigned := layerAssignment[dep]; !assigned {
					allDepsAssigned = false
					break
				}
			}

			if allDepsAssigned {
				nextLayer = append(nextLayer, pkgName)
				layerAssignment[pkgName] = layerNum + 1
			}
		}

		if len(nextLayer) == 0 {
			break
		}

		a.graph.Layers = append(a.graph.Layers, nextLayer)
		currentLayer = nextLayer
		layerNum++
	}

	// оставшиеся пакеты (с циклами) идут в последний слой
	if len(layerAssignment) < len(a.graph.Nodes) {
		var lastLayer []string
		for pkgName := range a.graph.Nodes {
			if _, assigned := layerAssignment[pkgName]; !assigned {
				lastLayer = append(lastLayer, pkgName)
			}
		}
		if len(lastLayer) > 0 {
			a.graph.Layers = append(a.graph.Layers, lastLayer)
		}
	}
}

// GetDependencyInfo возвращает информацию о зависимостях пакета
func (a *Analyzer) GetDependencyInfo(pkgName string) (string, error) {
	node, exists := a.graph.Nodes[pkgName]
	if !exists {
		return "", fmt.Errorf("package not found: %s", pkgName)
	}

	info := fmt.Sprintf("Package: %s\n", pkgName)
	info += fmt.Sprintf("Complexity: %d\n", node.Complexity)
	info += fmt.Sprintf("Exported Elements: %d\n", node.ExportedCount)
	info += fmt.Sprintf("Internal Elements: %d\n", node.InternalCount)
	info += fmt.Sprintf("Direct Dependencies: %d\n", len(node.Dependencies))
	info += fmt.Sprintf("Dependents: %d\n", len(node.Dependents))

	if len(node.CircularDeps) > 0 {
		info += fmt.Sprintf("⚠️  Circular Dependencies: %v\n", node.CircularDeps)
	}

	return info, nil
}

// GetGraph возвращает граф зависимостей
func (a *Analyzer) GetGraph() *DependencyGraph {
	return a.graph
}
