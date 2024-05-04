# docker build -t atcodersearch-backend:latest .
# docker container run --rm --name atcodersearch-backend -p 8080:8080 -e API_SERVER_LISTEN_PORT=8000 atcodersearch-backend:latest

FROM golang:1.21-bullseye AS stage

WORKDIR /app

COPY . .

RUN go build -o /app/atcodersearch

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

USER atcodersearch

RUN mkdir -p /var/tmp/atcoder/problem
RUN mkdir -p /var/tmp/atcoder/user
RUN mkdir -p /var/tmp/atcoder/submission
RUN mkdir -p /var/tmp/atcoder/solution

CMD ["atcodersearch"]
