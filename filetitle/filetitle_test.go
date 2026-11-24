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

	slugs "github.com/ZaparooProject/go-slugs"
	"github.com/stretchr/testify/assert"
)

func TestParseFilenameTitle_HTGDBPatchSuffixes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		filename string
		want     string
	}{
		{
			name:     "FastROM suffix",
			filename: "Super Castlevania IV (U) FastROM (Hack) v1.1 Vitor Vilela.sfc",
			want:     "Super Castlevania IV",
		},
		{
			name:     "hyphenated FastROM suffix",
			filename: "Arcana - Fastrom (Hack) v0.2 rainponcho.sfc",
			want:     "Arcana",
		},
		{
			name:     "uncensored suffix",
			filename: "Super Castlevania IV Uncensored (hack) v1 ShadowOne333.sfc",
			want:     "Super Castlevania IV",
		},
		{
			name:     "versionless restoration suffix",
			filename: "Contra III - The Alien Wars Restoration (Hack) Final SCD.sfc",
			want:     "Contra III - The Alien Wars",
		},
		{
			name:     "SA-1 suffix qualifiers",
			filename: "Gradius III SA-1 Origin USA (Hack) v1.6 Vitor Vilela.sfc",
			want:     "Gradius III",
		},
		{
			name: "version before marker",
			filename: "Final Fantasy III Ted Woolsey Uncensored Edition v1.3 by Rodimus Primal " +
				"(Hack).sfc",
			want: "Final Fantasy III",
		},
		{
			name: "composite suffix",
			filename: "Final Fantasy V GBA Script Port+Ginger Battle Galuf+Bugfixes, Sprite Touch-Ups " +
				"(Hack) v1.15 J121, v1.03 Chicken Knife.sfc",
			want: "Final Fantasy V",
		},
		{
			name:     "unknown hack suffix",
			filename: "Game Alternate Physics (Hack) v1.0 Author.sfc",
			want:     "Game Alternate Physics v1.0 Author",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, Parse(tt.filename, false, slugs.MediaTypeGame))
			assert.Equal(t, tt.want, ParseDisplay(tt.filename, false, slugs.MediaTypeGame))
		})
	}
}

