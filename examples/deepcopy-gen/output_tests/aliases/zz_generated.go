// +build !ignore_autogenerated

/*
Copyright 2018 The Kubernetes Authors.

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

// This file was autogenerated by deepcopy-gen. Do not edit it manually!

package aliases

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AliasMap) DeepCopyInto(out *AliasMap) {
	*out = make(AliasMap, len(*in))
	for key, val := range *in {
		(*out)[key] = val
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AliasMap.
func (in *AliasMap) DeepCopy() *AliasMap {
	if in == nil {
		return nil
	}
	out := new(AliasMap)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AliasSlice) DeepCopyInto(out *AliasSlice) {
	*out = make(AliasSlice, len(*in))
	copy(*out, *in)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AliasSlice.
func (in *AliasSlice) DeepCopy() *AliasSlice {
	if in == nil {
		return nil
	}
	out := new(AliasSlice)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AliasStruct) DeepCopyInto(out *AliasStruct) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AliasStruct.
func (in *AliasStruct) DeepCopy() *AliasStruct {
	if in == nil {
		return nil
	}
	out := new(AliasStruct)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Foo) DeepCopyInto(out *Foo) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Foo.
func (in *Foo) DeepCopy() *Foo {
	if in == nil {
		return nil
	}
	out := new(Foo)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FooAlias) DeepCopyInto(out *FooAlias) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FooAlias.
func (in *FooAlias) DeepCopy() *FooAlias {
	if in == nil {
		return nil
	}
	out := new(FooAlias)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FooMap) DeepCopyInto(out *FooMap) {
	*out = make(FooMap, len(*in))
	for key, val := range *in {
		(*out)[key] = val
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FooMap.
func (in *FooMap) DeepCopy() *FooMap {
	if in == nil {
		return nil
	}
	out := new(FooMap)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FooSlice) DeepCopyInto(out *FooSlice) {
	*out = make(FooSlice, len(*in))
	copy(*out, *in)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FooSlice.
func (in *FooSlice) DeepCopy() *FooSlice {
	if in == nil {
		return nil
	}
	out := new(FooSlice)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Map) DeepCopyInto(out *Map) {
	*out = make(Map, len(*in))
	for key, val := range *in {
		(*out)[key] = val
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Map.
func (in *Map) DeepCopy() *Map {
	if in == nil {
		return nil
	}
	out := new(Map)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Slice) DeepCopyInto(out *Slice) {
	*out = make(Slice, len(*in))
	copy(*out, *in)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Slice.
func (in *Slice) DeepCopy() *Slice {
	if in == nil {
		return nil
	}
	out := new(Slice)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Struct) DeepCopyInto(out *Struct) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Struct.
func (in *Struct) DeepCopy() *Struct {
	if in == nil {
		return nil
	}
	out := new(Struct)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Ttest) DeepCopyInto(out *Ttest) {
	*out = *in
	if in.Slice != nil {
		in, out := &in.Slice, &out.Slice
		*out = make(Slice, len(*in))
		copy(*out, *in)
	}
	if in.Pointer != nil {
		in, out := &in.Pointer, &out.Pointer
		if *in == nil {
			*out = nil
		} else {
			*out = new(int)
			**out = **in
		}
	}
	if in.PointerAlias != nil {
		in, out := &in.PointerAlias, &out.PointerAlias
		if *in == nil {
			*out = nil
		} else {
			*out = new(Builtin)
			**out = **in
		}
	}
	out.Struct = in.Struct
	if in.Map != nil {
		in, out := &in.Map, &out.Map
		*out = make(Map, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.SliceSlice != nil {
		in, out := &in.SliceSlice, &out.SliceSlice
		*out = make([]Slice, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = make(Slice, len(*in))
				copy(*out, *in)
			}
		}
	}
	if in.MapSlice != nil {
		in, out := &in.MapSlice, &out.MapSlice
		*out = make(map[string]Slice, len(*in))
		for key, val := range *in {
			if val == nil {
				(*out)[key] = nil
			} else {
				(*out)[key] = make([]int, len(val))
				copy((*out)[key], val)
			}
		}
	}
	out.FooAlias = in.FooAlias
	if in.FooSlice != nil {
		in, out := &in.FooSlice, &out.FooSlice
		*out = make(FooSlice, len(*in))
		copy(*out, *in)
	}
	if in.FooPointer != nil {
		in, out := &in.FooPointer, &out.FooPointer
		if *in == nil {
			*out = nil
		} else {
			*out = new(Foo)
			**out = **in
		}
	}
	if in.FooMap != nil {
		in, out := &in.FooMap, &out.FooMap
		*out = make(FooMap, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.AliasSlice != nil {
		in, out := &in.AliasSlice, &out.AliasSlice
		*out = make(AliasSlice, len(*in))
		copy(*out, *in)
	}
	if in.AliasPointer != nil {
		in, out := &in.AliasPointer, &out.AliasPointer
		if *in == nil {
			*out = nil
		} else {
			*out = new(int)
			**out = **in
		}
	}
	out.AliasStruct = in.AliasStruct
	if in.AliasMap != nil {
		in, out := &in.AliasMap, &out.AliasMap
		*out = make(AliasMap, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Ttest.
func (in *Ttest) DeepCopy() *Ttest {
	if in == nil {
		return nil
	}
	out := new(Ttest)
	in.DeepCopyInto(out)
	return out
}
