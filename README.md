# go-slugs

Media title slugs for the [Zaparoo Project](https://zaparoo.org).

`Slugify` normalizes a title so that different spellings of the same game,
movie, TV episode or song give the same slug. Zaparoo Core uses it to index
media and to resolve titles written on tokens, and other Zaparoo clients use
it to match titles the same way Core does.

```go
import "github.com/ZaparooProject/go-slugs"

slugs.Slugify(slugs.MediaTypeGame, "The Legend of Zelda: Ocarina of Time (USA)")
// "legendofzeldaocarinaoftime"

slugs.Slugify(slugs.MediaTypeTVShow, "Breaking Bad - 1x02 - Gray Matter")
// same slug as "Breaking Bad - S01E02 - Gray Matter"
```

Changing the output of `Slugify` changes how existing media indexes and
written tokens match, so every change to its output is a breaking change for
its users.

Two subpackages cover the steps before slugging when a client indexes files
itself. `filetitle` derives a title from a filename the way Core does, and
`browseprefix` detects and strips directory-wide list numbering first:

```go
import "github.com/ZaparooProject/go-slugs/filetitle"

title := filetitle.Parse("Super Mario Bros (USA) [!].sfc", false, slugs.MediaTypeGame)
// "Super Mario Bros"
slugs.Slugify(slugs.MediaTypeGame, title)
```

## Ports

`testdata/parity.json` holds this module's output for every string literal in
its tests: slugs, tokens and query words, filename titles, and list prefixes,
plus the `width.Fold` table. Ports to other languages test against it. Run
`task parity` after any change to output or tests; `task test` fails while the
file is stale.

## Development

```sh
task test   # go test -race ./...
task lint   # golangci-lint run ./...
task fuzz   # fuzz targets, 30s each (FUZZ_TIME to override)
```

## License

Apache License 2.0. See [LICENSE](LICENSE).
