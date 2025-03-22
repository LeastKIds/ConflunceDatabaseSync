package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"text/template"
)

type FieldInfo struct {
	Name      string
	Type      string
	JSONTag   string
	DBTag     string
	Reference string
}

type StructInfo struct {
	StructName  string
	DisplayName string // @name {유저} → 유저
	Fields      []FieldInfo
}

var markdownTemplate = `
### {{.StructName}} ({{.DisplayName}})

| 필드 이름 | 타입 | JSON 태그 | DB 태그 | 설명 |
|-----------|------|-----------|---------|------|
{{- range .Fields }}
| {{.Name}} | {{.Type}} | {{.JSONTag}} | {{.DBTag}} | {{.Reference}} |
{{- end }}

`

func main() {
	var results []StructInfo

	// entity 디렉터리 순회
	err := filepath.WalkDir("entity", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		fs := token.NewFileSet()
		node, err := parser.ParseFile(fs, path, nil, parser.ParseComments)
		if err != nil {
			return err
		}

		for _, decl := range node.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.TYPE {
				continue
			}

			var isDBStruct bool
			var displayName string

			if genDecl.Doc != nil {
				for _, comment := range genDecl.Doc.List {
					text := strings.TrimSpace(comment.Text)
					if strings.HasPrefix(text, "// @db") {
						isDBStruct = true
					}
					if strings.HasPrefix(text, "// @name") {
						displayName = extractNameFromComment(text)
					}
				}
			}

			if !isDBStruct {
				continue
			}

			for _, spec := range genDecl.Specs {
				typeSpec := spec.(*ast.TypeSpec)
				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					continue
				}

				s := StructInfo{
					StructName:  typeSpec.Name.Name,
					DisplayName: displayName,
				}

				for _, field := range structType.Fields.List {
					if len(field.Names) == 0 {
						continue
					}

					fieldName := field.Names[0].Name
					fieldType := exprToString(field.Type)

					jsonTag := ""
					dbTag := ""
					reference := ""
					if field.Tag != nil {
						tag := reflect.StructTag(strings.Trim(field.Tag.Value, "`"))
						jsonTag = tag.Get("json")
						dbTag = tag.Get("db")
						reference = tag.Get("reference")
					}

					s.Fields = append(s.Fields, FieldInfo{
						Name:      fieldName,
						Type:      fieldType,
						JSONTag:   jsonTag,
						DBTag:     dbTag,
						Reference: reference,
					})
				}
				results = append(results, s)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("에러:", err)
		os.Exit(1)
	}

	// Markdown 파일 열기/생성
	outFile, err := os.Create("output.md")
	if err != nil {
		fmt.Println("파일 생성 실패:", err)
		os.Exit(1)
	}
	defer outFile.Close()

	tmpl := template.Must(template.New("markdown").Parse(markdownTemplate))
	for _, r := range results {
		if r.DisplayName == "" {
			r.DisplayName = r.StructName
		}
		_ = tmpl.Execute(outFile, r)
	}

	fmt.Println("✅ output.md 파일 생성 완료!")
}

func exprToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + exprToString(t.X)
	case *ast.SelectorExpr:
		return exprToString(t.X) + "." + t.Sel.Name
	case *ast.ArrayType:
		return "[]" + exprToString(t.Elt)
	default:
		return fmt.Sprintf("%T", expr)
	}
}

func extractNameFromComment(comment string) string {
	re := regexp.MustCompile(`@name\s+\{(.+?)\}`)
	matches := re.FindStringSubmatch(comment)
	if len(matches) == 2 {
		return matches[1]
	}
	return ""
}
