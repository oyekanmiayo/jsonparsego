package main

import "fmt"

func parseTokens(TokenList []Token) (bool, error) {
	stack := NewStack()
	for _, t := range TokenList {
		switch t.tokenType {
		case LEFT_CURLY_BRACKET:
			// we should check if stack is empty or the top is a colon
			if stack.IsEmpty() || stack.Peek() == NAME_SEPARATOR {
				stack.Push(LEFT_CURLY_BRACKET)
			} else {
				return false, fmt.Errorf("{ can only come after a colon or at the beginning of an empty json file")
			}
		case RIGHT_CURLY_BRACKET:
			// This is for a case where we have just one name:value pair and value is a string
			// when we encounter the right curly brace, the value string will still be in the stack
			if stack.Peek() == VALUE_STRING || stack.Peek() == LITERAL || stack.Peek() == NUMBER {
				stack.Pop()
			}
			if stack.Peek() == LEFT_CURLY_BRACKET {
				stack.Pop()
			} else {
				return false, fmt.Errorf("there's no equivalent { for the }")
			}
		case VALUE_STRING:
			// This can happen if a name string is coming after a name:value pair (which follows
			// with a comma).
			if stack.Peek() == VALUE_SEPARATOR {
				stack.Pop()
			}
			stack.Push(VALUE_STRING)
		case LITERAL:
			// This can happen if a name string is coming after a name:value pair (which follows
			// with a comma).
			if stack.Peek() == VALUE_SEPARATOR {
				stack.Pop()
			} else {
				return false, fmt.Errorf("a literal must come after a comma (value separator)")
			}
			stack.Push(LITERAL)
		case NUMBER:
			if stack.Peek() == VALUE_SEPARATOR {
				stack.Pop()
			} else {
				return false, fmt.Errorf("a number must come after a comma (value separator)")
			}
			stack.Push(NUMBER)
		case NAME_STRING:
			stack.Push(NAME_STRING)
		case NAME_SEPARATOR:
			if stack.Peek() == NAME_STRING {
				stack.Pop()
			} else {
				return false, fmt.Errorf("a colon (name separator) must come after a name string")
			}
		case VALUE_SEPARATOR:
			if stack.Peek() == VALUE_STRING || stack.Peek() == LITERAL || stack.Peek() == NUMBER {
				stack.Pop()
			} else {
				return false, fmt.Errorf("a comma (value separator) must come after a value string or literal or number")
			}
			stack.Push(VALUE_SEPARATOR)
		default:
			return false, fmt.Errorf("unhandled token")
		}
	}

	if stack.IsEmpty() {
		return true, nil
	}

	return false, fmt.Errorf("unknown error")
}
