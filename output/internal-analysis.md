# API Documentation

This documentation was automatically generated from source code.

## Table of Contents

- [Project Overview](#project-overview)
- [Packages](#packages)
- [Architecture Analysis](#architecture-analysis)

## Project Overview

**Total Packages:** 3

**Total Elements:** 28

## Packages

### Package: `analyzer`

**Imports:**
- `awdoc/pkg/parser`
- `fmt`
- `sort`

#### Exported Elements

**Functions:**

- **`NewAnalyzer`** (function)
  ```go
  func NewAnalyzer(sourceInfo *parser.SourceInfo) *Analyzer
  ```
  NewAnalyzer создает новый анализатор

**Methods:**

- **`Analyze`** (method)
  ```go
  func (*Analyzer) Analyze() (*DependencyGraph, error)
  ```
  Analyze анализирует исходный код и строит граф

- **`GetDependencyInfo`** (method)
  ```go
  func (*Analyzer) GetDependencyInfo(pkgName string) (string, error)
  ```
  GetDependencyInfo возвращает информацию о зависимостях пакета

- **`GetGraph`** (method)
  ```go
  func (*Analyzer) GetGraph() *DependencyGraph
  ```
  GetGraph возвращает граф зависимостей

**Structs:**

- **`DependencyGraph`** (struct)
  DependencyGraph представляет граф зависимостей пакетов

- **`PackageNode`** (struct)
  PackageNode представляет узел графа (пакет)

- **`Analyzer`** (struct)
  Analyzer анализирует исходный код и строит граф зависимостей

#### Internal Elements

**Methods:**

- **`detectCycles`** (method)
  ```go
  func (*Analyzer) detectCycles()
  ```
  detectCycles находит циклические зависимости

- **`detectCyclesDFS`** (method)
  ```go
  func (*Analyzer) detectCyclesDFS(pkg string, visited map[string]bool, recStack map[string]bool, path []string)
  ```
  detectCyclesDFS вспомогательная функция для поиска циклов

- **`computeComplexity`** (method)
  ```go
  func (*Analyzer) computeComplexity()
  ```
  computeComplexity вычисляет сложность каждого пакета

- **`identifyGodObjects`** (method)
  ```go
  func (*Analyzer) identifyGodObjects()
  ```
  identifyGodObjects выявляет пакеты со слишком большой сложностью

- **`identifyLayers`** (method)
  ```go
  func (*Analyzer) identifyLayers()
  ```
  identifyLayers выявляет архитектурные слои


---

### Package: `generator`

**Imports:**
- `fmt`
- `sort`
- `strings`
- `awdoc/pkg/analyzer`
- `awdoc/pkg/parser`

#### Exported Elements

**Functions:**

- **`NewMarkdownGenerator`** (function)
  ```go
  func NewMarkdownGenerator(sourceInfo *parser.SourceInfo, graph *analyzer.DependencyGraph) *MarkdownGenerator
  ```
  NewMarkdownGenerator создает генератор Markdown

**Methods:**

- **`GenerateProjectDoc`** (method)
  ```go
  func (*MarkdownGenerator) GenerateProjectDoc() string
  ```
  GenerateProjectDoc генерирует документацию всего проекта

**Structs:**

- **`MarkdownGenerator`** (struct)
  MarkdownGenerator генерирует документацию в формате Markdown

#### Internal Elements

**Functions:**

- **`elementTypeLabel`** (function)
  ```go
  func elementTypeLabel(elemType parser.ElementType) string
  ```
  elementTypeLabel возвращает читаемое имя типа элемента

**Methods:**

- **`generatePackageDoc`** (method)
  ```go
  func (*MarkdownGenerator) generatePackageDoc(pkg *parser.Package) string
  ```
  generatePackageDoc генерирует документацию для пакета

- **`generateElementsDoc`** (method)
  ```go
  func (*MarkdownGenerator) generateElementsDoc(elements []parser.CodeElement) string
  ```
  generateElementsDoc генерирует документацию для элементов

- **`generateArchitectureAnalysis`** (method)
  ```go
  func (*MarkdownGenerator) generateArchitectureAnalysis() string
  ```
  generateArchitectureAnalysis генерирует анализ архитектуры

- **`countElements`** (method)
  ```go
  func (*MarkdownGenerator) countElements() int
  ```
  countElements подсчитывает общее количество элементов


---

### Package: `parser`

**Imports:**
- `fmt`
- `io/ioutil`
- `os`
- `path/filepath`
- `strings`

#### Exported Elements

**Functions:**

- **`NewParser`** (function)
  ```go
  func NewParser(language string) (Parser, error)
  ```
  Factory функция для создания парсера по расширению файла

- **`NewDirScanner`** (function)
  ```go
  func NewDirScanner(language string) *DirScanner
  ```
  NewDirScanner создает сканер для директории

- **`ParseProject`** (function)
  ```go
  func ParseProject(rootDir string, language string) (*SourceInfo, error)
  ```
  ParseProject анализирует весь проект в директории

- **`ReadFileContent`** (function)
  ```go
  func ReadFileContent(filePath string) (string, error)
  ```
  ReadFileContent возвращает содержимое файла

- **`GetLines`** (function)
  ```go
  func GetLines(filePath string) ([]string, error)
  ```
  GetLines возвращает строки файла

**Methods:**

- **`ScanFiles`** (method)
  ```go
  func (*DirScanner) ScanFiles(rootDir string) ([]string, error)
  ```
  ScanFiles находит все файлы нужного языка в директории

**Structs:**

- **`DirScanner`** (struct)
  DirScanner сканирует директорию и находит файлы нужного языка

**Interfaces:**

- **`Parser`** (interface)
  Parser - основной интерфейс для парсилки кода


---

## Architecture Analysis

### Architectural Layers

**Layer 0:**
- generator
- parser
- analyzer

### Complex Packages (God Objects)

Packages with high complexity that might need refactoring:

- **analyzer** (Complexity: 12, Dependencies: 0)

### Dependency Graph

```
```

