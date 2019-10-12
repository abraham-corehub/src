package main

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

func main() {
	testRun()
}

func testRun() {
	expressions := []string{
		"3+4-56-16",
		"3-4-6-16",
		"1--675+5",
		"--675+5",
		"^&*%(",
		"2345",
		"2-------4",
		"--------4",
		"36/12-2^8/2^6+6-12-16^3+5*9/3",
	}
	resultActual := []int{
		-65, -23, 681, 680, 0, 2345, -2, 4, -4088,
	}

	testPass := true
	expressionsFailed := make([]string, 0)
	errorsFailed := make([]error, 0)

	for indexExpression, expression := range expressions {
		_, resultTest, err := evaluateExpression(expression)

		if err == nil {
			if resultActual[indexExpression] != resultTest {
				testPass = false
				expressionsFailed = append(expressionsFailed, expression)
			}
			fmt.Print(expression, " = ", resultTest, "\n")
		} else {
			//fmt.Print(expression, " = ", err, "\n")
			testPass = false
			expressionsFailed = append(expressionsFailed, expression)
			errorsFailed = append(errorsFailed, err)
		}
	}
	if testPass {
		fmt.Print("Test Pass!")
	} else {
		for indexE, expressionFailed := range expressionsFailed {
			fmt.Print("Test Failed @ Expression :", expressionFailed, ", ", errorsFailed[indexE], "\n")
		}
	}
}

func evaluateExpression(expression string) (string, int, error) {
	expression = removeInvalidCharacters(expression)
	arrayIntegers, arrayOperators := performSingleOperandEvaluation(expression)
	result := 0
	var err error
	if len(arrayIntegers) <= 0 {
		return expression, result, errors.New("no numbers found")
	} else if len(arrayIntegers) == 1 {
		return expression, arrayIntegers[0], nil
	}
	operationOrder := "^/*-+"
	for _, operator := range operationOrder {
		for {
			indexOperator := strings.Index(arrayOperators, string(operator))
			if indexOperator >= 0 {
				operandA := arrayIntegers[indexOperator]
				operandB := arrayIntegers[indexOperator+1]
				result, err = evaluateExpressionTwoOperands(operandA, operandB, rune(operator))
				if err != nil {
					return expression, result, err
				} else if len(arrayIntegers) > 1 {
					arrayIntegers = replaceOperandsWithResult(arrayIntegers, indexOperator, result)
					arrayOperators = removeOperatorAtIndex(arrayOperators, indexOperator)
					//fmt.Println(operandA, string(operator), operandB, "=", result, "->", combineNumbersAndOperators(arrayIntegers, arrayOperators))
					//fmt.Println(combineNumbersAndOperatorsIntoExpression(arrayIntegers, arrayOperators))
				}
			} else {
				break
			}
		}
	}

	return expression, result, err
}

func removeInvalidCharacters(expression string) string {
	//fmt.Println(expression)
	validCharacters := "^*/-+1234567890"
	validExpression := expression
	flagValid := false
	for _, character := range expression {
		for _, validCharacter := range validCharacters {
			if character == validCharacter {
				flagValid = true
				break
			} else {
				flagValid = false
			}
		}
		if !flagValid {
			//fmt.Print("Removing invalid character '", string(character), "'\n")
			validExpression = strings.Replace(validExpression, string(character), "", -1)
		}
	}
	//fmt.Println(validExpression)
	return validExpression
}

func performSingleOperandEvaluation(expression string) ([]int, string) {
	arrayIntegers, arrayOperators, arrayPositionsOperators, arrayPositionsIntegers := separateNumbersOperatorsAndTheirPositionsFromExpression(expression)
	arrayIntegers, arrayOperators = evaluateExpressionSingleOperands(arrayIntegers, arrayOperators, arrayPositionsOperators, arrayPositionsIntegers)
	//fmt.Print(arrayIntegers, arrayOperators, "\n")
	return arrayIntegers, arrayOperators
}

func evaluateExpressionSingleOperands(arrayIntegers []int, arrayOperators string, arrayPositionsOperators []int, arrayPositionsIntegers []int) ([]int, string) {
	//fmt.Print(arrayOperators, arrayIntegers, arrayPositionsOperators, arrayPositionsIntegers, "\n")
	for _, positionOperator := range arrayPositionsOperators {
		positionIntegerSingleOperand := findPositionOFNextIntegerAfterOperatorPosition(positionOperator, arrayPositionsIntegers)
		indexOperatorSingleOperand := findIndexOfOperatorInArrayOperatorsAtPosition(positionIntegerSingleOperand-1, arrayPositionsOperators, arrayOperators)
		indexIntegerSingleOperand := findIndexOfIntegerInArrayIntegersAtPosition(positionIntegerSingleOperand, arrayPositionsIntegers, arrayIntegers)
		if indexOperatorSingleOperand >= 0 {
			distanceBetweenOperatorAndInteger := positionIntegerSingleOperand - positionOperator
			if (distanceBetweenOperatorAndInteger > 1 && positionIntegerSingleOperand >= 0) || (distanceBetweenOperatorAndInteger == 1 && indexIntegerSingleOperand == 0) {
				arrayIntegers[indexIntegerSingleOperand] = applySign(rune(arrayOperators[indexOperatorSingleOperand]), arrayIntegers[indexIntegerSingleOperand])
				arrayOperators = removeOperatorAtIndex(arrayOperators, indexOperatorSingleOperand)
				arrayPositionsOperators = updateArrayPositionsOperators(arrayPositionsOperators, indexOperatorSingleOperand)
				arrayPositionsIntegers = updateArrayPositionsIntegers(arrayPositionsIntegers, indexIntegerSingleOperand)
				arrayIntegers, arrayOperators = evaluateExpressionSingleOperands(arrayIntegers, arrayOperators, arrayPositionsOperators, arrayPositionsIntegers)
			}
		}
	}
	return arrayIntegers, arrayOperators
}

