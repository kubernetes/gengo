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
	"testing"
)

func TestTypedTagString(t *testing.T) {
	tests := []struct {
		name     string
		tag      Tag
		expected string
	}{
		{
			name: "simple tag name",
			tag: Tag{
				Name: "name",
			},
			expected: "name",
		},
		{
			name: "tag with single string arg",
			tag: Tag{
				Name: "tag",
				Args: []Arg{
					{Value: "value", Type: ArgTypeString},
				},
			},
			expected: `tag("value")`,
		},
		{
			name: "tag with multiple positional args",
			tag: Tag{
				Name: "tag",
				Args: []Arg{
					{Value: "value1", Type: ArgTypeString},
					{Value: "value2", Type: ArgTypeString},
				},
			},
			expected: `tag("value1", "value2")`,
		},
		{
			name: "tag with quoted string arg",
			tag: Tag{
				Name: "tag",
				Args: []Arg{
					{Value: "value with spaces", Type: ArgTypeString},
				},
			},
			expected: `tag("value with spaces")`,
		},
		{
			name: "tag with named args",
			tag: Tag{
				Name: "tag",
				Args: []Arg{
					{Name: "arg1", Value: "value1", Type: ArgTypeString},
					{Name: "arg2", Value: "value2", Type: ArgTypeString},
				},
			},
			expected: `tag(arg1: "value1", arg2: "value2")`,
		},
		{
			name: "tag with named arg of different types",
			tag: Tag{
				Name: "tag",
				Args: []Arg{
					{Name: "str", Value: "string value", Type: ArgTypeString},
					{Name: "int", Value: "42", Type: ArgTypeInt},
					{Name: "bool", Value: "true", Type: ArgTypeBool},
				},
			},
			expected: `tag(str: "string value", int: 42, bool: true)`,
		},
		{
			name: "tag with string value",
			tag: Tag{
				Name:      "tag",
				Value:     "value",
				ValueType: ValueTypeString,
			},
			expected: `tag="value"`,
		},
		{
			name: "tag with integer value",
			tag: Tag{
				Name:      "tag",
				Value:     "42",
				ValueType: ValueTypeInt,
			},
			expected: "tag=42",
		},
		{
			name: "tag with boolean value",
			tag: Tag{
				Name:      "tag",
				Value:     "true",
				ValueType: ValueTypeBool,
			},
			expected: "tag=true",
		},
		{
			name: "tag with raw value",
			tag: Tag{
				Name:      "tag",
				Value:     "some raw value // with comment",
				ValueType: ValueTypeRaw,
			},
			expected: "tag=some raw value // with comment",
		},
		{
			name: "tag with nested tag value",
			tag: Tag{
				Name:      "outer",
				ValueType: ValueTypeTag,
				ValueTag: &Tag{
					Name: "inner",
				},
			},
			expected: "outer=+inner",
		},
		{
			name: "complex nested tags",
			tag: Tag{
				Name: "level1",
				Args: []Arg{
					{Name: "arg", Value: "value", Type: ArgTypeString},
				},
				ValueType: ValueTypeTag,
				ValueTag: &Tag{
					Name: "level2",
					Args: []Arg{
						{Value: "42", Type: ArgTypeInt},
					},
					ValueType: ValueTypeTag,
					ValueTag: &Tag{
						Name:      "level3",
						Value:     "final",
						ValueType: ValueTypeString,
					},
				},
			},
			expected: `level1(arg: "value")=+level2(42)=+level3="final"`,
		},
		{
			name: "tag with args and string value",
			tag: Tag{
				Name: "tag",
				Args: []Arg{
					{Value: "arg", Type: ArgTypeString},
				},
				Value:     "value",
				ValueType: ValueTypeString,
			},
			expected: `tag("arg")="value"`,
		},
		{
			name: "tag with args and nested tag value",
			tag: Tag{
				Name: "outer",
				Args: []Arg{
					{Name: "param", Value: "value", Type: ArgTypeString},
				},
				ValueType: ValueTypeTag,
				ValueTag: &Tag{
					Name: "inner",
					Args: []Arg{
						{Value: "innerArg", Type: ArgTypeString},
					},
				},
			},
			expected: `outer(param: "value")=+inner("innerArg")`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.tag.String()
			if result != tt.expected {
				t.Errorf("got: %q, want: %q", result, tt.expected)
			}
		})
	}
}
