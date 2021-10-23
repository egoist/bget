**💛 You can help the author become a full-time open-source maintainer by [sponsoring him on GitHub](https://github.com/sponsors/egoist).**

---

# bget

[![npm version](https://badgen.net/npm/v/@egoist/bget)](https://npm.im/@egoist/bget)

> Download and install binaries from GitHub Releases.

## Preview

![Preview](https://user-images.githubusercontent.com/8784712/138543987-69075344-d781-4d9a-8a51-f6c33e37ba1c.gif)


## Install

One-off usage:

```bash
npx @egoist/bget user/repo
```

Or install as a global package:

```bash
npm i -g @egoist/bget

bget user/repo
```

## Usage

```bash
# From a repo's GitHub releases
bget user/repo

# Optionally a tag, default to `latest`
bget user/repo@v1.0.0
```

## TODO

- [ ] Support compressed files: `.zip` `.tar.gz` `.tar.xz`

## License

MIT &copy; [EGOIST](https://github.com/sponsors/egoist)
