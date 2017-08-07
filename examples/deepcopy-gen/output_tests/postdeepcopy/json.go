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

func deepCopyJSON(x interface{}) interface{} {
	switch x := x.(type) {
	case map[string]interface{}:
		clone := make(map[string]interface{}, len(x))
		for k, v := range x {
			clone[k] = deepCopyJSON(v)
		}
		return clone
	case []interface{}:
		clone := make([]interface{}, len(x))
		for i := range x {
			clone[i] = deepCopyJSON(x[i])
		}
		return clone
	default:
		// only non-pointer values (float64, int64, bool, string) are left. These can be copied by-value.
		return x
	}
}
