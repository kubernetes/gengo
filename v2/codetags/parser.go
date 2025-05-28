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

const (
	stBegin           = "stBegin"
	stTag             = "stTag"
	stArg             = "stArg"
	stNumber          = "stNumber"
	stPrefixedNumber  = "stPrefixedNumber"
	stQuotedString    = "stQuotedString"
	stNakedString     = "stNakedString"
	stEscape          = "stEscape"
	stEndOfToken      = "stEndOfToken"
	stMaybeValue      = "stMaybeValue"
	stValue           = "stValue"
	stMaybeComment    = "stMaybeComment"
	stTrailingSlash   = "stTrailingSlash"
	stTrailingComment = "stTrailingComment"
)

type tagKey struct {
	name  string
	args  []Arg
	value string
}

// Parse parses a comment tag into a TypedTag, or returns an error if the tag fails to parse.
//
// This function supports input of the following forms:
//
//	"key"
//	"key=value"
//	"key()=value"
//	"key(arg)=value"
//	"key(arg1: argValue1)=value"
//	"key(arg1: argValue1, arg2: argValue2)=value"
//
// When parsing Go comments, the Extract function it typically used to extract
// tags matching a prefix, when then can be parsed with this function.
//
// The tag may optionally contain function style arguments after the tag name.
// The arguments are optional. Arguments may either be a single positional
// argument or any number of named arguments, but not both. Argument values may
// be double-quoted strings, backtick-quoted strings, integers, booleans, or
// identifiers.
//
// The value is optional. If not specified, the resulting Tag will have "" as
// the value.
//
// A trailing comment is allowed if the tag has no value. That is, if tag does
// not end with "=<value>", then a " // <comment" is allowed.
//
// Examples:
//
//	"key("double-quoted") // comment is allowed here"
//	"key(`backtick-quoted`)"
//	"key(100)"
//	"key(true)"
//	"key(key1:`string value`)"
//	"key(key1: 1)"
//	"key(key1: true)"
//
// The tag grammar is:
//
// <tag> ::= <tagName> { "(" { <args> "}" ")" } { "=" <tagValue> }
// <args> ::= <argValue> | <namedArgs>
// <namedArgs> ::= <argNameAndValue> { "," <namedArgs> }
// <argNameAndValue> ::= <identifier> ":" <argValue>
// <argValue> ::= <identifier> | <string> | <int> | <bool>
//
// <tagName> ::= <identifier> { ":" <identifier> }
// <identifier> ::= [a-zA-Z_][a-zA-Z0-9_-.]*
// <string> ::= [`...` and "..." quoted strings with \\ and \" escaping]
// <int> ::= [decimal, hex (0x...), octal (0... or 0o...) or binary (0b...) notation with optional +/- prefix]
// <bool> ::= "true" | "false"
// <tagValue> ::= [all text after the = sign]
func Parse(tagText string) (TypedTag, error) {
	tagText = strings.TrimSpace(tagText)
	parsed, err := parseTagKey(tagText)
	if err != nil {
		return TypedTag{}, err
	}
	return TypedTag{Name: parsed.name, Args: parsed.args, Value: parsed.value}, nil
}

// ParseAll calls Parse on each tag in the input slice.
func ParseAll(tags []string) ([]TypedTag, error) {
	var out []TypedTag
	for _, tag := range tags {
		parsed, err := Parse(tag)
		if err != nil {
			return nil, err
		}
		out = append(out, parsed)
	}
	return out, nil
}

