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

// tagType, tagValue and canonicalTag mirror the tag model Zaparoo Core uses
// when it extracts tags from filenames. Title parsing only needs the
// structural subset of that model to decide which bracket groups to keep.
type tagType string

type tagValue string

type tagSource uint8

const (
	// tagSourceBracketed marks a tag that came from (), {}, <> or [].
	tagSourceBracketed tagSource = iota + 1
)

type canonicalTag struct {
	Type   tagType
	Value  tagValue
	Source tagSource
}

const (
	tagTypeDisc       tagType = "disc"
	tagTypeDiscTotal  tagType = "disctotal"
	tagTypeMedia      tagType = "media"
	tagTypePatch      tagType = "patch"
	tagTypeUnlicensed tagType = "unlicensed"
)

const (
	tagMediaDisc  tagValue = "disc"
	tagMediaDisk  tagValue = "disk"
	tagMediaFile  tagValue = "file"
	tagMediaPart  tagValue = "part"
	tagMediaSideA tagValue = "side-a"
	tagMediaSideB tagValue = "side-b"
	tagMediaSideC tagValue = "side-c"
	tagMediaSideD tagValue = "side-d"
	tagMediaTape  tagValue = "tape"
)

const tagUnlicensedHack tagValue = "hack"
