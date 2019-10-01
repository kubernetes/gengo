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
	"k8s.io/gengo/types"
)

// an unsafeConversionArbitrator decides whether conversions can be done using unsafe casts.
type unsafeConversionArbitrator struct {
	processedPairs           map[ConversionPair]bool
	manualConversionsTracker *ManualConversionsTracker
	functionTagName          string
}

func newUnsafeConversionArbitrator(manualConversionsTracker *ManualConversionsTracker) *unsafeConversionArbitrator {
	return &unsafeConversionArbitrator{
		processedPairs:           make(map[ConversionPair]bool),
		manualConversionsTracker: manualConversionsTracker,
	}
}

// canUseUnsafeConversion returns true iff x can be converted to y using an unsafe conversion.
func (a *unsafeConversionArbitrator) canUseUnsafeConversion(x, y *types.Type) bool {
	// alreadyVisitedTypes holds all the types that have already been checked in the structural type recursion.
	alreadyVisitedTypes := make(map[*types.Type]bool)
	return a.canUseUnsafeConversionWithCaching(x, y, alreadyVisitedTypes)
}

func (a *unsafeConversionArbitrator) canUseUnsafeConversionWithCaching(x, y *types.Type, alreadyVisitedTypes map[*types.Type]bool) bool {
	if x == y {
		return true
	}

	if equal, ok := a.processedPairs[ConversionPair{x, y}]; ok {
		return equal
	}
	if equal, ok := a.processedPairs[ConversionPair{y, x}]; ok {
		return equal
	}

	result := !a.nonCopyOnlyManualConversionFunctionExists(x, y) && a.canUseUnsafeRecursive(x, y, alreadyVisitedTypes)
	a.processedPairs[ConversionPair{x, y}] = result
	return result
}

// nonCopyOnlyManualConversionFunctionExists returns true iff the manual conversion tracker
// knows of a conversion function from x to y, that is not a copy-only conversion function.
func (a *unsafeConversionArbitrator) nonCopyOnlyManualConversionFunctionExists(x, y *types.Type) bool {
	conversionFunction, exists := a.manualConversionsTracker.preexists(x, y)
	return exists && !isCopyOnlyFunction(conversionFunction, a.functionTagName)
}

// setFunctionTagName sets the function tag name.
// That also invalidates the cache if the new function tag name is different than the previous one.
func (a *unsafeConversionArbitrator) setFunctionTagName(functionTagName string) {
	if a.functionTagName != functionTagName {
		a.functionTagName = functionTagName
		a.processedPairs = make(map[ConversionPair]bool)
	}
}

func (a *unsafeConversionArbitrator) canUseUnsafeRecursive(x, y *types.Type, alreadyVisitedTypes map[*types.Type]bool) bool {
	in, out := unwrapAlias(x), unwrapAlias(y)
	switch {
	case in == out:
		return true
	case in.Kind == out.Kind:
		// if the type exists already, return early to avoid recursion
		if alreadyVisitedTypes[in] {
			return true
		}
		alreadyVisitedTypes[in] = true

		switch in.Kind {
		case types.Struct:
			if len(in.Members) != len(out.Members) {
				return false
			}
			for i, inMember := range in.Members {
				outMember := out.Members[i]
				if !a.canUseUnsafeConversionWithCaching(inMember.Type, outMember.Type, alreadyVisitedTypes) {
					return false
				}
			}
			return true
		case types.Pointer:
			return a.canUseUnsafeConversionWithCaching(in.Elem, out.Elem, alreadyVisitedTypes)
		case types.Map:
			return a.canUseUnsafeConversionWithCaching(in.Key, out.Key, alreadyVisitedTypes) &&
				a.canUseUnsafeConversionWithCaching(in.Elem, out.Elem, alreadyVisitedTypes)
		case types.Slice:
			return a.canUseUnsafeConversionWithCaching(in.Elem, out.Elem, alreadyVisitedTypes)
		case types.Interface:
			// TODO: determine whether the interfaces are actually equivalent - for now, they must have the
			// same type.
			return false
		case types.Builtin:
			return in.Name.Name == out.Name.Name
		}
	}
	return false
}
