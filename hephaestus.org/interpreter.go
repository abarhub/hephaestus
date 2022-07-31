package main

import "fmt"

type Interpreter struct {
	functions []Function
}

type Valeur struct {
	valeurtype    Type
	valeurInt     int
	valeurString  string
	valeurBoolean bool
}

func NewInterpreter(functions []Function) *Interpreter {
	return &Interpreter{functions: functions}
}

func (interpreter *Interpreter) getIntValue(expression *Expression, symbolTable map[string]Valeur) (*Valeur, error) {
	if expression.code == EXPR_CODE_INT {
		return &Valeur{valeurtype: Type{code: TYPE_INT}, valeurInt: expression.valeurInt}, nil
	} else if expression.code == EXPR_CODE_STR {
		return &Valeur{valeurtype: Type{code: TYPE_STRING}, valeurString: expression.valeurString}, nil
	} else if expression.code == EXPR_CODE_TRUE || expression.code == EXPR_CODE_FALSE {
		return &Valeur{valeurtype: Type{code: TYPE_BOOLEAN}, valeurBoolean: expression.code == EXPR_CODE_TRUE}, nil
	} else if expression.code == EXPR_CODE_STR {
		return &Valeur{valeurtype: Type{code: TYPE_STRING}, valeurString: expression.valeurString}, nil
	} else if expression.code == EXPR_CODE_VAR {
		if val, ok := symbolTable[expression.variable]; ok {
			return &val, nil
		} else {
			return nil, fmt.Errorf("variable %s not declared", expression.variable)
		}
	} else if expression.code == EXPR_CODE_ADD || expression.code == EXPR_CODE_SUB ||
		expression.code == EXPR_CODE_EQU || expression.code == EXPR_CODE_LT || expression.code == EXPR_CODE_LTE ||
		expression.code == EXPR_CODE_GT || expression.code == EXPR_CODE_GTE {
		val, err := interpreter.getIntValue(expression.left, symbolTable)
		if err != nil {
			return nil, fmt.Errorf("error: %s", err)
		}
		val2, err2 := interpreter.getIntValue(expression.right, symbolTable)
		if err2 != nil {
			return nil, fmt.Errorf("error: %s", err2)
		}
		if expression.code == EXPR_CODE_ADD || expression.code == EXPR_CODE_SUB {
			if val.valeurtype.code == TYPE_INT && val2.valeurtype.code == TYPE_INT {
				var val3 int
				switch expression.code {
				case EXPR_CODE_ADD:
					val3 = val.valeurInt + val2.valeurInt
					break
				case EXPR_CODE_SUB:
					val3 = val.valeurInt - val2.valeurInt
					break
				default:
					return nil, fmt.Errorf("error: invalid opertator")
				}
				return &Valeur{valeurtype: Type{code: TYPE_INT}, valeurInt: val3}, nil
			} else {
				return nil, fmt.Errorf("error: var is not int")
			}
		} else if expression.code == EXPR_CODE_EQU || expression.code == EXPR_CODE_LT || expression.code == EXPR_CODE_LTE ||
			expression.code == EXPR_CODE_GT || expression.code == EXPR_CODE_GTE {
			if val.valeurtype.code == TYPE_INT && val2.valeurtype.code == TYPE_INT {
				var val3 bool
				switch expression.code {
				case EXPR_CODE_EQU:
					val3 = val.valeurInt == val2.valeurInt
					break
				case EXPR_CODE_LT:
					val3 = val.valeurInt < val2.valeurInt
					break
				case EXPR_CODE_LTE:
					val3 = val.valeurInt <= val2.valeurInt
					break
				case EXPR_CODE_GT:
					val3 = val.valeurInt > val2.valeurInt
					break
				case EXPR_CODE_GTE:
					val3 = val.valeurInt >= val2.valeurInt
					break
				default:
					return nil, fmt.Errorf("error: invalid opertator")
				}
				return &Valeur{valeurtype: Type{code: TYPE_BOOLEAN}, valeurBoolean: val3}, nil
			} else {
				return nil, fmt.Errorf("error: var is not int")
			}
		} else {
			return nil, fmt.Errorf("error: invalid opertator")
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
			if val.valeurtype.code == TYPE_INT {
				fmt.Printf("%d", val.valeurInt)
			} else if val.valeurtype.code == TYPE_STRING {
				fmt.Printf("%s", val.valeurString)
			} else if val.valeurtype.code == TYPE_BOOLEAN {
				fmt.Printf("%t", val.valeurBoolean)
			}
			symbolTable[instruction.Variable] = *val
			fmt.Printf("\n")
		}

		res = append(res, symbolTable)
	}

	return res, nil
}
