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
	"bytes"
	"regexp"
	"strings"
	"sync"

	slugs "github.com/ZaparooProject/go-slugs"
)

// Package-level compiled regexes for filename parsing.
// These are compiled once at initialization for optimal performance.
var (
	reHTGDBHackMarker   = regexp.MustCompile(`(?i)\((?:hack|patch)\)`)
	reHTGDBPatchVersion = regexp.MustCompile(`(?i)(?:^|[\s_+\-])v(\d+(?:\.\d+)*)`)
	reHTGDBPatchLabel   = regexp.MustCompile(`(?i)\b(?:` +
		`(?:ted[\s_-]+)?woolsey[\s_-]+uncensored[\s_-]+edition|` +
		`stereo\s*\+\s*bug[\s_-]*fix(?:es)?|gba[\s_-]+script[\s_-]+port|` +
		`splash[\s_-]+screen[\s_-]+removed|quality[\s_-]+of[\s_-]+life|` +
		`music[\s_-]+restore|sram[\s_-]+expansion|faster[\s_-]+rom|` +
		`fast[\s_-]*rom|slow[\s_-]*rom|all[\s_-]+bugs[\s_-]+fix|` +
		`bug[\s_-]*fix(?:es)?|uncensor(?:ed)?|relocalized|restoration|` +
		`performance|retouch|overhaul|tweak|wide[\s_-]*screen|` +
		`msu[\s_-]*1|sa[\s_-]*1|colori[sz](?:ation|ed)|redux|qol` +
		`)\b`)

	// Scene release artifact patterns (for modern media: movies, TV shows)
	// Note: These patterns match AFTER separator normalization, so hyphens have become spaces
	reResolution = regexp.MustCompile(`(?i)\b(480p|720p|1080p|2160p|4K|8K|UHD)\b`)
	reSource     = regexp.MustCompile(
		`(?i)\b(CAM|TS|TC|TELESYNC|TELECINE|SCREENER|SCR|HDTV|PDTV|WEB[\s-]?DL|WEBRip|WEBRIP|` +
			`BluRay|BRRip|BDRip|Blu[\s-]?ray|Remux|DVDRip|DVDSCR|R5)\b`,
	)
	reVideoCodec = regexp.MustCompile(`(?i)\b(x264|x265|h[\s.]?264|h[\s.]?265|HEVC|AVC|XviD|VP9|AV1)\b`)
	reAudioCodec = regexp.MustCompile(
		`(?i)\b(AAC|MP3|AC3|DTS|DD5[\s.]?1|DD\+?5[\s.]?1|TrueHD|Atmos|EAC3|FLAC|DTS[\s-]?HD)\b`,
	)
	reHDR         = regexp.MustCompile(`(?i)\b(HDR|HDR10|HDR10\+|DV|Dolby[\s.]?Vision)\b`)
	reSceneStatus = regexp.MustCompile(
		`(?i)\b(PROPER|REPACK|INTERNAL|LIMITED|READNFO|DIRFIX|NFOFIX|UNRATED|EXTENDED)\b`,
	)
	reReleaseGroup = regexp.MustCompile(`-([A-Z][A-Z0-9]{2,})$`)
)

// patchLabelAliases contains patch labels observed in ROM-hack collections.
// Unknown bracket labels remain available to the general filename parser.
var patchLabelAliases = map[string]tagValue{
	"all-bugs-fix":                   "bugfix",
	"bug-fix":                        "bugfix",
	"bug-fixes":                      "bugfix",
	"bugfix":                         "bugfix",
	"bugfixes":                       "bugfix",
	"color":                          "color",
	"colorisation":                   "color",
	"colorised":                      "color",
	"colorization":                   "color",
	"colorized":                      "color",
	"fast-rom":                       "fastrom",
	"faster-rom":                     "fastrom",
	"fastrom":                        "fastrom",
	"fix":                            "fix",
	"fixes":                          "fix",
	"gba-script-port":                "script-port",
	"hack":                           "hack",
	"msu-1":                          "msu1",
	"msu1":                           "msu1",
	"music":                          "music",
	"music-restore":                  "music",
	"no-sram":                        "no-sram",
	"overhaul":                       "overhaul",
	"performance":                    "performance",
	"quality-of-life":                "qol",
	"qol":                            "qol",
	"redux":                          "redux",
	"relocalized":                    "relocalized",
	"restoration":                    "restoration",
	"retouch":                        "retouch",
	"sa-1":                           "sa1",
	"sa1":                            "sa1",
	"script-port":                    "script-port",
	"slow-rom":                       "slowrom",
	"slowrom":                        "slowrom",
	"splash-screen-removed":          "splash-screen-removed",
	"sram":                           "sram",
	"sram-expansion":                 "sram",
	"stereo-bug-fixes":               "bugfix",
	"tweak":                          "tweak",
	"uncensor":                       "uncensored",
	"uncensored":                     "uncensored",
	"wide-screen":                    "widescreen",
	"widescreen":                     "widescreen",
	"woolsey-uncensored-edition":     "uncensored",
	"ted-woolsey-uncensored-edition": "uncensored",
}

