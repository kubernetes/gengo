package foo

type Blah[T any, U any, V any] struct {
	// V1 is the first field.
	V1 T `json:"v1"`
	// V2 is the second field.
	V2 U `json:"v2"`
	// V3 is the third field.
	V3 V `json:"v3"`
}
