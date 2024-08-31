mod atcoder;
mod cmd;
mod crawl;
mod history;

use clap::Parser;
use cmd::App;

#[tokio::main]
async fn main() {
    let app = App::parse();

    if let Err(err) = app.run().await {
        println!("command failed: {:#}", err)
    }
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
