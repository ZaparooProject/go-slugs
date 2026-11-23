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

package parity

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var update = flag.Bool("update", false, "rewrite testdata/parity.json")

// TestFixturesUpToDate fails when module output or test inputs changed
// without the fixtures being regenerated, so ports never test against stale
// expectations. Run "task parity" to regenerate.
func TestFixturesUpToDate(t *testing.T) {
	t.Parallel()

	root := filepath.Join("..", "..")
	got, err := Generate(root)
	require.NoError(t, err)

	path := filepath.Join(root, "testdata", "parity.json")
	if *update {
		require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o750))
		require.NoError(t, os.WriteFile(path, got, 0o644)) //nolint:gosec // checked-in, public fixture
		return
	}
	want, err := os.ReadFile(path) //nolint:gosec // fixed path under the module
	require.NoError(t, err, "run task parity to generate fixtures")
	assert.Equal(t, string(want), string(got), "fixtures are stale; run task parity")
}
