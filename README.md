# KF HTTP Redirect Server (kfrs)

## Overview
`kfrs` is a lightweight HTTP file server for the Killing Floor Dedicated Server. It serves .uz2 files from a specified directory while enforcing rate limits and IP bans to prevent excessive requests and ensure secure file access.

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
  --max-requests 10 \
  --ban-time 5
```

## License
[MIT][1]

[1]: https://