package analyzer

import (
	"github.com/reinbowARA/AWDoc/pkg/parser"
	"testing"
)

// TestAnalyzer тестирует анализатор зависимостей
func TestAnalyzer(t *testing.T) {
	// Создаем тестовую SourceInfo
	sourceInfo := &parser.SourceInfo{
		Files: []string{},
		Packages: map[string]*parser.Package{
			"pkg1": {
				Name: "pkg1",
				Path: "pkg1",
				Elements: []parser.CodeElement{
					{
						Name:     "Func1",
						Type:     parser.ElementFunc,
						Exported: true,
					},
				},
				Imports: map[string]bool{
					"pkg2": true,
				},
				ExportedAPI: []parser.CodeElement{
					{
						Name:     "Func1",
						Type:     parser.ElementFunc,
						Exported: true,
					},
				},
			},
			"pkg2": {
				Name: "pkg2",
				Path: "pkg2",
				Elements: []parser.CodeElement{
					{
						Name:     "Func2",
						Type:     parser.ElementFunc,
						Exported: true,
					},
				},
				Imports: map[string]bool{},
				ExportedAPI: []parser.CodeElement{
					{
						Name:     "Func2",
						Type:     parser.ElementFunc,
						Exported: true,
					},
				},
			},
		},
	}

	// Создаем анализатор
	analyzer := NewAnalyzer(sourceInfo)
	graph, err := analyzer.Analyze()
	if err != nil {
		t.Fatalf("Failed to analyze: %v", err)
	}

	if graph == nil {
		t.Fatal("Graph is nil")
	}

	// Проверяем что найдены пакеты
	if len(graph.Nodes) != 2 {
		t.Errorf("Expected 2 nodes, got %d", len(graph.Nodes))
	}

	// Проверяем что найдены зависимости
	if len(graph.Edges) == 0 {
		t.Error("No edges found in graph")
	}

	// Проверяем что pkg1 зависит от pkg2
	pkg1Edges := graph.Edges["pkg1"]
	hasPkg2 := false
	for _, dep := range pkg1Edges {
		if dep == "pkg2" {
			hasPkg2 = true
			break
		}
	}
	if !hasPkg2 {
		t.Error("pkg1 should depend on pkg2")
	}
}

// TestCycleDetection тестирует обнаружение циклов
func TestCycleDetection(t *testing.T) {
	// Создаем граф с циклом: A -> B -> C -> A
	sourceInfo := &parser.SourceInfo{
		Files: []string{},
		Packages: map[string]*parser.Package{
			"A": {
				Name:        "A",
				Path:        "A",
				Elements:    []parser.CodeElement{},
				Imports:     map[string]bool{"B": true},
				ExportedAPI: []parser.CodeElement{},
			},
			"B": {
				Name:        "B",
				Path:        "B",
				Elements:    []parser.CodeElement{},
				Imports:     map[string]bool{"C": true},
				ExportedAPI: []parser.CodeElement{},
			},
			"C": {
				Name:        "C",
				Path:        "C",
				Elements:    []parser.CodeElement{},
				Imports:     map[string]bool{"A": true},
				ExportedAPI: []parser.CodeElement{},
			},
		},
	}

	analyzer := NewAnalyzer(sourceInfo)
	graph, err := analyzer.Analyze()
	if err != nil {
		t.Fatalf("Failed to analyze: %v", err)
	}

	// Проверяем что найдены циклы
	if len(graph.Cycles) == 0 {
		t.Error("No cycles detected but one should exist")
	}
}
