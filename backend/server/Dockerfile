FROM rust:1.66 AS builder

WORKDIR /server

COPY Cargo.toml Cargo.toml
COPY ./src ./src

RUN cargo build --release

FROM debian:11.6

RUN apt-get update \
    && apt-get -y install curl jq \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /server/target/release/server /usr/local/bin/server
CMD ["server"]
