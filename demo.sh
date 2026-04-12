#!/bin/bash
# demo.sh - Demonstrate AWDoc capabilities

set -e

PROJECT_NAME="AWDoc"
VERSION="1.0.0"

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Functions
print_header() {
    echo -e "${BLUE}═════════════════════════════════════${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}═════════════════════════════════════${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_section() {
    echo -e "\n${YELLOW}► $1${NC}"
}

# Main demo
main() {
    clear
    print_header "$PROJECT_NAME - Demonstration"
    
    echo "Version: $VERSION"
    echo "Purpose: Analyze Go source code and generate API documentation"
    echo ""
    
    # Check if awdoc is built
    if [ ! -f "./awdoc" ]; then
        print_section "Building AWDoc..."
        go build -o awdoc main.go
        print_success "AWDoc built"
    else
        print_success "AWDoc already built"
    fi
    
    # Demo 1: Simple analysis
    print_section "1. Analyzing simple example (examples/sample/)"
    mkdir -p output
    ./awdoc -source ./examples/sample -lang go -output output/analysis.md
    print_success "Generated output/analysis.md"
    echo "  Elements found: $(grep -c "^\- \*\*" output/analysis.md || echo "multiple")"
    
    # Demo 2: Complex analysis
    print_section "2. Analyzing multi-package project (examples/complex/)"
    ./awdoc -source ./examples/complex -lang go -output output/complex-analysis.md
    print_success "Generated output/complex-analysis.md"
    echo "  Packages: graph"
    
    # Demo 3: Architecture analysis
    print_section "3. Analyzing entire AWDoc project"
    ./awdoc -source ./pkg -lang go -output output/internal-analysis.md
    print_success "Generated output/internal-analysis.md"
    
    # Show statistics
    print_section "4. Statistics"
    echo ""
    echo "Generated files:"
    echo "  - output/analysis.md: $(wc -l < output/analysis.md) lines"
    echo "  - output/complex-analysis.md: $(wc -l < output/complex-analysis.md) lines"
    echo "  - output/internal-analysis.md: $(wc -l < output/internal-analysis.md) lines"
    
    # Test examples
    print_section "5. Running unit tests"
    go test -v ./pkg/... -run TestAnalyzer
    print_success "Tests passed"
    
    # Show features
    print_section "6. Key Features Demonstrated"
    echo ""
    echo "  ✓ Code parsing and AST analysis"
    echo "  ✓ API documentation extraction"
    echo "  ✓ Dependency graph construction"
    echo "  ✓ Package complexity analysis"
    echo "  ✓ Markdown document generation"
    echo "  ✓ Architecture layer detection"
    echo ""
    
    print_section "7. Next Steps"
    echo ""
    echo "Try these commands:"
    echo "  - View generated docs: cat output/analysis.md"
    echo "  - Analyze your project: ./awdoc -source /path/to/project"
    echo "  - Custom output dir: ./awdoc -source . -output-dir ./docs"
    echo "  - Run tests: go test -v ./..."
    echo "  - Read docs: cat README.md"
    echo ""
    
    print_header "Demo Complete!"
}

# Run if executed directly
if [ "${BASH_SOURCE[0]}" == "${0}" ]; then
    main "$@"
fi
