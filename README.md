# chronytop

[![CI](https://github.com/geerlingguy/chronytop/actions/workflows/ci.yml/badge.svg)](https://github.com/geerlingguy/chronytop/actions/workflows/ci.yml)

A TUI for [chrony](https://chrony-project.org) NTP time server monitoring, in the style of [btop](https://github.com/aristocratos/btop).

![chronytop demo](resources/chronytop.gif)

It graphs NTP requests per second, and shows server counters, the clients chrony is rate limiting the most, sources, tracking, bandwidth, and CPU and memory use (including chronyd's).

Panels resize with the terminal, and can be shown or hidden.

## Requirements

- Linux, with a running chronyd 4.0 or later
- For Unix socket access: root (or `sudo`), or a user with access to the socket and write permission in its directory; see [Running without root](#running-without-root)

chronytop is a single static binary with no dependencies outside the Go standard library, so it runs on any distribution (Debian, Ubuntu, Fedora, RHEL, Arch, Alpine, etc.).

By default, it communicates directly with chronyd over the same Unix command socket normally used by `chronyc` for the server and client panels. Existing installations using that socket need no additional chronyd configuration. The `chronyc` executable is no longer required.

There is no automatic fallback to a different endpoint when the selected Unix socket is unavailable.

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
| `-socket` | `/var/run/chrony/chronyd.sock` | Absolute path to chronyd's Unix command socket |
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

Chronyd creates its Unix command socket by default; an explicit `bindcmdaddress`
directive is only needed to change its location or disable it. Chronytop's
upstream default is `/var/run/chrony/chronyd.sock`; on systems where `/var/run`
links to `/run`, this also reaches `/run/chrony/chronyd.sock`. Chronytop does not
parse `chrony.conf` or discover another binary's compiled-in defaults. Use
`-socket` for a different location.

Use `-socket /path/to/chronyd.sock` if chronyd uses a custom `bindcmdaddress` path. Two sockets are opened lazily and kept open: one for regular samples and one for the background client list. Each sample gets a fresh two-second deadline; client pagination gets a ten-second deadline for the entire list. A failed socket is closed and reopened on the next scheduled refresh. Valid command rejections do not close the connection. Connection and protocol errors appear in the affected panels. Counters not reported by older chrony versions are omitted from the server panel.

## Running without root

Unix socket mode requires permission to access chronyd's socket and to create
reply sockets in the same directory. On installations that grant these
permissions to chrony's group (`_chrony` on Debian and Ubuntu, `chrony` on
Fedora, RHEL, and Arch), add your user to that group, then log out and back in:

```
sudo usermod -aG _chrony $USER
```

Group membership alone is not sufficient if the socket directory is not
writable by the group. Check the distribution's socket directory permissions.
Access to chronyd's control socket permits more than reading its statistics:
group members can also step the clock or change sources with other clients.
Chronytop itself only sends monitoring commands.

## Security

Chronytop:

- Sends only read-only monitoring commands over the selected Unix socket and reads local `/proc` statistics. It never starts subprocesses, changes daemon settings, or resets counters.
- Writes its settings file only if you changed a setting. It skips the file and its directories unless they're owned by the user running chronytop (so a non-root `HOME` under `sudo` is ignored), and saves through a temporary file, so it never follows a symlink.
- Keeps reply sockets in private directories beside chronyd's socket for the lifetime of each connection, removing them on failure or normal shutdown. Permission changes and socket removal are anchored to a verified directory descriptor. Reply sockets allow chronyd's user to send responses; the connected socket accepts datagrams only from the selected daemon socket.
- Strips non-printable characters from everything it displays, so data from chrony can't inject terminal escape sequences.

CI runs `go vet`, `staticcheck`, `gosec`, and `govulncheck` on every change.

## Testing

Run `go test -race ./...` on Linux. The protocol tests use fake Unix daemons and require neither root nor an installed chronyd. To additionally query a running daemon without changing its configuration or counters:

```
sudo env CHRONYTOP_TEST_SOCKET=/run/chrony/chronyd.sock go test -run TestSocketLive -v
```

## License

MIT
