pub mod atcodersearch {
    tonic::include_proto!("atcodersearch");
}

use anyhow::Context;
use atcodersearch::{SearchProblemByKeywordRequest, SearchProblemResult};
use sqlx::postgres::PgPoolOptions;
use tonic::{transport::Server, Request, Response, Status};

use atcodersearch::problem_service_server::{ProblemService, ProblemServiceServer};
use clap::Parser;
use meilisearch_sdk::{client::Client, indexes::Index};
use sqlx::{Pool, Postgres};
use std::net::{IpAddr, Ipv4Addr, SocketAddr};
use std::sync::Arc;
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

#[derive(Debug)]
pub struct ProblemSearcher {
    pool: Pool<Postgres>,
    index: Arc<Index>,
}

impl ProblemSearcher {
    pub fn new(pool: Pool<Postgres>, index: Index) -> Self {
        Self {
            pool,
            index: Arc::new(index),
        }
    }
}

#[tonic::async_trait]
impl ProblemService for ProblemSearcher {
    async fn search_by_keyword(
        &self,
        request: Request<SearchProblemByKeywordRequest>,
    ) -> Result<Response<SearchProblemResult>, Status> {
        let req = request.into_inner();

        let res = SearchProblemResult {
            time: 0,
            total: 0,
            index: 0,
            count: 0,
            pages: 0,
            items: vec![],
            facet: None,
        };

        Ok(Response::new(res))
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
