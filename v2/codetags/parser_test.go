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

func TestParse(t *testing.T) {
	mksa := func(s ...string) []Arg {
		var args []Arg
		for _, v := range s {
			args = append(args, Arg{Value: v, Type: ArgTypeString})
		}
		return args
	}

	cases := []struct {
		input      string
		expectKey  string
		expectArgs []Arg
		err        bool
	}{
		// string args
		{input: "name", expectKey: "name"},
		{input: "name // comment", expectKey: "name"},
		{input: "name-dash", expectKey: "name-dash"},
		{input: "name.dot", expectKey: "name.dot"},
		{input: "name:colon", expectKey: "name:colon"},
		{input: "name()", expectKey: "name"},
		{input: "name(arg)", expectKey: "name", expectArgs: mksa("arg")},
		{input: "name(arg) // comment", expectKey: "name", expectArgs: mksa("arg")},
		{input: "name(ARG)", expectKey: "name", expectArgs: mksa("ARG")},
		{input: "name(ArG)", expectKey: "name", expectArgs: mksa("ArG")},
		{input: "name(has-dash)", expectKey: "name", expectArgs: mksa("has-dash")},
		{input: "name(has.dot)", expectKey: "name", expectArgs: mksa("has.dot")},
		{input: "name()", expectKey: "name"},
		{input: "name(lower)", expectKey: "name", expectArgs: mksa("lower")},
		{input: "name(CAPITAL)", expectKey: "name", expectArgs: mksa("CAPITAL")},
		{input: "name(MiXeD)", expectKey: "name", expectArgs: mksa("MiXeD")},
		{input: "name(mIxEd)", expectKey: "name", expectArgs: mksa("mIxEd")},
		{input: "name(_under)", expectKey: "name", expectArgs: mksa("_under")},
		{input: `name("hasQuotes")`, expectKey: "name", expectArgs: mksa("hasQuotes")},
		{input: "name(`hasRawQuotes`)", expectKey: "name", expectArgs: mksa("hasRawQuotes")},
		{input: "name(has space)", expectKey: "name", err: true},
		{input: "name(multiple, args)", expectKey: "name", err: true},
		{input: "name(noClosingParen", expectKey: "name", err: true},
		{input: "withRaw(`a = b`)", expectKey: "withRaw", expectArgs: mksa("a = b")},

		// invalid args
		{input: "name(arg1, arg2)", err: true},
		{input: "badRaw(missing`)", err: true},
		{input: "badMix(arg,`raw`)", err: true},

		// quotes
		{input: `quoted(s: "value \" \\")`, expectKey: "quoted", expectArgs: []Arg{
			{Name: "s", Value: "value \" \\", Type: ArgTypeString},
		}},
		{input: "backticks(s: `value`)", expectKey: "backticks", expectArgs: []Arg{
			{Name: "s", Value: `value`, Type: ArgTypeString},
		}},
		{input: "ident(k: value)", expectKey: "ident", expectArgs: []Arg{
			{Name: "k", Value: "value", Type: ArgTypeString},
		}},

		// numbers
		{input: "numbers(n1: 2, n2: -5, n3: 0xFF00B3, n4: 0o04167, n5: 0b10101)", expectKey: "numbers", expectArgs: []Arg{
			{Name: "n1", Value: "2", Type: ArgTypeInt},
			{Name: "n2", Value: "-5", Type: ArgTypeInt},
			{Name: "n3", Value: "0xFF00B3", Type: ArgTypeInt},
			{Name: "n4", Value: "0o04167", Type: ArgTypeInt},
			{Name: "n5", Value: "0b10101", Type: ArgTypeInt},
		}},

		// bools
		{input: "bools(t: true, f:false)", expectKey: "bools", expectArgs: []Arg{
			{Name: "t", Value: "true", Type: ArgTypeBool},
			{Name: "f", Value: "false", Type: ArgTypeBool},
		}},

		// mixed type args
		{input: "mixed(s: `value`, i: 2, b: true)", expectKey: "mixed", expectArgs: []Arg{
			{Name: "s", Value: "value", Type: ArgTypeString},
			{Name: "i", Value: "2", Type: ArgTypeInt},
			{Name: "b", Value: "true", Type: ArgTypeBool},
		}},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			parsed, err := parseTagKey(tc.input)
			if err != nil && tc.err == false {
				t.Errorf("[%q]: expected success, got: %v", tc.input, err)
				return
			}
			if err == nil {
				if tc.err == true {
					t.Errorf("[%q]: expected failure, got: %v(%v)", tc.input, parsed.Name, parsed.Args)
					return
				}
				if parsed.Name != tc.expectKey {
					t.Errorf("[%q]\nexpected key: %q, got: %q", tc.input, tc.expectKey, parsed.Name)
				}
				if len(parsed.Args) != len(tc.expectArgs) {
					t.Errorf("[%q]: expected %d args, got: %q", tc.input, len(tc.expectArgs), parsed.Args)
					return
				}
				for i := range tc.expectArgs {
					if want, got := tc.expectArgs[i], parsed.Args[i]; got != want {
						t.Errorf("[%q]\nexpected %q, got %q", tc.input, want, got)
					}
				}
			}
		})
	}
}
