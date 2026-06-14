# wolframfns

Browse mathematical functions from [functions.wolfram.com](https://functions.wolfram.com) on the command line.

`wolframfns` is a single pure-Go binary. No API key required.

## Install

```bash
go install github.com/tamnd/wolframfns-cli/cmd/wolframfns@latest
```

Or grab a prebuilt binary from the [releases](https://github.com/tamnd/wolframfns-cli/releases), or run the container image:

```bash
docker run --rm ghcr.io/tamnd/wolframfns:latest --help
```

## Usage

```bash
# List all mathematical functions (default 50)
wolframfns list

# Filter by category
wolframfns list --category GammaBetaErf

# Search by name or category
wolframfns search Bessel

# List all categories with function counts
wolframfns categories

# Output formats: table (default TTY), json, jsonl, csv, tsv, url, raw
wolframfns list -o json
wolframfns list --category ElementaryFunctions -o csv
```

## Commands

| Command | Description |
|---------|-------------|
| `list` | List mathematical functions (optional `--category` filter) |
| `search <query>` | Search by function name or category |
| `categories` | List all categories with function counts |
| `version` | Show version information |

## Global flags

```
-o, --output string    output format: table|json|jsonl|csv|tsv|url|raw (default "auto")
-n, --limit int        limit number of records (0 = command default)
    --fields strings   comma-separated columns to include
    --no-header        omit header row
    --template string  Go text/template per record
    --timeout duration per-request timeout (default 90s)
    --delay duration   minimum spacing between requests
```
