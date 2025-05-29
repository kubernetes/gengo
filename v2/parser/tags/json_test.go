/*
Copyright 2024 The Kubernetes Authors.

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

package tags

import (
	"testing"

	"k8s.io/gengo/v2/parser"
	"k8s.io/gengo/v2/types"
)

func TestJSON(t *testing.T) {
	p := parser.New()
	u := types.Universe{}
	// Proper packages with deps.
	pkgs, err := p.LoadPackagesTo(&u, "./testdata/tags")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	member := func(typeName, memberName string) types.Member {
		typ := pkgs[0].Type(typeName)
		for _, m := range typ.Members {
			if m.Name == memberName {
				return m
			}
		}
		t.Fatalf("member %s not found", memberName)
		return types.Member{}
	}

	tests := []struct {
		name     string
		member   types.Member
		expected JSON
	}{
		{
			name:     "name",
			member:   member("T1", "A"),
			expected: JSON{Name: "a"},
		},
		{
			name:     "omitempty",
			member:   member("T1", "B"),
			expected: JSON{Name: "b", Omitempty: true},
		},
		{
			name:     "inline",
			member:   member("T1", "C"),
			expected: JSON{Name: "", Inline: true},
		},
		{
			name:     "omit",
			member:   member("T1", "D"),
			expected: JSON{Name: "", Omit: true},
		},
		{
			name:     "empty",
			member:   member("T1", "E"),
			expected: JSON{Name: "E"},
		},
		{
			name:     "embedded struct",
			member:   member("T1", "T2"),
			expected: JSON{Name: "", Inline: true},
		},
		{
			name:     "embedded pointer",
			member:   member("T1", "T3"),
			expected: JSON{Name: "", Inline: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tags, ok := LookupJSON(tt.member)
			if !ok {
				t.Errorf("failed to lookup tags")
			}
			if tags != tt.expected {
				t.Errorf("expected %#+v, got %#+v", tt.expected, tags)
			}
		})
	}
}

type t1 struct {
	Name string `json:"name"`
}
