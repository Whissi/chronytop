# chronytop

[![CI](https://github.com/geerlingguy/chronytop/actions/workflows/ci.yml/badge.svg)](https://github.com/geerlingguy/chronytop/actions/workflows/ci.yml)

A TUI for [chrony](https://chrony-project.org) NTP time server monitoring, in the style of [btop](https://github.com/aristocratos/btop).

![chronytop screenshot](resources/screenshot.png)

It graphs NTP requests per second, and shows server counters, the clients chrony is rate limiting the most, sources, tracking, bandwidth, and CPU and memory use (including chronyd's).

Panels resize with the terminal, and can be shown or hidden.

## Requirements

- Linux, with chrony 4.0 or later (`/usr/bin/chronyc`)
- Root (or `sudo`) for the server and client panels, since chronyd only answers those over its root-owned socket

chronytop is a single static binary with no dependencies outside the Go standard library, so it runs on any distribution (Debian, Ubuntu, Fedora, RHEL, Arch, Alpine, etc.).

## Install

Download the binary for your architecture from [Releases](https://github.com/geerlingguy/chronytop/releases), and install it:

```
sudo install -o root -g root -m 755 chronytop-linux-amd64 /usr/local/bin/chronytop
```

Or build it from source with Go 1.24 or later:

```
CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o chronytop .
```

## Usage

```
sudo chronytop
```

| Option | Default | Description |
| --- | --- | --- |
| `-i` | `1s` | Sample interval (500ms to 10s) |
| `-p` | `100` | Minimum packets for a client to be listed in offenders |
| `-version` | | Print the version and exit |

| Key | Action |
| --- | --- |
| `1`-`7` | Show/hide a panel |
| `+` / `-` | Slower / faster sampling |
| `?` / `h` | Show/hide help |
| `Esc` | Close help, or quit |
| `q` | Quit |

Visible panels and the sample interval are saved when you quit, to `~/.config/chronytop/chronytop.conf` (`/root/.config/...` when run with `sudo`), and restored on the next run. `-i` overrides the saved interval.

The offenders list is refreshed every 10 seconds in the background. chronyd filters clients by `-p` itself, so the list stays cheap to fetch on busy public servers with hundreds of thousands of clients.

## Security

Since it runs as root, chronytop:

- Only reads from `/proc` and `chronyc`, and never opens network connections.
- Only writes its settings file, and only if you changed a setting. It skips the file and its directories unless they're owned by the user running chronytop (so a non-root `HOME` under `sudo` is ignored), and saves through a temporary file, so it never follows a symlink.
- Runs `/usr/bin/chronyc` by absolute path, with a minimal environment and no shell.
- Strips non-printable characters from everything it displays, so data from chrony can't inject terminal escape sequences.

CI runs `go vet`, `staticcheck`, `gosec`, and `govulncheck` on every change.

## License

MIT
