package main

import "fmt"

type Interpreter struct {
	functions []Function
}

func NewInterpreter(functions []Function) *Interpreter {
	return &Interpreter{functions: functions}
}

func (interpreter *Interpreter) getIntValue(expression *Expression, symbolTable map[string]int) (int, error) {
	if expression.code == EXPR_CODE_INT {
		return expression.valeurInt, nil
	} else if expression.code == EXPR_CODE_VAR {
		if val, ok := symbolTable[expression.variable]; ok {
			return val, nil
		} else {
			return 0, fmt.Errorf("variable %s not declared", expression.variable)
		}
	} else if expression.code == EXPR_CODE_ADD || expression.code == EXPR_CODE_SUB {
		val, err := interpreter.getIntValue(expression.left, symbolTable)
		if err != nil {
			return 0, fmt.Errorf("error: %s", err)
		}
		val2, err2 := interpreter.getIntValue(expression.right, symbolTable)
		if err2 != nil {
			return 0, fmt.Errorf("error: %s", err2)
		}
		if expression.code == EXPR_CODE_ADD {
			return val + val2, nil
		} else if expression.code == EXPR_CODE_SUB {
			return val - val2, nil
		}
	}

	return 0, nil
}

func (interpreter *Interpreter) interpreter() ([]map[string]int, error) {

	var res []map[string]int = nil
	for _, function := range interpreter.functions {
		fmt.Printf("function %s\n", function.name)

		symbolTable := make(map[string]int)

		for _, instruction := range function.instruction {
			fmt.Printf("%s=", instruction.variable)
			val, err := interpreter.getIntValue(instruction.valeur, symbolTable)
			if err != nil {
				return nil, fmt.Errorf("error: %s", err)
			}
			fmt.Printf("%d", val)
			symbolTable[instruction.variable] = val
			fmt.Printf("\n")
		}

		res = append(res, symbolTable)
	}

	return res, nil
}