// TestParse_SceneReleases tests scene release artifact stripping from movie and TV show filenames.
func TestParse_SceneReleases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Movie scene releases
		{
			name:     "Movie with full scene release tags",
			input:    "The.Dark.Knight.2008.1080p.BluRay.x264.DTS-GROUP.mkv",
			expected: "The Dark Knight 2008",
		},
		{
			name:     "Movie with WEB-DL and audio codec",
			input:    "Avatar.2009.2160p.WEB-DL.DD5.1.H264-TEAM.mkv",
			expected: "Avatar 2009",
		},
		{
			name:     "Movie with HDR and HEVC",
			input:    "Blade.Runner.1982.4K.UHD.HDR10.x265-GROUP.mkv",
			expected: "Blade Runner 1982",
		},
		{
			name:     "Movie with PROPER tag",
			input:    "Inception.2010.1080p.BluRay.PROPER.x264-SPARKS.mkv",
			expected: "Inception 2010",
		},
		{
			name:     "Movie with extended cut and codec",
			input:    "Lord.of.the.Rings.2001.EXTENDED.1080p.BluRay.x264.TrueHD-FGT.mkv",
			expected: "Lord of the Rings 2001",
		},
		{
			name:     "Movie with Dolby Vision and Atmos",
			input:    "Dune.2021.2160p.WEB-DL.DV.HDR10.Atmos.HEVC-CMRG.mkv",
			expected: "Dune 2021",
		},
		{
			name:     "Movie with underscores",
			input:    "The_Matrix_1999_1080p_BluRay_x264_AAC-YTS.mkv",
			expected: "The Matrix 1999",
		},

		// TV show scene releases
		{
			name:     "TV show with episode and quality",
			input:    "The.Office.US.S01E01.Pilot.1080p.WEB-DL.DD5.1.H264-GROUP.mkv",
			expected: "The Office US Pilot",
		},
		{
			name:     "TV show with season pack notation",
			input:    "Breaking.Bad.S05E16.1080p.BluRay.x264-ROVERS.mkv",
			expected: "Breaking Bad",
		},
		{
			name:     "TV show with HDTV source",
			input:    "Game.of.Thrones.S08E03.720p.HDTV.x264-AVS.mkv",
			expected: "Game of Thrones",
		},
		{
			name:     "TV show with PROPER and WEBRip",
			input:    "Stranger.Things.S01E01.PROPER.720p.WEBRip.x264-RARBG.mkv",
			expected: "Stranger Things",
		},

		// Edge cases
		{
			name:     "ROM filename should be unchanged",
			input:    "Super Mario Bros (USA) [!].sfc",
			expected: "Super Mario Bros",
		},
		{
			name:     "Game with year in parentheses strips the year",
			input:    "Elden Ring (2022).exe",
			expected: "Elden Ring",
		},
		{
			name:     "Filename with spaces (no scene tags)",
			input:    "My Movie Title.mkv",
			expected: "My Movie Title",
		},
		{
			name:     "Multiple quality indicators",
			input:    "Movie.2020.UHD.2160p.HDR.BluRay.x265-GRP.mkv",
			expected: "Movie 2020",
		},
		{
			name:     "Release group without hyphen",
			input:    "Movie.Title.2010.1080p.BluRay.x264.mkv",
			expected: "Movie Title 2010",
		},
		{
			name:     "Mixed case scene tags",
			input:    "Show.Name.s02e05.WEBDL.h264.AAC.mkv",
			expected: "Show Name",
		},

		// Critical edge cases from expert review
		{
			name:     "Movie titled 'Cam' should not have title stripped",
			input:    "Cam.2018.1080p.WEB-DL.x264-GROUP.mkv",
			expected: "Cam 2018",
		},
		{
			name:     "Movie titled 'TS' should not have title stripped",
			input:    "TS.2020.720p.BluRay.x264.mkv",
			expected: "TS 2020",
		},
		{
			name:     "Short title before year should be preserved",
			input:    "IT.2017.1080p.BluRay.x264-SPARKS.mkv",
			expected: "IT 2017",
		},
		{
			name:     "Game title with dash and number should be preserved",
			input:    "Mega-Man-X-4.sfc",
			expected: "Mega Man X 4",
		},
		{
			name:     "Game title ending in single letter-number should be preserved",
			input:    "F-Zero-X.n64",
			expected: "F Zero X",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := Parse(tt.input, false, slugs.MediaTypeMovie)
			assert.Equal(t, tt.expected, got, "Title parsing mismatch")
		})
	}
}

// All-caps game names, such as a cartridge header dumped by the Operator
// bridge, end in the same "-WORD" shape as a scene release group. Stripping it
// cut "F-ZERO" to "F" and "SUPER MARIO ALL-STARS" to "SUPER MARIO ALL".
func TestParse_GamesKeepTrailingUppercaseHyphenWord(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    string
		expected string
		slug     string
	}{
		{input: "SUPER MARIO ALL-STARS.sfc", expected: "SUPER MARIO ALL-STARS", slug: "supermarioallstars"},
		{input: "F-ZERO.sfc", expected: "F-ZERO", slug: "fzero"},
		{input: "X-MEN.md", expected: "X-MEN", slug: "xmen"},
		{input: "STREET FIGHTER II-TURBO.sfc", expected: "STREET FIGHTER II-TURBO", slug: "streetfighter2turbo"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
			for _, mediaType := range []slugs.MediaType{"", slugs.MediaTypeGame} {
				title := Parse(tt.input, false, mediaType)
				assert.Equal(t, tt.expected, title)
				assert.Equal(t, tt.slug, slugs.Slugify(slugs.MediaTypeGame, title))
			}
			assert.Equal(t, tt.expected, Parse(tt.input, false, slugs.MediaTypeGame))
			assert.Equal(t, tt.expected, ParseDisplay(tt.input, false, slugs.MediaTypeGame))
		})
	}
}

