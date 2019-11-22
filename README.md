# Import extraction

CLI for fast import extraction.

To build the project from source, run:

```bash
git clone https://github.com/src-d/imports.git
cd imports
GO111MODULE=off go build ./cmd/sourced-imports
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