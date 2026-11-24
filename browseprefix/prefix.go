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

package browseprefix

import (
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

const (
	KindNone Kind = ""
	KindRank Kind = "rank"
	KindDate Kind = "date"
)

const (
	DefaultThreshold = 0.5
	DefaultMinFiles  = 5
	// DefaultSampleLimit bounds how many paths are read to detect a directory's
	// prefix policy. The policy is a fraction-vs-threshold heuristic, so a sample of
	// this size estimates the ratio reliably without scanning directories that can
	// hold ~1M files on large libraries.
	DefaultSampleLimit = 2000
)

type Kind string

type Prefix struct {
	Kind Kind
	Rest string
}

type Policy struct {
	Kind    Kind
	Enabled bool
}

func StemFromPath(path string) string {
	base := filepath.Base(path)
	if idx := strings.LastIndex(base, "/"); idx >= 0 {
		base = base[idx+1:]
	}
	ext := filepath.Ext(base)
	if ext == "" {
		return base
	}
	return strings.TrimSuffix(base, ext)
}

func ParseStem(stem string) Prefix {
	trimmed := strings.TrimSpace(stem)
	if trimmed == "" {
		return Prefix{}
	}
	if prefix := parseDate(trimmed); prefix.Kind != KindNone {
		return prefix
	}
	return parseRank(trimmed)
}

func DetectPolicyForPaths(paths []string, threshold float64, minFiles int) Policy {
	if len(paths) < minFiles {
		return Policy{}
	}

	dateCount := 0
	rankCount := 0
	for _, path := range paths {
		prefix := ParseStem(StemFromPath(path))
		switch prefix.Kind {
		case KindDate:
			dateCount++
		case KindRank:
			rankCount++
		case KindNone:
		}
	}

	fileCount := float64(len(paths))
	if float64(dateCount)/fileCount > threshold {
		return Policy{Kind: KindDate, Enabled: true}
	}
	if float64(rankCount)/fileCount > threshold {
		return Policy{Kind: KindRank, Enabled: true}
	}
	return Policy{}
}

func StripWithPolicy(stem string, policy Policy) (string, bool) {
	if !policy.Enabled {
		return stem, false
	}
	prefix := ParseStem(stem)
	if prefix.Kind != policy.Kind || prefix.Rest == "" {
		return stem, false
	}
	return prefix.Rest, true
}

func parseRank(s string) Prefix {
	digitEnd := 0
	for digitEnd < len(s) && isASCIIDigit(s[digitEnd]) {
		digitEnd++
	}
	if digitEnd == 0 || digitEnd > 3 || digitEnd == len(s) {
		return Prefix{}
	}

	sepEnd := digitEnd
	for sepEnd < len(s) {
		ch := s[sepEnd]
		if ch != '.' && ch != '-' && ch != ')' && ch != '_' && !isSpaceByte(ch) {
			break
		}
		sepEnd++
	}
	if sepEnd == digitEnd || sepEnd == len(s) || !isListSeparator(s[digitEnd:sepEnd]) {
		return Prefix{}
	}

	rest := strings.TrimSpace(s[sepEnd:])
	if rest == "" {
		return Prefix{}
	}
	return Prefix{Kind: KindRank, Rest: rest}
}

// isListSeparator reports whether sep, what follows a leading number or a
// bare leading year, marks list numbering ("001 - Sonic", "42. Contra",
// "7) Zelda", "1991 - Sonic") rather than being part of the title. A bare
// space or a hyphen with no space beside it is a title's own: "720 Degrees",
// "10-Yard Fight" and "3-D WorldRunner" lost their numbers to a folder that
// was mostly such titles.
func isListSeparator(sep string) bool {
	if strings.ContainsAny(sep, ".)") {
		return true
	}
	return strings.Contains(sep, "-") && strings.ContainsFunc(sep, unicode.IsSpace)
}

func parseDate(s string) Prefix {
	if len(s) < 6 || !allDigits(s[:4]) {
		return Prefix{}
	}
	year, err := strconv.Atoi(s[:4])
	if err != nil || year < 1900 || year > 2099 {
		return Prefix{}
	}

	pos := 4
	if pos < len(s) && isDatePartSeparator(s[pos]) && pos+3 <= len(s) && allDigits(s[pos+1:pos+3]) {
		month, monthErr := strconv.Atoi(s[pos+1 : pos+3])
		if monthErr != nil || month < 1 || month > 12 {
			return Prefix{}
		}
		pos += 3
		if pos < len(s) && isDatePartSeparator(s[pos]) && pos+3 <= len(s) && allDigits(s[pos+1:pos+3]) {
			day, dayErr := strconv.Atoi(s[pos+1 : pos+3])
			if dayErr != nil || day < 1 || day > daysInMonth(year, month) {
				return Prefix{}
			}
			pos += 3
		}
	}

	if pos >= len(s) || !isSeparatorRune(rune(s[pos])) {
		return Prefix{}
	}
	rest := strings.TrimLeftFunc(s[pos:], isSeparatorRune)
	// A year alone is often the title's own ("1943 The Battle of Midway",
	// "2020 Super Baseball"), so it is a date prefix only when set off the
	// way a list number is. A month or day already marks a date.
	if pos == 4 && !isListSeparator(s[pos:len(s)-len(rest)]) {
		return Prefix{}
	}
	rest = strings.TrimSpace(rest)
	if rest == "" {
		return Prefix{}
	}
	return Prefix{Kind: KindDate, Rest: rest}
}

func allDigits(s string) bool {
	if s == "" {
		return false
	}
	for i := range s {
		if !isASCIIDigit(s[i]) {
			return false
		}
	}
	return true
}

func isASCIIDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isSpaceByte(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func isDatePartSeparator(ch byte) bool {
	return ch == '-' || ch == '.' || ch == '_'
}

func isSeparatorRune(r rune) bool {
	return r == '-' || r == '.' || r == '_' || r == ')' || unicode.IsSpace(r)
}

func daysInMonth(year, month int) int {
	switch month {
	case 4, 6, 9, 11:
		return 30
	case 2:
		if year%400 == 0 || (year%4 == 0 && year%100 != 0) {
			return 29
		}
		return 28
	default:
		return 31
	}
}
