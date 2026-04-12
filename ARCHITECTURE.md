# AWDoc - Architecture & Design

## System Architecture

```mermaid
graph TD
    A["🏗️ AWDoc - API Doc Generator"] --> B["Command Line Interface<br/>(main.go)<br/>Парсит аргументы,<br/>управляет потоком выполнения"]
    
    B --> C["Parser Module<br/>(pkg/parser/)"]
    B --> D["Analyzer Module<br/>(pkg/analyzer/)"]
    
    C --> E["GoParser<br/>Implementation<br/>go/ast, go/parser"]
    E --> F["SourceInfo<br/>Packages<br/>Elements<br/>Imports"]
    
    D --> G["DependencyGraphs<br/>Complexity Metrics<br/>Cycle Detection<br/>Layer Analysis"]
    G --> H["DependencyGraph<br/>Cycles<br/>Layers<br/>GodObjects"]
    
    F --> I["Generator Module<br/>(pkg/generator/)"]
    H --> I
    
    I --> J["Markdown Generator<br/>(implemented)"]
    I --> K["HTML Generator<br/>(TODO)"]
    
    J --> L["📄 Markdown Document<br/>(.md file)"]
    K --> M["🌐 HTML Document<br/>(with diagrams)"]
    
    style A fill:#4a90e2,stroke:#333,stroke-width:2px,color:#fff
    style B fill:#7cb342,stroke:#333,stroke-width:2px,color:#fff
    style C fill:#1976d2,stroke:#333,stroke-width:2px,color:#fff
    style D fill:#1976d2,stroke:#333,stroke-width:2px,color:#fff
    style E fill:#0097a7,stroke:#333,color:#fff
    style G fill:#0097a7,stroke:#333,color:#fff
    style I fill:#f57c00,stroke:#333,stroke-width:2px,color:#fff
    style J fill:#6a1b9a,stroke:#333,stroke-width:1px,color:#fff
    style K fill:#6a1b9a,stroke:#333,stroke-width:1px,color:#fff
    style L fill:#388e3c,stroke:#333,stroke-width:1px,color:#fff
    style M fill:#388e3c,stroke:#333,stroke-width:1px,color:#fff
```

## Component Breakdown

### 1. Parser Module (`pkg/parser/`)

**Responsibility:** Extract code structure from source files

**Classes:**

- `GoParser` - Parse Go source files using go/ast
- `Parser` (interface) - abstraction for different languages

**Key Exports:**

- `ParseProject(dir, lang)` - main entry point
- `CodeElement` - represents functions, types, etc.
- `Package` - represents a code package
- `SourceInfo` - container for all analysis data

**Dependencies:**

- Go standard library: `go/parser`, `go/ast`, `go/token`
- No external dependencies

**Complexity:** O(n) where n = number of files

### 2. Analyzer Module (`pkg/analyzer/`)

**Responsibility:** Analyze package relationships and complexity

**Classes:**

- `Analyzer` - main analysis engine
- `PackageNode` - node in dependency graph
- `DependencyGraph` - mathematical graph structure

**Key Exports:**

- `Analyze()` - build graph and metrics
- `GetDependencyInfo()` - query single package

**Algorithms:**

- **DFS Cycle Detection** - O(V+E) where V=packages, E=dependencies
- **Layer Assignment** - topological sort
- **Complexity Calculation** - weighted scoring

**Dependencies:**

- `pkg/parser` for data structures

**Complexity:** O(V²) worst case (with cycles)

### 3. Generator Module (`pkg/generator/`)

**Responsibility:** Create human-readable documentation

**Classes:**

- `MarkdownGenerator` - generate Markdown output
- `DocumentationBuilder` - facade for all generators

**Key Exports:**

- `GenerateProjectDoc()` - main generation method
- `DocumentationBuilder` - builder pattern

**Output Formats:**

1. Markdown - text-based, version control friendly
2. HTML (planned) - interactive, with diagrams

**Dependencies:**

- `pkg/parser` and `pkg/analyzer` for data
- Go standard library

**Complexity:** O(n) where n = number of elements

## Data Structures

### SourceInfo

```go
type SourceInfo struct {
    Files       []string              // files found
    RootDir     string                // analysis root
    Packages    map[string]*Package   // all packages
}
```

### Package

```go
type Package struct {
    Name        string              // package name
    Path        string              // import path
    Doc         string              // package docs
    Elements    []CodeElement       // all elements
    Imports     map[string]bool     // imported packages
    ExportedAPI []CodeElement       // only exported
}
```

### CodeElement

