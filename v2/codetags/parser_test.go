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
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParse(t *testing.T) {
	mksa := func(s ...string) []Arg {
		var args []Arg
		for _, v := range s {
			args = append(args, Arg{Value: v, Type: ArgTypeString})
		}
		return args
	}

	mkt := func(name string) Tag {
		return Tag{Name: name}
	}

	mkta := func(name string, args []Arg) Tag {
		return Tag{Name: name, Args: args}
	}

	mktv := func(name, value string, valueType ValueType) Tag {
		return Tag{Name: name, Value: value, ValueType: valueType}
	}

	mktt := func(name string, valueTag *Tag) Tag {
		return Tag{Name: name, ValueTag: valueTag, ValueType: ValueTypeTag}
	}

	cases := []struct {
		name         string
		input        string
		parseOptions []ParseOption
		expect       Tag
		wantError    string // substring natch
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

		// Error cases for basic tags
		{
			name:      "empty value",
			input:     "name=",
			wantError: "unexpected end of input",
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
			name:   "underscore in argument",
			input:  "name(_)",
			expect: mkta("name", mksa("_")),
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
		{
			name:   `"true" string as argument`,
			input:  `name("true")`,
			expect: mkta("name", mksa("true")),
		},
		{
			name:   `"true" and "false" strings as named arguments`,
			input:  `name(t: "true", f: "false")`,
			expect: mkta("name", []Arg{{Name: "t", Value: "true", Type: ArgTypeString}, {Name: "f", Value: "false", Type: ArgTypeString}}),
		},
		{
			name:   `strings containing "true" and "false"`,
			input:  `name("true")`,
			expect: mkta("name", mksa("true")),
		},

		// Error cases for arguments
		{
			name:      "space in argument",
			input:     "name(has space)",
			wantError: "unexpected character 's'",
		},
		{
			name:      "multiple positional args",
			input:     "name(multiple, args)",
			wantError: "multiple arguments must use 'name: value' syntax",
		},
		{
			name:      "unclosed parenthesis",
			input:     "name(noClosingParen",
			wantError: "unexpected end of input",
		},
		{
			name:      "unclosed raw string",
			input:     "badRaw(missing`)",
			wantError: "unterminated string",
		},
		{
			name:      "nested: comma-separated args",
			input:     "name=+name(arg1, arg2)",
			wantError: "multiple arguments must use 'name: value' syntax",
		},
		{
			name:      "nested: unclosed raw string",
			input:     "name=+badRaw(missing`)",
			wantError: "unterminated string",
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
			input: "numbers(n1: 2, n2: -5, n3: +5, n4: 0xFF, n5: -0x0000, n6: +0x00FF00, n7: 0o04167, n8: -0o04167, n9: +0o04167, n10: 0b10101, n11: -0b10101, n12: +0b10101)",
			expect: mkta("numbers", []Arg{
				{Name: "n1", Value: "2", Type: ArgTypeInt},
				{Name: "n2", Value: "-5", Type: ArgTypeInt},
				{Name: "n3", Value: "+5", Type: ArgTypeInt},
				{Name: "n4", Value: "0xFF", Type: ArgTypeInt},
				{Name: "n5", Value: "-0x0000", Type: ArgTypeInt},
				{Name: "n6", Value: "+0x00FF00", Type: ArgTypeInt},
				{Name: "n7", Value: "0o04167", Type: ArgTypeInt},
				{Name: "n8", Value: "-0o04167", Type: ArgTypeInt},
				{Name: "n9", Value: "+0o04167", Type: ArgTypeInt},
				{Name: "n10", Value: "0b10101", Type: ArgTypeInt},
				{Name: "n11", Value: "-0b10101", Type: ArgTypeInt},
				{Name: "n12", Value: "+0b10101", Type: ArgTypeInt},
			}),
		},

		// Boolean argument tests
		{
			name:  "boolean argument",
			input: "bools(true)",
			expect: mkta("bools", []Arg{
				{Value: "true", Type: ArgTypeBool},
			}),
		},
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
			name:      "space in value",
			input:     "key=one two",
			wantError: "unexpected character 't'",
		},
		{
			name:      "unclosed backtick quoted string",
			input:     "key=`unclosed",
			wantError: "unterminated string",
		},
		{
			name:      "unclosed double-quoted string",
			input:     `key="unclosed`,
			wantError: "unterminated string",
		},
		{
			name:      "illegal identifier",
			input:     `key=x@y`,
			wantError: "unexpected character '@'",
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
			expect: mktt("key", &Tag{Name: "key2"}),
		},
		{
			name:   "tag value with empty parentheses",
			input:  `key=+key2()`,
			expect: mktt("key", &Tag{Name: "key2"}),
		},
		{
			name:   "tag value with argument",
			input:  `key=+key2(arg)`,
			expect: mktt("key", &Tag{Name: "key2", Args: []Arg{{Value: "arg", Type: ArgTypeString}}}),
		},
		{
			name:  "tag value with named arguments",
			input: `key=+key2(k1: v1, k2: v2)`,
			expect: mktt("key", &Tag{Name: "key2", Args: []Arg{
				{Name: "k1", Value: "v1", Type: ArgTypeString},
				{Name: "k2", Value: "v2", Type: ArgTypeString},
			}}),
		},
		{
			name:  "nested tag values",
			input: `key=+key2=+key3`,
			expect: mktt("key", &Tag{
				Name:      "key2",
				ValueType: ValueTypeTag,
				ValueTag:  &Tag{Name: "key3"},
			}),
		},

		// Tag values with comments
		{
			name:   "simple tag value with comment",
			input:  `key=+key2 // comment`,
			expect: mktt("key", &Tag{Name: "key2"}),
		},
		{
			name:   "tag value with empty parentheses and comment",
			input:  `key=+key2() // comment`,
			expect: mktt("key", &Tag{Name: "key2"}),
		},
		{
			name:   "tag value with argument and comment",
			input:  `key=+key2(arg) // comment`,
			expect: mktt("key", &Tag{Name: "key2", Args: []Arg{{Value: "arg", Type: ArgTypeString}}}),
		},
		{
			name:  "tag value with named arguments and comment",
			input: `key=+key2(k1: v1, k2: v2) // comment`,
			expect: mktt("key", &Tag{Name: "key2", Args: []Arg{
				{Name: "k1", Value: "v1", Type: ArgTypeString},
				{Name: "k2", Value: "v2", Type: ArgTypeString},
			}}),
		},
		{
			name:  "3 level nested tag values with comment",
			input: `key=+key2=+key3 // comment`,
			expect: mktt("key", &Tag{
				Name:      "key2",
				ValueType: ValueTypeTag,
				ValueTag:  &Tag{Name: "key3"},
			}),
		},
		{
			name:  "4 level nested tag values with value",
			input: `key=+key2=+key3=+key4=value`,
			expect: mktt("key", &Tag{
				Name:      "key2",
				ValueType: ValueTypeTag,
				ValueTag: &Tag{
					Name:      "key3",
					ValueType: ValueTypeTag,
					ValueTag:  &Tag{Name: "key4", Value: "value", ValueType: ValueTypeString},
				},
			}),
		},
		{
			name:  "4 level nested tag values with args",
			input: `key(arg1)=+key2(arg2)=+key3(arg3)=+key4(arg4)`,
			expect: Tag{
				Name:      "key",
				Args:      []Arg{{Value: "arg1", Type: ArgTypeString}},
				ValueType: ValueTypeTag,
				ValueTag: &Tag{
					Name:      "key2",
					Args:      []Arg{{Value: "arg2", Type: ArgTypeString}},
					ValueType: ValueTypeTag,
					ValueTag: &Tag{
						Name:      "key3",
						Args:      []Arg{{Value: "arg3", Type: ArgTypeString}},
						ValueType: ValueTypeTag,
						ValueTag: &Tag{
							Name: "key4",
							Args: []Arg{{Value: "arg4", Type: ArgTypeString}},
						},
					},
				},
			},
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
			expect:       Tag{Name: "key"},
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
				if len(tc.wantError) == 0 {
					t.Errorf("Expected success, got error: %v", err)
				}
				if !strings.Contains(err.Error(), tc.wantError) {
					t.Errorf("Expected error to contain %q, got %q", tc.wantError, err.Error())
				}
				return
			}
			if len(tc.wantError) > 0 {
				t.Errorf("Expected error, got success: %v", parsed)
				return
			}
			if !reflect.DeepEqual(parsed, tc.expect) {
				t.Errorf("Parsed tag doesn't match expected.\nExpected: %s\n     Got: %s\n\n%s\n",
					tc.expect.String(), parsed.String(), cmp.Diff(tc.expect, parsed))
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
