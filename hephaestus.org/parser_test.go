package main

import (
	"reflect"
	"strings"
	"testing"
)

// Ensure the parser can parse strings into Statement ASTs.
func TestParser_ParseStatement(t *testing.T) {
	var tests = []struct {
		s     string
		funct []Function
		err   string
	}{
		// test ast
		{
			s: `void main () { x=5;y=18;}`,
			funct: []Function{{
				name: "main",
				instruction: []Instruction{
					{
						variable: "x",
						valeur:   5,
					}, {
						variable: "y",
						valeur:   18,
					},
				},
			},
			},
		},
		{
			s: `void test123() { abc=10; zzz=156;}`,
			funct: []Function{{
				name: "test123",
				instruction: []Instruction{
					{
						variable: "abc",
						valeur:   10,
					}, {
						variable: "zzz",
						valeur:   156,
					},
				},
			},
			},
		},
		// Errors
		{s: `void main()`, err: `found "", expected {`},
	}

	for i, tt := range tests {
		stmt, err := NewParser(strings.NewReader(tt.s)).Parse2()
		if !reflect.DeepEqual(tt.err, errstring(err)) {
			t.Errorf("%d. %q: error mismatch:\n  exp=%s\n  got=%s\n\n", i, tt.s, tt.err, err)
		} else if tt.err == "" && !reflect.DeepEqual(tt.funct, stmt) {
			t.Errorf("%d. %q\n\nstmt mismatch:\n\nexp=%#v\n\ngot=%#v\n\n", i, tt.s, tt.funct, stmt)
		}
	}
}

// errstring returns the string representation of an error.
func errstring(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}
