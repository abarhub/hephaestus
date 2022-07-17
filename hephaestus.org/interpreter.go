package main

import "fmt"

type Interpreter struct {
	functions []Function
}

func NewInterpreter(functions []Function) *Interpreter {
	return &Interpreter{functions: functions}
}

func (interpreter *Interpreter) interpreter() {

	for _, function := range interpreter.functions {
		fmt.Printf("function %s\n", function.name)

		for _, instruction := range function.instruction {
			fmt.Printf("%s=%d\n", instruction.variable, instruction.valeur)
		}

	}

}
