/*
Copyright 2015 The Kubernetes Authors.

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

package gengo

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
)

// ExtractCommentTags parses comments for lines of the form:
//
//	'marker' + "key=value".
//
// Values are optional; "" is the default.  A tag can be specified more than
// one time and all values are returned.  If the resulting map has an entry for
// a key, the value (a slice) is guaranteed to have at least 1 element.
//
// Example: if you pass "+" for 'marker', and the following lines are in
// the comments:
//
//	+foo=value1
//	+bar
//	+foo=value2
//	+baz="qux"
//
// Then this function will return:
//
//	map[string][]string{"foo":{"value1, "value2"}, "bar": {""}, "baz": {`"qux"`}}
//
// Deprecated: Use ExtractFunctionStyleCommentTags.
func ExtractCommentTags(marker string, lines []string) map[string][]string {
	out := map[string][]string{}
	for _, line := range lines {
		line = strings.Trim(line, " ")
		if len(line) == 0 {
			continue
		}
		if !strings.HasPrefix(line, marker) {
			continue
		}
		kv := strings.SplitN(line[len(marker):], "=", 2)
		if len(kv) == 2 {
			out[kv[0]] = append(out[kv[0]], kv[1])
		} else if len(kv) == 1 {
			out[kv[0]] = append(out[kv[0]], "")
		}
	}
	return out
}

// ExtractSingleBoolCommentTag parses comments for lines of the form:
//
//	'marker' + "key=value1"
//
// If the tag is not found, the default value is returned.  Values are asserted
// to be boolean ("true" or "false"), and any other value will cause an error
// to be returned.  If the key has multiple values, the first one will be used.
func ExtractSingleBoolCommentTag(marker string, key string, defaultVal bool, lines []string) (bool, error) {
	tags, err := ExtractFunctionStyleCommentTags(marker, []string{key}, lines)
	if err != nil {
		return false, err
	}
	values := tags[key]
	if values == nil {
		return defaultVal, nil
	}
	if values[0].Value == "true" {
		return true, nil
	}
	if values[0].Value == "false" {
		return false, nil
	}
	return false, fmt.Errorf("tag value for %q is not boolean: %q", key, values[0])
}

// ExtractFunctionStyleCommentTags parses comments for special metadata tags. The
// marker argument should be unique enough to identify the tags needed, and
// should not be a marker for tags you don't want, or else the caller takes
// responsibility for making that distinction.
//
// The tagNames argument is a list of specific tags being extracted. If this is
// nil or empty, all lines which match the marker are considered.  If this is
// specified, only lines with begin with marker + one of the tags will be
// considered.  This is useful when a common marker is used which may match
// lines which fail this syntax (e.g. which predate this definition).
//
// This function looks for input lines of the following forms:
//   - 'marker' + "key=value"
//   - 'marker' + "key()=value"
//   - 'marker' + "key(arg)=value"
//   - 'marker' + "key(`raw string`)=value"
//
// The arg is optional. It may be a Go identifier, a raw string literal enclosed
// in back-ticks or double-quotes, an integer (with support for hex, octal, or
// binary notation) or a boolean. If not specified (either as "key=value" or as
// "key()=value"), the resulting Tag will have an empty Args list.
//
// The value is optional.  If not specified, the resulting Tag will have "" as
// the value.
//
// Tag comment-lines may have a trailing end-of-line comment.
//
// The map returned here are keyed by the Tag's name without args.
//
// A tag can be specified more than one time and all values are returned.  If
// the resulting map has an entry for a key, the value (a slice) is guaranteed
// to have at least 1 element.
//
// Example: if you pass "+" as the marker, and the following lines are in
// the comments:
//
//	+foo=val1  // foo
//	+bar
//	+foo=val2  // also foo
//	+baz="qux"
//	+foo(arg)  // still foo
//
// Then this function will return:
//
//		map[string][]Tag{
//	 	"foo": []Tag{{
//				Name: "foo",
//				Args: nil,
//				Value: "val1",
//			}, {
//				Name: "foo",
//				Args: nil,
//				Value: "val2",
//			}, {
//				Name: "foo",
//				Args: []string{"arg"},
//				Value: "",
//			}, {
//				Name: "bar",
//				Args: nil,
//				Value: ""
//			}, {
//				Name: "baz",
//				Args: nil,
//				Value: "\"qux\""
//		   }}
//
// This function should be preferred instead of ExtractCommentTags.
func ExtractFunctionStyleCommentTags(marker string, tagNames []string, lines []string) (map[string][]Tag, error) {
	out := map[string][]Tag{}

	tagIdents := make([]TagIdentifier, len(tagNames))
	for i, name := range tagNames {
		tagIdents[i] = TagIdentifierFromString(name)
	}

	tags := ExtractTags(marker, nil, lines)
	for tagIdent, tagLines := range tags {
		if len(tagIdents) > 0 && !slices.Contains(tagIdents, tagIdent) {
			continue
		}
		for _, line := range tagLines {
			tag, err := ParseTagWithArgs(line)
			if err != nil {
				return nil, err
			}

			var stringArgs []string
			for _, arg := range tag.Args {
				if len(arg.Name) > 0 {
					return nil, fmt.Errorf("unexpected named argument: %q", arg.Name)
				}
				if s, ok := arg.Value.(string); !ok {
					return nil, fmt.Errorf("unexpected argument type: %T", arg.Value)
				} else {
					stringArgs = append(stringArgs, s)
				}
			}
			tag.Name = tagIdent.String()
			out[tag.Name] = append(out[tag.Name], Tag{
				Name:  tag.Name,
				Args:  stringArgs,
				Value: tag.Value,
			})
		}
	}

	return out, nil
}

// Tag represents a single comment tag.
type Tag struct {
	// Name is the name of the tag with no arguments.
	Name string
	// Args is a list of optional arguments to the tag.
	Args []string
	// Value is the value of the tag.
	Value string
}

func (t Tag) String() string {
	buf := bytes.Buffer{}
	buf.WriteString(t.Name)
	if len(t.Args) > 0 {
		buf.WriteString("(")
		for i, a := range t.Args {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(a)
		}
		buf.WriteString(")")
	}
	return buf.String()
}

// ExtractTags find comment lines that match the marker and group prefix. ExtractTags
// then returns comment lines, with the marker and group prefixes removed, grouped by group/name identifier keys.
// If group is nil, all groups are returned.  A group of "" represents comment lines without a group prefix.
//
// For example, the comment "+k8s:required" has a marker of "+" and a group of "k8s".
//
// Example: When called with "+" for 'marker', and "k8s" for group for these comment lines:
//
//	Comment line without marker
//	+k8s:required
//	+listType=set
//	+k8s:format=k8s-long-name
//
// Then this function will return:
//
//	map[TagIdentifier][]string{
//		{Group: "k8s", Name: "required"}: {"required"},
//		{Group: "k8s", Name: "required"}: {"format=k8s-long-name"},
//	}
func ExtractTags(marker string, group *string, lines []string) map[TagIdentifier][]string {
	out := map[TagIdentifier][]string{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, marker) {
			continue
		}
		line = line[len(marker):]

		nameEnd := len(line)
		if idx := strings.IndexAny(line, "(="); idx > 0 {
			nameEnd = idx
		}
		ident := TagIdentifierFromString(line[:nameEnd])

		if group != nil && ident.Group != *group {
			continue
		}
		if ident.Group != "" {
			line = line[len(ident.Group)+1:]
		}

		out[ident] = append(out[ident], line)
	}
	return out
}

// ParseTagWithArgs parses a single comment tag with typed args.
// The tagText must not have a marker or group prefix.
//
// The tag may optionally contain function style arguments after the tag name.
// Arguments may either be a single positional argument or any number of named
// arguments, but not both. The tag may optionally contain a value following "=".
// The value text is not parsed and is returned as a string.
// Argument values may be double-quoted strings, backtick-quoted strings,
// integers, booleans, or identifiers.
//
// Examples:
//
//	"tagName"
//	"tagName=value"
//	"tagName()"
//	"tagName(arg)=value"
//	"tagName("double-quoted")"
//	"tagName(`backtick-quoted`)"
//	"tagName(100)"
//	"tagName(true)"
//	"name(key1: value1)=value"
//	"name(key1: value1, key2: value2)"
//	"name(key1:`string value`)"
//	"name(key1: 1)"
//	"name(key1: true)"
//
// The tag grammar is:
//
// <tag> ::= <name> { "(" { <args> "}" ")" } { "=" <tagValue> }
// <args> ::= <argValue> | <namedArgs>
// <namedArgs> ::= <argNameAndValue> { "," <namedArgs> }
// <argNameAndValue> ::= <identifier> ":" <argValue>
// <argValue> ::= <identifier> | <string> | <int> | <bool>
//
// <identifier> ::= [a-zA-Z_][a-zA-Z0-9_]*
// <string> ::= [`...` and "..." quoted strings with \\ and \" escaping]
// <int> ::= [decimal, hex (0x), octal (0o) or binary (0b) notation with optional +/- prefix]
// <bool> ::= "true" | "false"
// <tagValue> ::= [all text after the = sign]
func ParseTagWithArgs(tagText string) (TypedTag, error) {
	tagText = strings.TrimSpace(tagText)
	parsed, err := parseTagKey(tagText)
	if err != nil {
		return TypedTag{}, err
	}
	return TypedTag{Name: parsed.name, Args: parsed.args, Value: parsed.value}, nil
}

// ExtractAndParseTagWithArgs combines ExtractTags and ParseTagWithArgs.
func ExtractAndParseTagWithArgs(marker string, group *string, lines []string) (map[TagIdentifier][]TypedTag, error) {
	out := map[TagIdentifier][]TypedTag{}
	for ident, lines := range ExtractTags(marker, group, lines) {
		for _, line := range lines {
			tag, err := ParseTagWithArgs(line)
			if err != nil {
				return nil, err
			}
			out[ident] = append(out[ident], tag)
		}
	}
	return out, nil
}

// TagIdentifier represents a tag name with an optional group prefix.
// Represented in string form as "<group>:<name>" or just "<name>".
type TagIdentifier struct {
	// Group is the optional group prefix of the tag.
	Group string
	Name  string
}

// TagIdentifierFromString parses a TagIdentifier from string form.
// tagName may be of form "<group>:<name>" or just "<name>".
func TagIdentifierFromString(tagName string) TagIdentifier {
	var ident TagIdentifier
	group, name, ok := strings.Cut(tagName, ":")
	if ok {
		ident = TagIdentifier{Group: group, Name: name}
	} else {
		ident = TagIdentifier{Group: "", Name: tagName}
	}
	return ident
}

func (t TagIdentifier) String() string {
	if len(t.Group) > 0 {
		return fmt.Sprintf("%s:%s", t.Group, t.Name)
	}
	return t.Name
}

// TypedTag represents a single comment tag with typed args.
type TypedTag struct {
	// Name is the name of the tag with no arguments.
	Name string
	// Args is a list of optional arguments to the tag.
	Args []Arg
	// Value is the value of the tag.
	Value string
}

func (t TypedTag) String() string {
	buf := strings.Builder{}
	buf.WriteString(t.Name)
	if len(t.Args) > 0 {
		buf.WriteString("(")
		for i, a := range t.Args {
			buf.WriteString(a.String())
			if i > 0 {
				buf.WriteString(", ")
			}
		}
		buf.WriteString(")")
	}
	if len(t.Value) > 0 {
		buf.WriteString("=")
		buf.WriteString(t.Value)
	}
	return buf.String()
}

// Arg represents a typed argument
type Arg struct {
	// Name is the name of a named argument. This is zero-valued for positional arguments.
	Name string
	// Value is the string value of an argument. It has been validated to match the Type.
	// Value may be a string, int, or bool.
	Value any
}

func (a Arg) String() string {
	if len(a.Name) > 0 {
		return fmt.Sprintf("%s: %v", a.Name, a.Value)
	}
	return fmt.Sprintf("%v", a.Value)
}
