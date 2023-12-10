package main

import (
	"fmt"
)

type Token struct {
	tokenType  TokenType
	tokenState TokenState
	tokenValue interface{}
	tokenDepth int
}

type TokenState int

const (
	INVALID TokenState = iota
	WITHIN_OBJECT
	WITHIN_ARRAY
)

type TokenType int

const (
	ILLEGAL TokenType = iota

	LEFT_CURLY_BRACKET  // '{'
	RIGHT_CURLY_BRACKET // '}'

	LEFT_SQUARE_BRACKET  // '['
	RIGHT_SQUARE_BRACKET // ']'

	NAME_SEPARATOR  // ':'
	VALUE_SEPARATOR // ','

	NAME_STRING  // for field titles which are always strings as in name:value
	VALUE_STRING // "string"

	LITERAL // false, true or null
	NUMBER
)

// (TODO) Write a post about this. It's a slice that works as a map. Why?
var tokens = [...]string{
	ILLEGAL:             "ILLEGAL",
	LEFT_CURLY_BRACKET:  "LEFT_CURLY_BRACKET",
	RIGHT_CURLY_BRACKET: "RIGHT_CURLY_BRACKET",

	LEFT_SQUARE_BRACKET:  "LEFT_SQUARE_BRACKET",
	RIGHT_SQUARE_BRACKET: "RIGHT_SQUARE_BRACKET",

	NAME_SEPARATOR:  "NAME_SEPARATOR",
	VALUE_SEPARATOR: "VALUE_SEPARATOR",

	NAME_STRING:  "NAME_STRING",
	VALUE_STRING: "VALUE_STRING",

	LITERAL: "LITERAL",
	NUMBER:  "NUMBER",
}

