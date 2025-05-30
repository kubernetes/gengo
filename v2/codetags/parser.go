/*
Copyright 2025 The Kubernetes Authors.

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

package codetags

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

// Parse parses a tag string into a TypedTag, or returns an error if the tag
// string fails to parse.
//
// ParseOption may be provided to modify the behavior of the parser. The below
// describes the default behavior.
//
// A tag consists of a name, optional arguments, and an optional scalar value or
// tag value. For example,
//
//	"name"
//	"name=50"
//	"name("featureX")=50"
//	"name(limit: 10, path: "/xyz")=text value"
//	"name(limit: 10, path: "/xyz")=+anotherTag(size: 100)"
//
// Arguments are optional and may be either:
//   - A single positional argument.
//   - One or more named arguments (in the format `name: value`).
//   - (Positional and named arguments cannot be mixed.)
//
// For example,
//
//	"name()"
//	"name(arg)"
//	"name(namedArg1: argValue1)"
//	"name(namedArg1: argValue1, namedArg2: argValue2)"
//
// Argument values may be strings, ints, booleans, or identifiers.
//
// For example,
//
//	"name("double-quoted")"
//	"name(`backtick-quoted`)"
//	"name(100)"
//	"name(true)"
//	"name(arg1: identifier)"
//	"name(arg1:`string value`)"
//	"name(arg1: 100)"
//	"name(arg1: true)"
//
// Note: When processing Go source code comments, the Extract function is
// typically used first to find and isolate tag strings matching a specific
// prefix. Those extracted strings can then be parsed using this function.
//
// The value part of the tag is optional and follows an equals sign "=". If a
// value is present, it must be a string, int, boolean, identifier, or tag.
//
// For example,
//
//	"name" // no value
//	"name=identifier"
//	"name="double-quoted value""
//	"name=`backtick-quoted value`"
//	"name(100)"
//	"name(true)"
//	"name=+anotherTag"
//	"name=+anotherTag(size: 100)"
//
// Trailing comments are ignored unless the RawValues option is enabled, in which
// case they are treated as part of the value.
//
// For example,
//
//	"key=value // This comment is ignored"
//
// Formal Grammar:
//
// <tag>             ::= <tagName> [ "(" [ <args> ] ")" ] [ ( "=" <value> | "=+" <tag> ) ]
// <args>            ::= <value> | <namedArgs>
// <namedArgs>       ::= <argNameAndValue> [ "," <namedArgs> ]*
// <argNameAndValue> ::= <identifier> ":" <value>
// <value>           ::= <identifier> | <string> | <int> | <bool>
//
// <tagName>       ::= [a-zA-Z_][a-zA-Z0-9_-.:]*
// <identifier>    ::= [a-zA-Z_][a-zA-Z0-9_-.]*
// <string>        ::= /* Go-style double-quoted or backtick-quoted strings,
// ...                    with standard Go escape sequences for double-quoted strings. */
// <int>           ::= /* Standard Go integer literals (decimal, 0x hex, 0o octal, 0b binary),
// ...                    with an optional +/- prefix. */
// <bool>          ::= "true" | "false"
func Parse(tag string, options ...ParseOption) (TypedTag, error) {
	opts := parseOpts{}
	for _, o := range options {
		o(&opts)
	}

	tag = strings.TrimSpace(tag)
	return parseTag(tag, opts)
}

// ParseAll calls Parse on each tag in the input slice.
func ParseAll(tags []string, options ...ParseOption) ([]TypedTag, error) {
	var out []TypedTag
	for _, tag := range tags {
		parsed, err := Parse(tag, options...)
		if err != nil {
			return nil, err
		}
		out = append(out, parsed)
	}
	return out, nil
}

type parseOpts struct {
	rawValues bool
}

// ParseOption provides a parser option.
type ParseOption func(*parseOpts)

// RawValues skips parsing of the value part of the tag. If enabled, the Value
// in the parse response will contain all text following the "=" sign, up to the last
// non-whitespace character, and ValueType will be set to ValueTypeRaw.
// Default: disabled
func RawValues(enabled bool) ParseOption {
	return func(opts *parseOpts) {
		opts.rawValues = enabled
	}
}

