package target1

import (
	"k8s.io/apimachinery/pkg/runtime"
)

type dummySchemeBuilder struct{}

func (dummySchemeBuilder) Register(func(s *runtime.Scheme) error) {}

var localSchemeBuilder = dummySchemeBuilder{}
