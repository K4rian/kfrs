# KF Redirect Server (KFRS)

A Docker image for the [Killing Floor Redirect Server (kfrs)][1] based on the official [Alpine Linux][2] [image][3].

---

## Environment variables
A few environment variables can be tweaked when creating a container to define the server configuration:

<details>
<summary>[Click to expand]</summary>

Variable          | Default value | Description 
---               | ---           | ---
KFRS_HOST         | 0.0.0.0       | IP/Host to bind to.
KFRS_PORT         | 9090          | TCP port to listen on.
KFRS_DIRECTORY    | ./redirect    | Directory to serve.
KFRS_MAX_REQUESTS | 20            | Max requests per IP/minute.
KFRS_BAN_TIME     | 15            | Ban duration (in minutes).

</details>

## Usage
Run the server using default configuration.<br>
Make sure that the `./redirect` directory exists before starting the container.
```bash
docker run -d \
  --name kfrs \
  -p 9090:9090/tcp \
  -e KFRS_HOST="0.0.0.0" \
  -e KFRS_PORT=9090 \
  -e KFRS_DIRECTORY="./redirect" \
  -e KFRS_MAX_REQUESTS=20 \
  -e KFRS_BAN_TIME=15 \
  -v ./redirect:/home/kfrs/redirect \
  -i k4rian/kfrs
```

## Using Compose
See the [docker-compose.yml][4] file.

## Manual build
__Requirements__:<br>
— Docker >= __18.09.0__<br>
— Git *(optional)*

Like any Docker image the building process is pretty straightforward: 

- Clone (or download) the GitHub repository to an empty folder on your local machine:
```bash
git clone https://github.com/K4rian/kfrs.git .
```

- Then run the following command inside the newly created folder:
```bash
cd docker/ && docker build --no-cache -t k4rian/kfrs .
```

[1]: https://github.com/K4rian/kfrs
[2]: https://www.alpinelinux.org/ "Alpine Linux Official Website"
[3]: https://hub.docker.com/_/alpine "Alpine Linux Docker Image"
[4]: https://github.com/K4rian/kfrs/blob/main/docker/docker-compose.yml