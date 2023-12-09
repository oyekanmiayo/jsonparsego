package main

import "fmt"

func parseTokens(TokenList []Token) (bool, error) {
	stack := NewStack()
	for i, t := range TokenList {
		i = i
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
			// Add one for }
			if stack.Peek() == VALUE_STRING || stack.Peek() == LITERAL || stack.Peek() == NUMBER {
				stack.Pop()
			} else if stack.Size() > 1 && stack.Peek() == NAME_SEPARATOR {
				// Valid for nested objects
				stack.Pop()
			}

			if stack.Peek() == LEFT_CURLY_BRACKET {
				stack.Pop()
			} else {
				return false, fmt.Errorf("there's no equivalent { for the }")
			}

			// this assumes that we have a nested
			if stack.Size() > 1 {
				stack.Push(RIGHT_CURLY_BRACKET)
			}
		case LEFT_SQUARE_BRACKET:
			if stack.IsEmpty() || (stack.Peek() != NAME_SEPARATOR) {
				return false, fmt.Errorf("[ can only come after a colon (:)")
			}

			if stack.Peek() != NAME_SEPARATOR {
				stack.Pop()
			}
			stack.Push(LEFT_SQUARE_BRACKET)
		case RIGHT_SQUARE_BRACKET:

			// Add one for }
			if stack.Peek() == VALUE_STRING || stack.Peek() == LITERAL || stack.Peek() == NUMBER {
				stack.Pop()
			}

			if stack.IsEmpty() || (stack.Peek() != LEFT_SQUARE_BRACKET) {
				return false, fmt.Errorf("] can only come after [")
			}

			if stack.Peek() == LEFT_SQUARE_BRACKET {
				stack.Pop()
			}

			stack.Push(RIGHT_SQUARE_BRACKET)
		case NAME_STRING:
			if stack.IsEmpty() || (stack.Peek() != VALUE_SEPARATOR && stack.Peek() != LEFT_CURLY_BRACKET) {
				return false, fmt.Errorf("a field name MUST come before a comma or a left curly bracket")
			}

			if stack.Peek() == VALUE_SEPARATOR {
				stack.Pop()
			}
			stack.Push(NAME_STRING)
		case VALUE_STRING:
			// This can happen if a name string is coming after a name:value pair (which follows
			// with a comma).
			if stack.Peek() == NAME_SEPARATOR {
				stack.Pop()
			} else if t.tokenState == WITHIN_ARRAY && (stack.Peek() == LEFT_SQUARE_BRACKET || stack.Peek() == VALUE_SEPARATOR) {
				if stack.Peek() == VALUE_SEPARATOR {
					stack.Pop()
				}
			} else {
				// error
			}
			stack.Push(VALUE_STRING)
		case LITERAL:
			// This can happen if a name string is coming after a name:value pair (which follows
			// with a comma).
			if stack.Peek() == NAME_SEPARATOR {
				stack.Pop()
			} else if t.tokenState == WITHIN_ARRAY && (stack.Peek() == LEFT_SQUARE_BRACKET || stack.Peek() == VALUE_SEPARATOR) {
				if stack.Peek() == VALUE_SEPARATOR {
					stack.Pop()
				}
			} else {
				return false, fmt.Errorf("a literal must come after a colon, a comma or a [. The last two only apply when within an array")
			}
			stack.Push(LITERAL)
		case NUMBER:
			if stack.Peek() == NAME_SEPARATOR {
				stack.Pop()
			} else if t.tokenState == WITHIN_ARRAY && (stack.Peek() == LEFT_SQUARE_BRACKET || stack.Peek() == VALUE_SEPARATOR) {
				if stack.Peek() == VALUE_SEPARATOR {
					stack.Pop()
				}
			} else {
				return false, fmt.Errorf("a number must come after a colon, a comma or a [. The last two only apply when within an array")
			}
			stack.Push(NUMBER)

		case NAME_SEPARATOR:
			if stack.Peek() == NAME_STRING {
				stack.Pop()
			} else {
				return false, fmt.Errorf("a colon (name separator) must come after a name string")
			}
			stack.Push(NAME_SEPARATOR)
		case VALUE_SEPARATOR:
			if stack.Peek() == VALUE_STRING || stack.Peek() == LITERAL || stack.Peek() == NUMBER || stack.Peek() == RIGHT_CURLY_BRACKET || stack.Peek() == RIGHT_SQUARE_BRACKET {
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
