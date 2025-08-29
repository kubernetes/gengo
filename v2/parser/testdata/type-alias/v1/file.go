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

// This is for back compat.
// It has the same number of lines as Blah's comment.
//
// Deprecated: use Blah instead.
type LocalBlahAliasDeprecated = Blah

// This is for back compat.
// It has the same number of lines as Blah's comment.
//
// Deprecated: use Blah instead.
// It's a whole paragraph of deprecated notes
type LocalBlahAliasDeprecatedLong = Blah
