# bget

> Download and install binary files from GitHub Releases.

## Usage

```bash
bget owner/repo [-b bin_name] [-d install_dir]
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
