package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strings"
)

type TableModel struct {
	Name     string
	Fields   []string
	Comments []string
}

func main() {
	m := parseFile("example.go")
	c := vueOutput(m)
	f := vueFormoutput(m)
	page := fmt.Sprintf(pageContentTpl, c, f, "[]")
	fmt.Println(page)
}

func parseFile(filename string) TableModel {
	fset := token.NewFileSet()
	f, _ := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	var t TableModel
	for name, obj := range f.Scope.Objects {
		t.Name = name
		tt, ok := obj.Decl.(*ast.TypeSpec)
		if !ok {
			continue
		}
		st, ok := tt.Type.(*ast.StructType)
		if !ok {
			continue
		}
		for _, field := range st.Fields.List {
			val := strings.Trim(field.Tag.Value, "`")
			tag := reflect.StructTag(val)
			t.Fields = append(t.Fields, tag.Get("json"))
			t.Comments = append(t.Comments, strings.TrimSpace(field.Comment.Text()))
		}
	}
	return t
}

const vueLineTpl = `<el-table-column prop="%s" label="%s"></el-table-column>`

func vueOutput(t TableModel) string {
	var buf bytes.Buffer
	for i := range t.Fields {
		name := t.Comments[i]
		if name == "" {
			name = t.Fields[i]
		}
		line := fmt.Sprintf(vueLineTpl, t.Fields[i], t.Comments[i])
		buf.WriteString(line)
		buf.WriteString("\n")
	}
	return buf.String()
}

const vueFormTpl = `<el-form-item label="%s">
          <el-input v-model="form.%s"></el-input>
      </el-form-item>`

func vueFormoutput(t TableModel) string {
	var buf bytes.Buffer
	for i := range t.Fields {
		name := t.Comments[i]
		if name == "" {
			name = t.Fields[i]
		}
		line := fmt.Sprintf(vueFormTpl, name, t.Fields[i])
		buf.WriteString(line)
		buf.WriteString("\n")
	}
	return buf.String()
}

const pageContentTpl = `<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <!-- 引入样式 -->
  <link rel="stylesheet" href="https://unpkg.com/element-ui/lib/theme-default/index.css">
</head>
<body>
  <div id="app">
      <el-table :data="tableData">
         %s
      </el-table>
      <el-form ref="form" :model="form" label-width="80px">  
        %s
      </el-form>
  </div>
</body>
  <!-- 先引入 Vue -->
  <script src="https://unpkg.com/vue/dist/vue.js"></script>
  <!-- 引入组件库 -->
  <script src="https://unpkg.com/element-ui/lib/index.js"></script>
  <script>
    new Vue({
      el: '#app',
      data: function() {
        return { 
          visible: false,
          tableData: %s,
          form:{}
        }
      }
    })
  </script>
</html>`
