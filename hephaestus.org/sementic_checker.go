package main

func (p *Parser) Checker(functionList []Function) error {

	for _, function := range functionList {
		for _, instr := range function.Instruction {
			if instr.Valeur != nil {

			}
		}
	}
	return nil
}
