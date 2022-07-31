package main

import "fmt"

type Interpreter struct {
	functions []Function
}

type ValeurCode int

const (
	// Special tokens
	CODE_INT ValeurCode = iota
	CODE_STRING
)

type Valeur struct {
	code         ValeurCode
	valeurInt    int
	valeurString string
}

func NewInterpreter(functions []Function) *Interpreter {
	return &Interpreter{functions: functions}
}

func (interpreter *Interpreter) getIntValue(expression *Expression, symbolTable map[string]Valeur) (*Valeur, error) {
	if expression.code == EXPR_CODE_INT {
		return &Valeur{code: CODE_INT, valeurInt: expression.valeurInt}, nil
	} else if expression.code == EXPR_CODE_STR {
		return &Valeur{code: CODE_STRING, valeurString: expression.valeurString}, nil
	} else if expression.code == EXPR_CODE_VAR {
		if val, ok := symbolTable[expression.variable]; ok {
			return &val, nil
		} else {
			return nil, fmt.Errorf("variable %s not declared", expression.variable)
		}
	} else if expression.code == EXPR_CODE_ADD || expression.code == EXPR_CODE_SUB {
		val, err := interpreter.getIntValue(expression.left, symbolTable)
		if err != nil {
			return nil, fmt.Errorf("error: %s", err)
		}
		val2, err2 := interpreter.getIntValue(expression.right, symbolTable)
		if err2 != nil {
			return nil, fmt.Errorf("error: %s", err2)
		}
		if expression.code == EXPR_CODE_ADD {
			if val.code == CODE_INT && val2.code == CODE_INT {
				val3 := val.valeurInt + val2.valeurInt
				return &Valeur{code: CODE_INT, valeurInt: val3}, nil
			} else {
				return nil, fmt.Errorf("error: var is not int")
			}
		} else if expression.code == EXPR_CODE_SUB {
			if val.code == CODE_INT && val2.code == CODE_INT {
				val3 := val.valeurInt - val2.valeurInt
				return &Valeur{code: CODE_INT, valeurInt: val3}, nil
			} else {
				return nil, fmt.Errorf("error: var is not int")
			}
		}
	}

	return nil, fmt.Errorf("expression not valid")
}

func (interpreter *Interpreter) interpreter() ([]map[string]Valeur, error) {

	var res []map[string]Valeur = nil
	for _, function := range interpreter.functions {
		fmt.Printf("function %s\n", function.Name)

		symbolTable := make(map[string]Valeur)

		for _, instruction := range function.Instruction {
			fmt.Printf("%s=", instruction.Variable)
			val, err := interpreter.getIntValue(instruction.Valeur, symbolTable)
			if err != nil {
				return nil, fmt.Errorf("error: %s", err)
			}
			if val.code == CODE_INT {
				fmt.Printf("%d", val.valeurInt)
			} else if val.code == CODE_STRING {
				fmt.Printf("%s", val.valeurString)
			}
			symbolTable[instruction.Variable] = *val
			fmt.Printf("\n")
		}

		res = append(res, symbolTable)
	}

	return res, nil
}
