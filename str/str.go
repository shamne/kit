package str

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/eyecuelab/kit/runeset"
)

//Runes removes any runes listed from the string.
//Note that this converts to runes and back to UTF-8, so RemoveRunes(s) == s
//does not necessarially hold, since the code points may differ.
func RemoveRunes(s string, toRemove ...rune) string {
	buf := bytes.Buffer{}
	set := runeset.FromRunes(toRemove...)
	for _, r := range s {
		if !set.Contains(r) {
			buf.WriteRune(r)
		}
	}
	return buf.String()
}

//SRemoveRunes removes any of the runes contained in the string toRemove from s.
//Equivalent to RemoveRunes(s, []rune(toRemove))
func SRemoveRunes(s string, toRemove string) string {
	buf := bytes.Buffer{}
	set := runeset.FromString(s)
	for _, r := range s {
		if !set.Contains(r) {
			buf.WriteRune(r)
		}
	}
	return buf.String()
}

//RemoveWhiteSpace is an alias for RemoveASCIIWhiteSpace.
var RemoveWhiteSpace = RemoveASCIIWhiteSpace

//MapErr maps f(string)string, err across a slice. If any error results, it
// it immediately returns an empty slice and the error.
func MapErr(f func(string) (string, error), strings []string) ([]string, error) {
	for i, str := range strings {
		result, err := f(str)
		if err != nil {
			return []string{}, fmt.Errorf("MapErr: %v", err)
		}
		strings[i] = result
	}
	return strings, nil
}

//Map maps f(string) across the remaining arguments, returning [f(str) for str in str]
func Map(f func(string) string, strings []string) []string {
	for i, str := range strings {
		strings[i] = f(str)
	}
	return strings
}

type RuneDiff struct {
	a, b     rune
	position int
}

func (rd RuneDiff) String() string {
	return fmt.Sprintf(`(%s, %s)`, subIfNoChar(rd.a), subIfNoChar(rd.b))
}

func subIfNoChar(r rune) string {
	if r == noChar {
		return "NO_CHAR"
	}
	return string(r)
}

//Diffs represents a difference in runes between strings.
type Diffs []RuneDiff

func (d Diffs) String() string {
	strs := make([]string, len(d))
	for i, rd := range d {
		strs[i] = rd.String()
	}
	return fmt.Sprintf("[%s]", strings.Join(strs, ", "))

}

//noCharacter is an code point from the unicode private use block which we use to represent one part of a diff having no character.
const noChar rune = '\uE011'

//RuneDiff returns a list of the differing rune-pairs in a string
func runeDiff(s, q []rune) Diffs {
	var diffs Diffs
	short := min(len(s), len(q))
	for i := 0; i < short; i++ {
		c, d := s[i], q[i]
		if c != d {
			diffs = append(diffs, RuneDiff{c, d, i})
		}
	}
	if len(s) == short {
		for i, d := range q[short:] {
			diffs = append(diffs, RuneDiff{a: noChar, b: d, position: short + i})
		}
	} else {
		for i, c := range s[short:] {
			diffs = append(diffs, RuneDiff{a: c, b: noChar, position: short + i})
		}
	}
	return diffs
}

//Diff compares the runes in the string for differences. Call Diff.String() for a pretty-printed list.
//Diff(a, b) == null if there is no difference.
func Diff(a, b string) Diffs {
	return runeDiff([]rune(a), []rune(b))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