var structuralMediaPrefixes = map[string]tagValue{
	"disc": tagMediaDisc,
	"disk": tagMediaDisk,
	"file": tagMediaFile,
	"part": tagMediaPart,
	"tape": tagMediaTape,
	"cd":   tagMediaDisc,
}

var structuralRomanNumbers = map[string]string{
	"i":    "1",
	"ii":   "2",
	"iii":  "3",
	"iv":   "4",
	"v":    "5",
	"vi":   "6",
	"vii":  "7",
	"viii": "8",
	"ix":   "9",
	"x":    "10",
}

// normTagCache and its mutex protect cached normalizeTag results for the duration
// of an indexing run. The bracket vocabulary is small (~200–400 unique strings) so
// the cache stabilises quickly. nil map means caching is disabled.
var (
	normTagCacheMu sync.RWMutex
	normTagCache   map[string]string
)

func cachedNormalizeTag(s string) string {
	normTagCacheMu.RLock()
	m := normTagCache
	if m != nil {
		if v, ok := m[s]; ok {
			normTagCacheMu.RUnlock()
			return v
		}
	}
	normTagCacheMu.RUnlock()

	v := normalizeTag(s)

	if m != nil {
		normTagCacheMu.Lock()
		if normTagCache != nil {
			normTagCache[s] = v
		}
		normTagCacheMu.Unlock()
	}
	return v
}

func normalizeStructuralNumber(s string) (string, bool) {
	if s == "" {
		return "", false
	}
	if roman, ok := structuralRomanNumbers[s]; ok {
		return roman, true
	}
	for i := range len(s) {
		if !isDigit(s[i]) {
			return "", false
		}
	}
	normalized := strings.TrimLeft(s, "0")
	if normalized == "" {
		normalized = "0"
	}
	return normalized, true
}

func structuralSideValue(s string) (tagValue, bool) {
	switch s {
	case "a", "1":
		return tagMediaSideA, true
	case "b", "2":
		return tagMediaSideB, true
	case "c", "3":
		return tagMediaSideC, true
	case "d", "4":
		return tagMediaSideD, true
	default:
		return "", false
	}
}

func compactStructuralMediaTag(normalized string, source tagSource) ([]canonicalTag, bool) {
	for prefix, mediaValue := range structuralMediaPrefixes {
		if !strings.HasPrefix(normalized, prefix) || len(normalized) == len(prefix) {
			continue
		}
		number, ok := normalizeStructuralNumber(normalized[len(prefix):])
		if !ok {
			continue
		}
		return []canonicalTag{
			{Type: tagTypeMedia, Value: mediaValue, Source: source},
			{Type: tagTypeDisc, Value: tagValue(number), Source: source},
		}, true
	}
	return nil, false
}

func parseStructuralSetTag(normalized string, source tagSource) ([]canonicalTag, bool) {
	if normalized == "" {
		return nil, false
	}

	if strings.HasPrefix(normalized, "side-") {
		if sideValue, ok := structuralSideValue(strings.TrimPrefix(normalized, "side-")); ok {
			return []canonicalTag{{Type: tagTypeMedia, Value: sideValue, Source: source}}, true
		}
	}

	tokens := strings.Split(normalized, "-")
	if len(tokens) < 2 {
		return compactStructuralMediaTag(normalized, source)
	}

	mediaValue, ok := structuralMediaPrefixes[tokens[0]]
	if !ok {
		return nil, false
	}
	number, ok := normalizeStructuralNumber(tokens[1])
	if !ok {
		return nil, false
	}

	results := []canonicalTag{
		{Type: tagTypeMedia, Value: mediaValue, Source: source},
		{Type: tagTypeDisc, Value: tagValue(number), Source: source},
	}

	next := 2
	if len(tokens) >= next+2 && tokens[next] == "of" {
		total, totalOK := normalizeStructuralNumber(tokens[next+1])
		if !totalOK {
			return nil, false
		}
		results = append(results, canonicalTag{Type: tagTypeDiscTotal, Value: tagValue(total), Source: source})
		next += 2
	}
	if len(tokens) >= next+2 && tokens[next] == "side" {
		sideValue, sideOK := structuralSideValue(tokens[next+1])
		if !sideOK {
			return nil, false
		}
		results = append(results, canonicalTag{Type: tagTypeMedia, Value: sideValue, Source: source})
		next += 2
	}
	if next != len(tokens) {
		return nil, false
	}

	return results, true
}

