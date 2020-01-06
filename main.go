package main

import (
	"bufio"
	"container/list"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/golang-collections/collections/stack"
)

//TODO:
// - add calculation with 'last' like 'last | 0xff'
// - change last when requesting a 'as' value so when 0x2 as dec -> last is 2

type expression struct {
	value      int
	expression string
}

var lastResult int = 0

func main() {
	fmt.Println("GoBitty - terminal bitwise operations calculator")

	input := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		ok := input.Scan()
		if ok == false {
			break
		}

		fmt.Println("-->", handleInput(input.Text()))
	}
}

// input rules:
// there exist 3 input numbers:
// decimal, hexadecimal, binary
// decimal ist just a normal input number 1234
// hexadecimal needs to have 0x before so -> 0x1234
// binary numbers need bx10101010 with MSB input
func handleInput(input string) string {
	switch {
	// quit command
	case input == "quit":
		fmt.Println("Quitting")
		os.Exit(0)
		return ""

	// help command
	case input == "help":
		return "This is the help command!"

	// return last element
	case input == "last":
		return strconv.Itoa(lastResult)

	// conversions
	case strings.Contains(input, " as "):
		params := strings.Split(input, " as ")
		//params[0] is the number argument f.e. '0x1234 as dec'
		if params[0] == "last" {
			return convertNumberToOutputFormat(lastResult, params[1])
		}

		if strings.ContainsAny(input, "&|^><~") {
			result, ok := calculateValue(params[0])
			if ok != true {
				return "Could not apply bitwise operation!"
			}
			return convertNumberToOutputFormat(result, params[1])

		}
		inputNumber, ok := convertInputToInt(params[0])
		if ok != true {
			return "Input number has wrong format!"
		}
		return convertNumberToOutputFormat(inputNumber, params[1])

	// calculations
	case strings.ContainsAny(input, "&|^><~"):
		result, ok := calculateValue(input)
		if ok != true {
			return "Could not apply bitwise operation!"
		}
		return strconv.Itoa(result)
	default:
		return "Unknown command!"
	}
}

func convertNumberToOutputFormat(inputNumber int, outputFormat string) string {
	switch outputFormat {
	case "dec":
		return strconv.Itoa(inputNumber)
	//add two-complement later
	case "bin":
		return fmt.Sprintf("bx%b", inputNumber)
	case "hex":
		return fmt.Sprintf("0x%X", inputNumber)

	default:
		return "Unknown conversion format!"
	}
}

func convertInputToInt(input string) (int, bool) {
	if input == "last" {
		return lastResult, true
	}
	//remove any (maybe) existing parenthesis
	input = strings.ReplaceAll(input, "(", "")
	input = strings.ReplaceAll(input, ")", "")

	//is a hexadecimal number
	if strings.HasPrefix(input, "0x") {
		number, err := strconv.ParseInt(strings.TrimPrefix(input, "0x"), 16, 64)
		return int(number), err == nil

	} else if strings.HasPrefix(input, "bx") {
		number, err := strconv.ParseInt(strings.TrimPrefix(input, "bx"), 2, 64)
		return int(number), err == nil

	} else if strings.HasPrefix(input, "~0x") {
		number, err := strconv.ParseInt(strings.TrimPrefix(input, "~0x"), 16, 64)
		return ^int(number), err == nil

	} else if strings.HasPrefix(input, "~bx") {
		number, err := strconv.ParseInt(strings.TrimPrefix(input, "~bx"), 2, 64)
		return ^int(number), err == nil

	} else if strings.HasPrefix(input, "~") {
		number, err := strconv.Atoi(input[1:])
		return ^number, err == nil
	} else {
		number, err := strconv.Atoi(input)
		return number, err == nil
	}
}

func calculateValue(input string) (int, bool) {
	input = strings.ReplaceAll(input, " ", "")
	expressions, ok := parseCalculationExpression(input)
	if ok != true {
		fmt.Println("Could not parse expression!")
		return 0, ok
	}

	//need so sort expression by their priority value
	sort.Slice(expressions, func(i int, j int) bool {
		return expressions[i].value > expressions[j].value
	})

	//fmt.Println("expressions: ", expressions)

	// now calculate and replace all expression step by step until we only have one number left
	for i, elem := range expressions {
		result, ok := calcExpression(elem.expression)
		if ok != true {
			fmt.Println("Could not calculate expression!")
			return 0, false
		}
		if i != len(expressions)-1 {
			expressions[i+1].expression = strings.Replace(expressions[i+1].expression, elem.expression, strconv.Itoa(result), -1)
			//fmt.Println("new expression: ", expressions[i+1])
		} else {
			lastResult = result
			return result, true
		}
	}
	return 0, false
}

func calcExpression(expression string) (int, bool) {
	operatorRegex := regexp.MustCompile("\\||&|\\^|>>|<<")
	operatorSplit := operatorRegex.Split(expression, -1)
	operatorSlices := operatorRegex.FindAllString(expression, -1)

	expressionsAsIntSlice := make([]int, len(operatorSplit))

	for i := range operatorSplit {
		value, ok := convertInputToInt(operatorSplit[i])
		if ok != true {
			fmt.Println("could not convert input to int")
			return 0, ok
		}
		expressionsAsIntSlice[i] = value
	}

	//dumbly parse from left to right
	for i, elem := range operatorSlices {
		switch elem {
		case "|":
			expressionsAsIntSlice[i] = expressionsAsIntSlice[i] | expressionsAsIntSlice[i+1]
			expressionsAsIntSlice[i+1] = expressionsAsIntSlice[i]
		case "&":
			expressionsAsIntSlice[i] = expressionsAsIntSlice[i] & expressionsAsIntSlice[i+1]
			expressionsAsIntSlice[i+1] = expressionsAsIntSlice[i]
		case ">>":
			expressionsAsIntSlice[i] = expressionsAsIntSlice[i] >> uint(expressionsAsIntSlice[i+1])
			expressionsAsIntSlice[i+1] = expressionsAsIntSlice[i]
		case "<<":
			expressionsAsIntSlice[i] = expressionsAsIntSlice[i] << uint(expressionsAsIntSlice[i+1])
			expressionsAsIntSlice[i+1] = expressionsAsIntSlice[i]
		case "^":
			expressionsAsIntSlice[i] = expressionsAsIntSlice[i] ^ expressionsAsIntSlice[i+1]
			expressionsAsIntSlice[i+1] = expressionsAsIntSlice[i]
		}
	}
	//return most "right" element in slice
	return expressionsAsIntSlice[len(expressionsAsIntSlice)-1], true
}

func parseCalculationExpression(input string) ([]expression, bool) {

	// add outer ( ) around every expression
	input = strings.Join([]string{"(", input, ")"}, "")

	expressionList := list.New()
	//the stack handles all starts of parenthesis
	parenthesisStack := stack.New()
	for i, char := range input {
		if char == '(' {
			parenthesisStack.Push(i)
		}
		if char == ')' {
			if parenthesisStack.Len() > 0 {
				expressionList.PushBack(expression{parenthesisStack.Len(), input[parenthesisStack.Peek().(int)+1 : i]})
				parenthesisStack.Pop()
			}
		}
	}

	//wrong amount of closing/ opening parenthesis
	if parenthesisStack.Len() != 0 {
		return nil, false
	}

	//parse list to slice
	counter := 0
	expressionSlice := make([]expression, expressionList.Len())
	for elem := expressionList.Front(); elem != nil; elem = elem.Next() {
		expressionSlice[counter] = elem.Value.(expression)
		counter++
	}
	return expressionSlice, true
}
