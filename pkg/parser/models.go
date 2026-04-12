package parser

import "fmt"

// CodeElement представляет элемент кода (функция, тип, константа и т.д.)
type CodeElement struct {
	Name        string      // имя элемента
	Type        ElementType // тип элемента
	Exported    bool        // экспортируется ли
	Doc         string      // документация из комментариев
	Signature   string      // сигнатура функции
	Params      []Parameter // параметры
	Returns     []Parameter // возвращаемые значения
	SourceFile  string      // файл источника
	StartLine   int         // начальная строка
	EndLine     int         // конечная строка
	Package     string      // пакет
	RelatedDocs []string    // примеры использования
}

// Parameter представляет параметр функции
type Parameter struct {
	Name string // имя параметра
	Type string // тип параметра
}

// ElementType - тип элемента кода
type ElementType string

const (
	ElementFunc      ElementType = "function"
	ElementMethod    ElementType = "method"
	ElementType_     ElementType = "type"
	ElementConst     ElementType = "const"
	ElementVar       ElementType = "var"
	ElementStruct    ElementType = "struct"
	ElementInterface ElementType = "interface"
)

// Package представляет пакет кода
type Package struct {
	Name        string
	Path        string          // путь к пакету
	Doc         string          // документация пакета
	Elements    []CodeElement   // элементы в пакете
	Imports     map[string]bool // импорты {путь -> bool}
	ExportedAPI []CodeElement   // только экспортируемые элементы
}

// Symbol представляет символ для удобства анализа
type Symbol struct {
	FullName    string // например: "fmt.Printf"
	LocalName   string // например: "Printf"
	Package     string // например: "fmt"
	ElementType ElementType
}

func (s Symbol) String() string {
	return fmt.Sprintf("%s.%s", s.Package, s.LocalName)
}

// SourceInfo содержит информацию об источнике
type SourceInfo struct {
	Files    []string // файлы в анализируемом проекте
	RootDir  string   // корневая директория проекта
	Packages map[string]*Package
}
