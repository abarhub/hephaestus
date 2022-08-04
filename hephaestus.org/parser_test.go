package main

import (
	"encoding/json"
	"fmt"
	"github.com/kr/pretty"
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
				ReturnType: Type{code: TYPE_VOID, position: &Position{
					line: 1, column: 1, pos: 0,
				}},
				Name: "main",
				Instruction: []Instruction{
					{
						Variable: "x",
						Valeur:   &Expression{code: EXPR_CODE_INT, valeurInt: 5, position: &Position{line: 1, column: 1, pos: 17}},
						position: &Position{line: 1, column: 1, pos: 15},
					}, {
						Variable: "y",
						Valeur:   &Expression{code: EXPR_CODE_INT, valeurInt: 18, position: &Position{line: 1, column: 1, pos: 21}},
					},
				},
			},
			},
		},
		{
			s: `void test123() { abc=10; zzz=156;}`,
			funct: []Function{{
				ReturnType: Type{code: TYPE_VOID, position: &Position{
					line: 1, column: 1, pos: 0,
				}},
				Name: "test123",
				Instruction: []Instruction{
					{
						Variable: "abc",
						Valeur:   &Expression{code: EXPR_CODE_INT, valeurInt: 10, position: &Position{line: 1, column: 1, pos: 21}},
						position: &Position{line: 1, column: 1, pos: 17},
					}, {
						Variable: "zzz",
						Valeur:   &Expression{code: EXPR_CODE_INT, valeurInt: 156, position: &Position{line: 1, column: 1, pos: 29}},
					},
				},
			},
			},
		},
		{
			s: `void test3() { x=10; y=x+15;}`,
			funct: []Function{{
				ReturnType: Type{code: TYPE_VOID, position: &Position{
					line: 1, column: 1, pos: 0,
				}},
				Name: "test3",
				Instruction: []Instruction{
					{
						Variable: "x",
						Valeur:   &Expression{code: EXPR_CODE_INT, valeurInt: 10, position: &Position{line: 1, column: 1, pos: 17}},
						position: &Position{line: 1, column: 1, pos: 15},
					}, {
						Variable: "y",
						Valeur: &Expression{code: EXPR_CODE_ADD,
							left:     &Expression{code: EXPR_CODE_VAR, variable: "x", position: &Position{line: 1, column: 1, pos: 23}},
							right:    &Expression{code: EXPR_CODE_INT, valeurInt: 15, position: &Position{line: 1, column: 1, pos: 25}},
							position: &Position{line: 1, column: 1, pos: 24},
						},
					},
				},
			},
			},
		},
		{
			s: `void test3() { x="abc"; y=x;}`,
			funct: []Function{{
				ReturnType: Type{code: TYPE_VOID, position: &Position{
					line: 1, column: 1, pos: 0,
				}},
				Name: "test3",
				Instruction: []Instruction{
					{
						Variable: "x",
						Valeur:   &Expression{code: EXPR_CODE_STR, valeurString: "abc", position: &Position{line: 1, column: 1, pos: 17}},
						position: &Position{line: 1, column: 1, pos: 15},
					}, {
						Variable: "y",
						Valeur: &Expression{code: EXPR_CODE_VAR,
							variable: "x", position: &Position{line: 1, column: 1, pos: 26}},
					},
				},
			},
			},
		},
		{
			s: `void test3() { x=true; y=false;z=5<=7;}`,
			funct: []Function{{
				ReturnType: Type{code: TYPE_VOID, position: &Position{
					line: 1, column: 1, pos: 0,
				}},
				Name: "test3",
				Instruction: []Instruction{
					{
						Variable: "x",
						Valeur:   &Expression{code: EXPR_CODE_TRUE, position: &Position{line: 1, column: 1, pos: 17}},
						position: &Position{line: 1, column: 1, pos: 15},
					}, {
						Variable: "y",
						Valeur:   &Expression{code: EXPR_CODE_FALSE, position: &Position{line: 1, column: 1, pos: 25}},
					}, {
						Variable: "z",
						Valeur: &Expression{code: EXPR_CODE_LTE,
							left:     &Expression{code: EXPR_CODE_INT, valeurInt: 5, position: &Position{line: 1, column: 1, pos: 33}},
							right:    &Expression{code: EXPR_CODE_INT, valeurInt: 7, position: &Position{line: 1, column: 1, pos: 36}},
							position: &Position{line: 1, column: 1, pos: 34},
						},
					},
				},
			},
			},
		},
		{
			s: `void test3() { x=10<3; y=14<=17;z=20>26;t=36>=50;v=40==63;}`,
			funct: []Function{{
				ReturnType: Type{code: TYPE_VOID, position: &Position{
					line: 1, column: 1, pos: 0,
				}},
				Name: "test3",
				Instruction: []Instruction{
					{
						Variable: "x",
						Valeur: &Expression{code: EXPR_CODE_LT,
							left:     &Expression{code: EXPR_CODE_INT, valeurInt: 10, position: &Position{line: 1, column: 1, pos: 17}},
							right:    &Expression{code: EXPR_CODE_INT, valeurInt: 3, position: &Position{line: 1, column: 1, pos: 20}},
							position: &Position{line: 1, column: 1, pos: 19},
						},
						position: &Position{line: 1, column: 1, pos: 15},
					}, {
						Variable: "y",
						Valeur: &Expression{code: EXPR_CODE_LTE,
							left:     &Expression{code: EXPR_CODE_INT, valeurInt: 14, position: &Position{line: 1, column: 1, pos: 25}},
							right:    &Expression{code: EXPR_CODE_INT, valeurInt: 17, position: &Position{line: 1, column: 1, pos: 29}},
							position: &Position{line: 1, column: 1, pos: 27},
						},
					}, {
						Variable: "z",
						Valeur: &Expression{code: EXPR_CODE_GT,
							left:     &Expression{code: EXPR_CODE_INT, valeurInt: 20, position: &Position{line: 1, column: 1, pos: 34}},
							right:    &Expression{code: EXPR_CODE_INT, valeurInt: 26, position: &Position{line: 1, column: 1, pos: 37}},
							position: &Position{line: 1, column: 1, pos: 36},
						},
					}, {
						Variable: "t",
						Valeur: &Expression{code: EXPR_CODE_GTE,
							left:     &Expression{code: EXPR_CODE_INT, valeurInt: 36, position: &Position{line: 1, column: 1, pos: 42}},
							right:    &Expression{code: EXPR_CODE_INT, valeurInt: 50, position: &Position{line: 1, column: 1, pos: 46}},
							position: &Position{line: 1, column: 1, pos: 44},
						},
					}, {
						Variable: "v",
						Valeur: &Expression{code: EXPR_CODE_EQU,
							left:     &Expression{code: EXPR_CODE_INT, valeurInt: 40, position: &Position{line: 1, column: 1, pos: 51}},
							right:    &Expression{code: EXPR_CODE_INT, valeurInt: 63, position: &Position{line: 1, column: 1, pos: 55}},
							position: &Position{line: 1, column: 1, pos: 53},
						},
					},
				},
			},
			},
		},
		{
			s: `void test3() { x=5; y=20;print(x,y);}`,
			funct: []Function{{
				ReturnType: Type{code: TYPE_VOID, position: &Position{
					line: 1, column: 1, pos: 0,
				}},
				Name: "test3",
				Instruction: []Instruction{
					{
						Code:     INSTRUCTION_AFFECTATION,
						Variable: "x",
						Valeur:   &Expression{code: EXPR_CODE_INT, valeurInt: 5, position: &Position{line: 1, column: 1, pos: 17}},
						position: &Position{line: 1, column: 1, pos: 15},
					}, {
						Code:     INSTRUCTION_AFFECTATION,
						Variable: "y",
						Valeur:   &Expression{code: EXPR_CODE_INT, valeurInt: 20, position: &Position{line: 1, column: 1, pos: 22}},
					}, {
						Code:         INSTRUCTION_CALL,
						FunctionName: "print",
						Parameter: []Expression{
							{code: EXPR_CODE_VAR, variable: "x", position: &Position{line: 1, column: 1, pos: 31}},
							{code: EXPR_CODE_VAR, variable: "y", position: &Position{line: 1, column: 1, pos: 33}},
						},
					},
				},
			},
			},
		},
		// Errors
		{s: `void main()`, err: `found "", expected { (pos=&{1 1 10})`},
	}

	for i, tt := range tests {
		stmt, err := NewParser(strings.NewReader(tt.s)).Parse2()

		//if diff := deep.Equal(t1, t2); diff != nil {
		//	t.Error(diff)
		//}

		if !reflect.DeepEqual(tt.err, errstring(err)) {
			t.Errorf("%d. %q: error mismatch:\n  exp=%s\n  got=%s\n\n", i, tt.s, tt.err, err)
			//} else if diff := deep.Equal(tt.funct, stmt); tt.err == "" && diff != nil {
			//	t.Error(diff)
		} else if tt.err == "" && !reflect.DeepEqual(tt.funct, stmt) {
			fmt.Print("abc\n")
			//fmt.Print("funct:" + spew.Sdump(tt.funct) + "\n")
			//fmt.Print("stmt:" + spew.Sdump(stmt) + "\n")
			//fmt.Printf("format:%# v\n", pretty.Formatter(tt.funct))
			//fmt.Printf("format2:%# v\n", pretty.Formatter(stmt))
			fmt.Printf("diff:%# v", pretty.Formatter(pretty.Diff(tt.funct, stmt)))
			fmt.Print("aaa\n")
			//assert.Errorf(t, fmt.Errorf("err"), "%d. %q\n\nstmt mismatch:\n\nexp=%#v\n\ngot=%#v\n\n", i, tt.s, tt.funct, stmt)
			t.Errorf("%d. %q\n\nstmt mismatch:\n\nexp=%#v\n\ngot=%#v\n\n", i, tt.s, tt.funct, stmt)
			//t.Errorf("%d. %q\n\nstmt mismatch:\n\nexp=%#v\n\ngot=%#v\n\n%#v", i, tt.s, tt.funct, stmt, spew.Sdump(tt.funct, stmt))
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

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "")
	return string(s)
}