func TestParse_VideoStripsReleaseGroup(t *testing.T) {
	t.Parallel()

	for _, mediaType := range []slugs.MediaType{slugs.MediaTypeMovie, slugs.MediaTypeTVShow, slugs.MediaTypeVideo} {
		assert.Equal(t, "Movie Title", Parse("Movie Title-GROUP.mkv", false, mediaType),
			"media type %s", mediaType)
	}
}

func TestStripStructuralSetMarkers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "disc number", input: "Final Fantasy VII (Disc 1)", want: "Final Fantasy VII"},
		{name: "disk of total", input: "Adventure (Disk 2 of 4)", want: "Adventure"},
		{name: "file number", input: "Magazine (File 3)", want: "Magazine"},
		{name: "side", input: "Cassette Game (Side A)", want: "Cassette Game"},
		{name: "compact cd", input: "Game [CD2]", want: "Game"},
		{name: "marker mid title", input: "Game (Disc 1) Special", want: "Game Special"},
		{name: "several markers", input: "Game (Disc 1) (Side B)", want: "Game"},
		{name: "other brackets kept", input: "Game (USA) (Disc 2) [!]", want: "Game (USA) [!]"},
		{name: "no brackets", input: "Disc Jam", want: "Disc Jam"},
		{name: "unclosed bracket kept", input: "Game (Disc 1", want: "Game (Disc 1"},
		{name: "unclosed bracket after marker", input: "Game (Disc 1) (USA", want: "Game (USA"},
		{name: "nested marker inside other group", input: "Game (Rev (Disc 1))", want: "Game (Rev (Disc 1))"},
		{name: "marker only", input: "(Disc 1)", want: "(Disc 1)"},
		{name: "empty", input: "", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, StripStructuralSetMarkers(tt.input))
		})
	}
}

func TestParseDisplay_PreservesOnlyStructuralSetMarkers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     string
		display   string
		canonical string
	}{
		{
			name:      "disc number preserved for display only",
			input:     "Final Fantasy VII (Disc 1).cue",
			display:   "Final Fantasy VII (Disc 1)",
			canonical: "Final Fantasy VII",
		},
		{
			name:      "disk total preserved and region stripped",
			input:     "Adventure (Disk 2 of 4) (USA).dsk",
			display:   "Adventure (Disk 2 of 4)",
			canonical: "Adventure",
		},
		{
			name:      "file number preserved and dump tag stripped",
			input:     "Magazine (File 3) [!].zip",
			display:   "Magazine (File 3)",
			canonical: "Magazine",
		},
		{
			name:      "side marker preserved",
			input:     "Cassette Game (Side A) (Europe).tap",
			display:   "Cassette Game (Side A)",
			canonical: "Cassette Game",
		},
		{
			name:      "compact cd marker preserved",
			input:     "Album (CD1) [FLAC].chd",
			display:   "Album (CD1)",
			canonical: "Album",
		},
		{
			name:      "normal metadata stripped",
			input:     "Super Mario Bros (USA) [!].sfc",
			display:   "Super Mario Bros",
			canonical: "Super Mario Bros",
		},
		{
			name:      "nonnumeric structural word stripped",
			input:     "Famicom Game (Disk Writer).fds",
			display:   "Famicom Game",
			canonical: "Famicom Game",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.display, ParseDisplay(tt.input, false, slugs.MediaTypeGame), "display title mismatch")
			assert.Equal(t, tt.canonical, Parse(tt.input, false, slugs.MediaTypeGame), "canonical title mismatch")
		})
	}
}

