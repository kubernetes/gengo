/*
Copyright 2017 The Kubernetes Authors.

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

package generator

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"k8s.io/gengo/namer"
	"k8s.io/gengo/parser"
	"k8s.io/gengo/types"
)

func TestExecuteCanonicalImport(t *testing.T) {
	// TODO: Generalize for more package customizations.

	tempdir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempdir)

	pkgPath := "k8s.io/bar/foo"
	fileName := "doc"

	p := DefaultPackage{
		PackageName: "foo",
		PackagePath: pkgPath,
		GeneratorList: []Generator{
			DefaultGen{OptionalName: fileName},
		},
	}

	// Create a bare bones context.
	it := namer.NewDefaultImportTracker(types.Name{})
	c, err := NewContext(
		parser.New(), namer.NameSystems{"": namer.NewRawNamer(pkgPath, &it)}, "",
	)
	if err != nil {
		t.Fatalf("failed to create new context: %v", err)
	}

	if err := c.ExecutePackage(tempdir, &p); err != nil {
		t.Fatalf("failed to execute template: %v", err)
	}

	// Execute appends ".go" to filename
	genFilePath := filepath.Join(tempdir, pkgPath, fileName) + ".go"
	data, err := ioutil.ReadFile(genFilePath)
	if err != nil {
		t.Fatalf("failed to read generated file %s: %v", genFilePath, err)
	}

	want := `package foo // import "k8s.io/bar/foo"
` // expect a newline.

	got := string(data)
	if got != want {
		t.Fatalf("expected generated file:\n%s\ngot:\n%s", want, got)
	}
}
