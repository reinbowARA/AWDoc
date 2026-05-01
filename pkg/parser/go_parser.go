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

	// обнаруживаем API запросы
	pkg.APIRequests = gp.DetectAPIRequests(file, fset, filePath, file.Name.Name)

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
				typeKind := gp.getTypeKind(typeSpec.Type)
				elem := &CodeElement{
					Name:       typeSpec.Name.Name,
					Type:       typeKind,
					Exported:   ast.IsExported(typeSpec.Name.Name),
					Doc:        formatDoc(decl.Doc),
					SourceFile: filePath,
					StartLine:  fset.Position(typeSpec.Pos()).Line,
					EndLine:    fset.Position(typeSpec.End()).Line,
				}

				// Extract fields for struct types
				if typeKind == ElementStruct {
					if structType, ok := typeSpec.Type.(*ast.StructType); ok {
						elem.Fields = gp.extractStructFields(structType)
					}
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

// extractStructFields извлекает поля из struct type
func (gp *GoParser) extractStructFields(st *ast.StructType) []StructField {
	var fields []StructField
	if st == nil || st.Fields == nil {
		return fields
	}

	for _, field := range st.Fields.List {
		// Пропускаем embedded fields (без имен) в базовом выводе
		if len(field.Names) == 0 {
			continue
		}

		typeName := gp.typeToString(field.Type)
		for _, name := range field.Names {
			fields = append(fields, StructField{
				Name:     name.Name,
				Type:     typeName,
				Doc:      formatDoc(field.Doc),
				Exported: ast.IsExported(name.Name),
			})
		}
	}

	return fields
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

// AnalyzeTestCoverage анализирует наличие тестов для элементов в пакете
func (gp *GoParser) AnalyzeTestCoverage(pkg *Package, testFunctions map[string]bool) {
	// testFunctions - это map из имён тестов (например: "TestFunction", "TestMethod", "TestGodObjects")

	testedCount := 0
	testableElementsCount := 0

	// Проходим по экспортируемым элементам
	for i, elem := range pkg.ExportedAPI {
		// Пропускаем структуры, интерфейсы и типы - они не требуют тестов
		if elem.Type == ElementStruct || elem.Type == ElementInterface || elem.Type == ElementType_ {
			continue
		}

		testableElementsCount++
		testName := gp.getExpectedTestName(elem.Name, elem.Type)

		// Проверяем есть ли соответствующий тест
		hasTest := false
		for testFunc := range testFunctions {
			if strings.EqualFold(testFunc, testName) ||
				strings.HasPrefix(strings.ToLower(testFunc), strings.ToLower("Test"+elem.Name)) {
				hasTest = true
				pkg.ExportedAPI[i].HasTests = true
				pkg.ExportedAPI[i].TestName = testFunc
				testedCount++
				break
			}
		}

		// Обновляем в Elements также
		for j, e := range pkg.Elements {
			if e.Name == elem.Name && e.Type == elem.Type {
				pkg.Elements[j].HasTests = hasTest
				if hasTest {
					pkg.Elements[j].TestName = pkg.ExportedAPI[i].TestName
				}
				break
			}
		}
	}

	// Вычисляем процент покрытия только для функций и методов
	pkg.TotalElements = testableElementsCount
	pkg.TestedElements = testedCount

	if pkg.TotalElements > 0 {
		pkg.Coverage = float64(testedCount) / float64(pkg.TotalElements) * 100
	} else {
		pkg.Coverage = 0
	}
}

// getExpectedTestName генерирует ожидаемое имя функции теста
func (gp *GoParser) getExpectedTestName(elementName string, elemType ElementType) string {
	switch elemType {
	case ElementFunc:
		return "Test" + elementName
	case ElementMethod:
		// для методов ищем Test<ReceiverType><MethodName>
		return "Test" + elementName
	default:
		return "Test" + elementName
	}
}

// DetectAPIRequests обнаруживает API запросы в файле
func (gp *GoParser) DetectAPIRequests(file *ast.File, fset *token.FileSet, filePath string, pkgName string) []APIRequest {
	var apiRequests []APIRequest

	// Ищем вызовы для различных фреймворков
	ast.Inspect(file, func(node ast.Node) bool {
		if call, ok := node.(*ast.CallExpr); ok {
			if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
				methodName := sel.Sel.Name
				httpMethods := map[string]string{
					"GET":     "GET",
					"POST":    "POST",
					"PUT":     "PUT",
					"DELETE":  "DELETE",
					"PATCH":   "PATCH",
					"HEAD":    "HEAD",
					"OPTIONS": "OPTIONS",
				}

				// gin gin.Context
				if ident, ok := sel.X.(*ast.Ident); ok {
					if ident.Name == "gin" || strings.Contains(ident.Name, "router") {
						if method, ok := httpMethods[methodName]; ok {
							if len(call.Args) >= 2 {
								path := gp.extractPathFromArg(call.Args[0])

								// Получаем документацию из комментариев
								docComment := gp.extractDocCommentSimple(file, fset, call.Pos())

								apiRequests = append(apiRequests, APIRequest{
									Name:        fmt.Sprintf("%s %s", method, path),
									Path:        path,
									Method:      method,
									Description: docComment,
									IsSwaggered: false,
									SourceFile:  filePath,
								})
							}
						}
					}
				}

				// gin.Context / *gin.Context вызовы
				if recv, ok := sel.X.(*ast.StarExpr); ok {
					if ident, ok := recv.X.(*ast.Ident); ok {
						if ident.Name == "c" || ident.Name == "ctx" || ident.Name == "ginCtx" {
							if method, ok := httpMethods[methodName]; ok {
								if len(call.Args) >= 2 {
									path := gp.extractPathFromArg(call.Args[0])
									docComment := gp.extractDocCommentSimple(file, fset, call.Pos())

									apiRequests = append(apiRequests, APIRequest{
										Name:        fmt.Sprintf("%s %s", method, path),
										Path:        path,
										Method:      method,
										Description: docComment,
										IsSwaggered: false,
										SourceFile:  filePath,
									})
								}
							}
						}
					}
				}

				// echo框架: router.GET, router.POST и т.д.
				if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == "e" {
					if method, ok := httpMethods[methodName]; ok {
						if len(call.Args) >= 2 {
							path := gp.extractPathFromArg(call.Args[0])
							docComment := gp.extractDocCommentSimple(file, fset, call.Pos())

							apiRequests = append(apiRequests, APIRequest{
								Name:        fmt.Sprintf("%s %s", method, path),
								Path:        path,
								Method:      method,
								Description: docComment,
								IsSwaggered: false,
								SourceFile:  filePath,
							})
						}
					}
				}
			}
		}

		// net/http: http.HandleFunc
		if call, ok := node.(*ast.CallExpr); ok {
			if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
				if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == "http" {
					if sel.Sel.Name == "HandleFunc" || sel.Sel.Name == "Handle" || sel.Sel.Name == "Get" || sel.Sel.Name == "Post" {
						if len(call.Args) >= 2 {
							path := gp.extractPathFromArg(call.Args[0])
							docComment := gp.extractDocCommentSimple(file, fset, call.Pos())

							apiRequests = append(apiRequests, APIRequest{
								Name:        fmt.Sprintf("GET %s", path),
								Path:        path,
								Method:      "GET",
								Description: docComment,
								IsSwaggered: false,
								SourceFile:  filePath,
							})
						}
					}
				}
			}
		}

		return true
	})

	// Проверяем наличие Swagger аннотаций
	for i := range apiRequests {
		apiRequests[i].IsSwaggered = gp.hasSwaggerAnnotationsSimple(file)
	}

	return apiRequests
}

