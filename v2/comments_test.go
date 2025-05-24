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
			"+simpleNoVal  // trailing comment",
			"+simpleWithVal=val  // trailing comment",
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
			"+simpleNoVal()  // trailing comment",
			"+simpleWithVal()=val  // trailing comment",
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
			"+k8s:union",
		},
		prefixes: []string{"pfx1Foo", "pfx2Foo", "k8s:union"},
		expect: map[string][]Tag{
			"pfx1Foo": mktags(
				Tag{"pfx1Foo", nil, ""},
				Tag{"pfx1Foo", mkstrs("arg"), ""}),
			"pfx2Foo": mktags(
				Tag{"pfx2Foo", nil, "val1"},
				Tag{"pfx2Foo", mkstrs("arg"), "val2"}),
			"k8s:union": mktags(
				Tag{Name: "k8s:union"}),
		},
	}, {
		name: "raw arg with =, ), and space",
		comments: []string{
			"+rawEq(`a=b c=d )`)=xyz",
		},
		expect: map[string][]Tag{
			"rawEq": mktags(Tag{"rawEq", mkstrs("a=b c=d )"), "xyz"}),
		},
	}, {
		name: "raw arg no value",
		comments: []string{
			"+onlyRaw(`zzz`)",
		},
		expect: map[string][]Tag{
			"onlyRaw": mktags(Tag{"onlyRaw", mkstrs("zzz"), ""}),
		},
	}, {
		name: "raw string arg complex",
		comments: []string{
			"+rawTag(`[self.foo==10, ()), {}}, \"foo\", 'foo']`)=val",
		},
		expect: map[string][]Tag{
			"rawTag": mktags(
				Tag{"rawTag", mkstrs("[self.foo==10, ()), {}}, \"foo\", 'foo']"), "val"}),
		},
	}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ExtractFunctionStyleCommentTags("+", tc.prefixes, tc.comments)
			if err != nil {
				t.Errorf("case %q: unexpected error: %v", tc.name, err)
				return
			}
			if !reflect.DeepEqual(result, tc.expect) {
				t.Errorf("case %q: wrong result:\n%v", tc.name, cmp.Diff(tc.expect, result))
			}
		})
	}
}

func TestExtractFunctionStyleCommentTypedTags(t *testing.T) {
	mktags := func(t ...TypedTag) []TypedTag { return t }

	cases := []struct {
		name     string
		comments []string
		prefixes []string
		expect   map[string][]TypedTag
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
			expect: map[string][]TypedTag{
				"quoted": mktags(
					TypedTag{Name: "quoted", Args: []Arg{
						{Value: String("value")},
					}},
				),
				"backticked": mktags(
					TypedTag{Name: "backticked", Args: []Arg{
						{Value: String("value")},
					}},
				),
				"ident": mktags(
					TypedTag{Name: "ident", Args: []Arg{
						{Value: String("value")},
					}},
				),
				"integer": mktags(
					TypedTag{Name: "integer", Args: []Arg{
						{Value: Int(2)},
					}}),
				"negative": mktags(
					TypedTag{Name: "negative", Args: []Arg{
						{Value: Int(-5)},
					}}),
				"hex": mktags(
					TypedTag{Name: "hex", Args: []Arg{
						{Value: Int(0xFF00B3)},
					}}),
				"octal": mktags(
					TypedTag{Name: "octal", Args: []Arg{
						{Value: Int(0o04167)},
					}}),
				"binary": mktags(
					TypedTag{Name: "binary", Args: []Arg{
						{Value: Int(0b10101)},
					}}),
				"true": mktags(
					TypedTag{Name: "true", Args: []Arg{
						{Value: Bool(true)},
					}}),
				"false": mktags(
					TypedTag{Name: "false", Args: []Arg{
						{Value: Bool(false)},
					}}),
			},
		},
		{
			name: "named params",
			comments: []string{
				"+strings(q: \"value\", b: `value`, i: value)",
				"+numbers(n1: 2, n2: -5, n3: 0xFF00B3, n4: 0o04167, n5: 0b10101)",
				"+bools(t: true, f:false)",
			},
			expect: map[string][]TypedTag{
				"strings": mktags(
					TypedTag{Name: "strings", Args: []Arg{
						{Name: "q", Value: String("value")},
						{Name: "b", Value: String(`value`)},
						{Name: "i", Value: String("value")},
					}}),
				"numbers": mktags(
					TypedTag{Name: "numbers", Args: []Arg{
						{Name: "n1", Value: Int(2)},
						{Name: "n2", Value: Int(-5)},
						{Name: "n3", Value: Int(0xFF00B3)},
						{Name: "n4", Value: Int(0o04167)},
						{Name: "n5", Value: Int(0b10101)},
					}}),
				"bools": mktags(
					TypedTag{Name: "bools", Args: []Arg{
						{Name: "t", Value: Bool(true)},
						{Name: "f", Value: Bool(false)},
					}}),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ExtractFunctionStyleCommentTypedTags("+", tc.prefixes, tc.comments)
			if err != nil {
				t.Errorf("case %q: unexpected error: %v", tc.name, err)
				return
			}
			if !reflect.DeepEqual(result, tc.expect) {
				t.Errorf("case %q: wrong result:\n%v", tc.name, cmp.Diff(tc.expect, result))
			}
		})
	}
}

