pub mod atcodersearch {
    tonic::include_proto!("atcodersearch");
}
mod problem;

use anyhow::Context;
use problem::ProblemSearcher;
use sqlx::postgres::PgPoolOptions;
use tonic::transport::Server;

use atcodersearch::problem_service_server::ProblemServiceServer;
use clap::Parser;
use meilisearch_sdk::client::Client;
use std::net::{IpAddr, Ipv4Addr, SocketAddr};
use std::{env, str::FromStr};
use tracing::level_filters::LevelFilter;
use tracing_subscriber::{
    fmt::{self, time::OffsetTime},
    EnvFilter,
};

#[derive(Parser)]
pub struct App {
    #[arg(long, env, hide_env_values = true)]
    database_url: String,
    #[arg(long, env, hide_env_values = true)]
    engine_url: String,
    #[arg(long, env, hide_env_values = true)]
    engine_master_key: String,
    #[arg(long, env, default_value_t = 8000)]
    port: u16,
}

impl App {
    #[tokio::main]
    pub async fn run(&self) -> anyhow::Result<()> {
        let addr = SocketAddr::new(IpAddr::V4(Ipv4Addr::new(127, 0, 0, 1)), self.port);

        let pool = PgPoolOptions::new()
            .connect(&self.database_url)
            .await
            .with_context(|| "connect to database")?;

        let client = Client::new(&self.engine_url, Some(&self.engine_master_key))
            .with_context(|| "create engine client")?;
        let index = client.index("problems");

        tracing::info!("start grpc server at port {}", self.port);
        Server::builder()
            .add_service(ProblemServiceServer::new(ProblemSearcher::new(pool, index)))
            .serve(addr)
            .await?;

        Ok(())
    }
}

fn main() -> anyhow::Result<()> {
    let level = env::var("RUST_LOG").unwrap_or(String::from("info"));
    let filter = EnvFilter::builder()
        .with_default_directive(LevelFilter::from_str(&level)?.into())
        .from_env_lossy();
    let format = fmt::format()
        .with_level(true)
        .with_target(true)
        .with_ansi(false)
        .with_thread_ids(false)
        .with_timer(OffsetTime::local_rfc_3339()?);
    let subscriber = tracing_subscriber::fmt()
        .with_env_filter(filter)
        .event_format(format)
        .json()
        .finish();
    tracing::subscriber::set_global_default(subscriber)?;

    App::parse().run()
}

#[cfg(test)]
pub mod testutil {
    use meilisearch_sdk::client::Client;
    use sqlx::{Pool, Postgres};
    use std::env;
    use testcontainers::{
        core::{IntoContainerPort, Mount, WaitFor},
        ContainerAsync, ContainerRequest, GenericImage, ImageExt,
    };

    pub fn create_db_container() -> anyhow::Result<ContainerRequest<GenericImage>> {
        let schema = env::current_dir()?
            .join("../db/schema.sql")
            .canonicalize()?;

        let request = GenericImage::new("postgres", "16-bullseye")
            .with_exposed_port(5432.tcp())
            .with_wait_for(WaitFor::message_on_stdout(
                "database system is ready to accept connections",
            ))
            .with_wait_for(WaitFor::message_on_stderr(
                "database system is ready to accept connections",
            ))
            // .with_wait_for(WaitFor::millis(1000))
            .with_mount(
                Mount::bind_mount(
                    schema
                        .into_os_string()
                        .into_string()
                        .map_err(|e| anyhow::anyhow!("{:?}", e))?,
                    "/docker-entrypoint-initdb.d/schema.sql",
                )
                .with_access_mode(testcontainers::core::AccessMode::ReadOnly),
            )
            .with_env_var("POSTGRES_PASSWORD", "atcodersearch")
            .with_env_var("POSTGRES_USER", "atcodersearch")
            .with_env_var("POSTGRES_DB", "atcodersearch")
            .with_env_var("POSTGRES_HOST_AUTH_METHOD", "password")
            .with_env_var("TZ", "Asia/Tokyo");

        Ok(request)
    }

    pub async fn create_pool_from_container(
        container: &ContainerAsync<GenericImage>,
    ) -> anyhow::Result<Pool<Postgres>> {
        let host = container.get_host().await?;
        let port = container.get_host_port_ipv4(5432).await?;

        let url = format!(
            "postgres://atcodersearch:atcodersearch@{}:{}/atcodersearch",
            host, port
        );
        let pool = Pool::connect(&url).await?;

        Ok(pool)
    }

    pub fn create_engine_container() -> anyhow::Result<ContainerRequest<GenericImage>> {
        todo!();
    }

    pub async fn create_client_from_container(
        container: &ContainerAsync<GenericImage>,
    ) -> anyhow::Result<Client> {
        todo!();
    }
}
