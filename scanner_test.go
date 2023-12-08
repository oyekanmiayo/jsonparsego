package main

import (
	"reflect"
	"testing"
)

func TestScanTokens(t *testing.T) {
	testCases := []struct {
		name              string
		data              []byte
		expectedTokenList []Token
	}{
		{
			name: "Valid json with only strings",
			data: []byte(`{"name":"ayo"}`),
			expectedTokenList: []Token{
				LEFT_CURLY_BRACKET, NAME_STRING, NAME_SEPARATOR, VALUE_STRING, RIGHT_CURLY_BRACKET,
			},
		},
		{
			name: "Valid json with numbers",
			data: []byte(`{"name":123}`),
			expectedTokenList: []Token{
				LEFT_CURLY_BRACKET, NAME_STRING, NAME_SEPARATOR, NUMBER, RIGHT_CURLY_BRACKET,
			},
		},
		{
			name: "Valid json with literals",
			data: []byte(`{"tbool":true, "fbool":false, "null":null }`),
			expectedTokenList: []Token{
				LEFT_CURLY_BRACKET, NAME_STRING, NAME_SEPARATOR, LITERAL, VALUE_SEPARATOR,
				NAME_STRING, NAME_SEPARATOR, LITERAL, VALUE_SEPARATOR, NAME_STRING, NAME_SEPARATOR,
				LITERAL, RIGHT_CURLY_BRACKET,
			},
		},
	}

	for _, tc := range testCases {
		actualTokenList := scanTokens(tc.data)
		if !reflect.DeepEqual(actualTokenList, tc.expectedTokenList) {
			t.Errorf("scanTokens expected %v, got %v", tc.expectedTokenList, actualTokenList)
		}
	}
}
