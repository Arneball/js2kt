package main

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
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
