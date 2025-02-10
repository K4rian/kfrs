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
```bash
./kfrs --host "0.0.0.0" \
  --port 9090 \
  --directory "./redirect" \
  --max-requests 20 \
  --ban-time 15
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
