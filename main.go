package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"math"
	"os"
	"text/template"
	"time"
)

type (
	Field struct {
		Name string
		Type string
	}
	Class struct {
		Name   string
		Fields []Field
	}
)

var (
	//go:embed template.tmpl
	tmpl string
	i    = 0
	t    *template.Template
)

func init() {
	tmpl, err := template.New("Negra").Parse(tmpl)
	if err != nil {
		panic(err)
	}
	t = tmpl
}

func main() {
	var m map[string]interface{}
	err := json.NewDecoder(os.Stdin).Decode(&m)
	if err != nil {
		panic(err)
	}

	err = write(os.Stdout, parseClass(m))
	if err != nil {
		panic(err)
	}
}

func write(w io.Writer, classes []Class) error {
	return t.Execute(w, classes)
}

func parseClass(decodedJson map[string]interface{}) []Class {
	data := Class{
		Name: fmt.Sprintf("ChangeMe%d", i),
	}
	i++
	var referencedClasses []Class
	for fieldName, fieldValue := range decodedJson {
		fieldType, classes := getType(fieldValue)
		referencedClasses = append(referencedClasses, classes...)
		data.Fields = append(data.Fields, Field{
			Name: fieldName,
			Type: fieldType,
		})
	}
	return append([]Class{data}, referencedClasses...)
}

func getType(v interface{}) (string, []Class) {
	switch v.(type) {
	case bool:
		return "Boolean", nil
	case int, int8, int16, int32, uint, uint8, uint16, uint32:
		return "Long", nil
	case float64:
		f := v.(float64)
		if f == math.Trunc(f) {
			return "Long", nil
		}
		return "Double", nil
	case string:
		if _, err := time.Parse(time.RFC3339, v.(string)); err == nil {
			return "java.util.Date", nil
		}
		if _, err := uuid.Parse(v.(string)); err == nil {
			return "java.util.UUID", nil
		}
		return "String", nil
	case time.Time:
		return "java.util.Date", nil
	case []interface{}:
		elems := v.([]interface{})
		if len(elems) == 0 {
			return "List<Nothing> = emptyList()", nil
		}
		m := v.([]interface{})[0]
		classes := parseClass(m.(map[string]interface{}))
		return fmt.Sprintf("List<%s>", classes[0].Name), classes
	case map[string]interface{}:
		elems := v.(map[string]interface{})
		classes := parseClass(elems)
		return classes[0].Name, classes
	}
	panic(fmt.Sprintf("Unknown shit: %s", v))
}
