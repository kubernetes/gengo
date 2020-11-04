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
	Field string `json:"Field,omitempty"`
	// +default=0
	OtherField int `json:"OtherField,omitempty"`

	// Default is forced to 0
	EmptyInt int

	// Default is forced to empty string
	EmptyString string

	// +default=["foo", "bar"]
	List []Item
	// +default={"s": "moo", "i": 5}
	Sub *SubStruct

	// Default is forced to empty struct
	OtherSub SubStruct

	// +default={"foo": "bar"}
	Map map[string]Item
}

// +default="foo"
type Item *string

type SubStruct struct {
	S string
	// +default=1
	I int `json:"I,omitempty"`
}
