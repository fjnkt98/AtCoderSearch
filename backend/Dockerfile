FROM golang:1.21-bullseye AS stage

ARG TARGETARCH

RUN apt-get update \
    && apt-get install -y \
    curl \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY ./batch /app/batch
COPY ./cmd /app/cmd
COPY ./pkg /app/pkg
COPY ./repository /app/repository
COPY ./server /app/server
COPY ./settings /app/settings
COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum
COPY ./main.go /app/main.go
COPY ./partitions.sql /app/partitions.sql
COPY ./schema.hcl /app/schema.hcl

RUN go build -o /app/atcodersearch

# ---------------------------------------------------------------------------------------------
FROM debian:bullseye-slim

RUN apt-get update \
    && apt-get install -y \
    dnsutils \
    curl \
    ca-certificates \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

RUN ln -sf /usr/share/zoneinfo/Asia/Tokyo /etc/localtime
ENV TZ="Asia/Tokyo"

RUN groupadd --gid 1000 atcodersearch \
    && useradd --uid 1000 --gid atcodersearch --shell /bin/bash --create-home atcodersearch

RUN curl -sSf https://atlasgo.sh | sh

COPY --from=stage --chown=atcodersearch:atcodersearch /app/atcodersearch /usr/local/bin/atcodersearch
COPY --from=stage --chown=atcodersearch:atcodersearch /app/schema.hcl /etc/schema.hcl

USER atcodersearch

RUN mkdir -p /var/tmp/atcoder/problem \
    && mkdir -p /var/tmp/atcoder/user \
    && mkdir -p /var/tmp/atcoder/submission 
ENV PROBLEM_SAVE_DIR=/var/tmp/atcoder/problem USER_SAVE_DIR=/var/tmp/atcoder/user SUBMISSION_SAVE_DIR=/var/tmp/atcoder/submission

CMD ["atcodersearch"]