// TestStripSceneArtifacts tests the scene artifact stripping function directly.
func TestStripSceneArtifacts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Strip resolution",
			input:    "Movie Title 1080p mkv",
			expected: "Movie Title   mkv",
		},
		{
			name:     "Strip source type",
			input:    "Movie Title BluRay mkv",
			expected: "Movie Title   mkv",
		},
		{
			name:     "Strip video codec",
			input:    "Movie Title x264 mkv",
			expected: "Movie Title   mkv",
		},
		{
			name:     "Strip audio codec",
			input:    "Movie Title DTS mkv",
			expected: "Movie Title   mkv",
		},
		{
			name:     "Strip HDR format",
			input:    "Movie Title HDR10 mkv",
			expected: "Movie Title   mkv",
		},
		{
			name:     "Strip scene status",
			input:    "Movie Title PROPER mkv",
			expected: "Movie Title   mkv",
		},
		{
			name:     "Strip release group",
			input:    "Movie Title-GROUP",
			expected: "Movie Title",
		},
		{
			name:     "Strip all artifacts",
			input:    "Movie 2020 1080p BluRay x264 DTS HDR10-GROUP",
			expected: "Movie 2020          ",
		},
		{
			name:     "Leave non-scene text alone",
			input:    "Regular Movie Title mkv",
			expected: "Regular Movie Title mkv",
		},
		{
			name:     "Case insensitive matching",
			input:    "Movie webdl H265 aac mkv",
			expected: "Movie       mkv",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := stripSceneArtifacts(tt.input, true)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestStripSceneArtifacts_KeepsReleaseGroupShapeWhenDisabled(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "SUPER MARIO ALL-STARS", stripSceneArtifacts("SUPER MARIO ALL-STARS", false))
	assert.Equal(t, "Movie 2020    -GROUP", stripSceneArtifacts("Movie 2020 1080p BluRay-GROUP", false),
		"other scene artifacts are still stripped")
}

// TestParse_AmigaVision verifies that titles derived from
// AmigaVision listing lines (the path basename) are cleaned correctly. The scanner
// leaves ScanResult.Name empty so the indexer derives the title from the path,
// which means the raw listing line — e.g. "1869 (AGA)[en]" — arrives as the
// filename argument here.
func TestParse_AmigaVision(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		filename  string
		wantTitle string
	}{
		{
			name:      "game with AGA chipset and single language",
			filename:  "1869 (AGA)[en]",
			wantTitle: "1869",
		},
		{
			name:      "game with OCS chipset and single language",
			filename:  "3D Pool (OCS)[en]",
			wantTitle: "3D Pool",
		},
		{
			name:      "game with OCS chipset and multi-language bracket",
			filename:  "7 Colors (OCS)[en-de-fr-it-es]",
			wantTitle: "7 Colors",
		},
		{
			name:      "demo with group paren and AGA chipset",
			filename:  "1001 Stolen Ideas (Airwalk)(AGA)",
			wantTitle: "1001 Stolen Ideas",
		},
		{
			name:      "demo with group paren and OCS chipset",
			filename:  "9 Fingers (Spaceballs)(OCS)",
			wantTitle: "9 Fingers",
		},
		{
			name:      "game title with leading number preserved",
			filename:  "1869 (OCS)[en]",
			wantTitle: "1869",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.wantTitle, Parse(tt.filename, false, slugs.MediaTypeGame))
		})
	}
}

func TestSetNormalizeCache_SameTitles(t *testing.T) {
	// Not parallel: the cache is package state.
	inputs := []string{
		"Final Fantasy VII (USA) (Disc 1 of 3).chd",
		"Game (Side B) [CD2] (Europe).adf",
		"Super Castlevania IV (U) FastROM (Hack) v1.1 Vitor Vilela.sfc",
	}
	want := make([]string, len(inputs))
	for i, in := range inputs {
		want[i] = ParseDisplay(in, false, slugs.MediaTypeGame)
	}

	cache := make(map[string]string)
	SetNormalizeCache(cache)
	t.Cleanup(func() { SetNormalizeCache(nil) })
	for range 2 {
		for i, in := range inputs {
			assert.Equal(t, want[i], ParseDisplay(in, false, slugs.MediaTypeGame))
		}
	}
	normTagCacheMu.RLock()
	defer normTagCacheMu.RUnlock()
	assert.NotEmpty(t, normTagCache)
}
