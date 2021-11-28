# bget

> Download and install binary files from GitHub Releases.

## Preview

![CleanShot 2021-11-29 at 00 54 42](https://user-images.githubusercontent.com/8784712/143778020-25b8de62-5b90-4097-8f11-d8ef2db172db.gif)

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

It will look for executable files (as well as compressed files) in the release assets, if it's a compressed file we simply use the largest file inside it.

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
