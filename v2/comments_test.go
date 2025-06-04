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
		"+FALSE=false # comment",
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

func TestExtractFunctionStyleCommentTags(t *testing.T) {
	mktags := func(t ...Tag) []Tag { return t }
	mkstrs := func(s ...string) []string { return s }

	cases := []struct {
		name        string
		comments    []string
		tagNames    []string
		parseValues bool
		expectError bool
		expect      map[string][]Tag
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
			"+k8s:union",
		},
		tagNames: []string{"pfx1Foo", "pfx2Foo", "k8s:union"},
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
	}, {
		name: "ParseValues - valid values",
		comments: []string{
			"+boolTag=true",
			"+intTag=42",
			"+stringTag=\"quoted string\"",
			"+rawStringTag=`raw string`",
			"+identTag=identifier",
		},
		parseValues: true,
		expect: map[string][]Tag{
			"boolTag":      mktags(Tag{"boolTag", nil, "true"}),
			"intTag":       mktags(Tag{"intTag", nil, "42"}),
			"stringTag":    mktags(Tag{"stringTag", nil, "quoted string"}),
			"rawStringTag": mktags(Tag{"rawStringTag", nil, "raw string"}),
			"identTag":     mktags(Tag{"identTag", nil, "identifier"}),
		},
	}, {
		name: "ParseValues - comments ignored",
		comments: []string{
			"+boolTag=true # this is a boolean",
			"+intTag=42 # this is an integer",
			"+stringTag=\"quoted string\" # this is a string",
		},
		parseValues: true,
		expect: map[string][]Tag{
			"boolTag":   mktags(Tag{"boolTag", nil, "true"}),
			"intTag":    mktags(Tag{"intTag", nil, "42"}),
			"stringTag": mktags(Tag{"stringTag", nil, "quoted string"}),
		},
	}, {
		name: "ParseValues enabled - invalid value",
		comments: []string{
			"+invalidTag=\"unclosed string",
		},
		parseValues: true,
		expectError: true,
	}, {
		name: "raw values with comments",
		comments: []string{
			"+boolTag=true // this is a boolean",
			"+intTag=42 // this is an integer",
			"+stringTag=\"quoted string\" // this is a string",
			"+invalidTag=\"unclosed string",
		},
		expect: map[string][]Tag{
			"boolTag":    mktags(Tag{"boolTag", nil, "true // this is a boolean"}),
			"intTag":     mktags(Tag{"intTag", nil, "42 // this is an integer"}),
			"stringTag":  mktags(Tag{"stringTag", nil, "\"quoted string\" // this is a string"}),
			"invalidTag": mktags(Tag{"invalidTag", nil, "\"unclosed string"}),
		},
	}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ExtractFunctionStyleCommentTags("+", tc.tagNames, tc.comments, ParseValues(tc.parseValues))

			if tc.expectError {
				if err == nil {
					t.Errorf("case %q: expected error but got none", tc.name)
				}
				return
			}

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
