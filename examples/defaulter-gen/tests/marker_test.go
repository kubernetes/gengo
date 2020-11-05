package main

import (
	"reflect"
	"testing"

	"k8s.io/gengo/examples/defaulter-gen/_output_tests/marker"
)

func getPointerFromString(s string) *string {
	return &s
}

func Test_Marker(t *testing.T) {
	testcases := []struct {
		name string
		in   marker.Defaulted
		out  marker.Defaulted
	}{
		{
			name: "default",
			in:   marker.Defaulted{},
			out: marker.Defaulted{
				Field:       "bar",
				OtherField:  0,
				EmptyInt:    0,
				EmptyString: "",
				List: []marker.Item{
					getPointerFromString("foo"),
					getPointerFromString("bar"),
				},
				Sub: &marker.SubStruct{
					S: "foo",
					I: 5,
				},
				OtherSub: marker.SubStruct{
					S: "",
					I: 1,
				},
				Map: map[string]marker.Item{
					"foo": getPointerFromString("bar"),
				},
			},
		},
		{
			name: "values-omitempty",
			in: marker.Defaulted{
				Field:      "changed",
				OtherField: 1,
			},
			out: marker.Defaulted{
				Field:       "changed",
				OtherField:  1,
				EmptyInt:    0,
				EmptyString: "",
				List: []marker.Item{
					getPointerFromString("foo"),
					getPointerFromString("bar"),
				},
				Sub: &marker.SubStruct{
					S: "foo",
					I: 5,
				},
				OtherSub: marker.SubStruct{
					S: "",
					I: 1,
				},
				Map: map[string]marker.Item{
					"foo": getPointerFromString("bar"),
				},
			},
		},
		{
			name: "lists",
			in: marker.Defaulted{
				List: []marker.Item{
					nil,
					getPointerFromString("bar"),
				},
			},
			out: marker.Defaulted{
				Field:       "bar",
				OtherField:  0,
				EmptyInt:    0,
				EmptyString: "",
				List: []marker.Item{
					getPointerFromString("apple"),
					getPointerFromString("bar"),
				},
				Sub: &marker.SubStruct{
					S: "foo",
					I: 5,
				},
				OtherSub: marker.SubStruct{
					S: "",
					I: 1,
				},
				Map: map[string]marker.Item{
					"foo": getPointerFromString("bar"),
				},
			},
		},
		{
			name: "stringmap",
			in: marker.Defaulted{
				Map: map[string]marker.Item{
					"foo": nil,
					"bar": getPointerFromString("banana"),
				},
			},
			out: marker.Defaulted{
				Field:       "bar",
				OtherField:  0,
				EmptyInt:    0,
				EmptyString: "",
				List: []marker.Item{
					getPointerFromString("foo"),
					getPointerFromString("bar"),
				},
				Sub: &marker.SubStruct{
					S: "foo",
					I: 5,
				},
				OtherSub: marker.SubStruct{
					S: "",
					I: 1,
				},
				Map: map[string]marker.Item{
					"foo": getPointerFromString("apple"),
					"bar": getPointerFromString("banana"),
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			marker.SetObjectDefaults_Defaulted(&tc.in)
			if !reflect.DeepEqual(tc.in, tc.out) {
				t.Errorf("Error: Expected and actual output are different \n actual: %+v\n expected: %+v\n", tc.in, tc.out)
			}
		})
	}
}

func Test_DefaultingFunction(t *testing.T) {
	in := marker.DefaultedWithFunction{}
	marker.SetObjectDefaults_DefaultedWithFunction(&in)
	out := marker.DefaultedWithFunction{
		S1: "default_function",
		S2: "default_marker",
	}
	if !reflect.DeepEqual(in, out) {
		t.Errorf("Error: Expected and actual output are different \n actual: %+v\n expected: %+v\n", in, out)
	}

}
