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

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollapseSpaces(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "hello world", collapseSpaces("hello world"))
	assert.Equal(t, "hello world", collapseSpaces("hello  world"))
	assert.Equal(t, "hello world", collapseSpaces("  hello   world  "))
	assert.Equal(t, "a b", collapseSpaces("a\t\nb"))
	assert.Empty(t, collapseSpaces(""))
	assert.Empty(t, collapseSpaces("   "))
}

func TestStripLeadingNumberPrefix(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "Song Title", stripLeadingNumberPrefix("01 - Song Title"))
	assert.Equal(t, "Game Name", stripLeadingNumberPrefix("1. Game Name"))
	assert.Equal(t, "track", stripLeadingNumberPrefix("001-track"))
	assert.Equal(t, "rest", stripLeadingNumberPrefix("42  rest"))
	assert.Equal(t, "nodigits", stripLeadingNumberPrefix("nodigits"))
	assert.Equal(t, "123", stripLeadingNumberPrefix("123"))
	assert.Equal(t, "123abc", stripLeadingNumberPrefix("123abc"))
	assert.Empty(t, stripLeadingNumberPrefix(""))
}

func TestIsYearValue(t *testing.T) {
	t.Parallel()
	assert.True(t, isYearValue("1970"))
	assert.True(t, isYearValue("1999"))
	assert.True(t, isYearValue("2000"))
	assert.True(t, isYearValue("2099"))
	assert.False(t, isYearValue("1969"))
	assert.False(t, isYearValue("2100"))
	assert.False(t, isYearValue("abcd"))
	assert.False(t, isYearValue("19"))
	assert.False(t, isYearValue(""))
}

func TestRemoveSeasonEpisode(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "Show   Title", removeSeasonEpisode("Show S01E02 Title"))
	assert.Equal(t, "No episodes", removeSeasonEpisode("No episodes"))
	assert.Equal(t, "Show Title", removeSeasonEpisode("ShowS01E02Title"))
}

func TestFindUnbracketedYear(t *testing.T) {
	t.Parallel()
	input := "Movie 2018 720p"
	m := findUnbracketedYear(input)
	assert.True(t, m.ok)
	assert.Equal(t, "2018", input[m.cap1s:m.cap1e])

	m = findUnbracketedYear("1995 title")
	assert.True(t, m.ok)

	m = findUnbracketedYear("title 2020")
	assert.True(t, m.ok)

	assert.False(t, findUnbracketedYear("abc2018def").ok)
	assert.False(t, findUnbracketedYear("title_2018_cut").ok)
	assert.False(t, findUnbracketedYear("title 1969 stuff").ok)
	assert.False(t, findUnbracketedYear("title 1234 stuff").ok)
	assert.False(t, findUnbracketedYear("").ok)
}

func TestWordBoundary(t *testing.T) {
	t.Parallel()
	s := "hello world"
	assert.True(t, isWordBoundaryBefore(s, 0))
	assert.False(t, isWordBoundaryBefore(s, 1))
	assert.True(t, isWordBoundaryBefore(s, 6))
	assert.True(t, isWordBoundaryAfter(s, 5))
	assert.False(t, isWordBoundaryAfter(s, 4))
	assert.True(t, isWordBoundaryAfter(s, len(s)))

	s = "hello_world"
	assert.False(t, isWordBoundaryBefore(s, 6))
	assert.False(t, isWordBoundaryAfter(s, 5))
}
