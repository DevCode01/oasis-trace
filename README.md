# Oasis Trace

**Go** · **v0.1.0**

_Access log analyzer: top IPs, paths, status codes and hourly traffic_

---

## Features

- Parses Apache/Nginx combined log format via regex
- Ranks top IPs and paths, breaks down status codes and bytes sent
- Hourly traffic histogram rendered as ASCII bars

## Installation

```bash
git clone https://github.com/DevCode01/oasis-trace.git
cd oasis-trace
```

## Usage

```bash
go run main.go access.log
```

## How it works

Parses each line against the Apache/Nginx combined log format with a single regex, then aggregates counts per status-code class, client IP, request path and hour-of-day bucket, printing the top 5 of each plus an ASCII histogram of traffic by hour.

## Options

- `<file>` - path to a log file (reads stdin if omitted)
- Unparseable lines are skipped and counted, never crash the run

## Example

```bash
$ go run main.go access.log
=== Access Log Report ===
Total requests:    4 (skipped 0 unparseable lines)
Total bytes sent:  2816

Status codes:
  2xx  3
  4xx  1

Top 5 IPs:
  192.168.1.10         2
  192.168.1.11         1
  192.168.1.12         1

Top 5 paths:
  /index.html                    2
  /about.html                    1
  /api/login                     1

Requests by hour:
  08h00  ##                                                 2
  09h00  ##                                                 2
```

## Testing

```bash
bash tests/run.sh
```

Includes a functional check against real captured output, plus the build/syntax check.

## License

MIT
