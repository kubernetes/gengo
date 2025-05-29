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
	"strconv"
	"strings"
	"unicode"
)

// Parse parses a tag string into a TypedTag, or returns an error if the tag
// string fails to parse.
//
// Any enabled ParseOptions modify the behavior of the parser. The below
// describes only the default behavior.
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
// Trailing comments are ignored unless opts.ParseValues=false, in which case they
// are treated as part of the value. See ParseValues.RawValues for details.
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
func Parse(tag string, opts ParseOptions) (TypedTag, error) {
	tag = strings.TrimSpace(tag)
	return parseTag(tag, opts)
}

// ParseOptions controls the behavior of Parse.
type ParseOptions struct {
	// RawValues disables parsing of the value part of the tag. If true, the Value
	// field will contain all text following the "=" sign, up to the last
	// non-whitespace character, and ValueType will be set to ValueTypeRaw.
	RawValues bool
}

// ParseAll calls Parse on each tag in the input slice.
func ParseAll(tags []string, opts ParseOptions) ([]TypedTag, error) {
	var out []TypedTag
	for _, tag := range tags {
		parsed, err := Parse(tag, opts)
		if err != nil {
			return nil, err
		}
		out = append(out, parsed)
	}
	return out, nil
}

const (
	stBegin = "stBegin"
	stTag   = "stTag"
	stArg   = "stArg"

	// arg value parsing states
	stArgNumber         = "stArgNumber"
	stArgPrefixedNumber = "stArgPrefixedNumber"
	stArgQuotedString   = "stArgQuotedString"
	stArgNakedString    = "stArgNakedString"
	stArgEscape         = "stArgEscape"
	stArgEndOfToken     = "stArgEndOfToken"

	// tag value parsing states
	stMaybeValue          = "stMaybeValue"
	stValue               = "stValue"
	stValueTagOrNumber    = "stValueTagOrNumber"
	stValueNumber         = "stValueNumber"
	stValuePrefixedNumber = "stValuePrefixedNumber"
	stValueQuotedString   = "stValueQuotedString"
	stValueNakedString    = "stValueNakedString"
	stValueEscape         = "stValueEscape"

	stMaybeComment    = "stMaybeComment"
	stTrailingSlash   = "stTrailingSlash"
	stTrailingComment = "stTrailingComment"
)