func canonicalPatchLabel(label string) (tagValue, bool) {
	value, ok := patchLabelAliases[label]
	return value, ok
}

func versionedPatchValue(value tagValue, version string) tagValue {
	if version == "" {
		return value
	}
	return tagValue(string(value) + ":" + version)
}

func hasHTGDBPatchMarker(filename string) bool {
	for offset := 0; offset < len(filename); {
		markerStart := strings.IndexByte(filename[offset:], '(')
		if markerStart < 0 {
			return false
		}
		markerStart += offset
		remaining := filename[markerStart:]
		if len(remaining) >= len("(hack)") && strings.EqualFold(remaining[:len("(hack)")], "(hack)") {
			return true
		}
		if len(remaining) >= len("(patch)") && strings.EqualFold(remaining[:len("(patch)")], "(patch)") {
			return true
		}
		offset = markerStart + 1
	}
	return false
}

func splitPatchFilenameExtension(filename string) (stem, extension string) {
	lastSlash := strings.LastIndexAny(filename, `/\\`)
	lastDot := strings.LastIndex(filename, ".")
	if lastDot <= lastSlash {
		return filename, ""
	}

	extensionBody := filename[lastDot+1:]
	if len(extensionBody) < 2 || len(extensionBody) > 4 || strings.ContainsAny(extensionBody, " \t") {
		return filename, ""
	}
	return filename[:lastDot], filename[lastDot:]
}

// parseHTGDBPatchSuffix handles HTGDB/SmokeMonster names such as
// "Axelay (U) FastROM (Hack) v1.0 Vitor Vilela". Only researched patch
// labels are recognized; arbitrary names containing "(Hack)" retain their
// existing title and revision behavior.
func parseHTGDBPatchSuffix(filename string) ([]canonicalTag, string, bool) {
	if !hasHTGDBPatchMarker(filename) {
		return nil, filename, false
	}

	stem, extension := splitPatchFilenameExtension(filename)
	markers := reHTGDBHackMarker.FindAllStringIndex(stem, -1)
	if len(markers) == 0 {
		return nil, filename, false
	}

	marker := markers[len(markers)-1]
	prefix := stem[:marker[0]]
	labelMatches := reHTGDBPatchLabel.FindAllStringIndex(prefix, -1)
	if len(labelMatches) == 0 {
		return nil, filename, false
	}

	stripStart := labelMatches[0][0]
	version := ""
	if match := reHTGDBPatchVersion.FindStringSubmatch(stem[stripStart:]); match != nil {
		version = strings.ReplaceAll(match[1], ".", "-")
	}

	tags := []canonicalTag{{
		Type:   tagTypeUnlicensed,
		Value:  tagUnlicensedHack,
		Source: tagSourceBracketed,
	}}
	seen := make(map[tagValue]struct{}, len(labelMatches))
	for _, match := range labelMatches {
		normalized := cachedNormalizeTag(prefix[match[0]:match[1]])
		patchValue, ok := canonicalPatchLabel(normalized)
		if !ok {
			continue
		}
		patchValue = versionedPatchValue(patchValue, version)
		if _, exists := seen[patchValue]; exists {
			continue
		}
		seen[patchValue] = struct{}{}
		tags = append(tags, canonicalTag{
			Type:   tagTypePatch,
			Value:  patchValue,
			Source: tagSourceBracketed,
		})
	}
	if len(tags) == 1 {
		return nil, filename, false
	}

	remaining := strings.TrimRight(prefix[:stripStart], " \t-_+") + extension
	return tags, remaining, true
}

