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

package generator_test

import (
	"bytes"
	"strings"
	"testing"

	"k8s.io/gengo/v2/generator"
	"k8s.io/gengo/v2/namer"
	"k8s.io/gengo/v2/parser"
)

func construct(t *testing.T, patterns ...string) *generator.Context {
	b := parser.New()

	if err := b.LoadPackages(patterns...); err != nil {
		t.Fatalf("unexpected error: %v", err)
		return nil
	}
	c, err := generator.NewContext(b, namer.NameSystems{
		"public":  namer.NewPublicNamer(0),
		"private": namer.NewPrivateNamer(0),
	}, "public")
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestSnippetWriter(t *testing.T) {
	c := construct(t, "./testdata/snippet_writer")
	b := &bytes.Buffer{}
	err := generator.NewSnippetWriter(b, c, "$", "$").
		Do("$.|public$$.|private$", c.Order[0]).
		Error()
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	if e, a := "Blahblah", b.String(); e != a {
		t.Errorf("Expected %q, got %q", e, a)
	}

	err = generator.NewSnippetWriter(b, c, "$", "$").
		Do("$.|public", c.Order[0]).
		Error()
	if err == nil {
		t.Errorf("expected error on invalid template")
	} else {
		// Dear reader, I apologize for making the worst change
		// detection test in the history of ever.
		if e, a := "snippet_writer_test.go", err.Error(); !strings.Contains(a, e) {
			t.Errorf("Expected %q but didn't find it in %q", e, a)
		}
	}
}

func TestArgsWith(t *testing.T) {
	orig := generator.Args{
		"a": "aaa",
		"b": "bbb",
	}

	withC := orig.With("c", "ccc")
	if len(orig) != 2 {
		t.Errorf("unexpected change to 'orig': %v", orig)
	}
	if len(withC) != 3 {
		t.Errorf("expected 'withC' to have 3 values: %v", withC)
	}

	withDE := withC.WithArgs(generator.Args{
		"d": "ddd",
		"e": "eee",
	})
	if len(orig) != 2 {
		t.Errorf("unexpected change to 'orig': %v", orig)
	}
	if len(withC) != 3 {
		t.Errorf("unexpected change to 'withC': %v", orig)
	}
	if len(withDE) != 5 {
		t.Errorf("expected 'withDE' to have 5 values: %v", withC)
	}

	withNewA := orig.With("a", "AAA")
	if orig["a"] != "aaa" {
		t.Errorf("unexpected change to 'orig': %v", orig)
	}
	if withNewA["a"] != "AAA" {
		t.Errorf("expected 'withNewA[\"a\"]' to be \"AAA\": %v", withNewA)
	}

	withNewAB := orig.WithArgs(generator.Args{
		"a": "AAA",
		"b": "BBB",
	})
	if orig["a"] != "aaa" || orig["b"] != "bbb" {
		t.Errorf("unexpected change to 'orig': %v", orig)
	}
	if withNewAB["a"] != "AAA" {
		t.Errorf("expected 'withNewAB[\"a\"]' to be \"AAA\": %v", withNewAB)
	}
	if withNewAB["b"] != "BBB" {
		t.Errorf("expected 'withNewAB[\"b\"]' to be \"BBB\": %v", withNewAB)
	}
}