func parseTag(input string, opts ParseOptions) (TypedTag, error) {
	var startTag, endTag *TypedTag // both ends of the chain when parsing chained tags

	tag := bytes.Buffer{}   // current tag name
	args := []Arg{}         // all tag arguments
	value := bytes.Buffer{} // current tag value
	var valueType ValueType // current value type
	var hasValue bool       // true if the tag has a value

	cur := Arg{}          // current argument accumulator
	buf := bytes.Buffer{} // string accumulator

	// These are defined outside the loop to make errors easier.
	var i int
	var r rune
	var incomplete bool
	var quote rune

	saveInt := func() error {
		s := buf.String()
		if _, err := strconv.ParseInt(s, 0, 64); err != nil {
			return fmt.Errorf("invalid number %q", s)
		} else {
			cur.Value = s
			cur.Type = ArgTypeInt
		}
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
		args = []Arg{}
		tag.Reset()
		return nil
	}
	saveValue := func() {
		endTag.Value = value.String()
		if opts.RawValues {
			endTag.ValueType = ValueTypeRaw
			return
		}
		endTag.ValueType = valueType
		if valueType == ValueTypeString && (endTag.Value == "true" || endTag.Value == "false") {
			endTag.ValueType = ValueTypeBool
		}
	}
	runes := []rune(input)
	st := stBegin
parseLoop:
	for i, r = range runes {
		switch st {
		case stBegin:
			switch {
			case unicode.IsSpace(r):
				continue
			case isIdentBegin(r):
				tag.WriteRune(r)
				st = stTag
			default:
				break parseLoop
			}
		case stTag:
			switch {
			case isIdentInterior(r) || r == ':':
				tag.WriteRune(r)
			case r == '(':
				incomplete = true
				st = stArg
			case r == '=':
				hasValue = true
				st = stValue
			case unicode.IsSpace(r):
				st = stMaybeComment
			default:
				break parseLoop
			}
		case stArg:
			switch {
			case unicode.IsSpace(r):
				continue
			case r == ')':
				incomplete = false
				st = stMaybeValue
			case r == '0':
				buf.WriteRune(r)
				st = stArgPrefixedNumber
			case r == '-' || r == '+' || unicode.IsDigit(r):
				buf.WriteRune(r)
				st = stArgNumber
			case r == '"' || r == '`':
				quote = r
				st = stArgQuotedString
			case isIdentBegin(r):
				buf.WriteRune(r)
				st = stArgNakedString
			default:
				break parseLoop
			}
		case stArgNumber:
			hexits := "abcdefABCDEF"
			switch {
			case unicode.IsDigit(r) || strings.Contains(hexits, string(r)):
				buf.WriteRune(r)
				continue
			case r == ',':
				if err := saveInt(); err != nil {
					return TypedTag{}, err
				}
				st = stArg
			case r == ')':
				if err := saveInt(); err != nil {
					return TypedTag{}, err
				}
				incomplete = false
				st = stMaybeValue
			case unicode.IsSpace(r):
				if err := saveInt(); err != nil {
					return TypedTag{}, err
				}
				st = stArgEndOfToken
			default:
				break parseLoop
			}
		case stArgPrefixedNumber:
			switch {
			case unicode.IsDigit(r):
				buf.WriteRune(r)
				st = stArgNumber
			case r == 'x' || r == 'o' || r == 'b':
				buf.WriteRune(r)
				st = stArgNumber
			default:
				break parseLoop
			}
		case stArgQuotedString:
			switch {
			case r == '\\':
				st = stArgEscape
			case r == quote:
				saveString()
				st = stArgEndOfToken
			default:
				buf.WriteRune(r)
			}
		case stArgEscape:
			switch {
			case r == quote || r == '\\':
				buf.WriteRune(r)
				st = stArgQuotedString
			default:
				return TypedTag{}, fmt.Errorf("unhandled escaped character %q", r)
			}
		case stArgNakedString:
			switch {
			case isIdentInterior(r):
				buf.WriteRune(r)
			case r == ',':
				saveBoolOrString()
				st = stArg
			case r == ')':
				saveBoolOrString()
				incomplete = false
				st = stMaybeValue
			case unicode.IsSpace(r):
				saveBoolOrString()
				st = stArgEndOfToken
			case r == ':':
				saveName()
				st = stArg
			default:
				break parseLoop
			}
		case stArgEndOfToken:
			switch {
			case unicode.IsSpace(r):
				continue
			case r == ',':
				st = stArg
			case r == ')':
				incomplete = false
				st = stMaybeValue
			default:
				break parseLoop
			}
		case stMaybeValue:
			switch {
			case r == '=':
				hasValue = true
				st = stValue
			case unicode.IsSpace(r):
				st = stMaybeComment
			default:
				break parseLoop
			}
		case stValue:
			switch {
			case opts.RawValues: // When enabled, consume all remaining chars
				value.WriteRune(r)
			case r == '+':
				st = stValueTagOrNumber // Might be a tag or a number so stValueTagOrNumber peeks
			case r == '0':
				value.WriteRune(r)
				valueType = ValueTypeInt
				st = stValuePrefixedNumber
			case r == '-' || unicode.IsDigit(r):
				value.WriteRune(r)
				valueType = ValueTypeInt
				st = stValueNumber
			case r == '"' || r == '`':
				quote = r
				valueType = ValueTypeString
				st = stValueQuotedString
			case isIdentBegin(r):
				value.WriteRune(r)
				valueType = ValueTypeString
				st = stValueNakedString
			default:
				break parseLoop
			}
		case stValueTagOrNumber: // Both tags and numbers can start with a +
			switch {
			case unicode.IsDigit(r):
				value.WriteRune(r)
				st = stValueNumber
			case isIdentBegin(r):
				if err := saveTag(); err != nil {
					return TypedTag{}, err
				}
				incomplete = false
				st = stMaybeValue
				tag.WriteRune(r)
				st = stTag
			default:
				break parseLoop
			}
		case stValueNumber:
			hexits := "abcdefABCDEF"
			switch {
			case unicode.IsDigit(r) || strings.Contains(hexits, string(r)):
				value.WriteRune(r)
				continue
			case unicode.IsSpace(r):
				st = stMaybeComment
			default:
				break parseLoop
			}
		case stValuePrefixedNumber:
			switch {
			case unicode.IsDigit(r):
				value.WriteRune(r)
				st = stValueNumber
			case r == 'x' || r == 'o' || r == 'b':
				value.WriteRune(r)
				st = stValueNumber
			default:
				break parseLoop
			}
		case stValueQuotedString:
			switch {
			case r == '\\':
				st = stValueEscape
			case r == quote:
				st = stMaybeComment
			default:
				value.WriteRune(r)
			}
		case stValueEscape:
			switch {
			case r == quote || r == '\\':
				value.WriteRune(r)
				st = stValueQuotedString
			default:
				return TypedTag{}, fmt.Errorf("unhandled escaped character %q", r)
			}
		case stValueNakedString:
			switch {
			case isIdentInterior(r):
				value.WriteRune(r)
			case unicode.IsSpace(r):
				st = stMaybeComment
			default:
				break parseLoop
			}
		case stMaybeComment:
			switch {
			case unicode.IsSpace(r):
				continue
			case r == '/':
				incomplete = true
				st = stTrailingSlash
			default:
				break parseLoop
			}
		case stTrailingSlash:
			switch {
			case r == '/':
				incomplete = false
				st = stTrailingComment
			default:
				break parseLoop
			}
		case stTrailingComment:
			i = len(runes) - 1
			break parseLoop
		default:
			return TypedTag{}, fmt.Errorf("unexpected internal parser error: unknown state: %s at position %d", st, i)
		}
	}
	if i != len(runes)-1 {
		return TypedTag{}, fmt.Errorf("unexpected character %q at position %d", r, i)
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
		return TypedTag{}, fmt.Errorf("unexpected internal parser error: no start tag")
	}
	return *startTag, nil
}

func isIdentBegin(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

func isIdentInterior(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '.' || r == '-'
}
