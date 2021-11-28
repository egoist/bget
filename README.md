# bget

> Download and install binary files from GitHub Releases.

## Install

For Mac users, you can use Homebrew to install it:

```bash
brew install egoist/tap/bget
```

For others:

```bash
curl -fsSL https://install.egoist.sh/bget.sh | bash -s -- -b /usr/local/bin
```

Or just grab the latest release from [GitHub Releases](https://github.com/egoist/bget/releases).

## Usage

```bash
bget owner/repo [-b bin_name] [-d install_dir]

# Download from a specific release
bget owner/repo#v1.2.3
```

- `bin_name` defaults to the name of the GitHub repo.
- `install_dir` defaults to `/usr/local/bin`.

## Example

```bash
bget egoist/doko
```

## Development

Build for release:

```bash
go build -o bget ./cmd
```

Development:

```bash
go run ./cmd
```

## License

MIT
