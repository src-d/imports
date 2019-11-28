# Import extraction

CLI for fast import extraction.

## Installation

To build the project from source, run:

```bash
git clone https://github.com/src-d/imports.git
cd imports
go install ./cmd/sourced-imports
```

Usage example:

```bash
sourced-imports --root ./project-dir/
```

Example output:

```
{"file":"LICENSE","lang":"Text"}
{"file":"README.md","lang":"Markdown"}
{"file":"cmd/sourced-imports/main.go","lang":"Go","imports":["flag","fmt","github.com/src-d/imports","github.com/src-d/imports/languages/all","os"]}
...
```

## Using as library

It is possible to use `imports` as a library, but you'll need to add this workaround to your `go.mod`:

```
replace github.com/smacker/go-tree-sitter => github.com/dennwc/go-tree-sitter dev
```

This workaround is temporary and won't be required in the future.