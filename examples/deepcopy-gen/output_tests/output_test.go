package output_tests

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/gofuzz"

	"k8s.io/gengo/examples/deepcopy-gen/output_tests/aliases"
	"k8s.io/gengo/examples/deepcopy-gen/output_tests/builtins"
	"k8s.io/gengo/examples/deepcopy-gen/output_tests/interfaces"
	"k8s.io/gengo/examples/deepcopy-gen/output_tests/maps"
	"k8s.io/gengo/examples/deepcopy-gen/output_tests/pointer"
	"k8s.io/gengo/examples/deepcopy-gen/output_tests/slices"
	"k8s.io/gengo/examples/deepcopy-gen/output_tests/structs"
)

func TestWithValueFuzzer(t *testing.T) {
	tests := []interface{}{
		aliases.Ttest{},
		builtins.Ttest{},
		interfaces.Ttest{},
		maps.Ttest{},
		pointer.Ttest{},
		slices.Ttest{},
		structs.Ttest{},
	}

	fuzzer := fuzz.New()
	fuzzer.NilChance(0.5)
	fuzzer.NumElements(0, 2)
	fuzzer.Funcs(interfaceFuzzers...)

	for _, test := range tests {
		t.Run(fmt.Sprintf("%T", test), func(t *testing.T) {
			N := 1000
			for i := 0; i < N; i++ {
				original := reflect.New(reflect.TypeOf(test)).Interface()

				fuzzer.Fuzz(original)

				reflectCopy := ReflectDeepCopy(original)

				if !reflect.DeepEqual(original, reflectCopy) {
					t.Errorf("original and reflectCopy are different:\n\n  original = %s\n\n  jsonCopy = %s", spew.Sdump(original), spew.Sdump(reflectCopy))
				}

				deepCopy := reflect.ValueOf(original).MethodByName("DeepCopy").Call(nil)[0].Interface()

				if !reflect.DeepEqual(original, deepCopy) {
					t.Fatalf("original and deepCopy are different:\n\n  original = %s\n\n  deepCopy() = %s", spew.Sdump(original), spew.Sdump(deepCopy))
				}

				ValueFuzz(original)

				if !reflect.DeepEqual(reflectCopy, deepCopy) {
					t.Fatalf("reflectCopy and deepCopy are different:\n\n  origin = %s\n\n  jsonCopy() = %s", spew.Sdump(original), spew.Sdump(deepCopy))
				}
			}
		})
	}
}

func BenchmarkReflectDeepCopy(b *testing.B) {
	fourtytwo := "fourtytwo"

	tests := []interface{}{
		maps.Ttest{
			Byte:      map[string]byte{"0": 42, "1": 42, "3": 42},
			Int16:     map[string]int16{"0": 42, "1": 42, "3": 42},
			Int32:     map[string]int32{"0": 42, "1": 42, "3": 42},
			Int64:     map[string]int64{"0": 42, "1": 42, "3": 42},
			Uint8:     map[string]uint8{"0": 42, "1": 42, "3": 42},
			Uint16:    map[string]uint16{"0": 42, "1": 42, "3": 42},
			Uint32:    map[string]uint32{"0": 42, "1": 42, "3": 42},
			Uint64:    map[string]uint64{"0": 42, "1": 42, "3": 42},
			Float32:   map[string]float32{"0": 42.0, "1": 42.0, "3": 42.0},
			Float64:   map[string]float64{"0": 42, "1": 42, "3": 42},
			String:    map[string]string{"0": "fourtytwo", "1": "fourtytwo", "3": "fourtytwo"},
			StringPtr: map[string]*string{"0": &fourtytwo, "1": &fourtytwo, "3": &fourtytwo},
		},
		slices.Ttest{
			Byte:      []byte{42, 42, 42},
			Int16:     []int16{42, 42, 42},
			Int32:     []int32{42, 42, 42},
			Int64:     []int64{42, 42, 42},
			Uint8:     []uint8{42, 42, 42},
			Uint16:    []uint16{42, 42, 42},
			Uint32:    []uint32{42, 42, 42},
			Uint64:    []uint64{42, 42, 42},
			Float32:   []float32{42.0, 42.0, 42.0},
			Float64:   []float64{42, 42, 42},
			String:    []string{"fourtytwo", "fourtytwo", "fourtytwo"},
			StringPtr: []*string{&fourtytwo, &fourtytwo, &fourtytwo},
		},
		pointer.Ttest{
			Types: map[string]*pointer.Ttest{},
		},
	}

	fuzzer := fuzz.New()
	fuzzer.NilChance(0.5)
	fuzzer.NumElements(0, 2)
	fuzzer.Funcs(interfaceFuzzers...)

	for _, test := range tests {
		b.Run(fmt.Sprintf("%T", test), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				switch t := test.(type) {
				case maps.Ttest:
					t.DeepCopy()
				case slices.Ttest:
					t.DeepCopy()
				case pointer.Ttest:
					t.DeepCopy()
				default:
					b.Fatalf("missing type case in switch for %T", t)
				}
			}
		})
	}
}
