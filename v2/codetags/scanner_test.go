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
	"testing"
)

func TestScannerBasics(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []rune
	}{
		{
			name:  "empty",
			input: "",
			want:  []rune{EOF},
		},
		{
			name:  "single character",
			input: "a",
			want:  []rune{'a', EOF},
		},
		{
			name:  "multiple characters",
			input: "abc",
			want:  []rune{'a', 'b', 'c', EOF},
		},
		{
			name:  "unicode characters",
			input: "αβγ",
			want:  []rune{'α', 'β', 'γ', EOF},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := scanner{buf: []rune(tt.input)}

			// next()
			for _, want := range tt.want {
				got := s.next()
				if got != want {
					t.Errorf("next() = %v, want %v", got, want)
				}
			}

			// peek()
			s = scanner{buf: []rune(tt.input)}
			for i, want := range tt.want {
				got := s.peek()
				if got != want {
					t.Errorf("peek() at position %d = %v, want %v", i, got, want)
				}
				s.next() // Advance scanner
			}

			// peekN()
			if len(tt.input) > 1 {
				s = scanner{buf: []rune(tt.input)}
				got := s.peekN(1)
				want := []rune(tt.input)[1]
				if got != want {
					t.Errorf("peekN(1) at position 0 = %v, want %v", got, want)
				}
			}
		})
	}
}

func TestNextNumber(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "simple integer",
			input: "123",
		},
		{
			name:  "negative integer",
			input: "-123",
		},
		{
			name:  "positive integer with plus sign",
			input: "+123",
			want:  "123",
		},
		{
			name:  "zero",
			input: "0",
		},
		{
			name:  "hex number",
			input: "0xFF",
		},
		{
			name:  "octal number",
			input: "0o77",
		},
		{
			name:  "binary number",
			input: "0b101",
		},
		{
			name:    "incomplete hex",
			input:   "0x",
			wantErr: true,
		},
		{
			name:    "incomplete octal",
			input:   "0o",
			wantErr: true,
		},
		{
			name:    "incomplete binary",
			input:   "0b",
			wantErr: true,
		},
		{
			name:    "number followed by non-digit",
			input:   "123abc",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.want == "" {
				tt.want = tt.input
			}
			s := scanner{buf: []rune(tt.input)}
			got, err := s.nextNumber()

			if (err != nil) != tt.wantErr {
				t.Errorf("nextNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != tt.want {
				t.Errorf("nextNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNextString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "double quoted string",
			input: `"hello"`,
			want:  "hello",
		},
		{
			name:  "backtick quoted string",
			input: "`hello`",
			want:  "hello",
		},
		{
			name:  "empty double quoted string",
			input: `""`,
			want:  "",
		},
		{
			name:  "empty backtick quoted string",
			input: "``",
			want:  "",
		},
		{
			name:  "string with escaped quote",
			input: `"hello \"world\""`,
			want:  `hello "world"`,
		},
		{
			name:  "string with escaped backslash",
			input: `"hello \\world"`,
			want:  `hello \world`,
		},
		{
			name:    "unterminated double quoted string",
			input:   `"hello`,
			wantErr: true,
		},
		{
			name:    "unterminated backtick quoted string",
			input:   "`hello",
			wantErr: true,
		},
		{
			name:    "invalid escape sequence",
			input:   `"hello \n world"`,
			wantErr: true,
		},
		{
			name:  "string followed by other content",
			input: `"hello" world`,
			want:  "hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := scanner{buf: []rune(tt.input)}
			got, err := s.nextString()

			if (err != nil) != tt.wantErr {
				t.Errorf("nextString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != tt.want {
				t.Errorf("nextString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNextIdent(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		isInteriorFn func(rune) bool
		want         string
		wantErr      bool
	}{
		{
			name:  "simple identifier with isIdentInterior",
			input: "abc",
			want:  "abc",
		},
		{
			name:  "identifier with underscore using isIdentInterior",
			input: "abc_def",
			want:  "abc_def",
		},
		{
			name:  "identifier with dash using isIdentInterior",
			input: "abc-def",
			want:  "abc-def",
		},
		{
			name:  "identifier with dot using isIdentInterior",
			input: "abc.def",
			want:  "abc.def",
		},
		{
			name:  "identifier with numbers using isIdentInterior",
			input: "abc123",
			want:  "abc123",
		},
		{
			name:  "identifier with colon using isIdentInterior",
			input: "abc:def",
			want:  "abc",
		},
		{
			name:  "identifier followed by invalid character",
			input: "abc@def",
			want:  "abc",
		},
		{
			name:  "identifier starting with underscore",
			input: "_abc",
			want:  "_abc",
		},
		{
			name:    "identifier starting with number",
			input:   "123abc",
			want:    "",
			wantErr: true,
		},

		{
			name:         "identifier with colon using isTagNameInterior",
			input:        "abc:def",
			isInteriorFn: isTagNameInterior,
			want:         "abc:def",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := scanner{buf: []rune(tt.input)}
			var got string
			var err error

			if tt.isInteriorFn == nil {
				tt.isInteriorFn = isIdentInterior
			}
			got, err = s.nextIdent(tt.isInteriorFn)

			if (err != nil) != tt.wantErr {
				t.Errorf("nextIdent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != tt.want {
				t.Errorf("nextIdent() = %v, want %v", got, tt.want)
			}
		})
	}
}
