package util

import (
	"reflect"
	"testing"
)

type Case struct {
	name     string
	sf       reflect.StructField
	expected string
}

func TestFieldUse(t *testing.T) {

	testCases := []Case{
		{
			"empty",
			reflect.StructField{
				Name: "Foo",
				Tag:  ``,
			},
			"foo",
		},
		{
			"only other",
			reflect.StructField{
				Name: "Foo",
				Tag:  `something:"other"`,
			},
			"foo",
		},
		{
			"json+other",
			reflect.StructField{
				Name: "Foo",
				Tag:  `json:"bar" something:"other"`,
			},
			"bar",
		},
		{
			"json,+other",
			reflect.StructField{
				Name: "Foo",
				Tag:  `json:",omitempty" something:"other"`,
			},
			"foo",
		},
		{
			"yaml",
			reflect.StructField{
				Name: "Foo",
				Tag:  `yaml:"bar,omitempty"`,
			},
			"bar",
		},
		{
			"yaml,",
			reflect.StructField{
				Name: "Foo",
				Tag:  `yaml:",omitempty"`,
			},
			"foo",
		},
		{
			"json,+yaml",
			reflect.StructField{
				Name: "Foo",
				Tag:  `json:",omitempty" yaml:"bar,omitempty"`,
			},
			"bar",
		},
	}

	for _, tc := range testCases {
		actual := getFieldUse(tc.sf)
		if actual != tc.expected {
			t.Errorf("name: %s, expected %s, got %s", tc.name, tc.expected, actual)
		}
	}
}

func TestOne(t *testing.T) {
	c := Case{
		"json,+other",
		reflect.StructField{
			Name: "Foo",
			Tag:  `json:",omitempty" something:"other"`,
		},
		"foo",
	}
	actual := getFieldUse(c.sf)
	if actual != c.expected {
		t.Errorf("name: %s, expected %s, got %s", c.name, c.expected, actual)
	}
}
