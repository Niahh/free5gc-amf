//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/packages"
	"log"
	"os"
	"reflect"
	"strings"
	"text/template"
)

// EnumMapping holds one enum constant and its generated error message.
type EnumMapping struct {
	ConstName string // e.g., "ngapType.CauseTransportPresentTransportResourceUnavailable"
	ErrorMsg  string // e.g., "transport resource unavailable"
}

// This is a representation of the information extracted from the Cause Type to be given to the template for generating
// the final file.
type CauseField struct {
	FieldName    string // e.g., "Transport" or "Nas"
	PointerType  string // e.g., "CauseTransport"
	EnumMappings []EnumMapping
	// The present constant follows a naming convention: "ngapType.CausePresent" + FieldName
	PresentConst string // e.g., "ngapType.CausePresentTransport"
}

// CauseData is used to pass data to the text template.
type CauseData struct {
	PackageImportPath string       // e.g., "github.com/free5gc/ngap"
	PackageName       string       // e.g., "ngapType"
	StructName        string       // e.g., "Cause"
	Fields            []CauseField // Pointer information for the different cause subtypes
}

type FileInfo struct {
	file     *ast.File
	typeSpec *ast.TypeSpec
}

// Struct to store the general information that the generator will need during the whole process.
type generatorCtx struct {
	pkgImportPath string
	packageName   string
	pkgFullName   string
	structToFind  string
	pkg           *packages.Package
}

// getPackage tries to find the package using its fullName, if not logs an error and exit.
func getPackage(pkgFullName string) *packages.Package {
	cfg := &packages.Config{
		// NeedSyntax -> Parse and load the ASTs  of each files
		Mode: packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo,
	}

	pkgs, err := packages.Load(cfg, pkgFullName)

	if err != nil {
		log.Fatalf("failed to load package %q: %v", pkgFullName, err)
	}

	if packages.PrintErrors(pkgs) > 0 {
		log.Fatal("Errors found when loading package")
	}

	// As we are giving the full path and not an ambiguous path, we can return directly the pkg found.
	return pkgs[0]
}

// This function is useful to retrieve the file for the subtypes of the Cause Type.
func getFileFromTypeName(pkg *packages.Package, typeName string) *ast.File {
	files := pkg.Syntax

	var fileFound *ast.File

	for _, file := range files {
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)

			if !ok {
				continue
			}

			for _, spec := range genDecl.Specs {
				// Same as for the genDecl
				if ts, typeOk := spec.(*ast.TypeSpec); typeOk {
					if typeName == ts.Name.Name {
						return file
					}
				} else {
					continue
				}
			}
		}
	}
	if fileFound == nil {
		log.Fatalf("no file found with %s type", typeName)
	}
	return nil
}

func getStructFromFile(file *ast.File, typeToFind string) *ast.StructType {
	var causeType *ast.TypeSpec

	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		// The type assertion failed so we cannot use the generic declaration (will result in nil ptr)
		if !ok {
			continue
		}

		for _, spec := range genDecl.Specs {
			// Same as for the genDecl
			if ts, typeOk := spec.(*ast.TypeSpec); typeOk {
				if ts.Name.Name == typeToFind {
					causeType = ts
					break
				}
			} else {
				continue
			}
		}
	}
	if causeType != nil {
		structType, ok := causeType.Type.(*ast.StructType)

		if !ok {
			log.Fatalf("Could not cast the type %s to a struct object", typeToFind)
		}
		return structType
	}

	log.Fatalf("no type found with %s type", typeToFind)
	return nil
}

func getConstListFromFile(f *ast.File) []string {
	var constNames []string

	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.CONST {
			continue
		}

		// Loop over all specifications in the const block.
		for _, spec := range genDecl.Specs {
			vspec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			// Append each constant's name to the result.
			for _, name := range vspec.Names {
				constNames = append(constNames, name.Name)
			}
		}
	}

	return constNames
}

func generateEnumMapping(consts []string, cf CauseField) []EnumMapping {
	var enumMappings []EnumMapping

	for _, constName := range consts {
		enumMapping := EnumMapping{
			ConstName: constName,
			ErrorMsg:  fmt.Sprintf("%s : %s", cf.FieldName, strings.Split(constName, "Present")[1]),
		}
		enumMappings = append(enumMappings, enumMapping)
	}
	return enumMappings
}

