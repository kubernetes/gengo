package target2

import (
	"k8s.io/apimachinery/pkg/conversion"

	"k8s.io/gengo/examples/conversion-gen/output_tests/peer2"
)

func Convert_peer2_Bar_To_target2_Bar(in *peer2.Bar, out *Bar, s conversion.Scope) error {
	return nil
}

// +k8s:conversion-fn=drop
func Convert_peer2_SubBarStruct_To_target2_SubBarStruct(in *peer2.SubBarStruct, out *SubBarStruct, s conversion.Scope) error {
	return nil
}

// Convert_Slice_byte_To_Slice_byte prevents recursing into every byte
func Convert_Slice_byte_To_Slice_byte(in *[]byte, out *[]byte, s conversion.Scope) error {
	if *in == nil {
		*out = nil
		return nil
	}
	*out = make([]byte, len(*in))
	copy(*out, *in)
	return nil
}
