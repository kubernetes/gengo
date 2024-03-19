package foo

type DeepCopyable[T any] interface {
	DeepCopy() T
}

type Blah[T DeepCopyable[T]] struct {
	// V is the first field.
	V T `json:"v"`
}
