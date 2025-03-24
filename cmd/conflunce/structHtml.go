package main

import (
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

func getStructTable() (string, error) {
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
		return "", err
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
			return "", err
		}
	}
	htmlContent := builder.String()

	return htmlContent, nil
}

func extractNameFromComment(comment string) string {
	re := regexp.MustCompile(`@name\s+\{(.+?)\}`)
	matches := re.FindStringSubmatch(comment)
	if len(matches) == 2 {
		return matches[1]
	}
	return ""
}
