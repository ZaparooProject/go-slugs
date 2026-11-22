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
	"strings"
)

// ParseGame normalizes game titles by applying game-specific transformations.
// This handles common game title patterns and variations to ensure consistent matching.
//
// Transformations applied (in order):
//  1. Split titles and strip articles: "The Zelda: Link's Awakening" → "Zelda Link's Awakening"
//  2. Strip trailing articles: "Legend, The" → "Legend"
//  3. Strip metadata brackets: (USA), [!], {Europe}, <Beta> → removed
//  4. Strip edition/version suffixes: "Edition", "Version", v1.0 → removed
//  5. Normalize separators: Convert periods to spaces (for abbreviation matching)
//  6. Expand abbreviations: "Bros" → "brothers", "vs" → "versus", "Dr" → "doctor"
//  7. Expand number words: "one" → "1", "two" → "2"
//  8. Normalize ordinals: "1st" → "1", "2nd" → "2"
//  9. Convert roman numerals: "VII" → "7", "II" → "2" (preserves "X" for games like Mega Man X)
//
// Examples:
//   - "Super Mario Bros. III (USA) [!]" → "super mario brothers 3"
//   - "Street Fighter II Version" → "street fighter 2"
//   - "Mega Man X" → "mega man x" (X preserved)
//   - "Final Fantasy VII" → "final fantasy 7"
func ParseGame(title string) string {
	s := title

	// Normalize width first so fullwidth separators are detected by SplitAndStripArticles
	// This converts "：" (fullwidth colon) to ":" (ASCII colon)
	// Skip for ASCII-only strings (no fullwidth characters possible)
	if !isASCII(s) {
		s = NormalizeWidth(s)
	}
	s = strings.TrimSpace(s)

	// MUST happen first: Split titles and strip articles
	// This preserves separators like ":" and "," which are needed for splitting
	s = SplitAndStripArticles(s)
	s = strings.TrimSpace(s)

	// Strip trailing articles (also needs "," intact)
	s = StripTrailingArticle(s)
	s = strings.TrimSpace(s)

	// Strip metadata brackets (region codes, dump info, tags)
	s = StripMetadataBrackets(s)
	s = strings.TrimSpace(s)

	// Strip edition and version suffixes
	s = StripEditionAndVersionSuffixes(s)
	s = strings.TrimSpace(s)

	// Collapse dotted initialisms (T.V. → TV) before TrimRight strips the trailing
	// period, leaving only one letter-period pair and defeating the {2,} threshold.
	// Must happen here so the Roman numeral pass below doesn't see isolated "V"→"5".
	s = CollapseDottedInitialisms(s)

	// Trim trailing separators before abbreviation expansion
	// This ensures "Bros-" → "Bros" so abbreviation matching works
	s = strings.TrimRight(s, "-:;.,_/\\ ")

	// Normalize symbols and separators (BUT NOT commas - needed for trailing articles)
	// This is similar to NormalizeSymbolsAndSeparators but preserves commas
	// Conjunction normalization
	s = strings.ReplaceAll(s, " & ", " and ")
	s = strings.ReplaceAll(s, "&", " and ")
	s = strings.ReplaceAll(s, " + ", " and ")
	s = strings.ReplaceAll(s, "+", " plus ")
	// Separator normalization (: _ / \ ; but NOT commas)
	s = strings.ReplaceAll(s, ":", " ")
	s = strings.ReplaceAll(s, "_", " ")
	s = strings.ReplaceAll(s, "/", " ")
	s = strings.ReplaceAll(s, "\\", " ")
	s = strings.ReplaceAll(s, ";", " ")
	// Convert periods to spaces
	s = strings.ReplaceAll(s, ".", " ")
	s = strings.TrimSpace(s)

	// Expand abbreviations and number words in a single pass to avoid
	// splitting/joining the string twice
	s = expandAbbreviationsAndNumbers(s)

	// Normalize ordinals
	s = NormalizeOrdinals(s)

	// Convert roman numerals (includes lowercasing)
	s = ConvertRomanNumerals(s)

	return s
}
