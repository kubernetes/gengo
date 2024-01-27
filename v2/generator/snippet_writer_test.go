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
	b := parser.New(nil)

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
