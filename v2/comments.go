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
	"strconv"
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
// in back-ticks, an integer (with support for hex, octal or binary notation) or
// a boolean. If not specified (either as "key=value" or as "key()=value"), the
// resulting Tag will have an empty Args list.
//
// The value is optional.  If not specified, the resulting Tag will have "" as
// the value.
//
// Tag comment-lines may have a trailing end-of-line comment.
//
// The map returned here is keyed by the Tag's name without args.
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
	typedTags, err := ExtractFunctionStyleCommentTypedTags(marker, tagNames, lines)
	if err != nil {
		return nil, err
	}
	out := map[string][]Tag{}
	for name, tags := range typedTags {
		for _, tag := range tags {
			var stringArgs []string
			for _, arg := range tag.Args {
				if len(arg.Name) > 0 {
					return nil, fmt.Errorf("unexpected named argument: %q", arg.Name)
				}
				if arg.Value.Type() != TypeString {
					return nil, fmt.Errorf("unexpected argument type: %q", arg.Value.Type())
				}
				stringArgs = append(stringArgs, arg.Value.String())
			}
			out[name] = append(out[name], Tag{
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

// ExtractFunctionStyleCommentTypedTags parses comments for special metadata tags.
// This function supports all the functionality of ExtractFunctionStyleCommentTags, and also supports named parameters
// and returns typed results.
func ExtractFunctionStyleCommentTypedTags(marker string, tagNames []string, lines []string) (map[string][]TypedTag, error) {
	stripTrailingComment := func(in string) string {
		idx := strings.LastIndex(in, "//")
		if idx == -1 {
			return in
		}
		return strings.TrimSpace(in[:idx])
	}

	out := map[string][]TypedTag{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, marker) {
			continue
		}
		line = line[len(marker):]
		parsed, err := parseTagKey(line)
		if err != nil {
			return nil, err
		}
		if parsed.name == "" {
			continue
		}
		if len(tagNames) > 0 && !slices.Contains(tagNames, parsed.name) {
			continue
		}

		var val string
		if parsed.valueStart > 0 {
			val = stripTrailingComment(line[parsed.valueStart:])
		}

		tag := TypedTag{Name: parsed.name, Args: parsed.args, Value: val}
		out[parsed.name] = append(out[parsed.name], tag)
	}
	return out, nil
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

// Arg represents a typed argument
type Arg struct {
	// Name is the name of a named argument. This is zero-valued for positional arguments.
	Name string
	// Value is the string value of an argument. It has been validated to match the Type.
	Value Value
}

type String string

func (String) Type() Type {
	return TypeString
}

func (s String) String() string {
	return string(s)
}

type Int int

func (Int) Type() Type {
	return TypeInt
}

func (i Int) String() string {
	return strconv.Itoa(int(i))
}

type Bool bool

func (Bool) Type() Type {
	return TypeBool
}

func (b Bool) String() string {
	return strconv.FormatBool(bool(b))
}

type Value interface {
	Type() Type
	String() string
}

// Type is the type of an arg.
type Type string

const (
	TypeString = "string"
	TypeInt    = "int"
	TypeBool   = "bool"
)
