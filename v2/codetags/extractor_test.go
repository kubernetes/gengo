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

	"github.com/google/go-cmp/cmp"
)

func TestExtract(t *testing.T) {
	tests := []struct {
		name   string
		prefix string
		lines  []string
		want   map[string][]string
	}{
		{
			name:   "example",
			prefix: "+k8s:",
			lines: []string{
				"Comment line without marker",
				"+k8s:noArgs // comment",
				"+withValue=value1",
				"+withValue=value2",
				"+k8s:withArg(arg1)=value1",
				"+k8s:withArg(arg2)=value2 // comment",
				"+k8s:withNamedArgs(arg1=value1, arg2=value2)=value",
			},
			want: map[string][]string{
				"noArgs":        {"noArgs // comment"},
				"withArg":       {"withArg(arg1)=value1", "withArg(arg2)=value2 // comment"},
				"withNamedArgs": {"withNamedArgs(arg1=value1, arg2=value2)=value"},
			},
		},
		{
			name:   "with args and values",
			prefix: "+",
			lines: []string{
				"+a:t1=tagValue",
				"+b:t1(arg)",
				"+b:t1(arg)=tagValue",
				"+a:t1(arg: value)",
				"+a:t1(arg: value)=tagValue",
			},
			want: map[string][]string{
				"a:t1": {"a:t1=tagValue", "a:t1(arg: value)", "a:t1(arg: value)=tagValue"},
				"b:t1": {"b:t1(arg)", "b:t1(arg)=tagValue"},
			},
		},
		{
			name:   "empty name",
			prefix: "+k8s:",
			lines:  []string{},
			want:   map[string][]string{},
		},
		{
			name:   "no matching lines",
			prefix: "+k8s:",
			lines: []string{
				"Comment line without marker",
				"Another comment line",
			},
			want: map[string][]string{},
		},
		{
			name:   "different marker",
			prefix: "@k8s:",
			lines: []string{
				"Comment line without marker",
				"@k8s:required",
				"@validation:required",
				"+k8s:format=k8s-long-name",
			},
			want: map[string][]string{
				"required": {"required"},
			},
		},
		{
			name:   "no group",
			prefix: "+",
			lines: []string{
				"+k8s:required",
				"+validation:required",
				"+validation:format=special",
			},
			want: map[string][]string{
				"k8s:required":        {"k8s:required"},
				"validation:required": {"validation:required"},
				"validation:format":   {"validation:format=special"},
			},
		},
		{
			name:   "no name",
			prefix: "+",
			lines: []string{
				"+",
				"+ ",
				"+ // comment",
			},
			want: map[string][]string{
				"": {"", " ", " // comment"},
			},
		},
		{
			name:   "whitespace",
			prefix: "+",
			lines: []string{
				" +name",
				" \t \t +name",
				"  +name",
				" +name ",
				" +name  ",
				" +name= value",
				"  +name = value",
				" +name =value ",
			},
			want: map[string][]string{
				"name": {"name", "name", "name", "name ", "name  ", "name= value", "name = value", "name =value "},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Extract(tt.prefix, tt.lines)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got:\n%#+v\nwant:\n%#+v\n", got, tt.want)
			}
		})
	}
}

func TestExtractAndParse(t *testing.T) {
	mktags := func(t ...Tag) []Tag { return t }

	cases := []struct {
		name     string
		comments []string
		expect   map[string][]Tag
	}{
		{
			name: "positional params",
			comments: []string{
				"+quoted(\"value\")",
				"+backticked(`value`)",
				"+ident(value)",
				"+integer(2)",
				"+negative(-5)",
				"+hex(0xFF00B3)",
				"+octal(0o04167)",
				"+binary(0b10101)",
				"+true(true)",
				"+false(false)",
			},
			expect: map[string][]Tag{
				"quoted": mktags(
					Tag{Name: "quoted", Args: []Arg{
						{Value: "value", Type: ArgTypeString},
					}},
				),
				"backticked": mktags(
					Tag{Name: "backticked", Args: []Arg{
						{Value: "value", Type: ArgTypeString},
					}},
				),
				"ident": mktags(
					Tag{Name: "ident", Args: []Arg{
						{Value: "value", Type: ArgTypeString},
					}},
				),
				"integer": mktags(
					Tag{Name: "integer", Args: []Arg{
						{Value: "2", Type: ArgTypeInt},
					}},
				),
				"negative": mktags(
					Tag{Name: "negative", Args: []Arg{
						{Value: "-5", Type: ArgTypeInt},
					}},
				),
				"hex": mktags(
					Tag{Name: "hex", Args: []Arg{
						{Value: "0xFF00B3", Type: ArgTypeInt},
					}},
				),
				"octal": mktags(
					Tag{Name: "octal", Args: []Arg{
						{Value: "0o04167", Type: ArgTypeInt},
					}},
				),
				"binary": mktags(
					Tag{Name: "binary", Args: []Arg{
						{Value: "0b10101", Type: ArgTypeInt},
					}},
				),
				"true": mktags(
					Tag{Name: "true", Args: []Arg{
						{Value: "true", Type: ArgTypeBool},
					}},
				),
				"false": mktags(
					Tag{Name: "false", Args: []Arg{
						{Value: "false", Type: ArgTypeBool},
					}},
				),
			},
		},
		{
			name: "named params",
			comments: []string{
				"+strings(q: \"value\", b: `value`, i: value)",
				"+numbers(n1: 2, n2: -5, n3: 0xFF00B3, n4: 0o04167, n5: 0b10101)",
				"+bools(t: true, f:false)",
			},
			expect: map[string][]Tag{
				"strings": mktags(
					Tag{Name: "strings", Args: []Arg{
						{Name: "q", Value: "value", Type: ArgTypeString},
						{Name: "b", Value: `value`, Type: ArgTypeString},
						{Name: "i", Value: "value", Type: ArgTypeString},
					}}),
				"numbers": mktags(
					Tag{Name: "numbers", Args: []Arg{
						{Name: "n1", Value: "2", Type: ArgTypeInt},
						{Name: "n2", Value: "-5", Type: ArgTypeInt},
						{Name: "n3", Value: "0xFF00B3", Type: ArgTypeInt},
						{Name: "n4", Value: "0o04167", Type: ArgTypeInt},
						{Name: "n5", Value: "0b10101", Type: ArgTypeInt},
					}}),
				"bools": mktags(
					Tag{Name: "bools", Args: []Arg{
						{Name: "t", Value: "true", Type: ArgTypeBool},
						{Name: "f", Value: "false", Type: ArgTypeBool},
					}}),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out := map[string][]Tag{}
			for name, matchedTags := range Extract("+", tc.comments) {
				parsed, err := ParseAll(matchedTags)
				if err != nil {
					t.Fatalf("case %q: unexpected error: %v", tc.name, err)
				}
				out[name] = parsed
			}
			if !reflect.DeepEqual(out, tc.expect) {
				t.Errorf("case %q: wrong result:\n%v", tc.name, cmp.Diff(tc.expect, out))
			}
		})
	}
}
