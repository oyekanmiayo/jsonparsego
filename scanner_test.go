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
				{tokenType: LEFT_CURLY_BRACKET, tokenState: INVALID},
				{tokenType: NAME_STRING, tokenState: WITHIN_OBJECT},
				{tokenType: NAME_SEPARATOR, tokenState: WITHIN_OBJECT},
				{tokenType: VALUE_STRING, tokenState: WITHIN_OBJECT},
				{tokenType: RIGHT_CURLY_BRACKET, tokenState: INVALID},
			},
		},
		{
			name: "Valid json with numbers",
			data: []byte(`{"name":123}`),
			expectedTokenList: []Token{
				{tokenType: LEFT_CURLY_BRACKET, tokenState: INVALID},
				{tokenType: NAME_STRING, tokenState: WITHIN_OBJECT},
				{tokenType: NAME_SEPARATOR, tokenState: WITHIN_OBJECT},
				{tokenType: NUMBER, tokenState: WITHIN_OBJECT},
				{tokenType: RIGHT_CURLY_BRACKET, tokenState: INVALID},
			},
		},
		{
			name: "Valid json with literals",
			data: []byte(`{"tbool":true, "fbool":false, "null":null }`),
			expectedTokenList: []Token{
				{tokenType: LEFT_CURLY_BRACKET, tokenState: INVALID},
				{tokenType: NAME_STRING, tokenState: WITHIN_OBJECT},
				{tokenType: NAME_SEPARATOR, tokenState: WITHIN_OBJECT},
				{tokenType: LITERAL, tokenState: WITHIN_OBJECT},
				{tokenType: VALUE_SEPARATOR, tokenState: WITHIN_OBJECT},
				{tokenType: NAME_STRING, tokenState: WITHIN_OBJECT},
				{tokenType: NAME_SEPARATOR, tokenState: WITHIN_OBJECT},
				{tokenType: LITERAL, tokenState: WITHIN_OBJECT},
				{tokenType: VALUE_SEPARATOR, tokenState: WITHIN_OBJECT},
				{tokenType: NAME_STRING, tokenState: WITHIN_OBJECT},
				{tokenType: NAME_SEPARATOR, tokenState: WITHIN_OBJECT},
				{tokenType: LITERAL, tokenState: WITHIN_OBJECT},
				{tokenType: RIGHT_CURLY_BRACKET, tokenState: INVALID},
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
		data              []byte
		expectedTokenList []Token
		errStr            string
	}{
		{
			data:              []byte(`"name":"ayo"}`),
			expectedTokenList: []Token{},
			errStr:            "json file is empty, so this is an illegal string",
		},
		{
			data:              []byte(`"A JSON payload should be an object or array, not a string."`),
			expectedTokenList: []Token{},
			errStr:            "json file is empty, so this is an illegal string",
		},
		{
			data:              []byte(`{unquoted_key: "keys must be quoted"}`),
			expectedTokenList: []Token{},
			errStr:            "invalid token: u",
		},
		{
			data:              []byte(`["Extra close"]]`),
			expectedTokenList: []Token{},
			errStr:            "not storing token states well (WITHIN_ARRAY) or invalid syntax",
		},
		{
			data:              []byte(`{"Extra value after close": true} "misplaced quoted value"`),
			expectedTokenList: []Token{},
			errStr:            "string token isn't WITHIN_OBJECT or WITHIN_ARRAY",
		},
		{
			data:              []byte(`{"Illegal expression": 1 + 2}`),
			expectedTokenList: []Token{},
			errStr:            "",
		},
		{
			data:              []byte(`{"Illegal invocation": alert()}`),
			expectedTokenList: []Token{},
			errStr:            "invalid token: a",
		},
		{
			data:              []byte(`{"Numbers cannot have leading zeroes": 013}`),
			expectedTokenList: []Token{},
			errStr:            "a number can't start with 0",
		},
		{
			data:              []byte(`{"Numbers cannot be hex": 0x14}`),
			expectedTokenList: []Token{},
			errStr:            "a number can't start with 0",
		},
		{
			data:              []byte(`["Illegal backslash escape: \x15"]`),
			expectedTokenList: []Token{},
			errStr:            "illegal escaped character: x",
		},
		{
			data:              []byte(`["Bad value", truth]`),
			expectedTokenList: []Token{},
			errStr:            "there's a t in the json that isn't in a string and isn't the true literal",
		},
		{
			data:              []byte(`['single quote']`),
			expectedTokenList: []Token{},
			errStr:            "invalid token: '",
		},
		{
			data:              []byte(`[2e]`),
			expectedTokenList: []Token{},
			errStr:            "character after the exponent sym (e|E) isn't a number",
		},
		{
			data:              []byte(`[2e+]`),
			expectedTokenList: []Token{},
			errStr:            "character after the plus or minus sym (+|-) isn't a number",
		},
		{
			data:              []byte(`[2e+-1]`),
			expectedTokenList: []Token{},
			errStr:            "character after the plus or minus sym (+|-) isn't a number",
		},
		{
			data:              []byte(`["mismatch"}`),
			expectedTokenList: []Token{},
			errStr:            "not WITHIN_OBJECT",
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

func TestGetTokenState_Panics(t *testing.T) {
	s := NewStack()
	s.Push(1)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	getTokenState(s)

	// 	assert
	// 	assert.Panics(t, func() { })
	//
	// 	s.Push(2)
	// 	s.Push(3)
}

func TestGetTokenState(t *testing.T) {
	tokenStateStack := NewStack()

	if state := getTokenState(tokenStateStack); state != INVALID {
		t.Errorf("Incorrect top element: expected %v, got %v", INVALID, state)
	}

	tokenStateStack.Push(WITHIN_OBJECT)
	if state := getTokenState(tokenStateStack); state != WITHIN_OBJECT {
		t.Errorf("Incorrect top element: expected %v, got %v", WITHIN_OBJECT, state)
	}

	tokenStateStack.Push(WITHIN_ARRAY)
	if state := getTokenState(tokenStateStack); state != WITHIN_ARRAY {
		t.Errorf("Incorrect top element: expected %v, got %v", WITHIN_ARRAY, state)
	}
}
