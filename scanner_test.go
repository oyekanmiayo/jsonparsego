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
