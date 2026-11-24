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

// Package parity generates the fixtures that ports of this module to other
// languages check themselves against. A port that reproduces every fixture
// gives the same titles and slugs as Zaparoo Core for the same inputs.
package parity

import (
	"bytes"
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"maps"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"

	slugs "github.com/ZaparooProject/go-slugs"
	"github.com/ZaparooProject/go-slugs/browseprefix"
	"github.com/ZaparooProject/go-slugs/filetitle"
	"golang.org/x/text/width"
)

// maxInputBytes skips long test literals such as generated stress inputs,
// which add fixture size without covering new rules.
const maxInputBytes = 300

// MediaTypes are the media types with distinct parsing. Every other media
// type takes the title as is, the same as "Image".
var MediaTypes = []slugs.MediaType{
	slugs.MediaTypeGame,
	slugs.MediaTypeMovie,
	slugs.MediaTypeTVShow,
	slugs.MediaTypeMusic,
	slugs.MediaTypeImage,
}

// Fixtures is the file a port tests against.
type Fixtures struct {
	// WidthFold lists every code point golang.org/x/text/width.Fold changes,
	// as [from, to] pairs, so ports need not reproduce the Unicode width tables.
	WidthFold [][2]rune `json:"widthFold"`
	Cases     []Case    `json:"cases"`
}

// Case holds every output this module gives for one input.
type Case struct {
	Input string `json:"input"`
	// Slugs, Tokens and QueryWords are keyed by media type. Every map holds
	// "Game"; another media type is listed only where its value differs, and
	// "Image" stands for every media type without its own parsing.
	Slugs      map[slugs.MediaType]string   `json:"slugs"`
	Tokens     map[slugs.MediaType][]string `json:"tokens"`
	QueryWords map[slugs.MediaType][]string `json:"queryWords"`
	// Titles and DisplayTitles are keyed by media type, then by
	// stripLeadingNumbers. "false" is always present; "true" is listed only
	// where it differs. Media types are omitted as for Slugs.
	Titles        map[slugs.MediaType]map[string]string `json:"titles"`
	DisplayTitles map[slugs.MediaType]map[string]string `json:"displayTitles"`
	Unstructured  string                                `json:"stripStructuralSetMarkers"`
	Stem          string                                `json:"stem"`
	Prefix        browseprefix.Prefix                   `json:"prefix"`
}

// Generate builds fixtures from every string literal in the module's tests
// under root, so each rule the tests cover is covered for ports too.
func Generate(root string) ([]byte, error) {
	inputs, err := collectInputs(root)
	if err != nil {
		return nil, err
	}

	fx := Fixtures{WidthFold: widthFoldTable(), Cases: make([]Case, 0, len(inputs))}
	for _, in := range inputs {
		fx.Cases = append(fx.Cases, buildCase(in))
	}
	// One case per line keeps diffs of regenerated fixtures readable.
	var buf bytes.Buffer
	table, err := json.Marshal(fx.WidthFold)
	if err != nil {
		return nil, err //nolint:wrapcheck // json errors are self-describing
	}
	buf.WriteString("{\"widthFold\":")
	buf.Write(table)
	buf.WriteString(",\n\"cases\":[\n")
	for i, c := range fx.Cases {
		line, err := json.Marshal(c)
		if err != nil {
			return nil, err //nolint:wrapcheck // json errors are self-describing
		}
		buf.Write(line)
		if i < len(fx.Cases)-1 {
			buf.WriteByte(',')
		}
		buf.WriteByte('\n')
	}
	buf.WriteString("]}\n")
	return buf.Bytes(), nil
}

func buildCase(in string) Case {
	c := Case{
		Input:         in,
		Slugs:         make(map[slugs.MediaType]string, len(MediaTypes)),
		Tokens:        make(map[slugs.MediaType][]string, len(MediaTypes)),
		QueryWords:    make(map[slugs.MediaType][]string, len(MediaTypes)),
		Titles:        make(map[slugs.MediaType]map[string]string, len(MediaTypes)),
		DisplayTitles: make(map[slugs.MediaType]map[string]string, len(MediaTypes)),
		Unstructured:  filetitle.StripStructuralSetMarkers(in),
		Stem:          browseprefix.StemFromPath(in),
		Prefix:        browseprefix.ParseStem(in),
	}
	for _, mt := range MediaTypes {
		res := slugs.SlugifyWithTokens(mt, in)
		words := slugs.SlugifyQueryWords(mt, in)
		if words == nil {
			words = []string{}
		}
		titles := map[string]string{}
		display := map[string]string{}
		for _, strip := range []bool{false, true} {
			key := strconv.FormatBool(strip)
			titles[key] = filetitle.Parse(in, strip, mt)
			display[key] = filetitle.ParseDisplay(in, strip, mt)
		}
		if titles["true"] == titles["false"] {
			delete(titles, "true")
		}
		if display["true"] == display["false"] {
			delete(display, "true")
		}

		game := slugs.MediaTypeGame
		if mt == game || res.Slug != c.Slugs[game] {
			c.Slugs[mt] = res.Slug
		}
		if mt == game || !slices.Equal(res.Tokens, c.Tokens[game]) {
			c.Tokens[mt] = res.Tokens
		}
		if mt == game || !slices.Equal(words, c.QueryWords[game]) {
			c.QueryWords[mt] = words
		}
		if mt == game || !maps.Equal(titles, c.Titles[game]) {
			c.Titles[mt] = titles
		}
		if mt == game || !maps.Equal(display, c.DisplayTitles[game]) {
			c.DisplayTitles[mt] = display
		}
	}
	return c
}

func widthFoldTable() [][2]rune {
	var table [][2]rune
	for r := rune(0); r <= utf8.MaxRune; r++ {
		if r >= 0xD800 && r <= 0xDFFF {
			continue
		}
		s := string(r)
		folded := width.Fold.String(s)
		if folded == s {
			continue
		}
		out, size := utf8.DecodeRuneInString(folded)
		if size != len(folded) {
			panic("width.Fold produced more than one rune for " + strconv.QuoteRune(r))
		}
		table = append(table, [2]rune{r, out})
	}
	return table
}

func collectInputs(root string) ([]string, error) {
	seen := map[string]struct{}{}
	fset := token.NewFileSet()
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			if path != root && (strings.HasPrefix(d.Name(), ".") || d.Name() == "testdata") {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, "_test.go") {
			return nil
		}
		file, err := parser.ParseFile(fset, path, nil, parser.SkipObjectResolution)
		if err != nil {
			return err //nolint:wrapcheck // parse errors carry the path
		}
		ast.Inspect(file, func(n ast.Node) bool {
			lit, ok := n.(*ast.BasicLit)
			if !ok || lit.Kind != token.STRING {
				return true
			}
			s, err := strconv.Unquote(lit.Value)
			if err != nil || s == "" || len(s) > maxInputBytes || !utf8.ValidString(s) {
				return true
			}
			seen[s] = struct{}{}
			return true
		})
		return nil
	})
	if err != nil {
		return nil, err //nolint:wrapcheck // walk errors carry the path
	}
	inputs := make([]string, 0, len(seen))
	for s := range seen {
		inputs = append(inputs, s)
	}
	slices.Sort(inputs)
	return inputs, nil
}
