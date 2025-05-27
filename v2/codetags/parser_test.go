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
	"fmt"
	"testing"
)

func TestParseTagKeyWithTagNames(t *testing.T) {
	mkss := func(s ...string) []string { return s }

	cases := []struct {
		input      string
		expectKey  string
		expectArgs []string
		err        bool
	}{
		{input: "name", expectKey: "name"},
		{input: "name()", expectKey: "name"},
		{input: "name(arg)", expectKey: "name", expectArgs: mkss("arg")},
		{input: "name()", expectKey: "name"},
		{input: "name(lower)", expectKey: "name", expectArgs: mkss("lower")},
		{input: "name(CAPITAL)", expectKey: "name", expectArgs: mkss("CAPITAL")},
		{input: "name(MiXeD)", expectKey: "name", expectArgs: mkss("MiXeD")},
		{input: "name(mIxEd)", expectKey: "name", expectArgs: mkss("mIxEd")},
		{input: "name(_under)", expectKey: "name", expectArgs: mkss("_under")},
		{input: `name("hasQuotes")`, expectKey: "name", expectArgs: mkss("hasQuotes")},
		{input: "name(`hasRawQuotes`)", expectKey: "name", expectArgs: mkss("hasRawQuotes")},
		{input: "name(has space)", expectKey: "name", err: true},
		{input: "name(has-dash)", expectKey: "name", err: true},
		{input: "name(multiple, args)", expectKey: "name", err: true},
		{input: "name(noClosingParen", expectKey: "name", err: true},
	}
	for _, tc := range cases {
		parsed, err := parseTagKey(tc.input)

		if err != nil && tc.err == false {
			t.Errorf("[%q]: expected success, got: %v", tc.input, err)
			continue
		}
		if err == nil {
			if tc.err == true {
				t.Errorf("[%q]: expected failure, got: %q", tc.input, parsed.name)
				continue
			}
			if parsed.name != tc.expectKey {
				t.Errorf("[%q]\nexpected key: %q, got: %q", tc.input, tc.expectKey, parsed.name)
			}
			if len(parsed.args) != len(tc.expectArgs) {
				t.Errorf("[%q]: expected %d args, got: %q", tc.input, len(tc.expectArgs), parsed.args)
				continue
			}
			for i := range tc.expectArgs {
				if want, got := tc.expectArgs[i], parsed.args[i]; fmt.Sprintf("%v", got.Value) != want {
					t.Errorf("[%q]\nexpected %q, got %q", tc.input, want, got)
				}
			}
		}
	}
}
