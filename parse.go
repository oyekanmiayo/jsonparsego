package main

func parseTokens(TokenList []Token) bool {
	stack := NewStack()
	for _, t := range TokenList {
		switch t {
		case LEFT_CURLY_BRACKET:
			stack.Push(LEFT_CURLY_BRACKET)
		case RIGHT_CURLY_BRACKET:
			// This is for a case where we have just one name:value pair and value is a string
			// when we encounter the right curly brace, the value string will still be in the stack
			if stack.Peek() == VALUE_STRING || stack.Peek() == LITERAL {
				stack.Pop()
			}
			if stack.Peek() == LEFT_CURLY_BRACKET {
				stack.Pop()
			} else {
				return false
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
			}
			stack.Push(LITERAL)
		case NAME_STRING:
			stack.Push(NAME_STRING)
		case NAME_SEPARATOR:
			if stack.Peek() == NAME_STRING {
				stack.Pop()
			} else {
				return false
			}
		case VALUE_SEPARATOR:
			if stack.Peek() == VALUE_STRING || stack.Peek() == LITERAL {
				stack.Pop()
			} else {
				return false
			}
			stack.Push(VALUE_SEPARATOR)
		default:
			panic("unhandled default case")
		}

	}

	if stack.IsEmpty() {
		return true
	}

	return false
}
