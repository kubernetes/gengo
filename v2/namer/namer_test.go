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

package namer

import (
	"go/token"
	"reflect"
	"testing"

	"k8s.io/gengo/v2/types"
)

func TestNameStrategy(t *testing.T) {
	u := types.Universe{}

	// Add some types.
	base := u.Type(types.Name{Package: "foo/bar", Name: "Baz"})
	base.Kind = types.Struct

	tmp := u.Type(types.Name{Package: "", Name: "[]bar.Baz"})
	tmp.Kind = types.Slice
	tmp.Elem = base

	tmp = u.Type(types.Name{Package: "", Name: "map[string]bar.Baz"})
	tmp.Kind = types.Map
	tmp.Key = types.String
	tmp.Elem = base

	tmp = u.Type(types.Name{Package: "foo/other", Name: "Baz"})
	tmp.Kind = types.Struct
	tmp.Members = []types.Member{{
		Embedded: true,
		Type:     base,
	}}

	tmp = u.Type(types.Name{Package: "", Name: "chan Baz"})
	tmp.Kind = types.Chan
	tmp.Elem = base

	tmp = u.Type(types.Name{Package: "", Name: "[4]Baz"})
	tmp.Kind = types.Array
	tmp.Elem = base
	tmp.Len = 4

	u.Type(types.Name{Package: "", Name: "string"})

	o := Orderer{NewPublicNamer(0)}
	order := o.OrderUniverse(u)
	orderedNames := make([]string, len(order))
	for i, t := range order {
		orderedNames[i] = o.Name(t)
	}
	expect := []string{"Array4Baz", "Baz", "Baz", "ChanBaz", "MapStringToBaz", "SliceBaz", "String"}
	if e, a := expect, orderedNames; !reflect.DeepEqual(e, a) {
		t.Errorf("Wanted %#v, got %#v", e, a)
	}

	o = Orderer{NewRawNamer("my/package", nil)}
	order = o.OrderUniverse(u)
	orderedNames = make([]string, len(order))
	for i, t := range order {
		orderedNames[i] = o.Name(t)
	}

	expect = []string{"[4]bar.Baz", "[]bar.Baz", "bar.Baz", "chan bar.Baz", "map[string]bar.Baz", "other.Baz", "string"}
	if e, a := expect, orderedNames; !reflect.DeepEqual(e, a) {
		t.Errorf("Wanted %#v, got %#v", e, a)
	}

	o = Orderer{NewRawNamer("foo/bar", nil)}
	order = o.OrderUniverse(u)
	orderedNames = make([]string, len(order))
	for i, t := range order {
		orderedNames[i] = o.Name(t)
	}

	expect = []string{"Baz", "[4]Baz", "[]Baz", "chan Baz", "map[string]Baz", "other.Baz", "string"}
	if e, a := expect, orderedNames; !reflect.DeepEqual(e, a) {
		t.Errorf("Wanted %#v, got %#v", e, a)
	}

	o = Orderer{NewPublicNamer(1)}
	order = o.OrderUniverse(u)
	orderedNames = make([]string, len(order))
	for i, t := range order {
		orderedNames[i] = o.Name(t)
	}
	expect = []string{"Array4BarBaz", "BarBaz", "ChanBarBaz", "MapStringToBarBaz", "OtherBaz", "SliceBarBaz", "String"}
	if e, a := expect, orderedNames; !reflect.DeepEqual(e, a) {
		t.Errorf("Wanted %#v, got %#v", e, a)
	}
}

// goKeywords is every keyword in the Go spec. A named type whose lowercased
// name is one of these breaks generators which emit the private name as an
// identifier.
var goKeywords = []string{
	"break", "case", "chan", "const", "continue", "default", "defer", "else",
	"fallthrough", "for", "func", "go", "goto", "if", "import", "interface",
	"map", "package", "range", "return", "select", "struct", "switch", "type",
	"var",
}

func TestNameStrategyKeywordNamedTypes(t *testing.T) {
	private := NewPrivateNamer(0)
	public := NewPublicNamer(0)

	for _, kw := range goKeywords {
		u := types.Universe{}
		typeName := IC(kw)
		typ := u.Type(types.Name{Package: "example.com/api/core/v1", Name: typeName})
		typ.Kind = types.Struct

		got := private.Name(typ)
		if e, a := "_"+kw, got; e != a {
			t.Errorf("private name of type %q: wanted %q, got %q", typeName, e, a)
		}
		if token.IsKeyword(got) {
			t.Errorf("private name of type %q is the keyword %q", typeName, got)
		}
		if !token.IsIdentifier(got) {
			t.Errorf("private name of type %q is not a legal identifier: %q", typeName, got)
		}

		// The public namer capitalizes the first character, so it can never
		// land on a keyword and must be left alone.
		if e, a := typeName, public.Name(typ); e != a {
			t.Errorf("public name of type %q: wanted %q, got %q", typeName, e, a)
		}
	}
}

func TestNameStrategyKeywordAnonymousTypes(t *testing.T) {
	u := types.Universe{}

	emptyInterface := u.Type(types.Name{Name: "interface{}"})
	emptyInterface.Kind = types.Interface

	emptyStruct := u.Type(types.Name{Name: "struct{}"})
	emptyStruct.Kind = types.Struct

	pkg := u.Type(types.Name{Package: "example.com/api/core/v1", Name: "Package"})
	pkg.Kind = types.Struct

	// A composite type embeds the name of its element, which is itself
	// keyword-escaped.
	slice := u.Type(types.Name{Name: "[]v1.Package"})
	slice.Kind = types.Slice
	slice.Elem = pkg

	// A name which merely contains a keyword is not a keyword and is untouched.
	iface := u.Type(types.Name{Package: "example.com/api/core/v1", Name: "Interfaces"})
	iface.Kind = types.Struct

	private := NewPrivateNamer(0)
	for _, tc := range []struct {
		typ    *types.Type
		expect string
	}{
		{emptyInterface, "_interface"},
		{emptyStruct, "_struct"},
		{slice, "slice_package"},
		{iface, "interfaces"},
	} {
		got := private.Name(tc.typ)
		if e, a := tc.expect, got; e != a {
			t.Errorf("private name of %q: wanted %q, got %q", tc.typ.Name, e, a)
		}
		if !token.IsIdentifier(got) || token.IsKeyword(got) {
			t.Errorf("private name of %q is not a legal identifier: %q", tc.typ.Name, got)
		}
	}
}
