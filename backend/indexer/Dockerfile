FROM rust:1.66 AS builder

WORKDIR /indexer

COPY Cargo.toml Cargo.toml
COPY ./src ./src

ENV DATABASE_URL postgres://atcoder:atcoder@localhost:5432/atcoder

RUN cargo build --release

FROM debian:11.6

COPY --from=builder /indexer/target/release/indexer /usr/local/bin/indexer
CMD ["indexer"]