// extractPathFromArg извлекает путь из аргумента вызова
func (gp *GoParser) extractPathFromArg(arg ast.Expr) string {
	if lit, ok := arg.(*ast.BasicLit); ok {
		return strings.Trim(lit.Value, "\"")
	}
	return ""
}

// extractDocCommentSimple извлекает документацию из комментариев (простая версия)
func (gp *GoParser) extractDocCommentSimple(file *ast.File, fset *token.FileSet, pos token.Pos) string {
	// Получаем номер строки для данной позиции
	line := fset.Position(pos).Line

	// Ищем комментарии, которые находятся непосредственно перед данной позицией
	for _, group := range file.Comments {
		for _, comment := range group.List {
			commentEndLine := fset.Position(comment.End()).Line
			if commentEndLine < line && line-commentEndLine <= 3 {
				return strings.TrimSpace(comment.Text)
			}
		}
	}

	return ""
}

// hasSwaggerAnnotationsSimple проверяет наличие Swagger аннотаций в файле (без fset)
func (gp *GoParser) hasSwaggerAnnotationsSimple(file *ast.File) bool {
	swaggerPatterns := []string{
		"@summary",
		"@description",
		"@tags",
		"@id",
		"@accept",
		"@produce",
		"@param",
		"@success",
		"@failure",
		"@router",
	}

	// Проверяем комментарии в файле
	for _, group := range file.Comments {
		text := group.Text()
		for _, pattern := range swaggerPatterns {
			if strings.Contains(text, pattern) {
				return true
			}
		}
	}

	return false
}