// stripSceneArtifacts removes common scene release artifacts from filenames.
// These patterns are common in movie and TV show filenames but should be removed
// from display titles for better readability.
//
// Removed patterns include:
//   - Resolution markers: 720p, 1080p, 2160p, 4K, UHD
//   - Source types: BluRay, WEB-DL, WEBRip, HDTV, DVDRip
//   - Video codecs: x264, x265, h264, HEVC, AVC
//   - Audio codecs: AAC, AC3, DTS, DD5.1, TrueHD, Atmos
//   - HDR formats: HDR, HDR10, DV (Dolby Vision)
//   - Scene status: PROPER, REPACK, INTERNAL, LIMITED
//   - Release groups: -GROUP suffix (uppercase alphanumeric), when releaseGroups
//
// Examples:
//   - "The.Dark.Knight.2008.1080p.BluRay.x264-GROUP" → "The.Dark.Knight.2008"
//   - "Movie.Title.2010.WEB-DL.DD5.1.H264" → "Movie.Title.2010"
//   - "Show.S01E01.PROPER.720p.HDTV" → "Show.S01E01"
//
// sceneMarkerBytes contains byte sequences that indicate scene release
// artifacts may be present. If none are found, stripSceneArtifacts can skip
// 7 regex evaluations. These are chosen to catch all patterns matched by the
// scene regexes while being absent from typical ROM filenames.
var sceneMarkerBytes = [][]byte{
	[]byte("264"), []byte("265"), // x264, h264, x265, h265
	[]byte("480"), []byte("720"), []byte("108"), []byte("216"), // resolutions
	[]byte("4K"), []byte("4k"), []byte("8K"), []byte("8k"), []byte("UHD"), []byte("uhd"),
	[]byte("BluRay"), []byte("bluray"), []byte("WEBDL"), []byte("webdl"),
	[]byte("WEB-DL"), []byte("web-dl"), []byte("WEBRip"), []byte("webrip"),
	[]byte("BDRip"), []byte("bdrip"), []byte("DVDRip"), []byte("dvdrip"),
	[]byte("HDTV"), []byte("hdtv"),
	[]byte("HEVC"), []byte("hevc"), []byte("AVC"), []byte("avc"),
	[]byte("XviD"), []byte("xvid"), []byte("AAC"), []byte("aac"),
	[]byte("DTS"), []byte("dts"), []byte("FLAC"), []byte("flac"),
	[]byte("HDR"), []byte("hdr"),
	[]byte("PROPER"), []byte("proper"), []byte("REPACK"), []byte("repack"),
}

// stripsReleaseGroups reports whether a media type's filenames follow scene
// release naming, where a trailing "-GROUP" names the release group. Game
// filenames do not: an all-caps title such as "F-ZERO" or "SUPER MARIO
// ALL-STARS", typical of a cartridge header, ends in the same shape.
func stripsReleaseGroups(mediaType slugs.MediaType) bool {
	switch mediaType {
	case slugs.MediaTypeMovie, slugs.MediaTypeTVShow, slugs.MediaTypeVideo:
		return true
	default:
		return false
	}
}

// hasSceneMarkers does a cheap byte scan for substrings that indicate scene
// release artifacts may be present. Returns false for typical ROM filenames,
// allowing stripSceneArtifacts to skip 7 regex evaluations.
func hasSceneMarkers(s string, releaseGroups bool) bool {
	sb := []byte(s)
	for _, marker := range sceneMarkerBytes {
		if bytes.Contains(sb, marker) {
			return true
		}
	}
	// Check for release group suffix: -GROUP at end of string
	if n := len(s); releaseGroups && n > 3 && s[n-1] >= 'A' && s[n-1] <= 'Z' {
		for i := n - 2; i >= 0; i-- {
			if s[i] == '-' {
				return true
			}
			if (s[i] < 'A' || s[i] > 'Z') && (s[i] < '0' || s[i] > '9') {
				break
			}
		}
	}
	return false
}