func scanTokens(data []byte) ([]Token, error) {

	var TokenList []Token
	tokenStateStack := NewStack()

	byteIdx := 0
	for ; byteIdx < len(data); byteIdx++ {
		b := data[byteIdx]

		switch b {
		case ' ':
		case '\n':
		case '\t':
		case '\r':
			continue
		case '{':
			TokenList = append(TokenList, Token{
				tokenType: LEFT_CURLY_BRACKET, tokenState: getTokenState(tokenStateStack),
			})
			tokenStateStack.Push(WITHIN_OBJECT)
		case '}':
			if tokenStateStack.Pop() != WITHIN_OBJECT {
				return []Token{}, fmt.Errorf("not WITHIN_OBJECT")
			}
			TokenList = append(TokenList, Token{
				tokenType: RIGHT_CURLY_BRACKET, tokenState: getTokenState(tokenStateStack),
			})
		case '[':
			TokenList = append(TokenList, Token{tokenType: LEFT_SQUARE_BRACKET})
			tokenStateStack.Push(WITHIN_ARRAY)
		case ']':
			if tokenStateStack.IsEmpty() || tokenStateStack.Pop() != WITHIN_ARRAY {
				return []Token{}, fmt.Errorf("not storing token states well (WITHIN_ARRAY) or invalid syntax")
			}
			TokenList = append(TokenList, Token{
				tokenType: RIGHT_SQUARE_BRACKET, tokenState: getTokenState(tokenStateStack),
			})
		// This assumes the beginning of a name or value string
		case '"':
			if len(TokenList) == 0 {
				return []Token{}, fmt.Errorf("json file is empty, so this is an illegal string")
			}

			i := byteIdx + 1
			for ; i < len(data); i++ {
				// If we find an escapes quotation \" continue
				if data[i-1] == '\\' {

					if data[i] == '"' || data[i] == '\\' || data[i] == '/' || data[i] == 'b' || data[i] == 'f' || data[i] == 'n' || data[i] == 'r' || data[i] == 't' {
						continue
					} else {
						// return an error
						return []Token{}, fmt.Errorf("illegal escaped character: %v\n", string(data[i]))
					}

				}

				if data[i] == '"' {
					// s := string(data[byteIdx+1 : i])
					// fmt.Println(s)
					break
				}
			}
			byteIdx = i

			// strToken := Token{}

			if tokenStateStack.IsEmpty() {
				return []Token{}, fmt.Errorf("string token isn't WITHIN_OBJECT or WITHIN_ARRAY")
			}

			if TokenList[len(TokenList)-1].tokenType == LEFT_CURLY_BRACKET {
				TokenList = append(TokenList, Token{
					tokenType: NAME_STRING, tokenState: getTokenState(tokenStateStack),
				})
			} else if tokenStateStack.Peek() == WITHIN_OBJECT && TokenList[len(TokenList)-1].tokenType == VALUE_SEPARATOR {
				TokenList = append(TokenList, Token{
					tokenType: NAME_STRING, tokenState: getTokenState(tokenStateStack),
				})
			} else if tokenStateStack.Peek() == WITHIN_ARRAY && TokenList[len(TokenList)-1].tokenType == VALUE_SEPARATOR {
				TokenList = append(TokenList, Token{
					tokenType: VALUE_STRING, tokenState: getTokenState(tokenStateStack),
				})
			} else if TokenList[len(TokenList)-1].tokenType == NAME_SEPARATOR {
				TokenList = append(TokenList, Token{
					tokenType: VALUE_STRING, tokenState: getTokenState(tokenStateStack),
				})
			} else if tokenStateStack.Peek() == WITHIN_ARRAY && TokenList[len(TokenList)-1].tokenType == LEFT_SQUARE_BRACKET {
				TokenList = append(TokenList, Token{
					tokenType: VALUE_STRING, tokenState: getTokenState(tokenStateStack),
				})
			}
		case ':':
			TokenList = append(TokenList, Token{
				tokenType: NAME_SEPARATOR, tokenState: getTokenState(tokenStateStack),
			})
		case ',':
			TokenList = append(TokenList, Token{
				tokenType: VALUE_SEPARATOR, tokenState: getTokenState(tokenStateStack),
			})
		case 'f':
			// For this case to remain valid, the case for quotation MUST always come before.

			// if TokenList[len(TokenList)-1].tokenType != NAME_SEPARATOR {
			// 	// Because we check for full strings earlier, the appearance of an f here means it has no quotation around it
			// 	// and the only two reasons for that are for false, or some invalid value. In any case, a colon must come before
			// 	// it.
			// 	return []Token{}, fmt.Errorf("there's a character f before a colon which means this is an invalid json")
			// }

			i := byteIdx
			if data[i+1] == 'a' && data[i+2] == 'l' && data[i+3] == 's' && data[i+4] == 'e' {
				// Found false
				TokenList = append(TokenList, Token{
					tokenType: LITERAL, tokenState: getTokenState(tokenStateStack),
				})
			} else {
				// Test
				return []Token{}, fmt.Errorf("there's an f in the json that isn't in a string and isn't the false literal")
			}
			byteIdx = i + 4
		case 't':
			// if TokenList[len(TokenList)-1].tokenType != NAME_SEPARATOR {
			// 	// Because we check for full strings earlier, the appearance of a t here means it has no quotation around it
			// 	// and the only two reasons for that are for true, or some invalid value. In any case, a colon must come before
			// 	// it.
			// 	return []Token{}, fmt.Errorf("there's a character t before a colon which means this is an invalid json")
			// }

			i := byteIdx
			if data[i+1] == 'r' && data[i+2] == 'u' && data[i+3] == 'e' {
				// Found true
				TokenList = append(TokenList, Token{
					tokenType: LITERAL, tokenState: getTokenState(tokenStateStack),
				})
			} else {
				// Test
				return []Token{}, fmt.Errorf("there's a t in the json that isn't in a string and isn't the true literal")
			}
			byteIdx = i + 3
		case 'n':
			// if TokenList[len(TokenList)-1].tokenType != NAME_SEPARATOR {
			// 	// Because we check for full strings earlier, the appearance of an n here means it has no quotation around it
			// 	// and the only two reasons for that are for null, or some invalid value. In any case, a colon must come before
			// 	// it.
			// 	return []Token{}, fmt.Errorf("there's a character n before a colon which means this is an invalid json")
			// }

			i := byteIdx
			if data[i+1] == 'u' && data[i+2] == 'l' && data[i+3] == 'l' {
				// Found true
				TokenList = append(TokenList, Token{
					tokenType: LITERAL, tokenState: getTokenState(tokenStateStack),
				})
			} else {
				// Test
				return []Token{}, fmt.Errorf("there's an n in the json that isn't in a string and isn't the null literal")
			}
			byteIdx = i + 3
		case '-':
			if byteIdx+1 == len(data) {
				// return error. nothing after fraction
				return []Token{}, fmt.Errorf("no character after the minus sym (+)")
			}
			byteIdx += 1
			b = data[byteIdx]
			fallthrough
		case '+':
			if byteIdx+1 == len(data) {
				// return error. nothing after fraction
				return []Token{}, fmt.Errorf("no character after the plus sym (+)")
			}
			byteIdx += 1
			b = data[byteIdx]
			fallthrough
		default:
			if IsDigit(b) {
				// if b == 0, it means number is starting with 0 and that's not allowed
				if b == '0' {
					return []Token{}, fmt.Errorf("a number can't start with 0")
				}

				i := byteIdx + 1
				for ; i < len(data); i++ {
					if IsDigit(data[i]) {
						continue
					} else if data[i] == '.' { // fraction
						if i+1 == len(data) {
							// return error. nothing after fraction
							return []Token{}, fmt.Errorf("no character after the fraction sym (.)")
						}

						if !IsDigit(data[i+1]) {
							// return error. what's after fraction isn't a number
							return []Token{}, fmt.Errorf("character after the fraction sym (.) isn't a number")

						}
						continue
					} else if data[i] == 'e' || data[i] == 'E' { // exponent
						if i+1 == len(data) {
							// return error. nothing after fraction
							return []Token{}, fmt.Errorf("no character after the exponent sym (e|E)")
						}

						if !IsDigit(data[i+1]) && data[i+1] != '-' && data[i+1] != '+' {
							// return error. what's after fraction isn't a number
							return []Token{}, fmt.Errorf("character after the exponent sym (e|E) isn't a number, minus or plsu")
						}
						continue
					} else if data[i] == '-' || data[i] == '+' {
						if i+1 == len(data) {
							// return error. nothing after fraction
							return []Token{}, fmt.Errorf("no character after the plus or minus sym (+|-)")
						}

						if IsDigit(data[i-1]) {
							// there's an issue
							// consider: +123, 1e-2
							// in either case, there must not be a number before + or -
							return []Token{}, fmt.Errorf("character before the plus or minus sym (+|-) IS a number and it shouldn't be")
						}

						if !IsDigit(data[i+1]) {
							// return error. what's after fraction isn't a number
							return []Token{}, fmt.Errorf("character after the plus or minus sym (+|-) isn't a number")
						}
						continue
					} else {
						break
					}
				}
				fmt.Println(string(data[byteIdx:i]))
				byteIdx = i - 1
				TokenList = append(TokenList, Token{
					tokenType: NUMBER, tokenState: getTokenState(tokenStateStack),
				})
			} else {
				fmt.Println(string(b))
				return []Token{}, fmt.Errorf("invalid token: %c", b)
			}
		}
	}

	if len(TokenList) == 0 {
		return []Token{}, fmt.Errorf("json file is empty")
	}

	return TokenList, nil
}

func IsDigit(b byte) bool {
	return '0' <= b && b <= '9'
}

// This is only valid for non-structural tokens that get added to the stack
// So invalid for {,},[,]
func getTokenState(tokenStateStack *Stack) TokenState {
	if tokenStateStack.IsEmpty() {
		return INVALID
	}

	tokenState, ok := tokenStateStack.Peek().(TokenState)
	if !ok {
		panic("Stored invalid value in tokenStateStack")
	}
	return tokenState
}
