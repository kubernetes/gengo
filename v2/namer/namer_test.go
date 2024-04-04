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

// NOTE: this test is intended to demostrate specific behavior
// described in https://github.com/kubernetes/kubernetes/issues/124192
// and should not be considered a good usage scenario of Namer
//
// Here we are trying to simulate the behavior when Namer is
// configured with invalid local package name and import tracker
// has valid configuration. In this case Namer shouldn't produce
// invalid syntax like '.<variableName>' and instead treat empty
// names from import tracker as types within the same package
func TestRawNamerEmptyImportInteraction(t *testing.T) {
	tracker := NewDefaultImportTracker(types.Name{
		Package: "resource.io/pkg",
		Name:    "name",
	})
	namer := NewRawNamer("resource/pkg", &tracker)

	gotName := namer.Name(&types.Type{
		Name: types.Name{
			Package: "resource.io/pkg",
			Name:    "flag",
			Path:    "/path/resource.io/pkg",
		},
	})

	if gotName != "flag" {
		t.Errorf("Wanted %#v, got %#v", "flag", gotName)
	}
}