func parseTagKey(input string) (tagKey, error) {
	tag := bytes.Buffer{}   // current tag name
	args := []Arg{}         // all tag arguments
	value := bytes.Buffer{} // current tag value

	cur := Arg{}          // current argument accumulator
	buf := bytes.Buffer{} // string accumulator

	// These are defined outside the loop to make errors easier.
	var i int
	var r rune
	var incomplete bool
	var quote rune

	saveInt := func() error {
		s := buf.String()
		if ival, err := strconv.ParseInt(s, 0, 64); err != nil {
			return fmt.Errorf("invalid number %q", s)
		} else {
			cur.Value = IntValue{s: s, i: ival}
		}
		args = append(args, cur)
		cur = Arg{}
		buf.Reset()
		return nil
	}
	saveString := func() {
		s := buf.String()
		cur.Value = StringValue(s)
		args = append(args, cur)
		cur = Arg{}
		buf.Reset()
	}
	saveBoolOrString := func() {
		s := buf.String()
		if s == "true" {
			cur.Value = BoolValue(true)
		} else if s == "false" {
			cur.Value = BoolValue(false)
		} else {
			cur.Value = StringValue(s)
		}
		args = append(args, cur)
		cur = Arg{}
		buf.Reset()
	}
	saveName := func() {
		cur.Name = buf.String()
		buf.Reset()
	}

	runes := []rune(input)
	st := stBegin
parseLoop:
	for i, r = range runes {
		//fmt.Printf("state: %s, char: %q\n", st, r)
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
				st = stPrefixedNumber
			case r == '-' || r == '+' || unicode.IsDigit(r):
				buf.WriteRune(r)
				st = stNumber
			case r == '"' || r == '`':
				quote = r
				st = stQuotedString
			case isIdentBegin(r):
				buf.WriteRune(r)
				st = stNakedString
			default:
				break parseLoop
			}
		case stNumber:
			hexits := "abcdefABCDEF"
			switch {
			case unicode.IsDigit(r) || strings.Contains(hexits, string(r)):
				buf.WriteRune(r)
				continue
			case r == ',':
				if err := saveInt(); err != nil {
					return tagKey{}, err
				}
				st = stArg
			case r == ')':
				if err := saveInt(); err != nil {
					return tagKey{}, err
				}
				incomplete = false
				st = stMaybeValue
			case unicode.IsSpace(r):
				if err := saveInt(); err != nil {
					return tagKey{}, err
				}
				st = stEndOfToken
			default:
				break parseLoop
			}
		case stPrefixedNumber:
			switch {
			case unicode.IsDigit(r):
				buf.WriteRune(r)
				st = stNumber
			case r == 'x' || r == 'o' || r == 'b':
				buf.WriteRune(r)
				st = stNumber
			default:
				break parseLoop
			}
		case stQuotedString:
			switch {
			case r == '\\':
				st = stEscape
			case r == quote:
				saveString()
				st = stEndOfToken
			default:
				buf.WriteRune(r)
			}
		case stEscape:
			switch {
			case r == quote || r == '\\':
				buf.WriteRune(r)
				st = stQuotedString
			default:
				return tagKey{}, fmt.Errorf("unhandled escaped character %q", r)
			}
		case stNakedString:
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
				st = stEndOfToken
			case r == ':':
				saveName()
				st = stArg
			default:
				break parseLoop
			}
		case stEndOfToken:
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
				st = stValue
			case unicode.IsSpace(r):
				st = stMaybeComment
			default:
				break parseLoop
			}
		case stValue: // This is a terminal state, it consumes the rest of the input as an opaque value.
			value.WriteRune(r)
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
			return tagKey{}, fmt.Errorf("unexpected character %q at position %d", r, i)
		}
	}
	if i != len(runes)-1 {
		return tagKey{}, fmt.Errorf("unexpected character %q at position %d", r, i)
	}
	if incomplete {
		return tagKey{}, fmt.Errorf("unexpected end of input")
	}
	usingNamedArgs := false
	for i, arg := range args {
		if (usingNamedArgs && arg.Name == "") || (!usingNamedArgs && arg.Name != "" && i > 0) {
			return tagKey{}, fmt.Errorf("can't mix named and positional arguments")
		}
		if arg.Name != "" {
			usingNamedArgs = true
		}
	}
	if !usingNamedArgs && len(args) > 1 {
		return tagKey{}, fmt.Errorf("multiple arguments must use 'name: value' syntax")
	}
	return tagKey{name: tag.String(), args: args, value: value.String()}, nil
}

func isIdentBegin(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

func isIdentInterior(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '.' || r == '-'
}
