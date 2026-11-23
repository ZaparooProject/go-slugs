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

package filetitle_test

import (
	"fmt"

	slugs "github.com/ZaparooProject/go-slugs"
	"github.com/ZaparooProject/go-slugs/filetitle"
)

func ExampleParse() {
	fmt.Println(filetitle.Parse("Super Mario Bros (USA) [!].sfc", false, slugs.MediaTypeGame))
	fmt.Println(filetitle.Parse("super_mario_bros.sfc", false, slugs.MediaTypeGame))
	fmt.Println(filetitle.Parse("1942 (USA)", false, slugs.MediaTypeGame))
	fmt.Println(filetitle.Parse("The.Dark.Knight.2008.1080p.BluRay.x264.mkv", false, slugs.MediaTypeMovie))
	// Output:
	// Super Mario Bros
	// super mario bros
	// 1942
	// The Dark Knight 2008
}

func ExampleParseDisplay() {
	fmt.Println(filetitle.ParseDisplay("Final Fantasy VII (USA) (Disc 1).chd", false, slugs.MediaTypeGame))
	fmt.Println(filetitle.Parse("Final Fantasy VII (USA) (Disc 1).chd", false, slugs.MediaTypeGame))
	// Output:
	// Final Fantasy VII (Disc 1)
	// Final Fantasy VII
}
