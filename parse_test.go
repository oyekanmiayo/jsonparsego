package main

import (
	"reflect"
	"testing"
)

func TestParseTokens_Valid(t *testing.T) {
	testCases := []struct {
		name         string
		tokenList    []Token
		expectedBool bool
	}{
		{
			name: "Valid json with one line",
			tokenList: []Token{
				LEFT_CURLY_BRACKET, NAME_STRING, NAME_SEPARATOR, VALUE_STRING, RIGHT_CURLY_BRACKET,
			},
			expectedBool: true,
		},
	}

	for _, tc := range testCases {
		actualBool, err := parseTokens(tc.tokenList)
		if !reflect.DeepEqual(actualBool, tc.expectedBool) {
			t.Errorf("expected: %v, got: %v. details: %v", tc.expectedBool, actualBool, err)
		}
	}
}
