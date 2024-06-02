FROM golang:1.21-bullseye AS stage

ARG TARGETARCH

RUN apt-get update \
    && apt-get install -y \
    curl \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY . .
RUN go build -o /app/atcodersearch

RUN curl -L "https://github.com/sqldef/sqldef/releases/download/v0.17.10/psqldef_linux_${TARGETARCH}.tar.gz" -o 'psqldef.tar.gz' \
    && tar zxvf psqldef.tar.gz

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
COPY --from=stage --chown=atcodersearch:atcodersearch /app/psqldef /usr/local/bin/sqldef
COPY --from=stage --chown=atcodersearch:atcodersearch /app/schema.sql /etc/schema.sql

USER atcodersearch

RUN mkdir -p /var/tmp/atcoder/problem \
    && mkdir -p /var/tmp/atcoder/user \
    && mkdir -p /var/tmp/atcoder/submission 
ENV PROBLEM_SAVE_DIR=/var/tmp/atcoder/problem USER_SAVE_DIR=/var/tmp/atcoder/user SUBMISSION_SAVE_DIR=/var/tmp/atcoder/submission

CMD ["atcodersearch"]
