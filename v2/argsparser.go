package gengo

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
	stMaybeComment    = "stMaybeComment"
	stTrailingSlash   = "stTrailingSlash"
	stTrailingComment = "stTrailingComment"
)

func parseTagKey(input string) (string, []Arg, error) {
	tag := bytes.Buffer{} // current tag name
	args := []Arg{}       // all tag arguments

	cur := Arg{}          // current argument accumulator
	buf := bytes.Buffer{} // string accumulator

	// These are defined outside the loop to make errors easier.
	var i int
	var r rune
	var incomplete bool
	var quote rune

	isIdentBegin := func(r rune) bool {
		return unicode.IsLetter(r) || r == '_'
	}
	isIdentInterior := func(r rune) bool {
		return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
	}
	saveInt := func() error {
		s := buf.String()
		if _, err := strconv.ParseInt(s, 0, 64); err != nil {
			return fmt.Errorf("invalid number %q", s)
		}
		cur.Typ = TypeInt
		cur.Value = s
		args = append(args, cur)
		cur = Arg{}
		buf.Reset()
		return nil
	}
	saveString := func() {
		s := buf.String()
		cur.Typ = TypeString
		cur.Value = s
		args = append(args, cur)
		cur = Arg{}
		buf.Reset()
	}
	saveBoolOrString := func() {
		s := buf.String()
		if s == "true" || s == "false" {
			cur.Typ = TypeBool
		} else {
			cur.Typ = TypeString
		}
		cur.Value = s
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
			case isIdentInterior(r):
				tag.WriteRune(r)
			case r == '(':
				incomplete = true
				st = stArg
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
				st = stMaybeComment
			case r == '0':
				buf.WriteRune(r)
				st = stPrefixedNumber
			case unicode.IsDigit(r):
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
					return "", nil, err
				}
				st = stArg
			case r == ')':
				if err := saveInt(); err != nil {
					return "", nil, err
				}
				incomplete = false
				st = stMaybeComment
			case unicode.IsSpace(r):
				if err := saveInt(); err != nil {
					return "", nil, err
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
				return "", nil, fmt.Errorf("unhandled escaped character %q", r)
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
				st = stMaybeComment
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
			panic("WTF: unknown state")
		}
	}
	if i != len(runes)-1 {
		return "", nil, fmt.Errorf("unexpected character %q at position %d", r, i)
	}
	if incomplete {
		return "", nil, fmt.Errorf("unexpected end of input")
	}
	usingNamedArgs := false
	for i, arg := range args {
		if (usingNamedArgs && arg.Name == "") || (!usingNamedArgs && arg.Name != "" && i > 0) {
			return "", nil, fmt.Errorf("can't mix named and positional arguments")
		}
		if arg.Name != "" {
			usingNamedArgs = true
		}
	}
	if !usingNamedArgs && len(args) > 1 {
		return "", nil, fmt.Errorf("multiple arguments must use 'name: value' syntax")
	}
	return tag.String(), args, nil
}
