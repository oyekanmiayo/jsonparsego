package main

import (
	"reflect"
	"strings"
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
				{tokenType: LEFT_CURLY_BRACKET}, {tokenType: NAME_STRING},
				{tokenType: NAME_SEPARATOR}, {tokenType: VALUE_STRING},
				{tokenType: RIGHT_CURLY_BRACKET},
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

func TestParseTokens_Invalid(t *testing.T) {
	testCases := []struct {
		name         string
		tokenList    []Token
		expectedBool bool
		errStr       string
	}{}
	for _, tc := range testCases {
		actualBool, err := parseTokens(tc.tokenList)
		if !reflect.DeepEqual(actualBool, tc.expectedBool) {
			t.Errorf("expected: %v, got: %v.", tc.expectedBool, actualBool)
		}

		if !strings.Contains(err.Error(), tc.errStr) {
			t.Errorf("expected: %v, got: %v.", tc.errStr, err.Error())
		}
	}
}
