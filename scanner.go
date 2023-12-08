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
				panic("Not storing token states well (WITHIN_OBJECT)")
			}
			TokenList = append(TokenList, Token{
				tokenType: RIGHT_CURLY_BRACKET, tokenState: getTokenState(tokenStateStack),
			})
		case '[':
			TokenList = append(TokenList, Token{tokenType: LEFT_SQUARE_BRACKET})
			tokenStateStack.Push(WITHIN_ARRAY)
		case ']':
			if tokenStateStack.Pop() != WITHIN_ARRAY {
				panic("Not storing token states well (WITHIN_ARRAY)")
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
				if data[i] == '"' && data[i-1] == '\\' {
					continue
				}

				if data[i] == '"' {
					s := string(data[byteIdx+1 : i])
					fmt.Println(s)
					break
				}
			}
			byteIdx = i

			// strToken := Token{}

			if TokenList[len(TokenList)-1].tokenType == LEFT_CURLY_BRACKET || TokenList[len(TokenList)-1].tokenType == VALUE_SEPARATOR {
				TokenList = append(TokenList, Token{
					tokenType: NAME_STRING, tokenState: getTokenState(tokenStateStack),
				})
			} else if TokenList[len(TokenList)-1].tokenType == NAME_SEPARATOR {
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
		default:
			if IsDigit(b) {
				i := byteIdx + 1
				for ; i < len(data); i++ {
					if IsDigit(data[i]) {
						continue
					} else if data[i] == '.' && IsDigit(data[i+1]) { // fraction
						continue
					} else if (data[i] == 'e' || data[i] == 'E') && IsDigit(data[i+1]) { // exponent
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
