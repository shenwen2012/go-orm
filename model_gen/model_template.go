package main

import (
	"text/template"
)

var header string = `// Code generated by model_gen
// {{.TableName}}.go contains model for the database table [{{.DbName}}.{{.TableName}}]

package {{.PkgName}}

import (
	"database/sql"
	{{if .ImportTime}}"time"{{end}}
)
`

var modelStruct string = `
func init() {
	addTable({{.Name}}{})
}

type {{.Name}} struct {
	{{range .Fields}}{{.Name}} {{.Type}} {{.Tag}}{{if .Comment}} // {{.Comment}}{{end}}
	{{end}}
}

func (obj {{.Name}}) TableName() string {
	return "{{.TableName}}"
}
`

var objApi string = `
// Start of the {{.Name}} APIs.

func (m *Model) Insert{{.Name}}({{.LowerName}} *{{.Name}}) error {
	return m.Insert({{.LowerName}})
}

func (m *Model) Get{{.Name}}ByPK(id {{.PrimaryField.Type}}) (*{{.Name}}, error) {
	var {{.LowerName}} {{.Name}}
	err := m.SelectByPK(&{{.LowerName}}, id)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err == nil {
		return &{{.LowerName}}, nil
	} else {
		return nil, err
	}
}
`
var testHeader string = `// Code generated by model_gen
// {{.TableName}}_test.go contains model test for the database table [{{.DbName}}.{{.TableName}}]

package {{.PkgName}}

import (
	"testing"
	{{if .ImportTime}}"time"{{end}}
)
`

var testCode string = `
func TestInsertAndGet{{.Name}}(t *testing.T) {
	OneTestScope(func(m *Model) {
		var {{.LowerName}} = &{{.Name}}{
		{{range .Fields}} {{if and (not .IsPrimaryKey) (.DefaultValueCode)}} {{.Name}}: {{.DefaultValueCode}},
		{{end}}{{end}}
		}
		err := m.Insert{{.Name}}({{.LowerName}})
		if err != nil {
			t.Fatalf("failed to Insert{{.Name}}, err: %+v", err)
		}
		loaded, err := m.Get{{.Name}}ByPK({{.LowerName}}.{{.PrimaryField.Name}})
		if err != nil {
			t.Fatalf("failed to Get{{.Name}}ByPK, err: %+v", err)
		}
		if loaded == nil {
			t.Fatalf("should have loaded one {{.Name}}")
		}
	})
}
`

var (
	tmHeader        *template.Template
	tmStruct        *template.Template
	tmObjApi        *template.Template
	tmTestHeader    *template.Template
	tmTestCode      *template.Template
)

func init() {
	tmHeader = template.Must(template.New("header").Parse(header))
	tmStruct = template.Must(template.New("modelStruct").Parse(modelStruct))
	tmObjApi = template.Must(template.New("objApi").Parse(objApi))
	tmTestHeader = template.Must(template.New("testHeader").Parse(testHeader))
	tmTestCode = template.Must(template.New("testCode").Parse(testCode))
}