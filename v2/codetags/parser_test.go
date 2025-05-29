/*
Copyright 2025 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package codetags

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	mksa := func(s ...string) []Arg {
		var args []Arg
		for _, v := range s {
			args = append(args, Arg{Value: v, Type: ArgTypeString})
		}
		return args
	}

	mkt := func(name string) TypedTag {
		return TypedTag{Name: name}
	}

	mkta := func(name string, args []Arg) TypedTag {
		return TypedTag{Name: name, Args: args}
	}

	mktv := func(name, value string, valueType ValueType) TypedTag {
		return TypedTag{Name: name, Value: value, ValueType: valueType}
	}

	mktt := func(name string, valueTag *TypedTag) TypedTag {
		return TypedTag{Name: name, ValueTag: valueTag, ValueType: ValueTypeTag}
	}

	cases := []struct {
		name         string
		input        string
		parseOptions []ParseOption
		expect       TypedTag
		err          bool
	}{
		// Basic tag name tests
		{
			name:   "simple name",
			input:  "name",
			expect: mkt("name"),
		},
		{
			name:   "name with comment",
			input:  "name // comment",
			expect: mkt("name"),
		},
		{
			name:   "name with dash",
			input:  "name-dash",
			expect: mkt("name-dash"),
		},
		{
			name:   "name with dot",
			input:  "name.dot",
			expect: mkt("name.dot"),
		},
		{
			name:   "name with colon",
			input:  "name:colon",
			expect: mkt("name:colon"),
		},

		// String arguments tests
		{
			name:   "empty parentheses",
			input:  "name()",
			expect: mkt("name"),
		},
		{
			name:   "single argument",
			input:  "name(arg)",
			expect: mkta("name", mksa("arg")),
		},
		{
			name:   "argument with comment",
			input:  "name(arg) // comment",
			expect: mkta("name", mksa("arg")),
		},
		{
			name:   "uppercase argument",
			input:  "name(ARG)",
			expect: mkta("name", mksa("ARG")),
		},
		{
			name:   "mixed case argument",
			input:  "name(ArG)",
			expect: mkta("name", mksa("ArG")),
		},
		{
			name:   "argument with dash",
			input:  "name(has-dash)",
			expect: mkta("name", mksa("has-dash")),
		},
		{
			name:   "argument with dot",
			input:  "name(has.dot)",
			expect: mkta("name", mksa("has.dot")),
		},
		{
			name:   "lowercase argument",
			input:  "name(lower)",
			expect: mkta("name", mksa("lower")),
		},
		{
			name:   "uppercase argument",
			input:  "name(CAPITAL)",
			expect: mkta("name", mksa("CAPITAL")),
		},
		{
			name:   "mixed case argument 1",
			input:  "name(MiXeD)",
			expect: mkta("name", mksa("MiXeD")),
		},
		{
			name:   "mixed case argument 2",
			input:  "name(mIxEd)",
			expect: mkta("name", mksa("mIxEd")),
		},
		{
			name:   "underscore in argument",
			input:  "name(_under)",
			expect: mkta("name", mksa("_under")),
		},
		{
			name:   "quoted argument",
			input:  `name("hasQuotes")`,
			expect: mkta("name", mksa("hasQuotes")),
		},
		{
			name:   "raw quoted argument",
			input:  "name(`hasRawQuotes`)",
			expect: mkta("name", mksa("hasRawQuotes")),
		},
		{
			name:   "raw string with spaces",
			input:  "withRaw(`a = b`)",
			expect: mkta("withRaw", mksa("a = b")),
		},

		// Error cases for arguments
		{
			name:  "space in argument",
			input: "name(has space)",
			err:   true,
		},
		{
			name:  "multiple comma-separated args",
			input: "name(multiple, args)",
			err:   true,
		},
		{
			name:  "unclosed parenthesis",
			input: "name(noClosingParen",
			err:   true,
		},
		{
			name:  "comma-separated args",
			input: "name(arg1, arg2)",
			err:   true,
		},
		{
			name:  "unclosed raw string",
			input: "badRaw(missing`)",
			err:   true,
		},
		{
			name:  "mixed arg formats",
			input: "badMix(arg,`raw`)",
			err:   true,
		},
		{
			name:  "nested: comma-separated args",
			input: "name=+name(arg1, arg2)",
			err:   true,
		},
		{
			name:  "nested: unclosed raw string",
			input: "name=+badRaw(missing`)",
			err:   true,
		},
		{
			name:  "nested: mixed arg formats",
			input: "name=+badMix(arg,`raw`)",
			err:   true,
		},

		// Named arguments tests
		{
			name:  "quoted named argument",
			input: `quoted(s: "value \" \\")`,
			expect: mkta("quoted", []Arg{
				{Name: "s", Value: "value \" \\", Type: ArgTypeString},
			}),
		},
		{
			name:  "backtick named argument",
			input: "backticks(s: `value`)",
			expect: mkta("backticks", []Arg{
				{Name: "s", Value: `value`, Type: ArgTypeString},
			}),
		},
		{
			name:  "identifier named argument",
			input: "ident(k: value)",
			expect: mkta("ident", []Arg{
				{Name: "k", Value: "value", Type: ArgTypeString},
			}),
		},

		// Numeric argument tests
		{
			name:  "numeric arguments",
			input: "numbers(n1: 2, n2: -5, n3: 0xFF00B3, n4: 0o04167, n5: 0b10101)",
			expect: mkta("numbers", []Arg{
				{Name: "n1", Value: "2", Type: ArgTypeInt},
				{Name: "n2", Value: "-5", Type: ArgTypeInt},
				{Name: "n3", Value: "0xFF00B3", Type: ArgTypeInt},
				{Name: "n4", Value: "0o04167", Type: ArgTypeInt},
				{Name: "n5", Value: "0b10101", Type: ArgTypeInt},
			}),
		},

		// Boolean argument tests
		{
			name:  "boolean arguments",
			input: "bools(t: true, f:false)",
			expect: mkta("bools", []Arg{
				{Name: "t", Value: "true", Type: ArgTypeBool},
				{Name: "f", Value: "false", Type: ArgTypeBool},
			}),
		},

		// Mixed type argument tests
		{
			name:  "mixed type arguments",
			input: "mixed(s: `value`, i: 2, b: true)",
			expect: mkta("mixed", []Arg{
				{Name: "s", Value: "value", Type: ArgTypeString},
				{Name: "i", Value: "2", Type: ArgTypeInt},
				{Name: "b", Value: "true", Type: ArgTypeBool},
			}),
		},

		// Scalar value tests
		{
			name:   "integer value",
			input:  `key=1`,
			expect: mktv("key", "1", ValueTypeInt),
		},
		{
			name:   "true value",
			input:  `key=true`,
			expect: mktv("key", "true", ValueTypeBool),
		},
		{
			name:   "false value",
			input:  `key=false`,
			expect: mktv("key", "false", ValueTypeBool),
		},
		{
			name:   "identifier value",
			input:  `key=ident`,
			expect: mktv("key", "ident", ValueTypeString),
		},
		{
			name:   "quoted string value",
			input:  `key="quoted"`,
			expect: mktv("key", "quoted", ValueTypeString),
		},
		{
			name:   "backtick string value",
			input:  "key=`quoted`",
			expect: mktv("key", "quoted", ValueTypeString),
		},

		// Error cases for scalar values
		{
			name:  "space in value",
			input: "key=one two",
			err:   true,
		},
		{
			name:  "unclosed backtick quoted string",
			input: "key=`unclosed",
			err:   true,
		},
		{
			name:  "unclosed double-quoted string",
			input: `key="unclosed`,
			err:   true,
		},
		{
			name:  "illegal identifier",
			input: `key=x@y`,
			err:   true,
		},

		// Scalar values with comments
		{
			name:   "integer value with comment",
			input:  `key=1 // comment`,
			expect: mktv("key", "1", ValueTypeInt),
		},
		{
			name:   "true value with comment",
			input:  `key=true // comment`,
			expect: mktv("key", "true", ValueTypeBool),
		},
		{
			name:   "false value with comment",
			input:  `key=false // comment`,
			expect: mktv("key", "false", ValueTypeBool),
		},
		{
			name:   "identifier value with comment",
			input:  `key=ident // comment`,
			expect: mktv("key", "ident", ValueTypeString),
		},
		{
			name:   "quoted string value with comment",
			input:  `key="quoted" // comment`,
			expect: mktv("key", "quoted", ValueTypeString),
		},
		{
			name:   "backtick string value with comment",
			input:  "key=`quoted` // comment",
			expect: mktv("key", "quoted", ValueTypeString),
		},

		// Tag values
		{
			name:   "simple tag value",
			input:  `key=+key2`,
			expect: mktt("key", &TypedTag{Name: "key2"}),
		},
		{
			name:   "tag value with empty parentheses",
			input:  `key=+key2()`,
			expect: mktt("key", &TypedTag{Name: "key2"}),
		},
		{
			name:   "tag value with argument",
			input:  `key=+key2(arg)`,
			expect: mktt("key", &TypedTag{Name: "key2", Args: []Arg{{Value: "arg", Type: ArgTypeString}}}),
		},
		{
			name:  "tag value with named arguments",
			input: `key=+key2(k1: v1, k2: v2)`,
			expect: mktt("key", &TypedTag{Name: "key2", Args: []Arg{
				{Name: "k1", Value: "v1", Type: ArgTypeString},
				{Name: "k2", Value: "v2", Type: ArgTypeString},
			}}),
		},
		{
			name:  "nested tag values",
			input: `key=+key2=+key3`,
			expect: mktt("key", &TypedTag{
				Name:      "key2",
				ValueType: ValueTypeTag,
				ValueTag:  &TypedTag{Name: "key3"},
			}),
		},

		// Tag values with comments
		{
			name:   "simple tag value with comment",
			input:  `key=+key2 // comment`,
			expect: mktt("key", &TypedTag{Name: "key2"}),
		},
		{
			name:   "tag value with empty parentheses and comment",
			input:  `key=+key2() // comment`,
			expect: mktt("key", &TypedTag{Name: "key2"}),
		},
		{
			name:   "tag value with argument and comment",
			input:  `key=+key2(arg) // comment`,
			expect: mktt("key", &TypedTag{Name: "key2", Args: []Arg{{Value: "arg", Type: ArgTypeString}}}),
		},
		{
			name:  "tag value with named arguments and comment",
			input: `key=+key2(k1: v1, k2: v2) // comment`,
			expect: mktt("key", &TypedTag{Name: "key2", Args: []Arg{
				{Name: "k1", Value: "v1", Type: ArgTypeString},
				{Name: "k2", Value: "v2", Type: ArgTypeString},
			}}),
		},
		{
			name:  "nested tag values with comment",
			input: `key=+key2=+key3 // comment`,
			expect: mktt("key", &TypedTag{
				Name:      "key2",
				ValueType: ValueTypeTag,
				ValueTag:  &TypedTag{Name: "key3"},
			}),
		},

		// Raw values
		{
			name:         "raw arbitrary content",
			input:        `key=this is \ arbitrary content !!`,
			parseOptions: []ParseOption{RawValues(true)},
			expect:       mktv("key", "this is \\ arbitrary content !!", ValueTypeRaw),
		},
		{
			name:         "raw value with comment",
			input:        `key=true // comment`,
			parseOptions: []ParseOption{RawValues(true)},
			expect:       mktv("key", "true // comment", ValueTypeRaw),
		},
		{
			name:         "raw value with tag syntax",
			input:        `key=+key2`,
			parseOptions: []ParseOption{RawValues(true)},
			expect:       mktv("key", "+key2", ValueTypeRaw),
		},
		{
			name:         "raw with key only",
			input:        `key`,
			parseOptions: []ParseOption{RawValues(true)},
			expect:       TypedTag{Name: "key"},
		},
		{
			name:         "raw with empty value",
			input:        `key=`,
			parseOptions: []ParseOption{RawValues(true)},
			expect:       mktv("key", "", ValueTypeRaw),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			parsed, err := Parse(tc.input, tc.parseOptions...)
			if err != nil {
				if !tc.err {
					t.Errorf("Expected success, got error: %v", err)
				}
				return
			}
			if tc.err {
				t.Errorf("Expected error, got success: %v", parsed)
				return
			}
			if !reflect.DeepEqual(parsed, tc.expect) {
				t.Errorf("Parsed tag doesn't match expected.\nExpected: %#v\nGot: %#v", tc.expect, parsed)
			}

			// round-trip testing
			roundTripped, err := Parse(parsed.String(), tc.parseOptions...)
			if err != nil {
				t.Errorf("Failed to reparse tag: %v", err)
				return
			}
			if !reflect.DeepEqual(roundTripped, parsed) {
				t.Errorf("Round-tripped tag doesn't match original.\nExpected: %#v\nGot: %#v", parsed, roundTripped)
			}
		})
	}
}
