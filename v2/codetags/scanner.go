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

type scanner struct {
	buf []rune
	pos int
}

func (s *scanner) next() rune {
	if s.pos >= len(s.buf) {
		return EOF
	}
	r := s.buf[s.pos]
	s.pos++
	return r
}

func (s *scanner) peek() rune {
	if s.pos >= len(s.buf) {
		return EOF
	}
	return s.buf[s.pos]
}

func (s *scanner) peekN(n int) rune {
	if s.pos >= len(s.buf) {
		return EOF
	}
	return s.buf[s.pos+n]
}

const (
	EOF = -1
)

const (
	stNumber       = "stNumber"
	stPrefixNumber = "stPrefixNumber"
)

func (s *scanner) nextNumber() (string, error) {
	var buf bytes.Buffer
	var incomplete bool
	st := stBegin

parseLoop:
	for r := s.peek(); r != EOF; r = s.peek() {
		switch st {
		case stBegin:
			switch {
			case r == '+':
				s.next() // consume +
			case r == '0':
				buf.WriteRune(s.next())
				st = stPrefixNumber
			case r == '-' || unicode.IsDigit(r):
				buf.WriteRune(s.next())
				st = stNumber
			default:
				break parseLoop
			}
		case stPrefixNumber:
			switch {
			case unicode.IsDigit(r):
				buf.WriteRune(s.next())
				st = stNumber
			case r == 'x' || r == 'o' || r == 'b':
				incomplete = true
				buf.WriteRune(s.next())
				st = stNumber
			default:
				break parseLoop
			}
		case stNumber:
			hexits := "abcdefABCDEF"
			switch {
			case unicode.IsDigit(r) || strings.Contains(hexits, string(r)):
				buf.WriteRune(s.next())
				incomplete = false
				continue
			default:
				break parseLoop
			}
		default:
			return "", fmt.Errorf("unexpected internal parser error: unknown state: %s at position %d", st, s.pos)
		}
	}
	if incomplete {
		return "", fmt.Errorf("unterminated number at position %d", s.pos)
	}
	return buf.String(), nil
}

const (
	stQuotedString = "stQuotedString"
	stEscape       = "stEscape"
)

func (s *scanner) nextString() (string, error) {
	var buf bytes.Buffer
	var quote rune
	var incomplete bool
	st := stBegin

parseLoop:
	for r := s.peek(); r != EOF; r = s.peek() {
		switch st {
		case stBegin:
			switch {
			case r == '"' || r == '`':
				incomplete = true
				quote = s.next() // consume quote
				st = stQuotedString
			default:
				return "", fmt.Errorf("expected string at position %d", s.pos)
			}
		case stQuotedString:
			switch {
			case r == '\\':
				s.next() // consume escape
				st = stEscape
			case r == quote:
				incomplete = false
				s.next()
				break parseLoop
			default:
				buf.WriteRune(s.next())
			}
		case stEscape:
			switch {
			case r == quote || r == '\\':
				buf.WriteRune(s.next())
				st = stQuotedString
			default:
				return "", fmt.Errorf("unhandled escaped character %q", r)
			}
		default:
			return "", fmt.Errorf("unexpected internal parser error: unknown state: %s at position %d", st, s.pos)
		}
	}
	if incomplete {
		return "", fmt.Errorf("unterminated string at position %d", s.pos)
	}
	return buf.String(), nil
}

const (
	stInterior = "stInterior"
)

func (s *scanner) nextIdent(isInteriorChar func(r rune) bool) (string, error) {
	var buf bytes.Buffer
	st := stBegin

parseLoop:
	for r := s.peek(); r != EOF; r = s.peek() {
		switch st {
		case stBegin:
			switch {
			case isIdentBegin(r):
				buf.WriteRune(s.next())
				st = stInterior
			}
		case stInterior:
			switch {
			case isInteriorChar(r):
				buf.WriteRune(s.next())
			default:
				break parseLoop
			}
		default:
			return "", fmt.Errorf("unexpected internal parser error: unknown state: %s at position %d", st, s.pos)
		}
	}
	return buf.String(), nil
}
