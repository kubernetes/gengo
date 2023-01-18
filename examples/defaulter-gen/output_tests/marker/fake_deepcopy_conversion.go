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
	runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func (in *Defaulted) DeepCopy() *Defaulted {
	if in == nil {
		return nil
	}
	out := new(Defaulted)
	in.DeepCopyInto(out)
	return out
}

func (in *Defaulted) DeepCopyInto(out *Defaulted) {
	*out = *in
	return
}

func (in *Defaulted) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *DefaultedWithFunction) DeepCopy() *DefaultedWithFunction {
	if in == nil {
		return nil
	}
	out := new(DefaultedWithFunction)
	in.DeepCopyInto(out)
	return out
}

func (in *DefaultedWithFunction) DeepCopyInto(out *DefaultedWithFunction) {
	*out = *in
	return
}

func (in *DefaultedWithFunction) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (obj *Defaulted) GetObjectKind() schema.ObjectKind             { return schema.EmptyObjectKind }
func (obj *DefaultedWithFunction) GetObjectKind() schema.ObjectKind { return schema.EmptyObjectKind }
