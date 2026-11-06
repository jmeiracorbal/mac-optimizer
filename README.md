# mac-optimizer

[![CI](https://github.com/jmeiracorbal/mac-optimizer/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/jmeiracorbal/mac-optimizer/actions/workflows/ci.yml)
[![Go](https://img.shields.io/badge/Go-1.22%2B-00ADD8?logo=go)](https://go.dev/)
[![Platform](https://img.shields.io/badge/platform-macOS%2012%2B-lightgrey?logo=apple)](https://www.apple.com/macos/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

A command-line tool for running maintenance tasks on macOS. It can clean selected caches and logs, restart a few system processes, and run other cleanup tasks on demand.

The program supports macOS Monterey (12) and later on Intel and Apple Silicon. It uses macOS system utilities and requires Go only when building from source.

## Install

### Homebrew

```sh
brew tap jmeiracorbal/tools
brew install mac-optimizer
```

### Installer script

The installer downloads the latest release for the current Mac and installs it to `/usr/local/bin` by default:

```sh
curl -fsSL https://raw.githubusercontent.com/jmeiracorbal/mac-optimizer/main/install.sh | sh
```

Set `INSTALL_DIR` to choose another destination. You may need to run the installer with permissions that allow writing to that directory.

### Build from source

```sh
git clone https://github.com/jmeiracorbal/mac-optimizer.git
cd mac-optimizer
go build -o mac-optimizer ./cmd/mac-optimizer
```

## Usage

Running without module names executes the default sequence: `cache`, `spotlight`, `permissions`, `services`, `memory`, and `logs`.

```sh
mac-optimizer                  # Run the default sequence
mac-optimizer --dry-run        # Preview supported file and command changes
mac-optimizer cache logs       # Run selected modules
mac-optimizer browsers downloads
mac-optimizer --help
mac-optimizer --version
```

### Options

| Option | Description |
| --- | --- |
| `-h`, `--help` | Show command help |
| `-v`, `--version` | Show version |
| `-d`, `--dry-run` | Show planned cleanup and skip commands that change the system |
| `-q`, `--quiet` | Suppress log messages; the header and statistics are still printed |
| `-f`, `--force` | Skip confirmation prompts |
| `--log-file FILE` | Set the log file path (default: `~/.mac-optimizer.log`) |

Available modules: `cache`, `spotlight`, `permissions`, `services`, `memory`, `logs`, `browsers`, and `downloads`. Use `all` to run every module.

## Modules

| Module | What it does |
| --- | --- |
| `cache` | Empties selected folders in `~/Library` and removes files older than one day from `/Library/Caches` and `/private/var/log/asl`. Also flushes the DNS cache. |
| `spotlight` | Disables Spotlight indexing, removes the existing index, and re-enables indexing. The new index is built in the background. |
| `permissions` | Attempts the legacy `diskutil repairPermissions` command, then sets `/usr/local/bin` and the home directory to mode `755`. |
| `services` | Sends termination signals to Finder, Dock, SystemUIServer, and ControlStrip so macOS can relaunch them. |
| `memory` | Runs macOS `purge` with administrator privileges. |
| `logs` | Deletes `.log` and `.out` files older than 30 days from `~/Library/Logs` and `/private/var/log`. |
| `browsers` | Clears selected Chrome, Firefox, Safari, and WebKit data. See the data warning below. |
| `downloads` | Deletes files older than 30 days under `~/Downloads`, after confirmation. |

Some modules call commands through `sudo` and may request an administrator password. The program writes a timestamped log to `~/.mac-optimizer.log` unless `--log-file` is set.

## Review before running

- `--dry-run` skips system-changing commands and previews supported file cleanup. Read-only system statistics and checks may still run.
- `--force` accepts confirmation prompts automatically. Use it only when you have reviewed the selected modules.
- `spotlight` removes the existing Spotlight index before macOS rebuilds it.
- `permissions` changes the permissions of your home directory and `/usr/local/bin`; it does not perform a general permissions repair on current macOS versions.
- `browsers` empties the Firefox `Profiles` directory and clears Safari local storage and databases as well as browser caches. This can remove browser profile data and sign you out of websites.
- `downloads` removes files based on age, not file type. Review the count and size shown before confirming.
- Cache and log cleanup permanently removes files from the listed locations.

Run only the modules you intend to use. For example, `mac-optimizer cache --dry-run` previews cache cleanup without applying it.

## Development

The project uses Go 1.22 or later. Build the CLI with:

```sh
go build ./cmd/mac-optimizer
```

## License

MIT. See [LICENSE](LICENSE).
