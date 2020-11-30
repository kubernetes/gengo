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
	"k8s.io/gengo/examples/defaulter-gen/_output_tests/empty"
)

type Defaulted struct {
	empty.TypeMeta

	// +default="bar"
	StringDefault string `json:"omitempty"`

	// Default is forced to empty string
	// Specifying the default is a no-op
	// +default=""
	StringEmptyDefault string

	// Not specifying a default still defaults for non-omitempty
	StringEmpty string

	// +default="default"
	StringPointer *string

	// +default=1
	IntDefault int `json:"omitempty"`

	// +default=0
	IntEmptyDefault int

	// Default is forced to 0
	IntEmpty int

	// +default=0.5
	FloatDefault float64 `json:"omitempty"`

	// +default=0.0
	FloatEmptyDefault float64

	FloatEmpty float64

	// +default=["foo", "bar"]
	List []Item
	// +default={"s": "foo", "i": 5}
	Sub *SubStruct

	//+default=[{"s": "foo1", "i": 1}, {"s": "foo2"}]
	StructList []SubStruct

	//+default=[{"s": "foo1", "i": 1}, {"s": "foo2"}]
	PtrStructList []*SubStruct

	//+default=["foo"]
	StringList []string

	// Default is forced to empty struct
	OtherSub SubStruct

	// +default={"foo": "bar"}
	Map map[string]Item

	// +default={"foo": {"S": "string", "I": 1}}
	StructMap map[string]SubStruct

	// +default={"foo": {"S": "string", "I": 1}}
	PtrStructMap map[string]*SubStruct

	// A default specified here overrides the default for the Item type
	// +default="banana"
	AliasPtr Item
}

// +default="apple"
type Item *string

type SubStruct struct {
	S string
	// +default=1
	I int `json:"I,omitempty"`
}

type DefaultedWithFunction struct {
	empty.TypeMeta
	// +default="default_marker"
	S1 string `json:"S1,omitempty"`
	// +default="default_marker"
	S2 string `json:"S2,omitempty"`
}