const (
	stBegin           = "stBegin"
	stTag             = "stTag"
	stMaybeArgs       = "stMaybeArgs"
	stArg             = "stArg"
	stArgEndOfToken   = "stArgEndOfToken"
	stMaybeValue      = "stMaybeValue"
	stValue           = "stValue"
	stMaybeComment    = "stMaybeComment"
	stTrailingSlash   = "stTrailingSlash"
	stTrailingComment = "stTrailingComment"
)

func parseTag(input string, opts parseOpts) (TypedTag, error) {
	var startTag, endTag *TypedTag // both ends of the chain when parsing chained tags

	tag := bytes.Buffer{}   // current tag name
	var args []Arg          // all tag arguments
	value := bytes.Buffer{} // current tag value
	var valueType ValueType // current value type
	var hasValue bool       // true if the tag has a value

	cur := Arg{}          // current argument accumulator
	buf := bytes.Buffer{} // string accumulator

	// These are defined outside the loop to make errors easier.
	s := scanner{buf: []rune(input)}
	var incomplete bool

	saveInt := func() error {
		s := buf.String()
		cur.Value = s
		cur.Type = ArgTypeInt
		args = append(args, cur)
		cur = Arg{}
		buf.Reset()
		return nil
	}
	saveString := func() {
		s := buf.String()
		cur.Value = s
		cur.Type = ArgTypeString
		args = append(args, cur)
		cur = Arg{}
		buf.Reset()
	}
	saveBoolOrString := func() {
		s := buf.String()
		if s == "true" || s == "false" {
			cur.Value = s
			cur.Type = ArgTypeBool
		} else {
			cur.Value = s
			cur.Type = ArgTypeString
		}
		args = append(args, cur)
		cur = Arg{}
		buf.Reset()
	}
	saveName := func() {
		cur.Name = buf.String()
		buf.Reset()
	}
	saveTag := func() error {
		usingNamedArgs := false
		for i, arg := range args {
			if (usingNamedArgs && arg.Name == "") || (!usingNamedArgs && arg.Name != "" && i > 0) {
				return fmt.Errorf("can't mix named and positional arguments")
			}
			if arg.Name != "" {
				usingNamedArgs = true
			}
		}
		if !usingNamedArgs && len(args) > 1 {
			return fmt.Errorf("multiple arguments must use 'name: value' syntax")
		}

		newTag := &TypedTag{Name: tag.String(), Args: args}
		if startTag == nil {
			startTag = newTag
			endTag = newTag
		} else {
			endTag.ValueTag = newTag
			endTag.ValueType = ValueTypeTag
			endTag = newTag
		}
		args = nil // Reset to nil instead of empty slice
		tag.Reset()
		return nil
	}
	saveValue := func() {
		endTag.Value = value.String()
		if opts.rawValues {
			endTag.ValueType = ValueTypeRaw
			return
		}
		endTag.ValueType = valueType
		if valueType == ValueTypeString && (endTag.Value == "true" || endTag.Value == "false") {
			endTag.ValueType = ValueTypeBool
		}
	}
	st := stBegin
parseLoop:
	for r := s.peek(); r != EOF; r = s.peek() {
		switch st {
		case stBegin:
			switch {
			case unicode.IsSpace(r):
				s.next()
				continue
			case isIdentBegin(r):
				st = stTag
			default:
				break parseLoop
			}
		case stTag:
			switch {
			case isIdentBegin(r):
				ident, err := s.nextIdent(isTagNameInterior)
				if err != nil {
					return TypedTag{}, err
				}
				tag.WriteString(ident)
				st = stMaybeArgs
			default:
				break parseLoop
			}
		case stMaybeArgs:
			switch {
			case r == '(':
				s.next()
				incomplete = true
				st = stArg
			case r == '=':
				s.next()
				hasValue = true
				st = stValue
			case unicode.IsSpace(r):
				s.next()
				st = stMaybeComment
			default:
				break parseLoop
			}
		case stArg:
			switch {
			case unicode.IsSpace(r):
				s.next()
				continue
			case r == ')':
				s.next()
				incomplete = false
				st = stMaybeValue
			case r == '-' || r == '+' || unicode.IsDigit(r):
				number, err := s.nextNumber()
				if err != nil {
					return TypedTag{}, err
				}
				buf.WriteString(number)
				if err := saveInt(); err != nil {
					return TypedTag{}, err
				}
				st = stArgEndOfToken
			case r == '"' || r == '`':
				str, err := s.nextString()
				if err != nil {
					return TypedTag{}, err
				}
				buf.WriteString(str)
				saveString()
				st = stArgEndOfToken
			case isIdentBegin(r):
				str, err := s.nextIdent(isIdentInterior)
				if err != nil {
					return TypedTag{}, err
				}
				buf.WriteString(str)
				r = s.peek()
				switch {
				case r == ',' || r == ')' || unicode.IsSpace(r):
					saveBoolOrString()
					st = stArgEndOfToken
				case r == ':':
					s.next()
					saveName()
					st = stArg
				}
			default:
				break parseLoop
			}
		case stArgEndOfToken:
			switch {
			case unicode.IsSpace(r):
				s.next()
				continue
			case r == ',':
				s.next()
				st = stArg
			case r == ')':
				s.next()
				incomplete = false
				st = stMaybeValue
			default:
				break parseLoop
			}
		case stMaybeValue:
			switch {
			case r == '=':
				s.next()
				hasValue = true
				st = stValue
			case unicode.IsSpace(r):
				s.next()
				st = stMaybeComment
			default:
				break parseLoop
			}
		case stValue:
			switch {
			case opts.rawValues: // When enabled, consume all remaining chars
				value.WriteRune(s.next())
			case r == '+' && isIdentBegin(s.peekN(1)): // tag value
				s.next() // consume +
				if err := saveTag(); err != nil {
					return TypedTag{}, err
				}
				st = stTag
			case r == '-' || r == '+' || unicode.IsDigit(r):
				number, err := s.nextNumber()
				valueType = ValueTypeInt
				if err != nil {
					return TypedTag{}, err
				}
				value.WriteString(number)
				st = stMaybeComment
			case r == '"' || r == '`':
				str, err := s.nextString()
				if err != nil {
					return TypedTag{}, err
				}
				value.WriteString(str)
				valueType = ValueTypeString
				st = stMaybeComment
			case isIdentBegin(r):
				str, err := s.nextIdent(isIdentInterior)
				if err != nil {
					return TypedTag{}, err
				}
				value.WriteString(str)
				valueType = ValueTypeString
				st = stMaybeComment
			default:
				break parseLoop
			}
		case stMaybeComment:
			switch {
			case unicode.IsSpace(r):
				s.next()
				continue
			case r == '/':
				s.next()
				incomplete = true
				st = stTrailingSlash
			default:
				break parseLoop
			}
		case stTrailingSlash:
			switch {
			case r == '/':
				s.next()
				incomplete = false
				st = stTrailingComment
			default:
				break parseLoop
			}
		case stTrailingComment:
			s.next()
			s.pos = len(s.buf)
			break parseLoop
		default:
			return TypedTag{}, fmt.Errorf("unexpected internal parser error: unknown state: %s at position %d", st, s.pos)
		}
	}
	if s.peek() != EOF {
		return TypedTag{}, fmt.Errorf("unexpected character %q at position %d", s.next(), s.pos)
	}
	if incomplete {
		return TypedTag{}, fmt.Errorf("unexpected end of input")
	}
	if err := saveTag(); err != nil {
		return TypedTag{}, err
	}
	if hasValue {
		saveValue()
	}
	if startTag == nil {
		return TypedTag{}, fmt.Errorf("unexpected internal parser error: no tags parsed")
	}
	return *startTag, nil
}

func isIdentBegin(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

func isIdentInterior(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '.' || r == '-'
}

func isTagNameInterior(r rune) bool {
	return isIdentInterior(r) || r == ':'
}
