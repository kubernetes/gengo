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
	"strings"
)

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
