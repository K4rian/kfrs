# KF Redirect Server (KFRS)

## Overview
`kfrs` is a lightweight HTTP file server for the Killing Floor Dedicated Server (KFDS).<br>
It serves .uz2 files from a specified directory while enforcing rate limits and IP bans to prevent excessive requests and ensure secure file access.

## Features
- Serves `.uz2` files from a specified directory.
- Blocks requests that exceed a configurable limit per IP.
- Enforces request filtering and IP banning.
- Allows only `GET` requests.

## Usage
### Basic Usage
```bash
./kfrs --host "0.0.0.0" \
  --port 9090 \
  --serve-dir "./redirect" \
  --max-requests 20 \
  --ban-time 15
```

### With File Logging
Enable logging to a file (text format):
```bash
./kfrs --host "0.0.0.0" \      # IP/Host to bind to
  --port 9090 \                # TCP port to listen on
  --serve-dir "./redirect" \   # Directory to serve files from
  --max-requests 20 \          # Max requests per IP/minute before banning
  --ban-time 15 \              # Ban duration (in minutes)
  --log-to-file \              # Enable file logging
  --log-level "info" \         # Log level (info, debug, warn, error)
  --log-file "./kfrs.log" \    # Log file path
  --log-file-format "text" \   # Log format (text or json)
  --log-max-size 10 \          # Max log file size (MB)
  --log-max-backups 5 \        # Max number of old log files to keep
  --log-max-age 28             # Max age of a log file (days)
```

### Using Environment Variables
You can also configure `kfrs` using environment variables:
```bash
export KFRS_HOST="0.0.0.0"         # IP/Host to bind to
export KFRS_PORT=9090              # TCP port to listen on
export KFRS_SERVE_DIR="./redirect" # Directory to serve files from
export KFRS_MAX_REQUESTS=20        # Max requests per IP/minute before banning
export KFRS_BAN_TIME=15            # Ban duration (in minutes)
export KFRS_LOG_TO_FILE=true       # Enable file logging
export KFRS_LOG_LEVEL="info"       # Log level (info, debug, warn, error)
export KFRS_LOG_FILE="./kfrs.log"  # Log file path
export KFRS_LOG_FILE_FORMAT="text" # Log format (text or json)
export KFRS_LOG_MAX_SIZE=10        # Max log file size (MB)
export KFRS_LOG_MAX_BACKUPS=5      # Max number of old log files to keep
export KFRS_LOG_MAX_AGE=28         # Max age of a log file (days)
./kfrs
```

You can add these export commands to a .env file and source it before running the server:
```bash
source kfrs.env && ./kfrs
```

## Using Docker
See [docker/][1]

## Building
Building is done with the `go` tool. If you have setup your `GOPATH` correctly, the following should work:
```bash
go get github.com/k4rian/kfrs
go build -ldflags "-w -s" github.com/k4rian/kfrs
```

## License
[MIT][2]

[1]: https://github.com/K4rian/kfrs/blob/main/docker
[2]: https://github.com/K4rian/kfrs/blob/main/LICENSE
