/*
Copyright 2019 The Kubernetes Authors.

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

package generator

import (
	"reflect"
	"testing"

	"k8s.io/gengo/v2/types"
)

func TestNewImportTracker(t *testing.T) {
	tests := []struct {
		name            string
		localPackage    string
		inputTypes      []*types.Type
		expectedImports []string
	}{
		{
			name:            "empty",
			inputTypes:      []*types.Type{},
			expectedImports: []string{},
		},
		{
			name: "builtin",
			inputTypes: []*types.Type{
				{Name: types.Name{Package: "context"}},
				{Name: types.Name{Package: "net/http"}},
				{Name: types.Name{Package: "x/net/http"}},
			},
			expectedImports: []string{
				`"context"`,
				`"net/http"`,
				`nethttp "x/net/http"`,
			},
		},
		{
			name: "sorting",
			inputTypes: []*types.Type{
				{Name: types.Name{Package: "foo/bar/pkg2"}},
				{Name: types.Name{Package: "foo/bar/pkg1"}},
			},
			expectedImports: []string{
				`"foo/bar/pkg1"`,
				`"foo/bar/pkg2"`,
			},
		},
		{
			name: "duplicate",
			inputTypes: []*types.Type{
				{Name: types.Name{Package: "foo/bar/pkg2/v1"}},
				{Name: types.Name{Package: "foo/bar/pkg3/v1"}},
				{Name: types.Name{Package: "foo/bar/pkg1/v1"}},
			},
			expectedImports: []string{
				`pkg1v1 "foo/bar/pkg1/v1"`,
				`"foo/bar/pkg2/v1"`,
				`pkg3v1 "foo/bar/pkg3/v1"`,
			},
		},
		{
			name: "reserved-keyword",
			inputTypes: []*types.Type{
				{Name: types.Name{Package: "my/reserved/pkg/struct"}},
			},
			expectedImports: []string{
				`_struct "my/reserved/pkg/struct"`,
			},
		},
		{
			name:         "local-symbol",
			localPackage: "bar.com/my/pkg",
			inputTypes: []*types.Type{
				{Name: types.Name{Package: "bar.com/my/pkg"}},
				{Name: types.Name{Package: "bar.com/external/pkg"}},
			},
			expectedImports: []string{
				`externalpkg "bar.com/external/pkg"`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var actualImports []string
			if len(tt.localPackage) == 0 {
				actualImports = NewImportTracker(tt.inputTypes...).ImportLines()
			} else {
				actualImports = NewImportTrackerForPackage(tt.localPackage, tt.inputTypes...).ImportLines()
			}
			if !reflect.DeepEqual(actualImports, tt.expectedImports) {
				t.Errorf("ImportLines(%v) = %v, want %v", tt.inputTypes, actualImports, tt.expectedImports)
			}
		})
	}
}
