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

// Package slugs turns media titles into slugs for matching. Two spellings of
// the same title, such as "The Legend of Zelda: Ocarina of Time (USA)" and
// "legend of zelda - ocarina of time", give the same slug, so a title written
// on a token can be found in a media index built from file names.
//
// Parsing depends on the media type: TV episode markers, movie and music
// scene tags and game metadata brackets are each handled by their own rules
// before the shared normalization pipeline runs.
package slugs