func getCauseFieldsInfo(causeStruct *ast.StructType, genCtx generatorCtx) []CauseField {
	causeStructName := genCtx.structToFind
	pkgName := genCtx.packageName

	var causeFields []CauseField

	for _, field := range causeStruct.Fields.List {
		// We check if the field is a pointer, as only the pointers are relevant for us
		star, ok := field.Type.(*ast.StarExpr)
		if !ok {
			continue
		} else {
			causeField := CauseField{}

			// In our case, the field have only one identifier, so we simplify by getting the first element.
			causeField.FieldName = field.Names[0].Name
			causeField.PresentConst = fmt.Sprintf("%s.CausePresent%s", pkgName, causeField.FieldName)

			// Name of the struct
			ident, ok := star.X.(*ast.Ident)
			if !ok {
				continue
			}

			causeField.PointerType = ident.Name

			file := getFileFromTypeName(genCtx.pkg, ident.Name)

			causeConst := getConstListFromFile(file)
			if len(causeConst) == 0 {
				causeField.EnumMappings = nil
			} else {
				causeField.EnumMappings = generateEnumMapping(causeConst, causeField)
			}

			causeFields = append(causeFields, causeField)
		}
	}
	if len(causeFields) == 0 {
		log.Fatalf("No cause fields found in struct %s", causeStructName)
	}
	return causeFields
}

// Initialize the context object that will be used during the whole generation process.
func initCauseGenCtxAndData() generatorCtx {
	genCtx := generatorCtx{
		pkgImportPath: "github.com/free5gc/ngap",
		packageName:   "ngapType",
		structToFind:  "Cause",
	}
	genCtx.pkgFullName = genCtx.pkgImportPath + "/" + genCtx.packageName

	pkg := getPackage(genCtx.pkgFullName)
	genCtx.pkg = pkg

	return genCtx
}

// The focus of this function is user readability, optimisation could have been made in pure information parsing.
// But for later maintenance, I chose to separate each steps even if it means re iterating on the same struct multiple
// times (i.e: the ast.File struct).
func getCauseData(ctx generatorCtx) CauseData {
	data := CauseData{
		PackageImportPath: ctx.pkgImportPath,
		PackageName:       ctx.packageName,
		StructName:        ctx.structToFind,
	}

	causeTypeFile := getFileFromTypeName(ctx.pkg, ctx.structToFind)
	causeStruct := getStructFromFile(causeTypeFile, ctx.structToFind)

	causeFields := getCauseFieldsInfo(causeStruct, ctx)
	data.Fields = causeFields

	return data
}

func generateFile(data CauseData) {
	// Template for the generated error message function.
	tmplText := `// Code generated by ngap_type_cause_generator.go; DO NOT EDIT.
package ngap

import (
	"{{ .PackageImportPath }}/{{ .PackageName }}"
	_ "golang.org/x/tools/go/packages"
)

{{- range .Fields }}
func get{{ $.StructName }}{{ .FieldName }}ErrorStr({{ $.StructName | toLower }} *{{ $.PackageName }}.{{ .PointerType }}) string {
	{{- if .EnumMappings | isNil }}
		return "{{ .FieldName }} : Unknown error" 
	{{- else }}
		switch {{ $.StructName | toLower }}.Value {
			{{- range .EnumMappings }}
			case {{ $.PackageName }}.{{ .ConstName }}:
				return "{{ .ErrorMsg }}"
			{{- end }}
			default:
				return "unknown cause"
		}
		return "unknown cause"
	{{- end }}
}
{{ end }}

func GetCauseErrorStr({{ .StructName | toLower }} *{{ .PackageName }}.{{ .StructName }}) string {
	if {{ .StructName | toLower }} != nil {
		switch {{ .StructName | toLower }}.Present {
			{{- range .Fields }}
			case {{ .PresentConst }}:
				return get{{ $.StructName }}{{ .FieldName }}ErrorStr({{ $.StructName | toLower }}.{{ .FieldName }})
			{{- end }}
			default:
				return "unknown {{ .PackageName }}.{{ $.StructName }}"
		}
	}

	return "unknown {{ .PackageName }}.{{ $.StructName }}"
}
	`

	funcMap := template.FuncMap{
		"toLower": strings.ToLower,
		"isNil": func(x interface{}) bool {
			if x == nil {
				return true
			}
			v := reflect.ValueOf(x)
			switch v.Kind() {
			case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
				return v.IsNil()
			}
			return false
		},
	}

	tmpl, err := template.New("errFunc").Funcs(funcMap).Parse(tmplText)
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	outFile, err := os.Create("error_message_gen.go")
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}
	defer outFile.Close()

	if err := tmpl.Execute(outFile, data); err != nil {
		log.Fatalf("Error executing template: %v", err)
	}
	log.Println("Generated error_message_gen.go successfully.")
}

func main() {
	// We set the ctx object with information that will be used during the creation of the CauseData Object.
	genCtx := initCauseGenCtxAndData()
	causeData := getCauseData(genCtx)

	generateFile(causeData)
}
