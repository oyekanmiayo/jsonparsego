package main

import (
	"fmt"
	"os"
)

type Token int

const (
	ILLEGAL Token = iota

	LEFT_CURLY_BRACKET  // '{'
	RIGHT_CURLY_BRACKET // '}'

	LEFT_SQUARE_BRACKET  // '['
	RIGHT_SQUARE_BRACKET // ']'

	NAME_SEPARATOR  // ':'
	VALUE_SEPARATOR // ','

	SPACE        // ' '
	NAME_STRING  // for field titles which are always strings as in name:value
	VALUE_STRING // "string"

	LITERAL // false, true or null
)

func scanTokens(data []byte) []Token {

	var TokenList []Token

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
			TokenList = append(TokenList, LEFT_CURLY_BRACKET)
		case '}':
			TokenList = append(TokenList, RIGHT_CURLY_BRACKET)
		// This assumes the beginning of a name or value string
		case '"':
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

			if len(TokenList) == 0 {
				panic("Illegal string")
			}

			if TokenList[len(TokenList)-1] == LEFT_CURLY_BRACKET || TokenList[len(TokenList)-1] == VALUE_SEPARATOR {
				TokenList = append(TokenList, NAME_STRING)
			} else if TokenList[len(TokenList)-1] == NAME_SEPARATOR {
				TokenList = append(TokenList, VALUE_STRING)
			}

		case ':':
			TokenList = append(TokenList, NAME_SEPARATOR)
		case ',':
			TokenList = append(TokenList, VALUE_SEPARATOR)
		case 'f':
			if TokenList[len(TokenList)-1] != NAME_SEPARATOR {
				panic("Invalid json")
			}

			i := byteIdx
			if data[i+1] == 'a' && data[i+2] == 'l' && data[i+3] == 's' && data[i+4] == 'e' {
				// Found false
				TokenList = append(TokenList, LITERAL)
			}
			byteIdx = i + 5
		case 't':
			if TokenList[len(TokenList)-1] != NAME_SEPARATOR {
				panic("Invalid json")
			}

			i := byteIdx
			if data[i+1] == 'r' && data[i+2] == 'u' && data[i+3] == 'e' {
				// Found true
				TokenList = append(TokenList, LITERAL)
			}
			byteIdx = i + 4
		case 'n':
			if TokenList[len(TokenList)-1] != NAME_SEPARATOR {
				panic("Invalid json")
			}

			i := byteIdx
			if data[i+1] == 'u' && data[i+2] == 'l' && data[i+3] == 'l' {
				// Found true
				TokenList = append(TokenList, LITERAL)
			}
			byteIdx = i + 4
		default:
			fmt.Println(string(b))
			os.Exit(1)
		}
	}

	return TokenList
}