func stripSceneArtifacts(input string, releaseGroups bool) string {
	// Fast path: skip 7 regex evaluations if no scene markers present
	if !hasSceneMarkers(input, releaseGroups) {
		return input
	}

	result := input

	// Remove resolution markers
	result = reResolution.ReplaceAllString(result, " ")

	// Remove source types
	result = reSource.ReplaceAllString(result, " ")

	// Remove video codecs
	result = reVideoCodec.ReplaceAllString(result, " ")

	// Remove audio codecs
	result = reAudioCodec.ReplaceAllString(result, " ")

	// Remove HDR formats
	result = reHDR.ReplaceAllString(result, " ")

	// Remove scene status tags
	result = reSceneStatus.ReplaceAllString(result, " ")

	// Remove release group suffix (at end of string)
	if releaseGroups {
		result = reReleaseGroup.ReplaceAllString(result, "")
	}

	return result
}

func stripMetadataBracketsForDisplay(s string) string {
	return filterBracketGroups(s, true)
}

// filterBracketGroups keeps the bracket groups that are structural set markers
// and drops the rest when keepStructural is set, and does the reverse when it
// is not. Text after an unclosed bracket is metadata to the first form and part
// of the name to the second.
func filterBracketGroups(s string, keepStructural bool) string {
	var result strings.Builder
	result.Grow(len(s))

	for i := 0; i < len(s); i++ {
		open := s[i]
		var closeBracket byte
		switch open {
		case '(':
			closeBracket = ')'
		case '[':
			closeBracket = ']'
		case '{':
			closeBracket = '}'
		case '<':
			closeBracket = '>'
		default:
			_ = result.WriteByte(open)
			continue
		}

		start := i
		depth := 1
		for i++; i < len(s); i++ {
			switch s[i] {
			case open:
				depth++
			case closeBracket:
				depth--
				if depth == 0 {
					content := s[start+1 : i]
					_, structural := parseStructuralSetTag(cachedNormalizeTag(content), tagSourceBracketed)
					if structural == keepStructural {
						_, _ = result.WriteString(s[start : i+1])
					}
				}
			}
			if depth == 0 {
				break
			}
		}
		if depth != 0 && !keepStructural {
			_, _ = result.WriteString(s[start:])
		}
	}

	return strings.TrimSpace(result.String())
}

// Parse extracts the canonical title of a media file from its filename. Zaparoo
// Core indexes this title and slugifies it, so other clients that index files
// must use it to get titles and slugs that match Core's.
//
// It strips the file extension, scene release artifacts (resolution, codecs,
// sources, and release groups for video), season and episode markers, and
// every bracket group: (), [], {} and <>. Filenames with no spaces but two or
// more dots, underscores or dashes have those separators turned into spaces.
//
// stripLeadingNumbers removes list numbering such as "1. " or "01 - ". Set it
// only when the directory as a whole is numbered; see the browseprefix package.
//
// An empty mediaType is treated as a game.
//
// Examples:
//   - "Super Mario Bros (USA) [!].sfc" → "Super Mario Bros"
//   - "super_mario_bros.sfc" → "super mario bros"
//   - "1942 (USA)" → "1942"
//   - "The.Dark.Knight.2008.1080p.BluRay.x264.mkv" → "The Dark Knight 2008"
func Parse(filename string, stripLeadingNumbers bool, mediaType slugs.MediaType) string {
	return parseFilenameTitle(filename, stripLeadingNumbers, false, mediaType)
}

// ParseDisplay is like Parse but keeps structural set markers such as
// "(Disc 1 of 2)" or "(Side A)", so files of a multi-part set stay
// distinguishable when shown one per row.
func ParseDisplay(filename string, stripLeadingNumbers bool, mediaType slugs.MediaType) string {
	return parseFilenameTitle(filename, stripLeadingNumbers, true, mediaType)
}

// StripStructuralSetMarkers removes structural set markers such as "(Disc 1)",
// "(Disk 2 of 4)" and "(Side A)" from a display title and leaves every other
// bracket group in place. Zaparoo Core indexes the same markers as media and
// disc tags, so a caller that shows tags alongside the name can drop them from
// the name. A title that is nothing but a marker is returned unchanged.
func StripStructuralSetMarkers(title string) string {
	if !strings.ContainsAny(title, "([{<") {
		return title
	}
	stripped := strings.TrimSpace(collapseSpaces(filterBracketGroups(title, false)))
	if stripped == "" {
		return title
	}
	return stripped
}

// SetNormalizeCache sets a map used to cache tag normalization during bulk
// parsing, such as indexing a media library. The bracket vocabulary is small,
// so the cache stays small. Pass nil to stop caching and release the map.
// The package guards the map; callers must not use it after passing it in.
func SetNormalizeCache(m map[string]string) {
	normTagCacheMu.Lock()
	normTagCache = m
	normTagCacheMu.Unlock()
}

