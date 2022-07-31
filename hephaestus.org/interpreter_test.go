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
				"x": {valeurtype: Type{code: TYPE_INT}, valeurInt: 5},
				"y": {valeurtype: Type{code: TYPE_INT}, valeurInt: 18},
			},
		},
		{
			s: `void main () { x=10;y=26;z=x+15;}`,
			symbolTable: map[string]Valeur{
				"x": {valeurtype: Type{code: TYPE_INT}, valeurInt: 10},
				"y": {valeurtype: Type{code: TYPE_INT}, valeurInt: 26},
				"z": {valeurtype: Type{code: TYPE_INT}, valeurInt: 25},
			},
		},
		{
			s: `void main () { x="abc";}`,
			symbolTable: map[string]Valeur{
				"x": {valeurtype: Type{code: TYPE_STRING}, valeurString: "abc"},
			},
		},
		{
			s: `void main () { x=5<=7;y=15<20;z=8>3;t=16>=12;w=36==42;}`,
			symbolTable: map[string]Valeur{
				"x": {valeurtype: Type{code: TYPE_BOOLEAN}, valeurBoolean: true},
				"y": {valeurtype: Type{code: TYPE_BOOLEAN}, valeurBoolean: true},
				"z": {valeurtype: Type{code: TYPE_BOOLEAN}, valeurBoolean: true},
				"t": {valeurtype: Type{code: TYPE_BOOLEAN}, valeurBoolean: true},
				"w": {valeurtype: Type{code: TYPE_BOOLEAN}, valeurBoolean: false},
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
			t.Errorf("%d. %q: error no program to execute (err:%s)\n", i, tt.s, err)
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
				t.Errorf("%d. %q\n\nstmt mismatch:\n\nexp=%#v\n\ngot=%#v\n\n", i, tt.s, tt.symbolTable, symbolTableList[0])
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
