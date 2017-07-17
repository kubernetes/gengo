/*
Copyright 2016 The Kubernetes Authors.

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

package wholepkg

type ManualDefault struct{}

func SetDefaults_ManualDefault(obj *ManualDefault) {
}

type Struct_With_Field_Has_Manual_Defaulter struct {
	TypeMeta      struct{}
	ManualDefault ManualDefault
	OtherField    string
}

type Struct_With_Field_With_Field_Has_Manual_Defaulter struct {
	TypeMeta struct{}
	Struct_With_Field_Has_Manual_Defaulter
	OtherField string
}

// No defaulter will be generated
type Struct_No_Field_Has_Manual_Defaulter struct {
	TypeMeta   struct{}
	OtherField string
}
