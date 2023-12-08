package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestScanTokens_Valid(t *testing.T) {
	testCases := []struct {
		name              string
		data              []byte
		expectedTokenList []Token
	}{
		{
			name: "Valid json with only strings",
			data: []byte(`{"name":"ayo"}`),
			expectedTokenList: []Token{
				{tokenType: LEFT_CURLY_BRACKET}, {tokenType: NAME_STRING},
				{tokenType: NAME_SEPARATOR}, {tokenType: VALUE_STRING},
				{tokenType: RIGHT_CURLY_BRACKET},
			},
		},
		{
			name: "Valid json with numbers",
			data: []byte(`{"name":123}`),
			expectedTokenList: []Token{
				{tokenType: LEFT_CURLY_BRACKET}, {tokenType: NAME_STRING},
				{tokenType: NAME_SEPARATOR}, {tokenType: NUMBER}, {tokenType: RIGHT_CURLY_BRACKET},
			},
		},
		{
			name: "Valid json with literals",
			data: []byte(`{"tbool":true, "fbool":false, "null":null }`),
			expectedTokenList: []Token{
				{tokenType: LEFT_CURLY_BRACKET}, {tokenType: NAME_STRING},
				{tokenType: NAME_SEPARATOR}, {tokenType: LITERAL}, {tokenType: VALUE_SEPARATOR},
				{tokenType: NAME_STRING}, {tokenType: NAME_SEPARATOR}, {tokenType: LITERAL},
				{tokenType: VALUE_SEPARATOR}, {tokenType: NAME_STRING}, {tokenType: NAME_SEPARATOR},
				{tokenType: LITERAL}, {tokenType: RIGHT_CURLY_BRACKET},
			},
		},
	}

	for _, tc := range testCases {
		actualTokenList, err := scanTokens(tc.data)
		if !reflect.DeepEqual(actualTokenList, tc.expectedTokenList) {
			t.Errorf("scanTokens expected: %v, got: %v. details: %v", tc.expectedTokenList, actualTokenList, err)
		}
	}
}

func TestScanTokens_Invalid(t *testing.T) {
	testCases := []struct {
		name              string
		data              []byte
		expectedTokenList []Token
		errStr            string
	}{
		{
			name:              "Invalid json with wrongly placed strings",
			data:              []byte(`"name":"ayo"}`),
			expectedTokenList: []Token{},
			errStr:            "json file is empty, so this is an illegal string",
		},
	}
	for _, tc := range testCases {
		actualTokenList, err := scanTokens(tc.data)
		if !reflect.DeepEqual(actualTokenList, tc.expectedTokenList) {
			t.Errorf("scanTokens expected: %v, got: %v. expected an error.", tc.expectedTokenList, actualTokenList)
		}

		if !strings.Contains(err.Error(), tc.errStr) {
			t.Errorf("expected: %v, got: %v.", tc.errStr, err.Error())
		}

	}
}