```go
type CodeElement struct {
    Name       string        // identifier name
    Type       ElementType   // function, type, method, etc.
    Exported   bool          // is it exported?
    Doc        string        // documentation
    Signature  string        // function signature
    Params     []Parameter   // function parameters
    Returns    []Parameter   // return values
    SourceFile string        // source file path
    StartLine  int           // line number
    EndLine    int           // line number
}
```

### DependencyGraph

```go
type DependencyGraph struct {
    Nodes      map[string]*PackageNode  // nodes
    Edges      map[string][]string      // adjacency list
    Cycles     [][]string               // circular deps
    Layers     [][]string               // topological layers
    GodObjects []string                 // complex packages
}
```

## Processing Pipeline

```mermaid
graph TD
    A["📂 Input Source Code"] --> B["Directory Scanner<br/>Scan for Go files<br/>Exclude vendor, .git"]
    B --> C["GoParser<br/>Parse each file<br/>Build AST"]
    C --> D["AST Analysis<br/>Extract elements<br/>Classify as exported/internal"]
    D --> E["Package Build<br/>Group elements by package<br/>Collect imports"]
    E --> F["Analyzer<br/>Build dependency graph<br/>Detect cycles DFS"]
    F --> G["Complexity Calculation<br/>Weighted scoring<br/>Identify god objects"]
    G --> H["Layer Analysis<br/>Topological sort<br/>Assign depth levels"]
    H --> I["Generator<br/>Format output<br/>Create documentation"]
    I --> J["📊 Output: Markdown/HTML"]
    
    style A fill:#1976d2,stroke:#333,stroke-width:2px,color:#fff
    style B fill:#0097a7,stroke:#333,stroke-width:1px,color:#fff
    style C fill:#0097a7,stroke:#333,stroke-width:1px,color:#fff
    style D fill:#0097a7,stroke:#333,stroke-width:1px,color:#fff
    style E fill:#0097a7,stroke:#333,stroke-width:1px,color:#fff
    style F fill:#f57c00,stroke:#333,stroke-width:1px,color:#fff
    style G fill:#f57c00,stroke:#333,stroke-width:1px,color:#fff
    style H fill:#f57c00,stroke:#333,stroke-width:1px,color:#fff
    style I fill:#6a1b9a,stroke:#333,stroke-width:1px,color:#fff
    style J fill:#388e3c,stroke:#333,stroke-width:2px,color:#fff
```

## Algorithm Details

### Cycle Detection (DFS)

```mermaid
graph TD
    A["detectCycles(graph)"] --> B["For each<br/>unvisited node"]
    B --> C["dfs(node,<br/>visited,<br/>recStack,<br/>path)"]
    C --> D["Mark node<br/>as visited"]
    D --> E["Add to<br/>recursive stack"]
    E --> F["Add to path"]
    F --> G{For each<br/>neighbor}
    G -->|not visited| H["dfs(neighbor)"]
    H --> G
    G -->|visited| I{In<br/>recStack?}
    I -->|Yes| J["🔴 Found Cycle<br/>Report cycle"]
    I -->|No| K["Continue"]
    J --> L["Remove from<br/>recursive stack"]
    K --> L
    L --> G
    G -->|done| M["✅ Cycle Detection<br/>Complete"]
    
    style A fill:#f57c00,stroke:#333,stroke-width:2px,color:#fff
    style B fill:#0097a7,stroke:#333,stroke-width:1px,color:#fff
    style C fill:#0097a7,stroke:#333,stroke-width:1px,color:#fff
    style G fill:#1976d2,stroke:#333,stroke-width:1px,color:#fff
    style I fill:#1976d2,stroke:#333,stroke-width:1px,color:#fff
    style J fill:#d32f2f,stroke:#333,stroke-width:2px,color:#fff
    style M fill:#388e3c,stroke:#333,stroke-width:2px,color:#fff
```

**Time Complexity:** O(V + E)
**Space Complexity:** O(V)

### Layer Assignment (Topological Sort)

```mermaid
graph TD
    A["assignLayers(graph)"] --> B["layer = 0"]
    B --> C["current = packages<br/>with no dependencies"]
    C --> D{current is<br/>not empty?}
    D -->|Yes| E["Assign all in current<br/>to layer"]
    E --> F["Find packages whose<br/>dependencies are assigned"]
    F --> G["current = next packages"]
    G --> H["layer++"]
    H --> D
    D -->|No| I["✅ All layers<br/>assigned"]
    
    style A fill:#f57c00,stroke:#333,stroke-width:2px,color:#fff
    style D fill:#1976d2,stroke:#333,stroke-width:1px,color:#fff
    style I fill:#388e3c,stroke:#333,stroke-width:2px,color:#fff
```

**Time Complexity:** O(V²) worst case
**Space Complexity:** O(V)

