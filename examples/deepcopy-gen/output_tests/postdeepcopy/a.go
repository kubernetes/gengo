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

package postdeepcopy

type Alias_WithPostDeepCopy string

type Struct_WithPostDeepCopy struct {
	marker int
}

type Struct_WithPostDeepCopyFields struct {
	A      Alias_WithPostDeepCopy
	APtr   *Alias_WithPostDeepCopy
	ASlice []Alias_WithPostDeepCopy
	AMap   map[string]Alias_WithPostDeepCopy

	S      Struct_WithPostDeepCopy
	SPtr   *Struct_WithPostDeepCopy
	SSlice []Struct_WithPostDeepCopy
	SMap   map[string]Struct_WithPostDeepCopy
}

type Struct_WithSkippedFields struct {
	S   string
	Ptr *string
	// +k8s:deepcopy-gen:skip-field
	//
	// JSON is what json.Unmarshal returns when an interface is passed.
	JSON interface{}
}

func (in *Alias_WithPostDeepCopy) postDeepCopy(out *Alias_WithPostDeepCopy) {
}

func (in *Struct_WithPostDeepCopy) postDeepCopy(out *Struct_WithPostDeepCopy) {
	out.marker = 42
}

func (in *Struct_WithSkippedFields) postDeepCopy(out *Struct_WithSkippedFields) {
	out.JSON = deepCopyJSON(in.JSON)
}
