package v1

// This is an alias within the same package.
// It comes before the type it aliases.
type LocalBlahAlias = Blah

// Blah is a test.
// A test, I tell you.
//
// Another paragraph of comments.
type Blah struct {
	// A is the first field.
	A int64 `json:"a"`

	// B is the second field.
	// Multiline comments work.
	B string `json:"b"`
}

// Blah is a test.
// A test, I tell you.
//
// Deprecated: use Blah instead. This is another alias within the same package.
//
// Another paragraph of comments.
type LocalBlahAliasDeprecated = Blah

// Blah is a test.
// A test, I tell you.
//
// Deprecated: use Blah instead. This is a third alias within the same package.
// It's a whole paragraph of deprecated notes
//
// Another paragraph of comments.
type LocalBlahAliasDeprecatedLong = Blah
