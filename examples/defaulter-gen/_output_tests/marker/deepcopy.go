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

func (in *SubStruct) DeepCopy() *SubStruct {
	if in == nil {
		return nil
	}
	out := new(SubStruct)
	in.DeepCopyInto(out)
	return out
}

func (in *SubStruct) DeepCopyInto(out *SubStruct) {
	*out = *in
	return
}

func (in *SubStruct) DeepCopyObject() runtime.Object {
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
func (obj *SubStruct) GetObjectKind() schema.ObjectKind             { return schema.EmptyObjectKind }
func (obj *DefaultedWithFunction) GetObjectKind() schema.ObjectKind { return schema.EmptyObjectKind }
