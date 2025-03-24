package main

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
