package main

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"sort"
	"strings"
	"testing"
)

func getClass(str string) Class {
	var m map[string]interface{}
	err := json.NewDecoder(strings.NewReader(str)).Decode(&m)
	if err != nil {
		panic(err)
	}
	classes := parseClass(m)
	return classes[0]
}

func assertJson(t *testing.T, str, wanted string) {
	var m map[string]interface{}
	err := json.NewDecoder(strings.NewReader(str)).Decode(&m)
	if err != nil {
		panic(err)
	}
	classes := parseClass(m)
	var b strings.Builder
	err = write(&b, classes)
	if err != nil {
		panic(err)
	}
	if b.String() != wanted {
		assert.Equal(t, wanted, b.String())
	}
}

func assertClasses(t *testing.T, str string, classes []Class) {
	var m map[string]interface{}
	err := json.NewDecoder(strings.NewReader(str)).Decode(&m)
	if err != nil {
		panic(err)
	}
	doSort := func(c []Class) {
		sort.Slice(c, func(i, j int) bool {
			return c[i].Name < c[j].Name
		})
		for _, class := range c {
			sort.Slice(class.Fields, func(i, j int) bool {
				return class.Fields[i].Name < class.Fields[j].Name
			})
		}
	}
	parsedJsonClasses := parseClass(m)
	doSort(parsedJsonClasses)
	doSort(classes)
	assert.Equal(t, classes, parsedJsonClasses)
}

func TestBasic(t *testing.T) {
	testCases := []struct {
		str    string
		wanted string
	}{
		{`"3"`, "String"},
		{"3", "Long"},
		{"true", "Boolean"},
		{"3.1", "Double"},
		{`"2021-02-25T13:45:24+00:00"`, "java.util.Date"},
		{`"f93e2204-91c7-4a86-bc6f-f65e0b892563"`, "java.util.UUID"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.wanted, func(t *testing.T) {
			c := getClass(fmt.Sprintf(`{"negra": %s}`, testCase.str))
			if c.Fields[0].Type != testCase.wanted {
				t.Errorf("Type of field was wrong")
			}
		})
	}
}

func TestRendering(t *testing.T) {
	testCases := []struct {
		str    string
		wanted string
	}{
		{`"3"`, "String"},
		{"3", "Long"},
		{"true", "Boolean"},
		{"3.1", "Double"},
		{`"2021-02-25T13:45:24+00:00"`, "java.util.Date"},
		{`"f93e2204-91c7-4a86-bc6f-f65e0b892563"`, "java.util.UUID"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.wanted, func(t *testing.T) {
			assertJson(
				t,
				fmt.Sprintf(`{"negra": %s}`, testCase.str),
				fmt.Sprintf(`data class ChangeMe%d(
    val negra: %s,
)

`, i, testCase.wanted))
		})
	}
}

// because someone found a special case
func TestEmptyList(t *testing.T) {
	assertJson(
		t,
		`{"negra": []}`,
		fmt.Sprintf(`data class ChangeMe%d(
    val negra: List<Nothing> = emptyList(),
)

`, i))
}

func TestNestedShit(t *testing.T) {
	i = 0
	assertClasses(
		t,
		`{
	"name": "honorificPrefix",
	"type": "select",
	"label": "Title",
	"options": [{
		"value": "",
		"label": "Mr"
	}]
}`,
		[]Class{
			{
				Name: "ChangeMe0",
				Fields: []Field{
					{"name", "String"},
					{"type", "String"},
					{"label", "String"},
					{"options", "List<ChangeMe1>"},
				},
			},
			{
				Name: "ChangeMe1",
				Fields: []Field{
					{"value", "String"},
					{"label", "String"},
				},
			},
		},
	)
}
