package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// GoParser парсит Go файлы и извлекает информацию
type GoParser struct{}

// Parse анализирует файл Go и возвращает информацию о пакете
func (gp *GoParser) Parse(filePath string) (*Package, error) {
	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parse error in %s: %w", filePath, err)
	}

	pkg := &Package{
		Name:        file.Name.Name,
		Path:        file.Name.Name,
		Doc:         formatDoc(file.Doc),
		Elements:    []CodeElement{},
		Imports:     make(map[string]bool),
		ExportedAPI: []CodeElement{},
	}

	// извлекаем импорты
	for _, decl := range file.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.IMPORT {
			for _, spec := range genDecl.Specs {
				if importSpec, ok := spec.(*ast.ImportSpec); ok {
					importPath := strings.Trim(importSpec.Path.Value, "\"")
					pkg.Imports[importPath] = true
				}
			}
		}
	}

	// извлекаем элементы кода
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			elem := gp.extractFunction(d, filePath, fset)
			if elem != nil {
				pkg.Elements = append(pkg.Elements, *elem)
				if elem.Exported {
					pkg.ExportedAPI = append(pkg.ExportedAPI, *elem)
				}
			}

		case *ast.GenDecl:
			elements := gp.extractGenDecl(d, filePath, fset)
			pkg.Elements = append(pkg.Elements, elements...)
			for _, elem := range elements {
				if elem.Exported {
					pkg.ExportedAPI = append(pkg.ExportedAPI, elem)
				}
			}
		}
	}

	return pkg, nil
}

// ParseDir парсит директорию (не реализовано для MVP)
func (gp *GoParser) ParseDir(dirPath string) (*SourceInfo, error) {
	return ParseProject(dirPath, "go")
}

// extractFunction извлекает информацию о функции/методе
func (gp *GoParser) extractFunction(fn *ast.FuncDecl, filePath string, fset *token.FileSet) *CodeElement {
	// определяем тип (функция или метод)
	elementType := ElementFunc
	if fn.Recv != nil {
		elementType = ElementMethod
	}

	// извлекаем параметры
	params := gp.extractParameters(fn.Type.Params)
	returns := gp.extractParameters(fn.Type.Results)

	elem := &CodeElement{
		Name:       fn.Name.Name,
		Type:       elementType,
		Exported:   ast.IsExported(fn.Name.Name),
		Doc:        formatDoc(fn.Doc),
		Signature:  gp.buildSignature(fn),
		Params:     params,
		Returns:    returns,
		SourceFile: filePath,
		StartLine:  fset.Position(fn.Pos()).Line,
		EndLine:    fset.Position(fn.End()).Line,
	}

	return elem
}

// extractGenDecl извлекает информацию из GenDecl (типы, константы, переменные)
func (gp *GoParser) extractGenDecl(decl *ast.GenDecl, filePath string, fset *token.FileSet) []CodeElement {
	var elements []CodeElement

	switch decl.Tok {
	case token.TYPE:
		for _, spec := range decl.Specs {
			if typeSpec, ok := spec.(*ast.TypeSpec); ok {
				elem := &CodeElement{
					Name:       typeSpec.Name.Name,
					Type:       gp.getTypeKind(typeSpec.Type),
					Exported:   ast.IsExported(typeSpec.Name.Name),
					Doc:        formatDoc(decl.Doc),
					SourceFile: filePath,
					StartLine:  fset.Position(typeSpec.Pos()).Line,
					EndLine:    fset.Position(typeSpec.End()).Line,
				}
				elements = append(elements, *elem)
			}
		}

	case token.CONST, token.VAR:
		elemType := ElementConst
		if decl.Tok == token.VAR {
			elemType = ElementVar
		}

		for _, spec := range decl.Specs {
			if valueSpec, ok := spec.(*ast.ValueSpec); ok {
				for _, name := range valueSpec.Names {
					elem := &CodeElement{
						Name:       name.Name,
						Type:       elemType,
						Exported:   ast.IsExported(name.Name),
						Doc:        formatDoc(decl.Doc),
						SourceFile: filePath,
						StartLine:  fset.Position(name.Pos()).Line,
						EndLine:    fset.Position(name.End()).Line,
					}
					elements = append(elements, *elem)
				}
			}
		}
	}

	return elements
}

// extractParameters извлекает параметры из FieldList
func (gp *GoParser) extractParameters(fl *ast.FieldList) []Parameter {
	var params []Parameter
	if fl == nil {
		return params
	}

	for _, field := range fl.List {
		typeName := gp.typeToString(field.Type)
		names := field.Names
		if len(names) == 0 {
			// безымянные параметры
			params = append(params, Parameter{
				Name: "",
				Type: typeName,
			})
		} else {
			for _, name := range names {
				params = append(params, Parameter{
					Name: name.Name,
					Type: typeName,
				})
			}
		}
	}

	return params
}

// typeToString преобразует AST тип в строку
func (gp *GoParser) typeToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return gp.typeToString(t.X) + "." + t.Sel.Name
	case *ast.StarExpr:
		return "*" + gp.typeToString(t.X)
	case *ast.ArrayType:
		return "[]" + gp.typeToString(t.Elt)
	case *ast.MapType:
		return "map[" + gp.typeToString(t.Key) + "]" + gp.typeToString(t.Value)
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.FuncType:
		return "func(...)"
	default:
		return "unknown"
	}
}

// buildSignature строит сигнатуру функции
func (gp *GoParser) buildSignature(fn *ast.FuncDecl) string {
	sig := "func "

	if fn.Recv != nil {
		sig += "("
		for _, field := range fn.Recv.List {
			sig += gp.typeToString(field.Type)
		}
		sig += ") "
	}

	sig += fn.Name.Name + "("

	if fn.Type.Params != nil {
		for i, field := range fn.Type.Params.List {
			if i > 0 {
				sig += ", "
			}
			if len(field.Names) > 0 {
				sig += field.Names[0].Name + " "
			}
			sig += gp.typeToString(field.Type)
		}
	}

	sig += ")"

	if fn.Type.Results != nil && len(fn.Type.Results.List) > 0 {
		sig += " "
		if len(fn.Type.Results.List) == 1 {
			sig += gp.typeToString(fn.Type.Results.List[0].Type)
		} else {
			sig += "("
			for i, field := range fn.Type.Results.List {
				if i > 0 {
					sig += ", "
				}
				sig += gp.typeToString(field.Type)
			}
			sig += ")"
		}
	}

	return sig
}

// getTypeKind определяет тип typedef
func (gp *GoParser) getTypeKind(expr ast.Expr) ElementType {
	switch expr.(type) {
	case *ast.StructType:
		return ElementStruct
	case *ast.InterfaceType:
		return ElementInterface
	default:
		return ElementType_
	}
}

// formatDoc форматирует документацию из комментариев
func formatDoc(doc *ast.CommentGroup) string {
	if doc == nil {
		return ""
	}
	return strings.TrimSpace(doc.Text())
}

// ExtractCallReferences находит вызовы функций (для анализа зависимостей)
// На MVP это упрощено - нужна инструмента дополнит анализа
func (gp *GoParser) ExtractCallReferences(filePath string) ([]string, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var callRefs []string
	ast.Inspect(file, func(node ast.Node) bool {
		if call, ok := node.(*ast.CallExpr); ok {
			if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
				if ident, ok := sel.X.(*ast.Ident); ok {
					callRefs = append(callRefs, ident.Name+"."+sel.Sel.Name)
				}
			}
		}
		return true
	})

	return callRefs, nil
}