func removeOperatorAtIndex(input string, index int) string {
	output := strings.Join([]string{input[:index], input[index+1:]}, "")
	return output
}

func removeValueAtIndex(arrayInput []int, index int) []int {
	arrayOutput := append(arrayInput[:index], arrayInput[index+1:]...)
	return arrayOutput
}

func findPositionOFNextIntegerAfterOperatorPosition(positionOperator int, arrayPositionsIntegers []int) int {
	for _, positionInteger := range arrayPositionsIntegers {
		if positionInteger-positionOperator >= 1 {
			return positionInteger
		}
	}
	return -1
}

func findIndexOfIntegerInArrayIntegersAtPosition(position int, arrayPositionsIntegers []int, arrayIntegers []int) int {
	for indexArrayPositionsIntegers, positionInteger := range arrayPositionsIntegers {
		if positionInteger == position {
			return indexArrayPositionsIntegers
		}
	}
	return -1
}

func findIndexOfOperatorInArrayOperatorsAtPosition(position int, arrayPositionsOperators []int, arrayOperators string) int {
	for indexArrayPositionsOperators, positionOperator := range arrayPositionsOperators {
		if positionOperator == position {
			return indexArrayPositionsOperators
		}
	}
	return -1
}
func updateArrayPositionsOperators(arrayPositionsOperators []int, startIndex int) []int {

	for i := startIndex; i < len(arrayPositionsOperators); i++ {
		arrayPositionsOperators[i]--
	}

	return removeValueAtIndex(arrayPositionsOperators, startIndex)
}

func updateArrayPositionsIntegers(arrayPositionsIntegers []int, startIndex int) []int {
	for i := startIndex; i < len(arrayPositionsIntegers); i++ {
		arrayPositionsIntegers[i]--
	}
	return arrayPositionsIntegers
}

func applySign(operator rune, integer int) int {
	switch operator {
	case '-':
		return integer * -1
	case '+':
		return integer
	}
	return integer
}

func replaceOperandsWithResult(arrayIntegers []int, indexOperator int, result int) []int {
	return append(append(arrayIntegers[:indexOperator], result), arrayIntegers[indexOperator+2:]...)
}

func combineNumbersAndOperatorsIntoExpression(arrayIntegers []int, operators string) string {
	var expression string
	for indexOperator, operator := range operators {
		expression = strings.Join([]string{expression, strconv.Itoa(arrayIntegers[indexOperator]), string(operator)}, "")
	}
	expression = strings.Join([]string{expression, strconv.Itoa(arrayIntegers[len(arrayIntegers)-1])}, "")
	return expression
}

func separateNumbersOperatorsAndTheirPositionsFromExpression(expression string) ([]int, string, []int, []int) {
	arrayIntegerSingleDigit := make([]int8, 0)
	arrayIntegerMultipleDigit := make([]int, 0)
	arrayOperators := make([]rune, 0)
	arrayPositionsOperators := make([]int, 0)
	arrayPositionsIntegers := make([]int, 0)
	positionOffset := 0
	for position, character := range expression {
		integerSingleDigit, err := strconv.Atoi(string(character))
		if err == nil {
			arrayIntegerSingleDigit = append(arrayIntegerSingleDigit, int8(integerSingleDigit))
			positionOffset++
		} else {
			arrayOperators = append(arrayOperators, character)
			arrayPositionsOperators = append(arrayPositionsOperators, position)
			if len(arrayIntegerSingleDigit) > 0 {
				arrayIntegerMultipleDigit = append(arrayIntegerMultipleDigit, convertArrayIntegerSingleDigitToIntegerMultipleDigit(arrayIntegerSingleDigit))
				arrayPositionsIntegers = append(arrayPositionsIntegers, position-positionOffset)
				positionOffset = 0
				arrayIntegerSingleDigit = make([]int8, 0)
			}
		}
	}
	if len(arrayIntegerSingleDigit) > 0 {
		arrayIntegerMultipleDigit = append(arrayIntegerMultipleDigit, convertArrayIntegerSingleDigitToIntegerMultipleDigit(arrayIntegerSingleDigit))
		arrayPositionsIntegers = append(arrayPositionsIntegers, len(expression)-positionOffset)
	}
	return arrayIntegerMultipleDigit, string(arrayOperators), arrayPositionsOperators, arrayPositionsIntegers
}

func convertArrayIntegerSingleDigitToIntegerMultipleDigit(arrayIntegerSingleDigit []int8) int {
	var integerMultipleDigit int
	for index, value := range arrayIntegerSingleDigit {
		integerMultipleDigit += int(value) * int(math.Pow10((len(arrayIntegerSingleDigit) - 1 - index)))
	}
	return integerMultipleDigit
}

func evaluateExpressionTwoOperands(operandA int, operandB int, operator rune) (int, error) {
	result := 0
	var err error
	switch string(operator) {
	case "^":
		result = int(math.Pow(float64(operandA), float64(operandB)))
	case "/":
		if operandB != 0 {
			result = operandA / operandB
		} else {
			err = errors.New("divide by zero")
		}
	case "*":
		result = operandA * operandB
	case "+":
		result = operandA + operandB
	case "-":
		result = operandA - operandB
	default:
		result = 0
	}
	return result, err
}

func strJoin(strA string, strB string) string {
	return strings.Join([]string{strA, strB}, "")
}

func bytesJoin(byteA []byte, byteB []byte) []byte {
	return bytes.Join([][]byte{byteA, byteB}, []byte(""))
}
