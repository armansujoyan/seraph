package parser

import (
	"errors"
	"fmt"
	"seraph/src/allocator"
	"seraph/src/common"
	"seraph/src/generator"
	"seraph/src/scanner"
	"seraph/src/utils"
)

var supportedOperands = []string{"+", "-", "*", "(", ")"}

func parseExpression(target *scanner.Token, iterator *scanner.TokenIterator) error {
	operatorStack := utils.NewStack[*common.Operator]()
	operandStack := utils.NewStack[*allocator.Register]()
	for iterator.HasMore() {
    next, _ := iterator.ViewNext()
    if next.Value == ";" {
      break;
    }
		token, err := iterator.Next()
		if errors.Is(err, scanner.ErrExhaustedInput) {
			return nil
		}
		switch token.Category {
		case "ident":
      token.Value = modularizeIdentifer(token.Value)
			if !symbolTable.Exists(token.Value) {
				return errors.New("Unknown identifier")
			}
			register, err := registerAllocator.Allocate()
			if err != nil {
				return fmt.Errorf("Can't parse experssion: %w", err)
			}
			register.Load(token.Value)
      register.SetType(allocator.VariableRegister)
			operandStack.Push(register)
		case "number":
			register, err := registerAllocator.Allocate()
			if err != nil {
				return fmt.Errorf("Can't parse experssion: %w", err)
			}
			register.Load(token.Value)
      register.SetType(allocator.NumberRegister)
			operandStack.Push(register)
		case "term":
			if !utils.Contains(supportedOperands, token.Value) {
				return errors.New("Invalid operand. Expected +, - or *")
			}
			operator, err := common.NewOperator(token.Value)
			if err != nil {
				return fmt.Errorf("Can't parse expression: %w", err)
			}
      if operator.Value == "(" {
        operatorStack.Push(operator)
        break;
      }
      if operator.Value == ")" {
        for operatorStack.Peek().Value != "(" {
          evaluateSingleOperand(operatorStack, operandStack)
        }
        operatorStack.Pop();
        break;
      }
			if operatorStack.Peek() != nil && operator.Precedence <= operatorStack.Peek().Precedence {
        evaluateSingleOperand(operatorStack, operandStack)
			}
			operatorStack.Push(operator)
		default:
      return errors.New("Invalid token. Expected operand, number, variable or ';', received: " + token.Value)
		}
	}

	calculatedRegister := exhaustStack(operatorStack, operandStack)
	generator.StoreRegister(target.Value, calculatedRegister)

	return nil
}

func exhaustStack(operators *utils.Stack[*common.Operator], operands *utils.Stack[*allocator.Register]) *allocator.Register {
	for operators.Peek() != nil {
    evaluateSingleOperand(operators, operands)
	}

	register := operands.Pop()
	if !register.GetIsLoaded() {
		generator.LoadRegister(register)
	}

	return register
}

func evaluateSingleOperand(operators *utils.Stack[*common.Operator], operands *utils.Stack[*allocator.Register]) {
    topOperator := operators.Pop()
    firstOperand := operands.Pop()
    secondOperand := operands.Pop()
    newOperand := generator.GenerateArithmeticStatement(firstOperand, secondOperand, topOperator, registerAllocator)
    operands.Push(newOperand)
}
