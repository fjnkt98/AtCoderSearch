# docker build -t atcodersearch-backend:0.1.0 .
# docker container run --rm --name atcodersearch-backend -p 8080:8080 -e API_SERVER_LISTEN_PORT=8000 atcodersearch-backend:0.1.0

FROM golang:1.20-bullseye AS stage

WORKDIR /app

COPY . .

RUN go build -o /app/atcodersearch

FROM debian:bullseye-slim

COPY --from=stage /app/atcodersearch /usr/local/bin/atcodersearch

RUN groupadd --gid 1000 atcodersearch \
    && useradd --uid 1000 --gid atcodersearch --shell /bin/bash --create-home atcodersearch
USER atcodersearch

CMD ["atcodersearch"]
