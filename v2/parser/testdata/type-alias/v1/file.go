package v1

// Blah is a test.
// A test, I tell you.
type Blah struct {
	// A is the first field.
	A int64 `json:"a"`

	// B is the second field.
	// Multiline comments work.
	B string `json:"b"`
}

// This is an alias within the same package.
type LocalBlahAlias = Blah
