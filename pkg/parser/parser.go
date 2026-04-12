package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Parser - основной интерфейс для парсилки кода
type Parser interface {
	Parse(filePath string) (*Package, error)
	ParseDir(dirPath string) (*SourceInfo, error)
}

// Factory функция для создания парсера по расширению файла
func NewParser(language string) (Parser, error) {
	switch strings.ToLower(language) {
	case "go":
		return &GoParser{}, nil
	default:
		return nil, fmt.Errorf("unsupported language: %s", language)
	}
}

// DirScanner сканирует директорию и находит файлы нужного языка
type DirScanner struct {
	language  string
	extension string
}

// NewDirScanner создает сканер для директории
func NewDirScanner(language string) *DirScanner {
	switch strings.ToLower(language) {
	case "go":
		return &DirScanner{language: "go", extension: ".go"}
	default:
		return &DirScanner{language: language, extension: "." + language}
	}
}

// ScanFiles находит все файлы нужного языка в директории
func (ds *DirScanner) ScanFiles(rootDir string) ([]string, error) {
	var files []string

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// пропускаем директории и файлы git
		if info.IsDir() {
			if info.Name() == ".git" || info.Name() == "vendor" || info.Name() == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}

		// проверяем расширение
		if filepath.Ext(path) == ds.extension {
			// для Go пропускаем тестовые файлы при первоначальном сканировании (но потом их обработаем)
			if ds.language == "go" && !strings.HasSuffix(path, "_test.go") {
				files = append(files, path)
			}
		}

		return nil
	})

	return files, err
}

// ParseProject анализирует весь проект в директории
func ParseProject(rootDir string, language string) (*SourceInfo, error) {
	parser, err := NewParser(language)
	if err != nil {
		return nil, err
	}

	scanner := NewDirScanner(language)
	files, err := scanner.ScanFiles(rootDir)
	if err != nil {
		return nil, err
	}

	sourceInfo := &SourceInfo{
		Files:    files,
		RootDir:  rootDir,
		Packages: make(map[string]*Package),
	}

	for _, file := range files {
		pkg, err := parser.Parse(file)
		if err != nil {
			// логируем ошибку но продолжаем обработку
			fmt.Printf("Warning: failed to parse %s: %v\n", file, err)
			continue
		}

		if pkg != nil {
			sourceInfo.Packages[pkg.Path] = pkg
		}
	}

	return sourceInfo, nil
}

// ReadFileContent возвращает содержимое файла
func ReadFileContent(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// GetLines возвращает строки файла
func GetLines(filePath string) ([]string, error) {
	content, err := ReadFileContent(filePath)
	if err != nil {
		return nil, err
	}
	return strings.Split(content, "\n"), nil
}
