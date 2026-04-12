# Demo PowerShell script - Demonstrate AWDoc capabilities

Clear-Host
Write-Host "===== AWDoc - Demonstration ====="
Write-Host "Version: 1.0.0"
Write-Host "Purpose: Analyze Go code and generate API documentation"
Write-Host ""

# Check if awdoc is built
if (-not (Test-Path ".\awdoc.exe")) {
    Write-Host "Building AWDoc..."
    go build -o awdoc.exe main.go
    Write-Host "AWDoc built"
} else {
    Write-Host "AWDoc already built"
}

# Demo 1: Simple analysis
Write-Host ""
Write-Host "1. Analyzing simple example (examples/sample/)"
.\awdoc.exe -source ./examples/sample -lang go -output-dir output
Write-Host "Generated output/analysis.md"

# Demo 2: Complex analysis
Write-Host ""
Write-Host "2. Analyzing multi-package project (examples/complex/)"
.\awdoc.exe -source ./examples/complex -lang go -output output/complex-analysis.md
Write-Host "Generated output/complex-analysis.md"

# Demo 3: Architecture analysis
Write-Host ""
Write-Host "3. Analyzing entire AWDoc project"
.\awdoc.exe -source ./pkg -lang go -output output/internal-analysis.md
Write-Host "Generated output/internal-analysis.md"

# Show statistics
Write-Host ""
Write-Host "4. Generated Documentation Files:"
Get-Item output\analysis.md, output\complex-analysis.md, output\internal-analysis.md -ErrorAction SilentlyContinue | Format-Table Name, Length

# Run tests
Write-Host ""
Write-Host "5. Running unit tests"
go test -v ./pkg/analyzer -run TestAnalyzer | Select-Object -First 30

Write-Host ""
Write-Host "6. Key Features Demonstrated:"
Write-Host "  * Code parsing and AST analysis"
Write-Host "  * API documentation extraction"
Write-Host "  * Dependency graph construction"
Write-Host "  * Package complexity analysis"
Write-Host "  * Markdown document generation"
Write-Host "  * Architecture layer detection"

Write-Host ""
Write-Host "===== Demo Complete! ====="
