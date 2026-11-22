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

package slugs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateBigrams(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Thai word",
			input:    "เพลง",
			expected: []string{"เพ", "พล", "ลง"},
		},
		{
			name:     "Empty string",
			input:    "",
			expected: []string{""},
		},
		{
			name:     "Single character",
			input:    "a",
			expected: []string{"a"},
		},
		{
			name:     "Two characters",
			input:    "ab",
			expected: []string{"ab"},
		},
		{
			name:     "Latin word",
			input:    "test",
			expected: []string{"te", "es", "st"},
		},
		{
			name:     "CJK characters",
			input:    "日本語",                //nolint:gosmopolitan // Japanese test data
			expected: []string{"日本", "本語"}, //nolint:gosmopolitan // Japanese test data
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := GenerateBigrams(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGenerateTrigrams(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Thai word",
			input:    "เพลงไทย",
			expected: []string{"เพล", "พลง", "ลงไ", "งไท", "ไทย"},
		},
		{
			name:     "Single character (falls back to input)",
			input:    "a",
			expected: []string{"a"},
		},
		{
			name:     "Two characters (falls back to bigram)",
			input:    "ab",
			expected: []string{"ab"},
		},
		{
			name:     "Three characters",
			input:    "abc",
			expected: []string{"abc"},
		},
		{
			name:     "Latin word",
			input:    "test",
			expected: []string{"tes", "est"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := GenerateTrigrams(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestJaccardSimilarity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		set1     []string
		set2     []string
		expected float64
	}{
		{
			name:     "Identical sets",
			set1:     []string{"a", "b", "c"},
			set2:     []string{"a", "b", "c"},
			expected: 1.0,
		},
		{
			name:     "No overlap",
			set1:     []string{"a", "b"},
			set2:     []string{"c", "d"},
			expected: 0.0,
		},
		{
			name:     "Partial overlap",
			set1:     []string{"a", "b", "c"},
			set2:     []string{"b", "c", "d"},
			expected: 0.5, // 2 common / 4 total
		},
		{
			name:     "Subset",
			set1:     []string{"a", "b"},
			set2:     []string{"a", "b", "c"},
			expected: 0.666666, // 2 common / 3 total (approximately)
		},
		{
			name:     "Both empty",
			set1:     []string{},
			set2:     []string{},
			expected: 1.0,
		},
		{
			name:     "One empty",
			set1:     []string{"a"},
			set2:     []string{},
			expected: 0.0,
		},
		{
			name:     "Single element match",
			set1:     []string{"a"},
			set2:     []string{"a"},
			expected: 1.0,
		},
		{
			name:     "Thai bigrams example",
			set1:     []string{"เพ", "พล", "ลง"},
			set2:     []string{"เพ", "พล", "ลง", "งไ", "ไท", "ทย"},
			expected: 0.5, // 3 common / 6 total
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := JaccardSimilarity(tt.set1, tt.set2)
			assert.InDelta(t, tt.expected, result, 0.001)
		})
	}
}

func TestNGramMatching_ThaiExample(t *testing.T) {
	t.Parallel()

	// Full title
	fullTitle := "เพลงไทย"
	fullBigrams := GenerateBigrams(fullTitle)

	// Query: just the first word
	query := "เพลง"
	queryBigrams := GenerateBigrams(query)

	// Calculate similarity
	similarity := JaccardSimilarity(queryBigrams, fullBigrams)

	// Should have good similarity (50% - 3 bigrams match out of 6 total)
	assert.InDelta(t, 0.5, similarity, 0.01)
	assert.GreaterOrEqual(t, similarity, 0.5, "Thai substring should match with n-gram similarity >= 0.5")
}

func TestNGramMatching_ThaiSecondWord(t *testing.T) {
	t.Parallel()

	// Full title
	fullTitle := "เพลงไทย"
	fullBigrams := GenerateBigrams(fullTitle)

	// Query: just the second word
	query := "ไทย"
	queryBigrams := GenerateBigrams(query)

	// Calculate similarity
	similarity := JaccardSimilarity(queryBigrams, fullBigrams)

	// Should still have some similarity
	assert.Greater(t, similarity, 0.0, "Thai second word should have some similarity")
	assert.GreaterOrEqual(t, similarity, 0.25, "Expected at least 25% similarity for second word")
}

func TestIsThai(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Thai text", "เพลงไทย", true},
		{"Latin text", "Thai music", false},
		{"Mixed Thai-Latin", "Thai เพลง", true},
		{"CJK text", "日本語", false}, //nolint:gosmopolitan // Japanese test data
		{"Empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := IsThai(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNeedsNGramMatching_Function(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Thai", "เพลงไทย", true},
		{"Burmese", "မြန်မာ", true},
		{"Khmer", "ភាសាខ្មែរ", true},
		{"Lao", "ພາສາລາວ", true},
		{"Latin", "English", false},
		{"CJK", "日本語", false}, //nolint:gosmopolitan // Japanese test data
		{"Arabic", "موسيقى", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := needsNGramMatching(DetectScript(tt.input))
			assert.Equal(t, tt.expected, result)
		})
	}
}