### Complexity Scoring

```mermaid
graph LR
    A["Elements × 1"] --> S["Complexity<br/>Score"]
    B["Dependencies × 3"] --> S
    C["Dependents × 2"] --> S
    D["CircularDeps × 5"] --> S
    
    style A fill:#4CAF50,stroke:#333,stroke-width:1px,color:#fff
    style B fill:#FF9800,stroke:#333,stroke-width:1px,color:#fff
    style C fill:#2196F3,stroke:#333,stroke-width:1px,color:#fff
    style D fill:#F44336,stroke:#333,stroke-width:2px,color:#fff
    style S fill:#9C27B0,stroke:#333,stroke-width:2px,color:#fff
```

**Factors:**

1. **Elements (weight 1):** More code = more complexity
2. **Dependencies (weight 3):** External dependencies increase risk
3. **Dependents (weight 2):** More impact if broken
4. **Circular (weight 5):** Major architecture issue

## Design Patterns Used

### 1. Parser Pattern

- **Interface-based:** `Parser` interface allows multiple implementations
- **Factory Method:** `NewParser()` creates language-specific parsers
- **Strategy Pattern:** Each `Parser` implements different parsing strategy

### 2. Builder Pattern

- `DocumentationBuilder` constructs complex documents
- Separates construction from representation
- Allows different output formats

### 3. Adapter Pattern

- `GoParser` adapts Go's `ast` package to our `CodeElement` model
- Hides language-specific details

### 4. Facade Pattern

- `DocumentationBuilder` provides simplified interface
- Hides complexity of multiple generators

## Extension Points

### Adding New Language Support

```mermaid
graph LR
    A["1. Create new Parser"] --> B["Implement Parser interface<br/>Parse to CodeElement[]<br/>Handle language-specific comments"]
    C["2. Register in Factory"] --> D["Add case in NewParser<br/>Update DirScanner"]
    E["3. Testing"] --> F["Add language_parser_test.go<br/>Test element extraction"]
    
    B --> C
    D --> E
    F --> G["✅ New Language Support<br/>Ready"]
    
    style A fill:#1976d2,stroke:#333,stroke-width:1px,color:#fff
    style C fill:#1976d2,stroke:#333,stroke-width:1px,color:#fff
    style E fill:#1976d2,stroke:#333,stroke-width:1px,color:#fff
    style G fill:#388e3c,stroke:#333,stroke-width:2px,color:#fff
```

### Adding New Output Format

```mermaid
graph LR
    A["1. Create new Generator"] --> B["Implement Generator interface<br/>Generate format-specific output"]
    C["2. Register in Builder"] --> D["Add method to DocumentationBuilder<br/>Update main.go"]
    E["3. Testing"] --> F["Add generator_test.go<br/>Validate output structure"]
    
    B --> C
    D --> E
    F --> G["✅ New Output Format<br/>Ready"]
    
    style A fill:#f57c00,stroke:#333,stroke-width:1px,color:#fff
    style C fill:#f57c00,stroke:#333,stroke-width:1px,color:#fff
    style E fill:#f57c00,stroke:#333,stroke-width:1px,color:#fff
    style G fill:#388e3c,stroke:#333,stroke-width:2px,color:#fff
```

## Performance Characteristics

| Operation | Complexity | Notes |
| --- | --- | --- |
| Parse file | O(n) | where n = file size |
| Parse project | O(m) | where m = total lines of code |
| Detect cycles | O(V+E) | V = packages, E = dependencies |
| Assign layers | O(V²) | worst case with many deps |
| Generate docs | O(n) | where n = number of elements |
| **Total** | **O(V² + m)** | dominated by analysis |

**Typical times:**

- Small project (10 packages, 1K elements): < 100ms
- Medium project (50 packages, 5K elements): 200-500ms
- Large project (100+ packages, 10K+ elements): 1-3 seconds

## Dependencies

### External

- None for parser and analyzer
- Standard Go library only

### Internal

- `parser` used by `analyzer` and `generator`
- `analyzer` used by `generator`
- No circular dependencies

## Error Handling

### Parser Errors

- File not found → log and continue
- Invalid syntax → log with line number
- Unsupported construct → skip element

### Analyzer Errors

- Invalid graph → return error
- Infinite loops in cycles → timeout protection

### Generator Errors

- File write errors → return error
- Template errors → return error

## Concurrency

Current implementation is **single-threaded**:

- Parser processes files sequentially
- Analyzer runs single-threaded
- Generator creates output sequentially

**Future optimization:**

- Parallel file parsing with sync.WaitGroup
- Concurrent analysis of independent packages
- Parallelized document generation
