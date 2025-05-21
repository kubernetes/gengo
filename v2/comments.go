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
	"errors"
	"fmt"
	"strings"
	"text/scanner"
	"unicode"
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
// The arg is optional.  It may be a Go identifier or a raw string literal
// enclosed in back-ticks.  If not specified (either as "key=value" or as
// "key()=value"), the resulting Tag will have an empty Args list.
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
//		}}
//
// This function should be preferred instead of ExtractCommentTags.
func ExtractFunctionStyleCommentTags(marker string, tagNames []string, lines []string) (map[string][]Tag, error) {
	stripTrailingComment := func(in string) string {
		parts := strings.SplitN(in, "//", 2)
		return strings.TrimSpace(parts[0])
	}

	out := map[string][]Tag{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if !strings.HasPrefix(line, marker) {
			continue
		}
		body := stripTrailingComment(line[len(marker):])
		key, val, err := splitKeyValScanner(body)
		if err != nil {
			return nil, err
		}

		tag := Tag{}
		if name, args, err := parseTagKey(key, tagNames); err != nil {
			return nil, err
		} else if name != "" {
			tag.Name, tag.Args = name, args
			tag.Value = val
			out[tag.Name] = append(out[tag.Name], tag)
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

// parseTagKey parses the key part of an extended comment tag, including
// optional arguments. The input is assumed to be the entire text of the
// original input after the marker, up to the '=' or end-of-line.
//
// The tags argument is an optional list of tag names to match. If it is nil or
// empty, all tags match.
//
// At the moment, arguments are very strictly formatted (see parseTagArgs) and
// whitespace is not allowed.
//
// This function returns the key name and arguments, unless tagNames was
// specified and the input did not match, in which case it returns "".
func parseTagKey(input string, tagNames []string) (string, []string, error) {
	parts := strings.SplitN(input, "(", 2)
	key := parts[0]

	if len(tagNames) > 0 {
		found := false
		for _, tn := range tagNames {
			if key == tn {
				found = true
				break
			}
		}
		if !found {
			return "", nil, nil
		}
	}

	var args []string
	if len(parts) == 2 {
		if ret, err := parseTagArgs(parts[1]); err != nil {
			return key, nil, fmt.Errorf("failed to parse tag args: %v", err)
		} else {
			args = ret
		}
	}
	return key, args, nil
}

// parseTagArgs parses the arguments part of an extended comment tag. The input
// is assumed to be the entire text of the original input after the opening
// '(', including the trailing ')'.
//
// At the moment this assumes that the entire string between the opening '('
// and the trailing ')' is a single Go-style identifier token OR a raw string
// literal. The single Go-style token may consist only of letters and digits
// and whitespace is not allowed.
func parseTagArgs(input string) ([]string, error) {
	s := initArgScanner(input)
	var args []string
	if s.Peek() != ')' {
		// Arg found.
		arg, err := parseArg(s)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}
	// Expect one closing ')' after the arg.
	if s.Scan() != ')' {
		return nil, fmt.Errorf("no closing ')' found: %q", input)
	}
	// Expect no whitespace, etc. after the one ')'.
	if s.Scan() != scanner.EOF {
		pos := s.Pos().Offset - len(s.TokenText())
		return nil, fmt.Errorf("unexpected characters after ')': %q", input[pos:])
	}
	return args, nil
}

type argScanner struct {
	*scanner.Scanner
	errs []error
}

func initArgScanner(input string) *argScanner {
	s := &argScanner{Scanner: &scanner.Scanner{}}

	s.Init(strings.NewReader(input))
	s.Mode = scanner.ScanIdents | scanner.ScanRawStrings
	s.Whitespace = 0

	s.Error = func(_ *scanner.Scanner, msg string) {
		s.errs = append(s.errs,
			fmt.Errorf("error parsing %q at %v: %s", input, s.Position, msg))
	}
	return s
}

func (s *argScanner) unexpectedTokenError(expected string, token string) error {
	s.Error(s.Scanner, fmt.Sprintf("expected %s but got (%q)", expected, token))
	return errors.Join(s.errs...)
}

func parseArg(s *argScanner) (string, error) {
	switch tok := s.Scan(); tok {
	case scanner.RawString:
		return s.TokenText(), nil
	case scanner.Ident:
		txt := s.TokenText()
		for _, r := range txt {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
				return "", s.unexpectedTokenError("letter or digit", txt)
			}
		}
		return txt, nil
	case ',':
		return "", fmt.Errorf("multiple arguments are not supported")
	default:
		return "", s.unexpectedTokenError("Go-style identifier or raw string", s.TokenText())
	}
}

// splitKeyValScanner parses a tag body of the form key[=val]. It parses left to
// right and stops at the first "=" that is not inside a quoted or raw
// string literal. Text before that point becomes the key (trimmed of spaces).
// Text after becomes the val. If no "=" is found, the whole input
// is returned as key and val is empty. The parsing understands Go-style identifiers,
// and raw strings. Any other token or scanner error is
// reported to the caller.
func splitKeyValScanner(input string) (key, val string, err error) {
	var s scanner.Scanner
	s.Init(strings.NewReader(input))
	s.Mode = scanner.ScanIdents | scanner.ScanRawStrings
	for {
		switch tok := s.Scan(); tok {
		case scanner.EOF:
			return strings.TrimSpace(input), "", nil
		case '=':
			// Split at the first top-level '='.  Everything before (trimmed) is the
			// key, everything after (not trimmed) is the value.
			start := s.Pos().Offset - len(s.TokenText())
			key = strings.TrimSpace(input[:start])
			if start+len(s.TokenText()) < len(input) {
				val = input[start+len(s.TokenText()):]
			}
			return key, val, nil
		}
	}
}
