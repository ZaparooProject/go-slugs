// Copyright 2026 The Zaparoo Project Contributors.
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package filetitle

import "strings"

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func isAlphanumeric(b byte) bool {
	return isDigit(b) || (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z')
}

func isWhitespace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r' || b == '\f' || b == '\v'
}

// isWordBoundary checks if the position i in s is at a word boundary.
// A word boundary exists between a regexp word char and a non-word char (or start/end).
func isWordBoundaryBefore(s string, i int) bool {
	if i == 0 {
		return true
	}
	return !isRegexWordByte(s[i-1])
}

func isWordBoundaryAfter(s string, i int) bool {
	if i >= len(s) {
		return true
	}
	return !isRegexWordByte(s[i])
}

// isRegexWordByte matches Go regexp's ASCII `\w` class used by `\b`.
func isRegexWordByte(b byte) bool {
	return isAlphanumeric(b) || b == '_'
}

// collapseSpaces replaces runs of whitespace with a single space.
// Also trims leading/trailing whitespace (unlike the regex it replaces).
// Callers always call strings.TrimSpace afterward, so this is equivalent.
// Replaces: reMultiSpace = regexp.MustCompile(`\s+`)
func collapseSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// stripLeadingNumberPrefix removes a leading number followed by separators.
// Replaces: reLeadingNum = regexp.MustCompile(`^\d+[.\s\-]+`)
// Used as: reLeadingNum.ReplaceAllString(title, "")
func stripLeadingNumberPrefix(s string) string {
	i := 0
	for i < len(s) && isDigit(s[i]) {
		i++
	}
	if i == 0 {
		return s
	}
	j := i
	for j < len(s) && (s[j] == '.' || s[j] == '-' || isWhitespace(s[j])) {
		j++
	}
	if j == i {
		return s
	}
	return s[j:]
}

// isYearValue checks if a 4-character string is a valid year in range 1970-2099.
func isYearValue(s string) bool {
	if len(s) != 4 || !isDigit(s[0]) || !isDigit(s[1]) || !isDigit(s[2]) || !isDigit(s[3]) {
		return false
	}
	if s[0] == '1' && s[1] == '9' && (s[2] >= '7') {
		return true
	}
	return s[0] == '2' && s[1] == '0'
}

// parseMatch holds the result of a string-based pattern match, replacing []int
// from FindStringSubmatchIndex. Fields are indices into the original string.
type parseMatch struct {
	start, end   int // full match bounds
	cap1s, cap1e int // first capture group
	ok           bool
}

// removeSeasonEpisode removes all S##E### patterns, replacing with space.
// Replaces: reSeasonEpisode.ReplaceAllString(title, " ")
func removeSeasonEpisode(s string) string {
	var b strings.Builder
	i := 0
	changed := false
	for i < len(s) {
		if (s[i] == 'S' || s[i] == 's') && i+3 < len(s) {
			pos := i + 1
			dStart := pos
			for pos < len(s) && pos-dStart < 2 && isDigit(s[pos]) {
				pos++
			}
			if pos > dStart && pos < len(s) && (s[pos] == 'E' || s[pos] == 'e') {
				pos++
				eStart := pos
				for pos < len(s) && pos-eStart < 3 && isDigit(s[pos]) {
					pos++
				}
				if pos > eStart {
					if !changed {
						b.Grow(len(s))
						_, _ = b.WriteString(s[:i])
						changed = true
					}
					_ = b.WriteByte(' ')
					i = pos
					continue
				}
			}
		}
		if changed {
			_ = b.WriteByte(s[i])
		}
		i++
	}
	if !changed {
		return s
	}
	return b.String()
}

// findUnbracketedYear finds a year (1970-2099) at word boundaries.
// Replaces: reYearScene = regexp.MustCompile(`\b(19[789]\d|20\d{2})\b`)
func findUnbracketedYear(s string) parseMatch {
	for i := 0; i <= len(s)-4; i++ {
		if !isWordBoundaryBefore(s, i) {
			continue
		}
		if !isYearValue(s[i : i+4]) {
			continue
		}
		if !isWordBoundaryAfter(s, i+4) {
			continue
		}
		return parseMatch{
			start: i, end: i + 4,
			cap1s: i, cap1e: i + 4,
			ok: true,
		}
	}
	return parseMatch{}
}
