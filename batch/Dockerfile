FROM golang:1.23-bullseye AS stage

ARG TARGETARCH

RUN apt-get update \
    && apt-get install -y \
    curl \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY ./crawl /app/crawl
COPY ./internal /app/internal
COPY ./pkg /app/pkg
COPY ./queries /app/queries
COPY ./repository /app/repository
COPY ./update /app/update
COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum
COPY ./main.go /app/main.go

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

COPY --from=stage --chown=atcodersearch:atcodersearch /app/atcodersearch /usr/local/bin/atcodersearch

USER atcodersearch

CMD ["atcodersearch"]
