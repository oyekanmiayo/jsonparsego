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

	SPACE  // ' '
	STRING // "string"
)

func scanTokens(data []byte) []Token {

	var TokenList []Token

	byteIdx := 0
	for ; byteIdx < len(data); byteIdx++ {
		b := data[byteIdx]

		switch b {
		case ' ':
		case '\n':
			continue
		case '{':
			TokenList = append(TokenList, LEFT_CURLY_BRACKET)
		case '}':
			TokenList = append(TokenList, RIGHT_CURLY_BRACKET)
		case '"':
			i := byteIdx + 1
			for ; i < len(data); i++ {
				if data[i] == '"' {
					s := string(data[byteIdx+1 : i])
					fmt.Println(s)
					break
				}
			}
			byteIdx = i
			TokenList = append(TokenList, STRING)
		case ':':
			TokenList = append(TokenList, NAME_SEPARATOR)
		case ',':
			TokenList = append(TokenList, VALUE_SEPARATOR)
		default:
			fmt.Println(string(b))
			os.Exit(1)
		}
	}

	return TokenList
}
