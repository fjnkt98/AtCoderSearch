# docker build -t atcoder-search-indexer .
# docker run --rm --name indexer atcoder-search-indexer indexer generate

FROM rust:1.66 AS builder

WORKDIR /indexer

COPY Cargo.toml Cargo.toml
COPY ./src ./src

RUN cargo build --release

FROM debian:11.6

RUN apt-get update \
    && apt-get install -y \
    lsb-release \
    curl \
    gnupg2 \
    wget

RUN echo "deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list
RUN wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add -
RUN apt-get update \
    && apt-get -y install postgresql-client-15 \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

RUN cd /usr/local/bin && curl -L https://github.com/k0kubun/sqldef/releases/download/v0.15.12/psqldef_linux_amd64.tar.gz | tar xfz -

COPY --from=builder /indexer/target/release/indexer /usr/local/bin/indexer
COPY schema.sql /schema.sql
CMD ["indexer"]
