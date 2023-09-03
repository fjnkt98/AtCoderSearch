# docker build -t atcodersearch-backend:0.1.0 .
# docker container run --rm --name atcodersearch-backend -p 8080:8080 -e API_SERVER_LISTEN_PORT=8000 atcodersearch-backend:0.1.0

FROM golang:1.20-bullseye AS stage

WORKDIR /app

COPY . .

RUN go build -o /app/atcodersearch

FROM debian:buster-slim

COPY --from=stage /app/atcodersearch /usr/local/bin/atcodersearch

CMD ["atcodersearch"]
