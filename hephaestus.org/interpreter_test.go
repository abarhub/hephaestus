package main

import (
	"reflect"
	"strings"
	"testing"
)

// Ensure the parser can parse strings into Statement ASTs.
func TestParser_interpreter(t *testing.T) {
	var tests = []struct {
		s           string
		symbolTable map[string]int
		err         string
	}{
		// test interpreter
		{
			s: `void main () { x=5;y=18;}`,
			symbolTable: map[string]int{
				"x": 5,
				"y": 18,
			},
		},
		{
			s: `void main () { x=10;y=26;z=x+15;}`,
			symbolTable: map[string]int{
				"x": 10,
				"y": 26,
				"z": 25,
			},
		},
		/*{
			s: `void test123() { abc=10; zzz=156;}`,
			funct: []Function{{
				name: "test123",
				instruction: []Instruction{
					{
						variable: "abc",
						valeur:   &Expression{code: EXPR_CODE_INT, valeurInt: 10},
					}, {
						variable: "zzz",
						valeur:   &Expression{code: EXPR_CODE_INT, valeurInt: 156},
					},
				},
			},
			},
		},
		{
			s: `void test3() { x=10; y=x+15;}`,
			funct: []Function{{
				name: "test3",
				instruction: []Instruction{
					{
						variable: "x",
						valeur:   &Expression{code: EXPR_CODE_INT, valeurInt: 10},
					}, {
						variable: "y",
						valeur: &Expression{code: EXPR_CODE_ADD,
							left:  &Expression{code: EXPR_CODE_VAR, variable: "x"},
							right: &Expression{code: EXPR_CODE_INT, valeurInt: 15}},
					},
				},
			},
			},
		},*/
		// Errors
		//{s: `void main()`, err: `found "", expected {`},
	}

	for i, tt := range tests {
		funct, err := NewParser(strings.NewReader(tt.s)).Parse2()

		interpreter := NewInterpreter(funct)
		var symbolTableList []map[string]int
		symbolTableList, err = interpreter.interpreter()

		if !reflect.DeepEqual(tt.err, errstring2(err)) {
			t.Errorf("%d. %q: error mismatch:\n  exp=%s\n  got=%s\n\n", i, tt.s, tt.err, err)
		} else if tt.err == "" && !reflect.DeepEqual(tt.symbolTable, symbolTableList[0]) {
			t.Errorf("%d. %q\n\nstmt mismatch:\n\nexp=%#v\n\ngot=%#v\n\n", i, tt.s, tt.symbolTable, symbolTableList)
		}
	}
}

// errstring returns the string representation of an error.
func errstring2(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}
