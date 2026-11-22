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

func TestSlugify_Cyrillic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Russian game title",
			input:    "Тетрис",
			expected: "тетрис",
		},
		{
			name:     "Russian with spaces",
			input:    "Супер Марио",
			expected: "супермарио",
		},
		{
			name:     "Cyrillic with diacritics (ё→е)",
			input:    "Ёлки",
			expected: "елки",
		},
		{
			name:     "Cyrillic with metadata",
			input:    "Тетрис (Russia) [!]",
			expected: "тетрис",
		},
		{
			name:     "Mixed Cyrillic and Latin",
			input:    "Super Тетрис Bros",
			expected: "superтетрисbrothers",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := Slugify(MediaTypeGame, tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSlugify_Greek(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Greek game title",
			input:    "Παιχνίδι",
			expected: "παιχνιδι",
		},
		{
			name:     "Greek with diacritics removed",
			input:    "Σούπερ Μάριο",
			expected: "σουπερμαριο",
		},
		{
			name:     "Greek question mark normalized and removed",
			input:    "Τι είναι;",
			expected: "τιειναι", // Diacritics removed, punctuation removed in final slugification
		},
		{
			name:     "Mixed Greek and Latin",
			input:    "Super Ελληνικό Game",
			expected: "superελληνικοgame", // Diacritics removed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := Slugify(MediaTypeGame, tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSlugify_Indic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Hindi (Devanagari)",
			input:    "दिलवाले",
			expected: "दिलवाले",
		},
		{
			name:     "Bengali",
			input:    "বাংলা",
			expected: "বাংলা",
		},
		{
			name:     "Tamil",
			input:    "தமிழ்",
			expected: "தமிழ்",
		},
		{
			name:     "Telugu",
			input:    "తెలుగు",
			expected: "తెలుగు",
		},
		{
			name:     "Mixed Hindi and Latin",
			input:    "Dilwale दिलवाले",
			expected: "dilwaleदिलवाले",
		},
		{
			name:     "Hindi with metadata",
			input:    "गाना (India) [!]",
			expected: "गाना",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := Slugify(MediaTypeGame, tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSlugify_Arabic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Arabic music title",
			input:    "موسيقى",
			expected: "موسيقى",
		},
		{
			name:     "Arabic with spaces",
			input:    "سوبر ماريو",
			expected: "سوبرماريو",
		},
		{
			name:     "Arabic with vowel marks stripped",
			input:    "مُوسِيقَى",
			expected: "موسيقى",
		},
		{
			name:     "Arabic punctuation normalized and removed",
			input:    "ما هذا؟",
			expected: "ماهذا", // Punctuation normalized then removed in final slugification
		},
		{
			name:     "Mixed Arabic and Latin",
			input:    "Super موسيقى",
			expected: "superموسيقى",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := Slugify(MediaTypeGame, tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSlugify_Hebrew(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Hebrew word",
			input:    "שיר",
			expected: "שיר",
		},
		{
			name:     "Hebrew with vowel marks (niqqud) stripped",
			input:    "שִׁיר",
			expected: "שיר",
		},
		{
			name:     "Mixed Hebrew and Latin",
			input:    "Super שיר",
			expected: "superשיר",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := Slugify(MediaTypeGame, tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSlugify_Thai(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Thai music title",
			input:    "เพลงไทย",
			expected: "เพลงไทย",
		},
		{
			name:     "Thai with spaces",
			input:    "เพลง ไทย",
			expected: "เพลงไทย",
		},
		{
			name:     "Mixed Thai and Latin",
			input:    "Thai เพลง Music",
			expected: "thaiเพลงmusic",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := Slugify(MediaTypeGame, tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSlugify_Amharic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Amharic music title",
			input:    "ሙዚቃ",
			expected: "ሙዚቃ",
		},
		{
			name:     "Amharic with punctuation normalized and removed",
			input:    "ሙዚቃ።",
			expected: "ሙዚቃ", // Punctuation normalized then removed in final slugification
		},
		{
			name:     "Mixed Amharic and Latin",
			input:    "Ethiopian ሙዚቃ",
			expected: "ethiopianሙዚቃ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := Slugify(MediaTypeGame, tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSlugify_Turkish(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Turkish dotted I lowercased correctly",
			input:    "İstanbul",
			expected: "istanbul",
		},
		{
			name:     "Turkish dotless i transliterated",
			input:    "Kışlık",
			expected: "kislik", // Turkish special chars removed (Latin script = ASCII slug)
		},
		{
			name:     "Turkish special characters removed",
			input:    "Şişli Güneş",
			expected: "sisligunes", // ş→s, ğ→g, ü→u, ö→o (diacritics removed)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := Slugify(MediaTypeGame, tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestScriptDetection(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected ScriptType
	}{
		{"Latin", "Super Mario", ScriptLatin},
		{"CJK - Japanese", "ドラゴンクエスト", ScriptCJK},
		{"CJK - Chinese", "街头霸王", ScriptCJK}, //nolint:gosmopolitan // Chinese test data
		{"CJK - Korean", "스트리트파이터", ScriptCJK},
		{"Cyrillic", "Тетрис", ScriptCyrillic},
		{"Greek", "Παιχνίδι", ScriptGreek},
		{"Indic - Devanagari", "दिलवाले", ScriptIndic},
		{"Indic - Bengali", "বাংলা", ScriptIndic},
		{"Indic - Tamil", "தமிழ்", ScriptIndic},
		{"Arabic", "موسيقى", ScriptArabic},
		{"Hebrew", "שיר", ScriptHebrew},
		{"Thai", "เพลงไทย", ScriptThai},
		{"Burmese", "မြန်မာ", ScriptBurmese},
		{"Khmer", "ភាសាខ្មែរ", ScriptKhmer},
		{"Lao", "ພາສາລາວ", ScriptLao},
		{"Amharic", "ሙዚቃ", ScriptAmharic},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := DetectScript(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNeedsUnicodeSlug(t *testing.T) {
	t.Parallel()

	tests := []struct {
		script   ScriptType
		expected bool
	}{
		{ScriptLatin, false},
		{ScriptCJK, true},
		{ScriptCyrillic, true},
		{ScriptGreek, true},
		{ScriptIndic, true},
		{ScriptArabic, true},
		{ScriptHebrew, true},
		{ScriptThai, true},
		{ScriptBurmese, true},
		{ScriptKhmer, true},
		{ScriptLao, true},
		{ScriptAmharic, true},
	}

	for _, tt := range tests {
		result := needsUnicodeSlug(tt.script)
		assert.Equal(t, tt.expected, result)
	}
}

func TestNeedsNGramMatching(t *testing.T) {
	t.Parallel()

	tests := []struct {
		script   ScriptType
		expected bool
	}{
		{ScriptLatin, false},
		{ScriptCJK, false},
		{ScriptCyrillic, false},
		{ScriptThai, true},
		{ScriptBurmese, true},
		{ScriptKhmer, true},
		{ScriptLao, true},
	}

	for _, tt := range tests {
		result := needsNGramMatching(tt.script)
		assert.Equal(t, tt.expected, result)
	}
}
