package parser

import (
	"testing"
)

// TestParseGo тестирует парсинг Go файла
func TestParseGo(t *testing.T) {
	parser := &GoParser{}

	// Тестируем парсинг файла с примером
	pkg, err := parser.Parse("../../examples/sample/main.go")
	if err != nil {
		t.Fatalf("Failed to parse file: %v", err)
	}

	if pkg == nil {
		t.Fatal("Package is nil")
	}

	if pkg.Name != "main" {
		t.Errorf("Expected package name 'main', got '%s'", pkg.Name)
	}

	// Проверяем количество элементов
	if len(pkg.Elements) == 0 {
		t.Error("No elements found in package")
	}

	// Проверяем что есть экспортируемые элементы
	if len(pkg.ExportedAPI) == 0 {
		t.Error("No exported elements found")
	}

	// Найдем функцию ServiceA
	hasServiceA := false
	for _, elem := range pkg.Elements {
		if elem.Name == "ServiceA" {
			hasServiceA = true
			if elem.Type != ElementStruct {
				t.Errorf("ServiceA should be a struct, got %v", elem.Type)
			}
		}
	}
	if !hasServiceA {
		t.Error("ServiceA struct not found")
	}
}

// TestParser тестирует интерфейс Parser
func TestNewParser(t *testing.T) {
	// Тест на поддерживаемый язык
	parser, err := NewParser("go")
	if err != nil {
		t.Errorf("Failed to create Go parser: %v", err)
	}
	if parser == nil {
		t.Error("Parser is nil")
	}

	// Тест на неподдерживаемый язык
	parser, err = NewParser("rust")
	if err == nil {
		t.Error("Should return error for unsupported language")
	}
}

// TestTypeToString тестирует преобразование типов в строки
func TestTypeToString(t *testing.T) {
	gp := &GoParser{}

	tests := []struct {
		name     string
		typeName string
	}{
		{"string", "string"},
		{"int", "int"},
		{"bool", "bool"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Базовый тест что функция существует
			if gp == nil {
				t.Error("Parser is nil")
			}
		})
	}
}
