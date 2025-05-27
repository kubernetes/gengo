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
	mkss := func(s ...string) []Arg {
		var args []Arg
		for _, v := range s {
			args = append(args, Arg{Value: ArgString(v)})
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
		{input: "name(arg)", expectKey: "name", expectArgs: mkss("arg")},
		{input: "name(arg) // comment", expectKey: "name", expectArgs: mkss("arg")},
		{input: "name(ARG)", expectKey: "name", expectArgs: mkss("ARG")},
		{input: "name(ArG)", expectKey: "name", expectArgs: mkss("ArG")},
		{input: "name(has-dash)", expectKey: "name", expectArgs: mkss("has-dash")},
		{input: "name(has.dot)", expectKey: "name", expectArgs: mkss("has.dot")},
		{input: "name()", expectKey: "name"},
		{input: "name(lower)", expectKey: "name", expectArgs: mkss("lower")},
		{input: "name(CAPITAL)", expectKey: "name", expectArgs: mkss("CAPITAL")},
		{input: "name(MiXeD)", expectKey: "name", expectArgs: mkss("MiXeD")},
		{input: "name(mIxEd)", expectKey: "name", expectArgs: mkss("mIxEd")},
		{input: "name(_under)", expectKey: "name", expectArgs: mkss("_under")},
		{input: `name("hasQuotes")`, expectKey: "name", expectArgs: mkss("hasQuotes")},
		{input: "name(`hasRawQuotes`)", expectKey: "name", expectArgs: mkss("hasRawQuotes")},
		{input: "name(has space)", expectKey: "name", err: true},
		{input: "name(multiple, args)", expectKey: "name", err: true},
		{input: "name(noClosingParen", expectKey: "name", err: true},

		// invalid args
		{input: "name(arg1, arg2)", err: true},
		{input: "withRaw(`a = b`)", expectKey: "withRaw", expectArgs: mkss("a = b")},
		{input: "badRaw(missing`)", err: true},
		{input: "badMix(arg,`raw`)", err: true},

		// quotes
		{input: `quoted(s: "value \" \\")`, expectKey: "quoted", expectArgs: []Arg{
			{Name: "s", Value: ArgString("value \" \\")},
		}},
		{input: "backticks(s: `value`)", expectKey: "backticks", expectArgs: []Arg{
			{Name: "s", Value: ArgString(`value`)},
		}},
		{input: "ident(k: value)", expectKey: "ident", expectArgs: []Arg{
			{Name: "k", Value: ArgString("value")},
		}},

		// numbers
		{input: "numbers(n1: 2, n2: -5, n3: 0xFF00B3, n4: 0o04167, n5: 0b10101)", expectKey: "numbers", expectArgs: []Arg{
			{Name: "n1", Value: MustArgInt("2")},
			{Name: "n2", Value: MustArgInt("-5")},
			{Name: "n3", Value: MustArgInt("0xFF00B3")},
			{Name: "n4", Value: MustArgInt("0o04167")},
			{Name: "n5", Value: MustArgInt("0b10101")},
		}},

		// bools
		{input: "bools(t: true, f:false)", expectKey: "bools", expectArgs: []Arg{
			{Name: "t", Value: ArgBool(true)},
			{Name: "f", Value: ArgBool(false)},
		}},

		// mixed type args
		{input: "mixed(s: `value`, i: 2, b: true)", expectKey: "mixed", expectArgs: []Arg{
			{Name: "s", Value: ArgString("value")},
			{Name: "i", Value: MustArgInt("2")},
			{Name: "b", Value: ArgBool(true)},
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
					t.Errorf("[%q]: expected failure, got: %v(%v)", tc.input, parsed.name, parsed.args)
					return
				}
				if parsed.name != tc.expectKey {
					t.Errorf("[%q]\nexpected key: %q, got: %q", tc.input, tc.expectKey, parsed.name)
				}
				if len(parsed.args) != len(tc.expectArgs) {
					t.Errorf("[%q]: expected %d args, got: %q", tc.input, len(tc.expectArgs), parsed.args)
					return
				}
				for i := range tc.expectArgs {
					if want, got := tc.expectArgs[i], parsed.args[i]; got != want {
						t.Errorf("[%q]\nexpected %q, got %q", tc.input, want, got)
					}
				}
			}
		})
	}
}
