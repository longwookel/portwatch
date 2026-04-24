# portwatch

A lightweight CLI daemon that monitors open ports and alerts on unexpected changes.

---

## Installation

```bash
go install github.com/yourusername/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portwatch.git && cd portwatch && go build -o portwatch .
```

---

## Usage

Start the daemon with default settings:

```bash
portwatch start
```

Specify a custom polling interval and alert on any new or closed ports:

```bash
portwatch start --interval 30s --notify
```

Take a snapshot of the current port state to use as a baseline:

```bash
portwatch snapshot
```

View detected changes since the last snapshot:

```bash
portwatch diff
```

### Example Output

```
[2024-01-15 10:32:01] ALERT: New port opened   → TCP :8080
[2024-01-15 10:45:17] ALERT: Port closed       → TCP :3000
[2024-01-15 11:00:00] OK: No changes detected
```

---

## Configuration

portwatch looks for a config file at `~/.portwatch.yaml`:

```yaml
interval: 60s
notify: true
ignore:
  - 22
  - 80
  - 443
```

---

## License

MIT © 2024 [yourusername](https://github.com/yourusername)