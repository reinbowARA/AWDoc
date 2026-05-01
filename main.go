package main

import (
	"github.com/reinbowARA/AWDoc/pkg/analyzer"
	"github.com/reinbowARA/AWDoc/pkg/generator"
	"github.com/reinbowARA/AWDoc/pkg/parser"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// флаги командной строки
	sourceDir := flag.String("source", ".", "Source code directory to analyze")
	language := flag.String("lang", "go", "Programming language (go, etc.)")
	outputFile := flag.String("output", "", "Output file path (default: output/analysis.md)")
	outputDir := flag.String("output-dir", "output", "Output directory for documentation")
	format := flag.String("format", "markdown", "Output format (markdown, html)")
	help := flag.Bool("help", false, "Show help message")

	flag.Parse()

	// показываем справку если требуется
	if *help {
		printHelp()
		os.Exit(0)
	}

	// валидация
	if *sourceDir == "" {
		fmt.Fprintf(os.Stderr, "Error: source directory is required\n")
		os.Exit(1)
	}

	// Если outputFile не указан, используем значение по умолчанию
	finalOutputFile := *outputFile
	if finalOutputFile == "" {
		// Определяем правильное расширение в зависимости от формата
		ext := ".md"
		if *format == "html" {
			ext = ".html"
		}
		finalOutputFile = filepath.Join(*outputDir, "analysis"+ext)
	} else if finalOutputFile != "" && *outputFile == "" {
		// Если outputFile пуст но не явно передан, обновляем расширение на основе формата
		baseName := filepath.Base(finalOutputFile)
		basePath := filepath.Dir(finalOutputFile)

		// Удаляем старое расширение
		name := strings.TrimSuffix(baseName, filepath.Ext(baseName))

		// Добавляем правильное расширение
		ext := ".md"
		if *format == "html" {
			ext = ".html"
		}
		finalOutputFile = filepath.Join(basePath, name+ext)
	}

	fmt.Printf("🔍 Analyzing %s code in: %s\n", *language, *sourceDir)

	// парсим проект
	sourceInfo, err := parser.ParseProject(*sourceDir, *language)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Parse error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Found %d packages\n", len(sourceInfo.Packages))

	// анализируем зависимости
	fmt.Println("📊 Analyzing dependencies...")
	analyzer := analyzer.NewAnalyzer(sourceInfo)
	graph, err := analyzer.Analyze()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Analysis error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Found %d cycles\n", len(graph.Cycles))
	fmt.Printf("✓ Identified %d god objects\n", len(graph.GodObjects))

	// генерируем документацию
	fmt.Printf("📝 Generating %s documentation...\n", *format)

	builder := generator.NewDocumentationBuilder(sourceInfo, graph)
	var docs string

	switch *format {
	case "markdown":
		docs = builder.BuildMarkdown()
	case "html":
		docs = builder.BuildHTML()
	default:
		fmt.Fprintf(os.Stderr, "❌ Unknown format: %s\n", *format)
		os.Exit(1)
	}

	// сохраняем документацию
	err = saveDocumentation(finalOutputFile, docs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Save error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Documentation saved to: %s\n", finalOutputFile)

	// печатаем статистику
	printStatistics(sourceInfo, graph, analyzer)
}

// saveDocumentation сохраняет документацию в файл
func saveDocumentation(filePath string, content string) error {
	// создаем директорию если её нет
	dir := filepath.Dir(filePath)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	// пишем файл
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return err
	}

	return nil
}

// printStatistics выводит статистику
func printStatistics(sourceInfo *parser.SourceInfo, graph *analyzer.DependencyGraph, an *analyzer.Analyzer) {
	fmt.Println("\n📈 Statistics:")
	fmt.Println("═══════════════════════════════════")
	fmt.Printf("Total packages: %d\n", len(sourceInfo.Packages))

	totalElements := 0
	totalExported := 0
	for _, pkg := range sourceInfo.Packages {
		totalElements += len(pkg.Elements)
		totalExported += len(pkg.ExportedAPI)
	}

	fmt.Printf("Total code elements: %d\n", totalElements)
	fmt.Printf("Exported elements: %d\n", totalExported)

	fmt.Printf("\nArchitectural layers: %d\n", len(graph.Layers))
	fmt.Printf("Circular dependencies: %d\n", len(graph.Cycles))
	fmt.Printf("Complex packages: %d\n", len(graph.GodObjects))

	if len(graph.Cycles) > 0 {
		fmt.Println("\n⚠️  Detected dependency cycles:")
		for _, cycle := range graph.Cycles {
			fmt.Printf("  - %v\n", cycle)
		}
	}

	if len(graph.GodObjects) > 0 {
		fmt.Println("\n⚠️  Complex packages (god objects):")
		for _, pkg := range graph.GodObjects {
			node := graph.Nodes[pkg]
			fmt.Printf("  - %s (complexity: %d, deps: %d)\n", pkg, node.Complexity, len(node.Dependencies))
		}
	}
}

// printHelp выводит справку
func printHelp() {
	fmt.Println(`AWDoc - Go API Documentation Generator

USAGE:
  awdoc [flags]

FLAGS:
  -source string
        Source code directory to analyze (default: ".")
  
  -lang string
        Programming language (default: "go")
  
  -output-dir string
        Output directory for documentation (default: "output")
  
  -output string
        Output file path (overrides -output-dir)
  
  -format string
        Output format: markdown or html (default: "markdown")
  
  -help
        Show this help message

EXAMPLES:
  # Analyze current directory, save to output/analysis.md
  awdoc

  # Analyze specific directory
  awdoc -source ./pkg -lang go

  # Custom output directory
  awdoc -source . -output-dir ./docs

  # Custom output file
  awdoc -source . -output ./documentation/api.md

  # HTML format
  awdoc -format html -output output/api.html

VERSION:
  AWDoc 1.0.0`)
}
