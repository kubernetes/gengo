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

// +k8s:deepcopy-gen=package

// This is a test package.
package aliases

type Foo struct {
	X int
}

type Builtin int
type Slice []int
type Pointer *int
type PointerAlias *Builtin
type Struct Foo
type Map map[string]int

type FooAlias Foo
type FooSlice []Foo
type FooPointer *Foo
type FooMap map[string]Foo

type AliasBuiltin Builtin
type AliasSlice Slice
type AliasPointer Pointer
type AliasStruct Struct
type AliasMap Map

// Aliases
type Ttest struct {
	Builtin      Builtin
	Slice        Slice
	Pointer      Pointer
	PointerAlias PointerAlias
	Struct       Struct
	Map          Map
	SliceSlice   []Slice
	MapSlice     map[string]Slice

	FooAlias   FooAlias
	FooSlice   FooSlice
	FooPointer FooPointer
	FooMap     FooMap

	AliasBuiltin AliasBuiltin
	AliasSlice   AliasSlice
	AliasPointer AliasPointer
	AliasStruct  AliasStruct
	AliasMap     AliasMap
}