func TestParseTagKey(t *testing.T) {
	mkss := func(s ...string) []Arg {
		var args []Arg
		for _, v := range s {
			args = append(args, Arg{Value: String(v)})
		}
		return args
	}

	cases := []struct {
		input      string
		expectKey  string
		expectArgs []Arg
		err        bool
	}{
		{"simple", "simple", nil, false},
		{"parens()", "parens", nil, false},
		{"withArgLower(arg)", "withArgLower", mkss("arg"), false},
		{"withArgUpper(ARG)", "withArgUpper", mkss("ARG"), false},
		{"withArgMixed(ArG)", "withArgMixed", mkss("ArG"), false},
		{"withArgs(arg1, arg2)", "", nil, true},
		{"argWithDash(arg-name) ", "", nil, true},
		{"withRaw(`a = b`)", "withRaw", mkss("a = b"), false},
		{"badRaw(missing`)", "", nil, true},
		{"badMix(arg,`raw`)", "", nil, true},
		{`quoted(s: "value \" \\")`, "quoted", []Arg{
			{Name: "s", Value: String("value \" \\")},
		}, false},
		{"backticks(s: `value`)", "backticks", []Arg{
			{Name: "s", Value: String(`value`)},
		}, false},
		{"ident(k: value)", "ident", []Arg{
			{Name: "k", Value: String("value")},
		}, false},
		{"numbers(n1: 2, n2: -5, n3: 0xFF00B3, n4: 0o04167, n5: 0b10101)", "numbers", []Arg{
			{Name: "n1", Value: Int(2)},
			{Name: "n2", Value: Int(-5)},
			{Name: "n3", Value: Int(0xFF00B3)},
			{Name: "n4", Value: Int(0o04167)},
			{Name: "n5", Value: Int(0b10101)},
		}, false},
		{"bools(t: true, f:false)", "bools", []Arg{
			{Name: "t", Value: Bool(true)},
			{Name: "f", Value: Bool(false)},
		}, false},
		{"mixed(s: `value`, i: 2, b: true)", "mixed", []Arg{
			{Name: "s", Value: String("value")},
			{Name: "i", Value: Int(2)},
			{Name: "b", Value: Bool(true)},
		}, false},
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
				if want, got := tc.expectArgs[i], parsed.args[i]; got.Value.String() != want {
					t.Errorf("[%q]\nexpected %q, got %q", tc.input, want, got)
				}
			}
		}
	}
}
