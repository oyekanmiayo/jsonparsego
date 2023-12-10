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
			name: "Valid string json",
			tokenList: func() []Token {
				tokens, _ := scanTokens([]byte(`{"name":"ayo"}`))
				return tokens
			}(),
			expectedBool: true,
		},
		{
			name: "Valid num json",
			tokenList: func() []Token {
				tokens, _ := scanTokens([]byte(
					`{
					  "num1": 123,
					  "num2": 987
					}`))
				return tokens
			}(),
			expectedBool: true,
		},
		{
			name: "Valid array json with string at idx 0",
			tokenList: func() []Token {
				tokens, _ := scanTokens([]byte(
					`{
					  "array": ["testing", "hello", 123, {"name":"hello"}]
					}`))
				return tokens
			}(),
			expectedBool: true,
		},
		{
			tokenList: func() []Token {
				tokens, _ := scanTokens([]byte(`{}`))
				return tokens
			}(),
			expectedBool: true,
		},
		{
			tokenList: func() []Token {
				tokens, _ := scanTokens([]byte(`
				[
					"JSON Test Pattern pass1",
					{"object with 1 member":["array with 1 element"]},
					{},
					[],
					-42,
					true,
					false,
					null,
					{
						"integer": 1234567890,
						"real": -9876.543210,
						"e": 0.123456789e-12,
						"E": 1.234567890E+34,
						"":  23456789012E66,
						"zero": 0,
						"one": 1,
						"space": " ",
						"quote": "\"",
						"backslash": "\\",
						"controls": "\b\f\n\r\t",
						"slash": "/ & \/",
						"alpha": "abcdefghijklmnopqrstuvwyz",
						"ALPHA": "ABCDEFGHIJKLMNOPQRSTUVWYZ",
						"digit": "0123456789",
						"0123456789": "digit",
						"special": "1~!@#$%^&*()_+-={':[,]}|;.</>?",
									"hex": "\u0123\u4567\u89AB\uCDEF\uabcd\uef4A",
										"true": true,
										"false": false,
										"null": null,
										"array":[  ],
							"object":{  },
							"address": "50 St. James Street",
							"url": "http://www.JSON.org/",
							"comment": "// /* <!-- --",
							"# -- --> */": " ",
							" s p a c e d " :[1,2 , 3
				
							,
				
							4 , 5        ,          6           ,7        ],"compact":[1,2,3,4,5,6,7],
							"jsontext": "{\"object with 1 member\":[\"array with 1 element\"]}",
							"quotes": "&#34; \u0022 %22 0x22 034 &#x22;",
							"\/\\\"\uCAFE\uBABE\uAB98\uFCDE\ubcda\uef4A\b\f\n\r\t1~!@#$%^&*()_+-=[]{}|;:',./<>?"
							: "A key can be any string"
						},
						0.5 ,98.6
						,
						99.44
						,
				
						1066,
						1e1,
						0.1e1,
						1e-1,
						1e00,2e+00,2e-00
						,"rosebud"]
				`))
				return tokens
			}(),
			expectedBool: true,
		},
		{
			tokenList: func() []Token {
				tokens, _ := scanTokens([]byte(`[[[[[[[[[[[[[[[[[[["Not too deep"]]]]]]]]]]]]]]]]]]]`))
				return tokens
			}(),
			expectedBool: true,
		},
		{
			tokenList: func() []Token {
				tokens, _ := scanTokens([]byte(`
				{
					"JSON Test Pattern pass3": {
						"The outermost value": "must be an object or array.",
						"In this test": "It is an object."
					}
				}`))
				return tokens
			}(),
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
	}{
		{
			tokenList: func() []Token {
				tokens, _ := scanTokens([]byte(`["Unclosed array"`))
				return tokens
			}(),
			expectedBool: false,
		},
		{
			tokenList: func() []Token {
				tokens, _ := scanTokens([]byte(`["extra comma",]`))
				return tokens
			}(),
			expectedBool: false,
		},
		{
			tokenList: func() []Token {
				tokens, _ := scanTokens([]byte(`["double extra comma",,]`))
				return tokens
			}(),
			expectedBool: false,
		},
		{
			tokenList: func() []Token {
				tokens, _ := scanTokens([]byte(`[   , "<-- missing value"]`))
				return tokens
			}(),
			expectedBool: false,
		},
		{
			tokenList: func() []Token {
				tokens, _ := scanTokens([]byte(`["Comma after the close"],`))
				return tokens
			}(),
			expectedBool: false,
		},
		{
			tokenList: func() []Token {
				tokens, _ := scanTokens([]byte(`{"Extra comma": true,}`))
				return tokens
			}(),
			expectedBool: false,
		},
		{
			tokenList: func() []Token {
				tokens, _ := scanTokens([]byte(`{"Missing colon" null}`))
				return tokens
			}(),
			expectedBool: false,
		},
		{
			tokenList: func() []Token {
				tokens, _ := scanTokens([]byte(`{"Double colon":: null}`))
				return tokens
			}(),
			expectedBool: false,
		},
		{
			tokenList: func() []Token {
				tokens, _ := scanTokens([]byte(`{"Comma instead of colon", null}`))
				return tokens
			}(),
			expectedBool: false,
		},
		{
			tokenList: func() []Token {
				tokens, _ := scanTokens([]byte(`["Colon instead of comma": false]`))
				return tokens
			}(),
			expectedBool: false,
		},
		{
			tokenList: func() []Token {
				tokens, _ := scanTokens([]byte(`{"Comma instead if closing brace": true,`))
				return tokens
			}(),
			expectedBool: false,
		},
	}
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
