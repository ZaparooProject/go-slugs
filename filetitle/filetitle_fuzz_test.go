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
	"strings"
	"testing"
	"unicode/utf8"

	slugs "github.com/ZaparooProject/go-slugs"
)

// FuzzParse tests title extraction with random inputs.
func FuzzParse(f *testing.F) {
	// ROM filenames
	f.Add("Super Mario Bros. (USA).nes", false)
	f.Add("Sonic the Hedgehog (Europe).md", false)
	f.Add("1942 (USA).nes", false)     // Year-like title
	f.Add("1. Game Title (USA)", true) // Leading number
	f.Add("01 - Song Name.mp3", true)  // Track number

	// Scene releases
	f.Add("The.Dark.Knight.2008.1080p.BluRay.x264-GROUP.mkv", false)
	f.Add("Movie.Title.2010.WEB-DL.DD5.1.H264.mp4", false)
	f.Add("Show.S01E01.PROPER.720p.HDTV.avi", false)

	// Edge cases
	f.Add("", false)
	f.Add("   ", false)
	f.Add("()", false)
	f.Add(".....", false)
	f.Add("_____", false)
	f.Add("-----", false)

	// Unicode
	f.Add("日本ゲーム (Japan).rom", false) //nolint:gosmopolitan // testing Unicode
	f.Add("Spëcîål Chàrs (World)", false)

	// Very long filename
	f.Add(strings.Repeat("A", 500)+"(USA).rom", false)

	// Control characters
	f.Add("Game\x00Name (USA)", false)
	f.Add("Game\tName\n(USA)", false)

	// Multiple extensions
	f.Add("File.rom.backup.zip", false)
	f.Add("Game..zip", false)

	// Structural set markers
	f.Add("Game (Disc 1) (USA).cue", false)
	f.Add("Game (Disk 2 of 4) [Side B", false)
	f.Add("Game ((Disc 1)) [CD2]", false)

	f.Fuzz(func(t *testing.T, filename string, stripLeadingNumbers bool) {
		stripped := StripStructuralSetMarkers(filename)
		if len(stripped) > len(filename) {
			t.Errorf("Stripping markers grew the title: %q from %q", stripped, filename)
		}
		if utf8.ValidString(filename) && !utf8.ValidString(stripped) {
			t.Errorf("Stripping markers broke UTF-8: %q from %q", stripped, filename)
		}
		if stripped != StripStructuralSetMarkers(filename) {
			t.Errorf("Non-deterministic marker stripping for: %q", filename)
		}

		result := Parse(filename, stripLeadingNumbers, slugs.MediaTypeGame)

		if !utf8.ValidString(result) {
			t.Errorf("Invalid UTF-8 in result: %q from filename: %q", result, filename)
		}

		// Invalid UTF-8 bytes get replaced with 3-byte replacement chars,
		// so output can be up to 3x input plus some overhead.
		if len(result) > len(filename)*3+10 {
			t.Errorf("Result unexpectedly longer: input=%d, output=%d for filename: %q",
				len(filename), len(result), filename)
		}

		if strings.Contains(result, "  ") {
			t.Errorf("Result contains multiple consecutive spaces: %q from filename: %q",
				result, filename)
		}

		if result != strings.TrimSpace(result) {
			t.Errorf("Result has leading/trailing whitespace: %q from filename: %q",
				result, filename)
		}

		result2 := Parse(filename, stripLeadingNumbers, slugs.MediaTypeGame)
		if result != result2 {
			t.Errorf("Non-deterministic result for filename: %q", filename)
		}
	})
}
