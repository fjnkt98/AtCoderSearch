# docker build -t atcodersearch-backend:latest .

FROM golang:1.21-bullseye AS stage

RUN apt-get update \
    && apt-get install -y \
    curl \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY . .
RUN go build -o /app/atcodersearch

RUN curl -L 'https://github.com/sqldef/sqldef/releases/download/v0.17.10/psqldef_linux_amd64.tar.gz' -o 'psqldef.tar.gz' \
    && tar zxvf psqldef.tar.gz

# ---------------------------------------------------------------------------------------------
FROM debian:bullseye-slim

RUN apt-get update \
    && apt-get install -y \
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

RUN mkdir -p /var/tmp/atcoder/problem
RUN mkdir -p /var/tmp/atcoder/user
RUN mkdir -p /var/tmp/atcoder/submission
RUN mkdir -p /var/tmp/atcoder/solution

CMD ["atcodersearch"]
