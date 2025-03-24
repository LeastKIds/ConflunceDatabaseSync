package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"text/template"
)

type FieldInfo struct {
	No        string
	Column    string
	Name      string
	Type      string
	PK        string
	FK        string
	Null      string
	Index     string
	Reference string
}

type StructInfo struct {
	StructName  string
	DisplayName string // 예: @name {유저} → "유저"
	Fields      []FieldInfo
}

// HTML 템플릿 예제
var htmlTemplateStr = `
<h3>{{.StructName}} ({{.DisplayName}})</h3>
<table border="1" cellspacing="0" cellpadding="4">
  <thead>
    <tr>
      <th>No</th>
      <th>칼럼명</th>
      <th>칼럼뜻</th>
      <th>타입</th>
      <th>PK</th>
	  <th>FK</th>
	  <th>Null</th>
	  <th>Index</th>
	  <th>참고</th>
    </tr>
  </thead>
  <tbody>
  {{- range .Fields }}
    <tr>
      <td>{{.No}}</td>
      <td>{{.Column}}</td>
      <td>{{.Name}}</td>
      <td>{{.Type}}</td>
      <td>{{.PK}}</td>
	  <td>{{.FK}}</td>
	  <td>{{.Null}}</td>
	  <td>{{.Index}}</td>
	  <td>{{.Reference}}</td>
    </tr>
  {{- end }}
  </tbody>
</table>
<hr/>
`

func main() {
	var results []StructInfo

	// entity 디렉터리 순회: .go 파일에서 @db 주석이 있는 구조체만 추출
	err := filepath.WalkDir("entity", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// .go 파일만 처리
		if d.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
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
				// 주석 내에 @db, @name 등을 확인
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
					// 임베디드 필드는 무시
					if len(field.Names) == 0 {
						continue
					}
					var no, column, name, typ, pk, fk, null, index, reference string
					if field.Tag != nil {
						tag := reflect.StructTag(strings.Trim(field.Tag.Value, "`"))
						no = tag.Get("no")
						column = tag.Get("column")
						name = tag.Get("name")
						typ = tag.Get("type")
						pk = tag.Get("pk")
						fk = tag.Get("fk")
						null = tag.Get("null")
						index = tag.Get("index")
						reference = tag.Get("reference")
					}
					s.Fields = append(s.Fields, FieldInfo{
						No:        no,
						Column:    column,
						Name:      name,
						Type:      typ,
						PK:        pk,
						FK:        fk,
						Null:      null,
						Index:     index,
						Reference: reference,
					})
				}
				results = append(results, s)
			}
		}
		return nil
	})

	if err != nil {
		log.Fatalf("파일 순회 중 에러 발생: %v", err)
	}

	// 메모리상에서 HTML 콘텐츠 생성
	var builder strings.Builder
	tmpl := template.Must(template.New("htmlTemplate").Parse(htmlTemplateStr))
	for _, r := range results {
		if r.DisplayName == "" {
			r.DisplayName = r.StructName
		}
		if err := tmpl.Execute(&builder, r); err != nil {
			log.Fatalf("템플릿 실행 에러: %v", err)
		}
	}
	htmlContent := builder.String()

	fmt.Println("생성된 HTML 내용:")
	fmt.Println(htmlContent)
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
