FROM rust:1.66 AS builder

WORKDIR /server

COPY Cargo.toml Cargo.toml
COPY ./src ./src

RUN cargo build --release

FROM debian:11.6

COPY --from=builder /server/target/release/server /usr/local/bin/server
CMD ["server"]
