/*
Copyright 2020 The Kubernetes Authors.

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

package marker

import (
	"reflect"
	"testing"
)

func getPointerFromString(s string) *string {
	return &s
}

func Test_Marker(t *testing.T) {
	testcases := []struct {
		name string
		in   Defaulted
		out  Defaulted
	}{
		{
			name: "default",
			in:   Defaulted{},
			out: Defaulted{
				StringDefault:      "bar",
				StringEmptyDefault: "",
				StringEmpty:        "",
				StringPointer:      getPointerFromString("default"),
				IntDefault:         1,
				IntEmptyDefault:    0,
				IntEmpty:           0,
				FloatDefault:       0.5,
				FloatEmptyDefault:  0.0,
				FloatEmpty:         0.0,
				List: []Item{
					getPointerFromString("foo"),
					getPointerFromString("bar"),
				},
				Sub: &SubStruct{
					S: "foo",
					I: 5,
				},
				OtherSub: SubStruct{
					S: "",
					I: 1,
				},
				StructList: []SubStruct{
					{
						S: "foo1",
						I: 1,
					},
					{
						S: "foo2",
						I: 1,
					},
				},
				PtrStructList: []*SubStruct{
					{
						S: "foo1",
						I: 1,
					},
					{
						S: "foo2",
						I: 1,
					},
				},
				StringList: []string{
					"foo",
				},
				Map: map[string]Item{
					"foo": getPointerFromString("bar"),
				},
				StructMap: map[string]SubStruct{
					"foo": SubStruct{
						S: "string",
						I: 1,
					},
				},
				PtrStructMap: map[string]*SubStruct{
					"foo": &SubStruct{
						S: "string",
						I: 1,
					},
				},
				AliasPtr: getPointerFromString("banana"),
			},
		},
		{
			name: "values-omitempty",
			in: Defaulted{
				StringDefault: "changed",
				IntDefault:    5,
			},
			out: Defaulted{
				StringDefault:      "changed",
				StringEmptyDefault: "",
				StringEmpty:        "",
				StringPointer:      getPointerFromString("default"),
				IntDefault:         5,
				IntEmptyDefault:    0,
				IntEmpty:           0,
				FloatDefault:       0.5,
				FloatEmptyDefault:  0.0,
				FloatEmpty:         0.0,
				List: []Item{
					getPointerFromString("foo"),
					getPointerFromString("bar"),
				},
				Sub: &SubStruct{
					S: "foo",
					I: 5,
				},
				StructList: []SubStruct{
					{
						S: "foo1",
						I: 1,
					},
					{
						S: "foo2",
						I: 1,
					},
				},
				PtrStructList: []*SubStruct{
					{
						S: "foo1",
						I: 1,
					},
					{
						S: "foo2",
						I: 1,
					},
				},
				StringList: []string{
					"foo",
				},
				OtherSub: SubStruct{
					S: "",
					I: 1,
				},
				Map: map[string]Item{
					"foo": getPointerFromString("bar"),
				},
				StructMap: map[string]SubStruct{
					"foo": SubStruct{
						S: "string",
						I: 1,
					},
				},
				PtrStructMap: map[string]*SubStruct{
					"foo": &SubStruct{
						S: "string",
						I: 1,
					},
				},
				AliasPtr: getPointerFromString("banana"),
			},
		},
		{
			name: "lists",
			in: Defaulted{
				List: []Item{
					nil,
					getPointerFromString("bar"),
				},
			},
			out: Defaulted{
				StringDefault:      "bar",
				StringEmptyDefault: "",
				StringEmpty:        "",
				StringPointer:      getPointerFromString("default"),
				IntDefault:         1,
				IntEmptyDefault:    0,
				IntEmpty:           0,
				FloatDefault:       0.5,
				FloatEmptyDefault:  0.0,
				FloatEmpty:         0.0,
				List: []Item{
					getPointerFromString("apple"),
					getPointerFromString("bar"),
				},
				Sub: &SubStruct{
					S: "foo",
					I: 5,
				},
				StructList: []SubStruct{
					{
						S: "foo1",
						I: 1,
					},
					{
						S: "foo2",
						I: 1,
					},
				},
				PtrStructList: []*SubStruct{
					{
						S: "foo1",
						I: 1,
					},
					{
						S: "foo2",
						I: 1,
					},
				},
				StringList: []string{
					"foo",
				},
				OtherSub: SubStruct{
					S: "",
					I: 1,
				},
				Map: map[string]Item{
					"foo": getPointerFromString("bar"),
				},
				StructMap: map[string]SubStruct{
					"foo": SubStruct{
						S: "string",
						I: 1,
					},
				},
				PtrStructMap: map[string]*SubStruct{
					"foo": &SubStruct{
						S: "string",
						I: 1,
					},
				},
				AliasPtr: getPointerFromString("banana"),
			},
		},
		{
			name: "stringmap",
			in: Defaulted{
				Map: map[string]Item{
					"foo": nil,
					"bar": getPointerFromString("banana"),
				},
			},
			out: Defaulted{
				StringDefault:      "bar",
				StringEmptyDefault: "",
				StringEmpty:        "",
				StringPointer:      getPointerFromString("default"),
				IntDefault:         1,
				IntEmptyDefault:    0,
				IntEmpty:           0,
				FloatDefault:       0.5,
				FloatEmptyDefault:  0.0,
				FloatEmpty:         0.0,
				List: []Item{
					getPointerFromString("foo"),
					getPointerFromString("bar"),
				},
				Sub: &SubStruct{
					S: "foo",
					I: 5,
				},
				StructList: []SubStruct{
					{
						S: "foo1",
						I: 1,
					},
					{
						S: "foo2",
						I: 1,
					},
				},
				PtrStructList: []*SubStruct{
					{
						S: "foo1",
						I: 1,
					},
					{
						S: "foo2",
						I: 1,
					},
				},
				StringList: []string{
					"foo",
				},
				OtherSub: SubStruct{
					S: "",
					I: 1,
				},
				Map: map[string]Item{
					"foo": getPointerFromString("apple"),
					"bar": getPointerFromString("banana"),
				},
				StructMap: map[string]SubStruct{
					"foo": SubStruct{
						S: "string",
						I: 1,
					},
				},
				PtrStructMap: map[string]*SubStruct{
					"foo": &SubStruct{
						S: "string",
						I: 1,
					},
				},
				AliasPtr: getPointerFromString("banana"),
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			SetObjectDefaults_Defaulted(&tc.in)
			if !reflect.DeepEqual(tc.in, tc.out) {
				t.Errorf("Error: Expected and actual output are different \n actual: %+v\n expected: %+v\n", tc.in, tc.out)
			}
		})
	}
}

func Test_DefaultingFunction(t *testing.T) {
	in := DefaultedWithFunction{}
	SetObjectDefaults_DefaultedWithFunction(&in)
	out := DefaultedWithFunction{
		S1: "default_function",
		S2: "default_marker",
	}
	if !reflect.DeepEqual(in, out) {
		t.Errorf("Error: Expected and actual output are different \n actual: %+v\n expected: %+v\n", in, out)
	}

}
