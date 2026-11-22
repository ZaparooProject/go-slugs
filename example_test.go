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

package slugs_test

import (
	"fmt"

	"github.com/ZaparooProject/go-slugs"
)

func ExampleSlugify() {
	fmt.Println(slugs.Slugify(slugs.MediaTypeGame, "The Legend of Zelda: Ocarina of Time (USA)"))

	a := slugs.Slugify(slugs.MediaTypeTVShow, "Breaking Bad - 1x02 - Gray Matter")
	b := slugs.Slugify(slugs.MediaTypeTVShow, "Breaking Bad - S01E02 - Gray Matter")
	fmt.Println(a == b)
	// Output:
	// legendofzeldaocarinaoftime
	// true
}
