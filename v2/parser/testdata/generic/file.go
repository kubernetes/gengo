package foo

type Blah[T any] struct {
	// V is the first field.
	V T `json:"v"`
}
