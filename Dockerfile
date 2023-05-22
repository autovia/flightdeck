# syntax=docker/dockerfile:1

# Copyright (c) Autovia GmbH
# SPDX-License-Identifier: Apache-2.0
# 
# This Dockerfile contains multiple targets.
# e.g. `docker build --target=release .`
#
# The release target needs a RELEASE_VERSION argument that must be provided 
# e.g. `docker build --build-arg RELEASE_VERSION=0.0.1 --target=release -t autovia/flightdeck:0.0.1 .`

FROM ubuntu AS dev

RUN apt update
RUN apt install -y nodejs npm wget

WORKDIR /usr/local/

ENV GO_VERSION=1.19.8
RUN wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz
ENV PATH "$PATH:/usr/local/go/bin"
RUN go version

RUN mkdir /usr/local/flightdeck
WORKDIR /usr/local/flightdeck

COPY . .

FROM dev AS build

WORKDIR /usr/local/flightdeck
RUN ./api/bin/build

WORKDIR /usr/local/flightdeck/ui
RUN npm install
RUN npm run build

ENTRYPOINT [ "/usr/local/flightdeck/api/bin/flightdeck.linux.amd64", "-fileserver=true", "-fileserverpath=/usr/local/flightdeck/ui/build", "-incluster=true", "-addr=0.0.0.0:3000" ]

FROM alpine AS release

LABEL name="flightdeck" \
      maintainer="Autovia Team <team@autovia.io>" \
      vendor="Autovia GmbH" \
      version=$RELEASE_VERSION \
      release=$RELEASE_VERSION \
      summary="flightdeck is a web-based interactive diagram UI for Kubernetes to observe and troubleshoot applications running in the cluster." \
      description="flightdeck is a web-based interactive diagram UI for Kubernetes to observe and troubleshoot applications running in the cluster. Please submit issues to https://github.com/autovia/flightdeck/issues"

RUN mkdir /usr/local/flightdeck

COPY --from=build /usr/local/flightdeck/api/bin/flightdeck.linux.amd64 /usr/local/flightdeck/flightdeck.linux.amd64
COPY --from=build /usr/local/flightdeck/ui/build /usr/local/flightdeck/dist/

ENTRYPOINT [ "/usr/local/flightdeck/flightdeck.linux.amd64", "-fileserver=true", "-fileserverpath=/usr/local/flightdeck/dist", "-incluster=true", "-addr=0.0.0.0:3000" ]

FROM dev