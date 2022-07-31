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
		symbolTable map[string]Valeur
		err         string
	}{
		// test interpreter
		{
			s: `void main () { x=5;y=18;}`,
			symbolTable: map[string]Valeur{
				"x": {code: CODE_INT, valeurInt: 5},
				"y": {code: CODE_INT, valeurInt: 18},
			},
		},
		{
			s: `void main () { x=10;y=26;z=x+15;}`,
			symbolTable: map[string]Valeur{
				"x": {code: CODE_INT, valeurInt: 10},
				"y": {code: CODE_INT, valeurInt: 26},
				"z": {code: CODE_INT, valeurInt: 25},
			},
		},
		{
			s: `void main () { x="abc";}`,
			symbolTable: map[string]Valeur{
				"x": {code: CODE_STRING, valeurString: "abc"},
			},
		},
		// Errors
		{
			s:   `void main () { x=y;}`,
			err: "error: variable y not declared",
		},
	}

	for i, tt := range tests {
		funct, err := NewParser(strings.NewReader(tt.s)).Parse2()

		if funct == nil {
			t.Errorf("%d. %q: error no program to execute\n", i, tt.s)
		} else {
			interpreter := NewInterpreter(funct)
			var symbolTableList []map[string]Valeur
			symbolTableList, err = interpreter.interpreter()

			if !reflect.DeepEqual(tt.err, errstring2(err)) {
				t.Errorf("%d. %q: error mismatch:\n  exp=%s\n  got=%s\n\n", i, tt.s, tt.err, err)
			} else if err != nil && tt.err == "" {
				t.Errorf("%d. %q: error mismatch:\n  exp=%s\n  got=%s\n\n", i, tt.s, tt.err, err)
			} else if symbolTableList == nil && tt.symbolTable != nil {
				t.Errorf("%d. %q: error no symbol table\n", i, tt.s)
			} else if tt.err == "" && !reflect.DeepEqual(tt.symbolTable, symbolTableList[0]) {
				t.Errorf("%d. %q\n\nstmt mismatch:\n\nexp=%#v\n\ngot=%#v\n\n", i, tt.s, tt.symbolTable, symbolTableList)
			}
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
