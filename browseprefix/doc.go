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

// Package browseprefix detects list numbering such as "01 - " or date stamps
// at the start of filenames in a directory, and strips them the way Zaparoo
// Core does before it parses a title. Numbering is detected per directory, so
// a lone "1942.nes" keeps its number while a numbered playlist loses it.
package browseprefix
