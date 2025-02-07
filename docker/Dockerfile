FROM golang:alpine3.20 AS builder

WORKDIR /tmp/rserver

RUN apk update \
    && apk -U add --no-cache \
        git \
    && mkdir -p /server \
    && git clone "https://github.com/K4rian/kfrs.git" /tmp/rserver \
    && go mod download \
    && CGO_ENABLED=0 go build -a -o kfrs . \
    && cp /tmp/rserver/kfrs /server/kfrs \
    && chmod +x /server/kfrs \
    && rm -R /tmp/rserver

FROM alpine:3.20

ENV USERNAME=kfrs
ENV USERHOME=/home/$USERNAME

ENV KFRS_HOST="0.0.0.0"
ENV KFRS_PORT=9090
ENV KFRS_DIRECTORY="./redirect"
ENV KFRS_MAX_REQUESTS=20
ENV KFRS_BAN_TIME=15

RUN apk update \
    && adduser --disabled-password $USERNAME \
    && rm -rf /tmp/* /var/tmp/*

COPY --from=builder --chown=$USERNAME /server/* $USERHOME/

USER $USERNAME
WORKDIR $USERHOME

RUN mkdir -p "${KFRS_DIRECTORY}"

ENTRYPOINT ["./kfrs"]