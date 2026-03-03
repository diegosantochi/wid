# wid — What Was I Doing?

A minimal CLI/TUI application for keeping track of tasks and notes.

## Installation

Binaries are available in the releases/tags section

### With Go

```sh
go install github.com/diegosantochi/wid@latest
```

Or clone and build manually:

```sh
git clone https://github.com/diegosantochi/wid.git
cd wid
go build -o wid .
```

## Usage

### Open the TUI

```sh
wid
```

### Initialize a project-local store

Creates a `.wid.yaml` in the current directory (useful for per-directory task lists):

```sh
wid init
```

## Storage

`wid` looks for a `.wid.yaml` file starting from the current directory and walking up through parent directories. If none is found, it falls back to `~/.wid.yaml`.

You can:
- Keep a **global** task list in your home directory (automatic, no setup needed)
- Keep a **directory-specific** list by running `wid init` inside the directory

## License

[GPL-3.0](LICENSE)
