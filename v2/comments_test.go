/*
Copyright 2015 The Kubernetes Authors.

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

package gengo

import (
	"reflect"
	"strings"
	"testing"
	"text/scanner"

	"github.com/google/go-cmp/cmp"
)

func TestExtractCommentTags(t *testing.T) {
	commentLines := []string{
		"Human comment that is ignored.",
		"+foo=value1",
		"+bar",
		"+foo=value2",
		"+baz=qux,zrb=true",
		"+bip=\"value3\"",
	}

	a := ExtractCommentTags("+", commentLines)
	e := map[string][]string{
		"foo": {"value1", "value2"},
		"bar": {""},
		"baz": {"qux,zrb=true"},
		"bip": {`"value3"`},
	}
	if !reflect.DeepEqual(e, a) {
		t.Errorf("Wrong result:\n%v", cmp.Diff(e, a))
	}
}

func TestExtractSingleBoolCommentTag(t *testing.T) {
	commentLines := []string{
		"Human comment that is ignored.",
		"+TRUE=true",
		"+FALSE=false",
		"+MULTI=true",
		"+MULTI=false",
		"+MULTI=multi",
		"+NOTBOOL=blue",
		"+EMPTY",
	}

	testCases := []struct {
		key string
		def bool
		exp bool
		err string // if set, ignore exp.
	}{
		{"TRUE", false, true, ""},
		{"FALSE", true, false, ""},
		{"MULTI", false, true, ""},
		{"NOTBOOL", false, true, "is not boolean"},
		{"EMPTY", false, true, "is not boolean"},
		{"ABSENT", true, true, ""},
		{"ABSENT", false, false, ""},
	}

	for i, tc := range testCases {
		v, err := ExtractSingleBoolCommentTag("+", tc.key, tc.def, commentLines)
		if err != nil && tc.err == "" {
			t.Errorf("[%d]: unexpected failure: %v", i, err)
		} else if err == nil && tc.err != "" {
			t.Errorf("[%d]: expected failure: %v", i, tc.err)
		} else if err != nil {
			if !strings.Contains(err.Error(), tc.err) {
				t.Errorf("[%d]: unexpected error: expected %q, got %q", i, tc.err, err)
			}
		} else if v != tc.exp {
			t.Errorf("[%d]: unexpected value: expected %t, got %t", i, tc.exp, v)
		}
	}
}

func TestExtractExtendedCommentTags(t *testing.T) {
	mktags := func(t ...Tag) []Tag { return t }
	mkstrs := func(s ...string) []string { return s }

	cases := []struct {
		name     string
		comments []string
		prefixes []string
		expect   map[string][]Tag
	}{{
		name: "no args",
		comments: []string{
			"Human comment that is ignored",
			"+simpleNoVal",
			"+simpleWithVal=val",
			"+duplicateNoVal",
			"+duplicateNoVal",
			"+duplicateWithVal=val1",
			"+duplicateWithVal=val2",
		},
		expect: map[string][]Tag{
			"simpleNoVal":   mktags(Tag{"simpleNoVal", nil, ""}),
			"simpleWithVal": mktags(Tag{"simpleWithVal", nil, "val"}),
			"duplicateNoVal": mktags(
				Tag{"duplicateNoVal", nil, ""},
				Tag{"duplicateNoVal", nil, ""}),
			"duplicateWithVal": mktags(
				Tag{"duplicateWithVal", nil, "val1"},
				Tag{"duplicateWithVal", nil, "val2"}),
		},
	}, {
		name: "empty parens",
		comments: []string{
			"Human comment that is ignored",
			"+simpleNoVal()",
			"+simpleWithVal()=val",
			"+duplicateNoVal()",
			"+duplicateNoVal()",
			"+duplicateWithVal()=val1",
			"+duplicateWithVal()=val2",
		},
		expect: map[string][]Tag{
			"simpleNoVal":   mktags(Tag{"simpleNoVal", nil, ""}),
			"simpleWithVal": mktags(Tag{"simpleWithVal", nil, "val"}),
			"duplicateNoVal": mktags(
				Tag{"duplicateNoVal", nil, ""},
				Tag{"duplicateNoVal", nil, ""}),
			"duplicateWithVal": mktags(
				Tag{"duplicateWithVal", nil, "val1"},
				Tag{"duplicateWithVal", nil, "val2"}),
		},
	}, {
		name: "mixed no args and empty parens",
		comments: []string{
			"Human comment that is ignored",
			"+noVal",
			"+withVal=val1",
			"+noVal()",
			"+withVal()=val2",
		},
		expect: map[string][]Tag{
			"noVal": mktags(
				Tag{"noVal", nil, ""},
				Tag{"noVal", nil, ""}),
			"withVal": mktags(
				Tag{"withVal", nil, "val1"},
				Tag{"withVal", nil, "val2"}),
		},
	}, {
		name: "with args",
		comments: []string{
			"Human comment that is ignored",
			"+simpleNoVal(arg)",
			"+simpleWithVal(arg)=val",
			"+duplicateNoVal(arg1)",
			"+duplicateNoVal(arg2)",
			"+duplicateWithVal(arg1)=val1",
			"+duplicateWithVal(arg2)=val2",
		},
		expect: map[string][]Tag{
			"simpleNoVal":   mktags(Tag{"simpleNoVal", mkstrs("arg"), ""}),
			"simpleWithVal": mktags(Tag{"simpleWithVal", mkstrs("arg"), "val"}),
			"duplicateNoVal": mktags(
				Tag{"duplicateNoVal", mkstrs("arg1"), ""},
				Tag{"duplicateNoVal", mkstrs("arg2"), ""}),
			"duplicateWithVal": mktags(
				Tag{"duplicateWithVal", mkstrs("arg1"), "val1"},
				Tag{"duplicateWithVal", mkstrs("arg2"), "val2"}),
		},
	}, {
		name: "mixed no args and empty parens",
		comments: []string{
			"Human comment that is ignored",
			"+noVal",
			"+withVal=val1",
			"+noVal(arg)",
			"+withVal(arg)=val2",
		},
		expect: map[string][]Tag{
			"noVal": mktags(
				Tag{"noVal", nil, ""},
				Tag{"noVal", mkstrs("arg"), ""}),
			"withVal": mktags(
				Tag{"withVal", nil, "val1"},
				Tag{"withVal", mkstrs("arg"), "val2"}),
		},
	}, {
		name: "prefixes",
		comments: []string{
			"Human comment that is ignored",
			"+pfx1Foo",
			"+pfx2Foo=val1",
			"+pfx3Bar",
			"+pfx4Bar=val",
			"+pfx1Foo(arg)",
			"+pfx2Foo(arg)=val2",
			"+pfx3Bar(arg)",
			"+pfx4Bar(arg)=val",
		},
		prefixes: []string{"pfx1Foo", "pfx2Foo"},
		expect: map[string][]Tag{
			"pfx1Foo": mktags(
				Tag{"pfx1Foo", nil, ""},
				Tag{"pfx1Foo", mkstrs("arg"), ""}),
			"pfx2Foo": mktags(
				Tag{"pfx2Foo", nil, "val1"},
				Tag{"pfx2Foo", mkstrs("arg"), "val2"}),
		},
	}, {
		name: "raw arg with =, ), and space",
		comments: []string{
			"+rawEq(`a=b c=d )`)=xyz",
		},
		expect: map[string][]Tag{
			"rawEq": mktags(Tag{"rawEq", mkstrs("`a=b c=d )`"), "xyz"}),
		},
	}, {
		name: "raw arg no value",
		comments: []string{
			"+onlyRaw(`zzz`)",
		},
		expect: map[string][]Tag{
			"onlyRaw": mktags(Tag{"onlyRaw", mkstrs("`zzz`"), ""}),
		},
	}, {
		name: "raw string arg complex",
		comments: []string{
			"+rawTag(`[self.foo==10, ()), {}}, \"foo\", 'foo']`)=val",
		},
		expect: map[string][]Tag{
			"rawTag": mktags(
				Tag{"rawTag", mkstrs("`[self.foo==10, ()), {}}, \"foo\", 'foo']`"), "val"}),
		},
	}, {
		name: "JSON object",
		comments: []string{
			`+rawTag({"a": 1, "b": "x"})=val`,
		},
		expect: map[string][]Tag{
			"rawTag": mktags(
				Tag{"rawTag", mkstrs(`{"a": 1, "b": "x"}`), "val"}),
		},
	}}

	for _, tc := range cases {
		result, _ := ExtractFunctionStyleCommentTags("+", tc.prefixes, tc.comments)
		if !reflect.DeepEqual(result, tc.expect) {
			t.Errorf("case %q: wrong result:\n%v", tc.name, cmp.Diff(tc.expect, result))
		}
	}
}

func TestParseTagKey(t *testing.T) {
	mkss := func(s ...string) []string { return s }

	cases := []struct {
		input      string
		expectKey  string
		expectArgs []string
		err        bool
	}{
		{"simple", "simple", nil, false},
		{"parens()", "parens", nil, false},
		{"withArgLower(arg)", "withArgLower", mkss("arg"), false},
		{"withArgUpper(ARG)", "withArgUpper", mkss("ARG"), false},
		{"withArgMixed(ArG)", "withArgMixed", mkss("ArG"), false},
		{"name({})", "name", mkss("{}"), false},
		{`name({"a":1})`, "name", mkss(`{"a":1}`), false},
		{"name([])", "name", mkss("[]"), false},
		{"withArgs(arg1, arg2)", "", nil, true},
		{"argWithDash(arg-name) ", "", nil, true},
		{"withRaw(`a = b`)", "withRaw", mkss("`a = b`"), false},
		{"badRaw(missing`)", "", nil, true},
		{"badMix(arg,`raw`)", "", nil, true},
	}
	for _, tc := range cases {
		s := initTagKeyScanner(tc.input)
		key, args, err := s.parseTagKey(nil)
		if err != nil && tc.err == false {
			t.Errorf("[%q]: expected success, got: %v", tc.input, err)
			continue
		}
		if err == nil {
			if tc.err == true {
				t.Errorf("[%q]: expected failure, got: %v(%v)", tc.input, key, args)
				continue
			}
			if key != tc.expectKey {
				t.Errorf("[%q]\nexpected key: %q, got: %q", tc.input, tc.expectKey, key)
			}
			if len(args) != len(tc.expectArgs) {
				t.Errorf("[%q]: expected %d args, got: %q", tc.input, len(tc.expectArgs), args)
				continue
			}
			for i := range tc.expectArgs {
				if want, got := tc.expectArgs[i], args[i]; got != want {
					t.Errorf("[%q]\nexpected %q, got %q", tc.input, want, got)
				}
			}
		}
	}
}

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
		{input: "nameNoMatch", expectKey: ""},
		{input: "nameNoMatch()", expectKey: ""},
		{input: "nameNoMatch(arg)", expectKey: ""},

		{input: "name()", expectKey: "name"},
		{input: "name(lower)", expectKey: "name", expectArgs: mkss("lower")},
		{input: "name(CAPITAL)", expectKey: "name", expectArgs: mkss("CAPITAL")},
		{input: "name(MiXeD)", expectKey: "name", expectArgs: mkss("MiXeD")},
		{input: "name(mIxEd)", expectKey: "name", expectArgs: mkss("mIxEd")},
		{input: "name(_under)", expectKey: "name", expectArgs: mkss("_under")},
		{input: `name("hasQuotes")`, expectKey: "name", expectArgs: mkss("\"hasQuotes\"")},
		{input: "name(`hasRawQuotes`)", expectKey: "name", expectArgs: mkss("`hasRawQuotes`")},
		{input: "name(has space)", expectKey: "name", err: true},
		{input: "name(has-dash)", expectKey: "name", err: true},
		{input: "name(multiple, args)", expectKey: "name", err: true},
		{input: "name(noClosingParen", expectKey: "name", err: true},
	}
	for _, tc := range cases {
		s := initTagKeyScanner(tc.input)
		key, args, err := s.parseTagKey([]string{"name"})

		if err != nil && tc.err == false {
			t.Errorf("[%q]: expected success, got: %v", tc.input, err)
			continue
		}
		if err == nil {
			if tc.err == true {
				t.Errorf("[%q]: expected failure, got: %q", tc.input, key)
				continue
			}
			if key != tc.expectKey {
				t.Errorf("[%q]\nexpected key: %q, got: %q", tc.input, tc.expectKey, key)
			}
			if len(args) != len(tc.expectArgs) {
				t.Errorf("[%q]: expected %d args, got: %q", tc.input, len(tc.expectArgs), args)
				continue
			}
			for i := range tc.expectArgs {
				if want, got := tc.expectArgs[i], args[i]; got != want {
					t.Errorf("[%q]\nexpected %q, got %q", tc.input, want, got)
				}
			}
		}
	}
}

func TestParseJSON(t *testing.T) {
	cases := []struct {
		input      string
		err        bool
		incomplete bool
	}{
		{
			input: `[]`,
		},
		{
			input: `{}`,
		},
		{
			input: `[1]`,
		},
		{
			input: `{"a":1}`,
		},
		{
			input: `[1, 2]`,
		},
		{
			input: `{"a": 1, "b": 2}`,
		},
		{
			input: `1.1`,
		},
		{
			input: `-4`,
		},
		{
			input: `true`,
		},
		{
			input: `false`,
		},
		{
			input: `"string"`,
		},
		{
			input: "null",
		},
		{
			input: `{"key":"value" }`,
		},
		{
			input: `[1 ]`,
		},
		{
			input: `[1 ,2]`,
		},

		// invalid
		{
			input: `[1,]`,
			err:   true,
		},
		{
			input: `[1,]`,
			err:   true,
		},
		{
			input: `{"a":1,}`,
			err:   true,
		},
		{
			input: `{"a"`,
			err:   true,
		},
		{
			input: `"a`,
			err:   true,
		},
		{
			input: `UNKNOWN`,
			err:   true,
		},
		{
			input: `1.4e-10`, // parse consumes 1.4, not the full number
		},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			s := initTagKeyScanner(tc.input)
			out, err := s.scanJSONFlavoredValue()
			if err != nil && tc.err == false {
				t.Errorf("[%q]: expected success, got: %v", tc.input, err)
				return
			}
			if err == nil {
				if tc.err == true {
					t.Errorf("[%q]: expected failure, got: %q", tc.input, out)
					return
				}
				if out != tc.input {
					t.Errorf("expected %q got %q", tc.input, out)
				}
			}

			gotIncomplete := s.Scan() != scanner.EOF
			if tc.incomplete != gotIncomplete {
				t.Errorf("Expected incomplete=%t but got %t", tc.incomplete, gotIncomplete)
			}
		})
	}
}
