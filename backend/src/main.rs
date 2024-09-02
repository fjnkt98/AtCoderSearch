mod atcoder;
mod cmd;
mod crawl;
mod history;

use clap::Parser;
use cmd::App;
use std::{env, str::FromStr};
use tracing::level_filters::LevelFilter;
use tracing_subscriber::{
    fmt::{self, time::OffsetTime},
    EnvFilter,
};

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

    let app = App::parse();
    app.run()?;

    Ok(())
}

#[cfg(test)]
pub mod testutil {
    use sqlx::{Pool, Postgres};
    use std::env;
    use testcontainers::{
        core::{IntoContainerPort, Mount, WaitFor},
        ContainerAsync, ContainerRequest, GenericImage, ImageExt,
    };

    pub fn create_container() -> anyhow::Result<ContainerRequest<GenericImage>> {
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
}
