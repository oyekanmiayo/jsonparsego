package main

import (
	"fmt"
)

func parseTokensV2(TokenList []Token) (bool, error) {
	// This stack should only contain {, }, [ or ] at any point
	stack := NewStack()
	for i, t := range TokenList {
		switch t.tokenType {
		case LEFT_CURLY_BRACKET:
			if i == 0 {
				stack.Push(LEFT_CURLY_BRACKET)
				continue
			}

			prevTkn := TokenList[i-1].tokenType
			if prevTkn != NAME_SEPARATOR && (t.tokenState == WITHIN_ARRAY && prevTkn != VALUE_SEPARATOR && prevTkn != LEFT_SQUARE_BRACKET) {
				return false, fmt.Errorf("NAME_SEPARATOR should preceed LEFT_CURLY_BRACKET if this is a nested object or else nothing should preceed it, instead got: %v\n", tokens[prevTkn])
			}
			stack.Push(LEFT_CURLY_BRACKET)
		case RIGHT_CURLY_BRACKET:
			prevTkn := TokenList[i-1].tokenType
			if prevTkn != VALUE_STRING && prevTkn != NUMBER && prevTkn != LITERAL && prevTkn != RIGHT_CURLY_BRACKET && prevTkn != RIGHT_SQUARE_BRACKET {
				return false, fmt.Errorf("VALUE_STRING or RIGHT_CURLY_BRACKET or LITERAL or RIGHT_SQUARE_BRACKET should preceed RIGHT_CURLY_BRACKET, instead got: %v\n", tokens[prevTkn])
			}
			stack.Pop()
		case LEFT_SQUARE_BRACKET:
			prevTkn := TokenList[i-1].tokenType
			if prevTkn != NAME_SEPARATOR {
				return false, fmt.Errorf("NAME_SEPARATOR should preceed LEFT_SQUARE_BRACKET, instead got: %v\n", tokens[prevTkn])
			}
			stack.Push(LEFT_SQUARE_BRACKET)
		case RIGHT_SQUARE_BRACKET:
			prevTkn := TokenList[i-1].tokenType
			if prevTkn != VALUE_STRING && prevTkn != NUMBER && prevTkn != LITERAL && prevTkn != RIGHT_CURLY_BRACKET {
				return false, fmt.Errorf("VALUE_STRING or RIGHT_CURLY_BRACKET or LITERAL should preceed RIGHT_SQUARE_BRACKET, instead got: %v\n", tokens[prevTkn])
			}
			stack.Pop()
		case NAME_STRING:
			prevTkn := TokenList[i-1].tokenType
			if prevTkn != LEFT_CURLY_BRACKET && prevTkn != VALUE_SEPARATOR {
				return false, fmt.Errorf("LEFT_CURLY_BRACKET or VALUE_SEPARATOR should preceed NAME_STRING, instead got: %v\n", tokens[prevTkn])
			}
		case NAME_SEPARATOR:
			prevTkn := TokenList[i-1].tokenType
			if prevTkn != NAME_STRING {
				return false, fmt.Errorf("NAME_STRING should preceed NAME_SEPARATOR, instead got: %v\n", tokens[prevTkn])
			}
		case VALUE_STRING:
			prevTkn := TokenList[i-1].tokenType
			if prevTkn != NAME_SEPARATOR && (t.tokenState == WITHIN_ARRAY && prevTkn != LEFT_SQUARE_BRACKET && prevTkn != VALUE_SEPARATOR) {
				return false, fmt.Errorf("NAME_SEPARATOR or LEFT_SQUARE_BRACKET (when within an array) or VALUE_SEPARATOR (when within an array) should preceed VALUE_STRING, instead got: %v\n", tokens[prevTkn])
			}
		case VALUE_SEPARATOR:
			prevTkn := TokenList[i-1].tokenType
			prevTknState := TokenList[i-1].tokenState
			if prevTknState == INVALID {
				return false, fmt.Errorf("VALUE_SEPARATOR should not come after LEFT_CURLY_BRACKET or  RIGHT_CURLY_BRACKET (when the object isn't nested, got: %v\n", tokens[prevTkn])
			}

			if prevTkn != RIGHT_CURLY_BRACKET && prevTkn != VALUE_STRING && prevTkn != NUMBER && prevTkn != LITERAL {
				return false, fmt.Errorf("RIGHT_CURLY_BRACKET or VALUE_STRING or NUMBER or LITERAL must precede VALUE_SEPARATOR, got: %v\n", tokens[prevTkn])
			}
		case NUMBER:
			prevTkn := TokenList[i-1].tokenType
			if prevTkn != NAME_SEPARATOR && (t.tokenState == WITHIN_ARRAY && prevTkn != VALUE_SEPARATOR && prevTkn != LEFT_SQUARE_BRACKET) {
				return false, fmt.Errorf("NAME_SEPARATOR or VALUE_SEPARATOR (within an array) or LEFT_SQUARE_BRACKET (within an array) should preceed NUMBER, instead got: %v\n", tokens[prevTkn])
			}
		case LITERAL:
			prevTkn := TokenList[i-1].tokenType
			if prevTkn != NAME_SEPARATOR {
				return false, fmt.Errorf("NAME_SEPARATOR should preceed LITERAL, instead got: %v\n", tokens[prevTkn])
			}
		default:
			return false, fmt.Errorf("invalid token: %v\n", tokens[t.tokenType])
		}
	}

	if stack.IsEmpty() {
		return true, nil
	}

	return false, fmt.Errorf("unknown error")
}
