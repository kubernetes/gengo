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

import (
	"reflect"
	"testing"

	fuzz "github.com/google/gofuzz"
)

func TestNoCopyOptimizationForAssignablesWithPostDeepCopy(t *testing.T) {
	x := &Struct_WithPostDeepCopyFields{}

	fuzzer := fuzz.New()
	fuzzer.Funcs(func(s *Struct_WithPostDeepCopy, c fuzz.Continue) {
		s.marker = 0
	})
	fuzzer.Fuzz(x)

	y := x.DeepCopy()

	if y.S.marker != 42 {
		t.Errorf("postDeepCopy not called on x.S")
	}
	y.S.marker = 0

	if y.SPtr != nil && y.SPtr.marker != 42 {
		t.Errorf("postDeepCopy not called on x.SPtr")

	}
	y.SPtr.marker = 0

	for i := range y.SSlice {
		if y.SSlice[i].marker != 42 {
			t.Errorf("postDeepCopy not called on x.SSlice")
			break
		}
		y.SSlice[i].marker = 0
	}

	for k := range y.SMap {
		if y.SMap[k].marker != 42 {
			t.Errorf("postDeepCopy not called on x.SMap")
			break
		}
		clone := y.SMap[k]
		clone.marker = 0
		y.SMap[k] = clone
	}

	if !reflect.DeepEqual(y, x) {
		t.Errorf("objects should be equal, but are not: %#v, %#v", x, y)
	}
}