func parseFilenameTitle(
	filename string, stripLeadingNumbers, preserveStructuralTags bool, mediaType slugs.MediaType,
) string {
	title := filename
	if mediaType == "" || mediaType == slugs.MediaTypeGame {
		if _, withoutPatchSuffix, matched := parseHTGDBPatchSuffix(title); matched {
			title = withoutPatchSuffix
		}
	}

	// Step 1: Remove file extension first (simplifies later processing)
	// Find the last dot that's likely an extension (after the last slash if any)
	if lastDot := strings.LastIndex(title, "."); lastDot > 0 {
		// Only remove if the extension looks reasonable (2-4 characters after the dot)
		ext := title[lastDot+1:]
		if len(ext) >= 2 && len(ext) <= 4 && !strings.Contains(ext, " ") {
			title = title[:lastDot]
		}
	}

	// Step 2: Strip release group BEFORE separator normalization
	// Release groups are typically "-GROUP" at the end, and the hyphen will be converted to a space
	// So we need to remove it early before that conversion happens. Only
	// scene-named media has them; see stripsReleaseGroups.
	// reReleaseGroup is anchored at end-of-string and needs a '-', so a title
	// with no hyphen at all cannot match. ReplaceAllString allocates a new
	// string even on no match, and this runs twice per file, so skip the engine
	// for the common hyphen-free case.
	releaseGroups := stripsReleaseGroups(mediaType)
	if releaseGroups && strings.IndexByte(title, '-') >= 0 {
		title = reReleaseGroup.ReplaceAllString(title, "")
	}

	// Step 3: Normalize filename separators (before scene artifact stripping)
	// Heuristic: If filename has no spaces AND has 2+ separators (dots, underscores, or dashes),
	// treat them all as space substitutes. Common in scene releases and ROM filenames.
	// Examples: "The.Dark.Knight.2008.mkv", "super_mario_bros.sfc", "mega-man-x.nes"
	dotCount := strings.Count(title, ".")
	sepCount := strings.Count(title, "_") + strings.Count(title, "-")

	if !strings.Contains(title, " ") && (dotCount >= 2 || sepCount >= 2) {
		title = strings.ReplaceAll(title, ".", " ")
		title = strings.ReplaceAll(title, "_", " ")
		title = strings.ReplaceAll(title, "-", " ")
	}

	// Step 4: Strip scene release artifacts CONTEXTUALLY (AFTER separator normalization)
	// CRITICAL: Only strip artifacts from the portion AFTER the year to avoid corrupting titles
	// Example: "Cam (2018)" should NOT have "Cam" removed as a scene tag
	yearMatch := findUnbracketedYear(title)

	var extractedYearForStripping string
	if yearMatch.ok {
		yearIdx := yearMatch.start
		extractedYearForStripping = title[yearMatch.cap1s:yearMatch.cap1e]

		beforeYear := title[:yearIdx]
		afterYear := title[yearIdx+len(extractedYearForStripping):]

		// Only strip artifacts from AFTER the year to protect the title
		cleanedAfterYear := stripSceneArtifacts(afterYear, releaseGroups)
		title = beforeYear + extractedYearForStripping + cleanedAfterYear
	} else {
		// No year found - apply stripping to entire title (safer for non-media files)
		title = stripSceneArtifacts(title, releaseGroups)
	}

	// Step 5: Strip S##E## patterns (TV show episode markers) from titles
	// These are useful as metadata but clutter the display title
	title = removeSeasonEpisode(title)

	// Step 6: Optionally strip leading number prefixes (e.g., "1. ", "01 - ")
	// Only done when contextual detection confirms list-based numbering
	if stripLeadingNumbers {
		title = stripLeadingNumberPrefix(title)
		title = strings.TrimSpace(title)
	}

	// Step 7: Remove bracket metadata. Per-file display titles retain
	// structural set markers, while canonical titles strip every bracket group.
	if preserveStructuralTags {
		title = stripMetadataBracketsForDisplay(title)
	} else {
		title = slugs.StripMetadataBrackets(title)
	}
	title = strings.TrimSpace(title)

	// Step 8: Normalize multiple spaces to single space
	// This handles cases where bracket removal creates gaps like "Game [USA] [!]" → "Game  "
	title = collapseSpaces(title)

	return strings.TrimSpace(title)
}
