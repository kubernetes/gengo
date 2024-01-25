/*
Copyright YEAR The Kubernetes Authors.

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

package a

import "k8s.io/gengo/testdata/a/b"

// A is a type for testing.
type A string

// AA is a struct type for testing.
type AA struct {
	FieldA string
}

// AFunc is a member function for testing.
func (a *AA) AFunc(i *int, j int) (*A, b.ITest, error) {
	return nil, nil, nil
}
