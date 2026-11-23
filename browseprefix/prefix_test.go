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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseStem(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		stem string
		kind Kind
		rest string
	}{
		{name: "rank dash", stem: "001 - Sonic", kind: KindRank, rest: "Sonic"},
		{name: "rank period", stem: "42. Contra", kind: KindRank, rest: "Contra"},
		{name: "rank paren", stem: "7) Zelda", kind: KindRank, rest: "Zelda"},
		// ReplayOS, 2026-10-26: titles that start with a number lost it.
		{name: "hyphenated title", stem: "10-Yard Fight (USA, Europe)", kind: KindNone},
		{name: "hyphenated title no space", stem: "3-D WorldRunner (USA)", kind: KindNone},
		{name: "number word title", stem: "720 Degrees (USA)", kind: KindNone},
		{name: "date full", stem: "1991-06-23 - Sonic", kind: KindDate, rest: "Sonic"},
		{name: "date dotted", stem: "1991.06.23 Sonic", kind: KindDate, rest: "Sonic"},
		{name: "date year", stem: "1991 - Sonic", kind: KindDate, rest: "Sonic"},
		{name: "date year period", stem: "1991. Sonic", kind: KindDate, rest: "Sonic"},
		{name: "date year month space", stem: "1989-01 Bomberman (USA)", kind: KindDate, rest: "Bomberman (USA)"},
		// MiSTer, 2026-11-03: a year that starts the title is not a date prefix.
		{name: "year title", stem: "1943 The Battle of Midway (US)", kind: KindNone},
		{name: "year title short", stem: "2020 Super Baseball", kind: KindNone},
		{name: "year title hyphen", stem: "1999-Nights", kind: KindNone},
		{name: "bare numeric title", stem: "1942", kind: KindNone, rest: ""},
		{name: "number plus letter title", stem: "3D Worldrunner", kind: KindNone, rest: ""},
		{name: "invalid date", stem: "1991-13-01 - Sonic", kind: KindNone, rest: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := ParseStem(tt.stem)
			assert.Equal(t, tt.kind, got.Kind)
			assert.Equal(t, tt.rest, got.Rest)
		})
	}
}

func TestStemFromPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		path string
		stem string
	}{
		{
			name: "filename with extension",
			path: filepath.Join("roms", "genesis", "001 - Sonic.gen"),
			stem: "001 - Sonic",
		},
		{
			name: "nested date file",
			path: filepath.Join("roms", "genesis", "history", "1991-06-23 - Sonic.zip"),
			stem: "1991-06-23 - Sonic",
		},
		{name: "no extension", path: filepath.Join("roms", "genesis", "Sonic"), stem: "Sonic"},
		{
			name: "dot in title",
			path: filepath.Join("roms", "genesis", "Sonic. The Hedgehog.gen"),
			stem: "Sonic. The Hedgehog",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.stem, StemFromPath(tt.path))
		})
	}
}

func TestStripWithPolicy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		path    string
		rest    string
		parsed  Kind
		policy  Policy
		strips  bool
		enabled bool
	}{
		{
			name:    "enabled rank policy strips rank prefix",
			path:    filepath.Join("roms", "genesis", "001 - Sonic.gen"),
			policy:  Policy{Kind: KindRank, Enabled: true},
			rest:    "Sonic",
			strips:  true,
			parsed:  KindRank,
			enabled: true,
		},
		{
			name:    "disabled policy keeps stem",
			path:    filepath.Join("roms", "genesis", "001 - Sonic.gen"),
			policy:  Policy{Kind: KindRank},
			rest:    "001 - Sonic",
			parsed:  KindRank,
			enabled: false,
		},
		{
			name:    "kind mismatch keeps stem",
			path:    filepath.Join("roms", "genesis", "Sonic.gen"),
			policy:  Policy{Kind: KindDate, Enabled: true},
			rest:    "Sonic",
			parsed:  KindNone,
			enabled: true,
		},
		{
			name:    "empty rest keeps stem",
			path:    filepath.Join("roms", "genesis", "001.gen"),
			policy:  Policy{Kind: KindRank, Enabled: true},
			rest:    "001",
			parsed:  KindNone,
			enabled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			stem := StemFromPath(tt.path)
			prefix := ParseStem(stem)
			got, ok := StripWithPolicy(stem, tt.policy)

			assert.Equal(t, tt.parsed, prefix.Kind)
			assert.Equal(t, tt.rest, got)
			assert.Equal(t, tt.strips, ok)
			assert.Equal(t, tt.enabled, tt.policy.Enabled)
		})
	}
}

func TestDetectPolicyForPathsUsesStems(t *testing.T) {
	t.Parallel()

	paths := []string{
		filepath.Join("roms", "nes", "1942.nes"),
		filepath.Join("roms", "nes", "007.nes"),
		filepath.Join("roms", "nes", "3D Worldrunner.nes"),
		filepath.Join("roms", "nes", "720 Degrees.nes"),
		filepath.Join("roms", "nes", "8 Eyes.nes"),
	}

	policy := DetectPolicyForPaths(paths, DefaultThreshold, DefaultMinFiles)
	assert.False(t, policy.Enabled)
}
