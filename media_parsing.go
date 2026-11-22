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

// ParseWithMediaType is the entry point for media-type-aware parsing.
// It delegates to the appropriate parser based on media type.
// Each parser applies media-specific normalization BEFORE the universal pipeline.
//
// Media-specific parsers are implemented in separate files:
//   - ParseTVShow → media_parsing_tv.go
//   - ParseGame → media_parsing_game.go
//   - ParseMovie → media_parsing_movie.go
//   - ParseMusic → media_parsing_music.go
func ParseWithMediaType(mediaType MediaType, title string) string {
	switch mediaType {
	case MediaTypeTVShow:
		return ParseTVShow(title)
	case MediaTypeGame:
		return ParseGame(title)
	case MediaTypeMovie:
		return ParseMovie(title)
	case MediaTypeMusic:
		return ParseMusic(title)
	case MediaTypeAudio:
		// TODO: Implement ParseAudio in media_parsing_audio.go
		// Future work needed:
		//   - Audiobooks: Extract chapter/part numbers, strip "Unabridged"/"Abridged" markers
		//   - Podcasts: Normalize episode numbers, strip date formats
		//   - General: Strip narrator names if in metadata format
		// Example: "Book Title (Unabridged) - Chapter 1" → "Book Title Chapter 1"
		return title
	case MediaTypeVideo:
		// TODO: Implement ParseVideo in media_parsing_video.go
		// Future work needed:
		//   - Music videos: Normalize "Artist - Song (Official Video)" formats
		//   - Strip quality tags (1080p, 4K, etc.) similar to movies
		//   - Handle "feat." and artist collaborations
		// Example: "Artist - Song Title (Official Music Video) [4K]" → "Artist Song Title"
		return title
	case MediaTypeImage, MediaTypeApplication:
		// No special parsing needed for images and applications
		return title
	default:
		return title
	}
}
