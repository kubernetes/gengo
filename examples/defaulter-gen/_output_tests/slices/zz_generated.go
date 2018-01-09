// +build !ignore_autogenerated

/*
Copyright The Kubernetes Authors.

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

// This file was autogenerated by defaulter-gen. Do not edit it manually!

package slices

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// RegisterDefaults adds defaulters functions to the given scheme.
// Public to allow building arbitrary schemes.
// All generated defaulters are covering - they call all nested defaulters.
func RegisterDefaults(scheme *runtime.Scheme) error {
	scheme.AddTypeDefaultingFunc(&Ttest{}, func(obj interface{}) { SetObjectDefaults_Ttest(obj.(*Ttest)) })
	scheme.AddTypeDefaultingFunc(&TtestList{}, func(obj interface{}) { SetObjectDefaults_TtestList(obj.(*TtestList)) })
	scheme.AddTypeDefaultingFunc(&TtestPointerList{}, func(obj interface{}) { SetObjectDefaults_TtestPointerList(obj.(*TtestPointerList)) })
	return nil
}

func SetObjectDefaults_Ttest(in *Ttest) {
	SetDefaults_Ttest(in)
}

func SetObjectDefaults_TtestList(in *TtestList) {
	for i := range in.Items {
		a := &in.Items[i]
		SetObjectDefaults_Ttest(a)
	}
}

func SetObjectDefaults_TtestPointerList(in *TtestPointerList) {
	for i := range in.Items {
		a := &in.Items[i]
		if a != nil {
			SetObjectDefaults_Ttest(a)
		}
	}
}
