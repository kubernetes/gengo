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
	"fmt"
	"strconv"
	"strings"
)

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
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(a.String())
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
	Value ArgValue
}

// String returns the string representation of the argument.
func (a Arg) String() string {
	if len(a.Name) > 0 {
		return fmt.Sprintf("%s: %v", a.Name, a.Value)
	}
	return fmt.Sprintf("%v", a.Value)
}

// ArgValue represents a typed value of an argument.
// It may be a string, int, or bool.
type ArgValue interface {
	// String returns the string representation of the value that was parsed from the tag.
	String() string
	// Type returns the type of the value.
	Type() ArgType
}

// ArgType represents the type of an argument.
type ArgType string

const (
	ArgTypeString ArgType = "string"
	ArgTypeInt    ArgType = "int"
	ArgTypeBool   ArgType = "bool"
)

type ArgString string

func (a ArgString) String() string {
	return string(a)
}

func (a ArgString) Type() ArgType {
	return ArgTypeString
}

// ArgInt represents an integer argument. The string representation is the original value
// that was parsed from the tag and may be in hex, octal, binary, or decimal notation.
type ArgInt struct {
	s string
	i int64
}

// MustArgInt parses an int argument. Panics if the argument is invalid.
func MustArgInt(s string) ArgInt {
	i, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		panic(fmt.Sprintf("invalid int value: %q", s))
	}
	return ArgInt{s: s, i: i}
}

func (a ArgInt) String() string {
	return a.s
}

// Int returns the integer value of the argument.
func (a ArgInt) Int() int64 {
	return a.i
}

func (a ArgInt) Type() ArgType {
	return ArgTypeInt
}

type ArgBool bool

func (a ArgBool) String() string {
	if a == true {
		return "true"
	} else {
		return "false"
	}
}

// Bool returns the boolean value of the argument.
func (a ArgBool) Bool() bool {
	return bool(a)
}

func (a ArgBool) Type() ArgType {
	return ArgTypeBool
}
