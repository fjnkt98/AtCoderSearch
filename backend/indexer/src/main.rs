mod models;
mod utils;

use std::ffi::OsString;
use std::str::FromStr;

use crate::utils::crawlers::{ContestCrawler, ProblemCrawler};
use crate::utils::generator::DocumentGenerator;
use crate::utils::uploader::DocumentUploader;
use anyhow::{bail, Context, Result};
use clap::{Args, Parser, Subcommand};
use dotenvy::dotenv;
use solrust::client::solr::SolrClient;
use sqlx::postgres::Postgres;
use sqlx::Pool;
use std::env;
use std::path::PathBuf;
use tracing_subscriber::{
    filter::{EnvFilter, LevelFilter},
    fmt,
};

#[derive(Debug, Parser)]
#[command(name = "indexer")]
#[command(
    about = "All you need to indexing",
    long_about = "Crawl and problems, generate document json, and post json into Solr core."
)]
struct Cli {
    #[command(subcommand)]
    command: Commands,
}
#[derive(Debug, Subcommand)]
enum Commands {
    Crawl(CrawlArgs),
    Generate(GenerateArgs),
    Post(PostArgs),
}

#[derive(Debug, Args)]
struct CrawlArgs {
    #[arg(long)]
    all: bool,
}

#[derive(Debug, Args)]
struct GenerateArgs {
    path: Option<OsString>,
}

#[derive(Debug, Args)]
struct PostArgs {
    path: Option<OsString>,
    #[arg(short, long)]
    optimize: bool,
}

#[tokio::main(flavor = "multi_thread", worker_threads = 10)]
async fn main() -> Result<()> {
    dotenv().ok();
    let args = Cli::parse();

    let log_level = env::var("RUST_LOG").unwrap_or(String::from("info"));
    let filter = EnvFilter::builder()
        .with_default_directive(
            LevelFilter::from_str(&log_level)
                .expect("Cannot parse specified log level!")
                .into(),
        )
        .from_env_lossy();

    let format = fmt::format()
        .with_level(true)
        .with_target(true)
        .with_ansi(false)
        .with_thread_ids(true);

    let subscriber = tracing_subscriber::fmt()
        .with_env_filter(filter)
        .event_format(format)
        .finish();
    tracing::subscriber::set_global_default(subscriber)
        .with_context(|| "Failed to set subscriber.")?;

    let database_url: String = env::var("DATABASE_URL").with_context(|| {
        let message = "DATABASE_URL must be configured.";
        tracing::error!(message);
        message
    })?;

    let pool: Pool<Postgres> = sqlx::postgres::PgPoolOptions::new()
        .max_connections(5)
        .connect(&database_url)
        .await
        .with_context(|| {
            let message = "Failed to create database connection pool.";
            tracing::error!(message);
            message
        })?;

    match args.command {
        Commands::Crawl(args) => {
            let crawler = ContestCrawler::new(&pool);
            crawler.run().await.with_context(|| {
                let message = "Failed to crawl and save contest information.";
                tracing::error!(message);
                message
            })?;

            let crawler = ProblemCrawler::new(&pool);
            crawler.run(args.all).await.with_context(|| {
                let message = "Failed to crawl and save problem information [{}]";
                tracing::error!(message);
                message
            })?;
        }
        Commands::Generate(args) => {
            let savedir = args.path.and_then(|path| Some(PathBuf::from(path))).or_else(|| {
                let path = env::var("DOCUMENT_SAVE_DIRECTORY").with_context(|| {
                    let message = "Documents save directory does not configured. Check and make sure DOCUMENT_SAVE_DIRECTORY environment variable.";
                    tracing::error!(message);
                    message
                }).ok()?;
                Some(PathBuf::from(path))
            }).with_context(|| "Failed to set document save dir.")?;
            tracing::info!("Documents will be saved into {:?}.", savedir);

            let generator = DocumentGenerator::new(&pool, &savedir);
            tracing::info!("Delete existing documents");
            generator.truncate().await.or_else(|e| {
                tracing::error!(
                    "Failed to delete existing json documents [{}]",
                    e.to_string()
                );
                bail!(e.to_string());
            })?;
            match generator.generate(1000).await {
                Ok(()) => {
                    tracing::info!("Successfully generate documents");
                }
                Err(e) => {
                    tracing::error!("Failed to generate documents");
                    bail!(e.to_string())
                }
            }
        }
        Commands::Post(args) => {
            let solr_host = env::var("SOLR_HOST").unwrap_or_else(|_| {
                tracing::info!("SOLR_HOST environment variable is not set. Default value `http://localhost` will be used.");
                String::from("http://localhost")
            });
            let solr_port = match env::var("SOLR_PORT") {
                Ok(v) => match v.parse::<u32>() {
                    Ok(port) => port,
                    Err(e) => {
                        bail!(e.to_string());
                    }
                },
                Err(_) => {
                    tracing::info!("SOLR_PORT environment variable is not set. Default value `8983` will be used.");
                    8983u32
                }
            };

            let core_name = env::var("CORE_NAME").with_context(|| {
                let message = "CORE_NAME must be configured";
                tracing::error!(message);
                message
            })?;

            let solr = SolrClient::new(&solr_host, solr_port).with_context(|| {
                let message = "Failed to create Solr client.";
                tracing::error!(message);
                message
            })?;
            let core = solr.core(&core_name).await.with_context(|| {
                let message = "Failed to create Solr core client";
                tracing::error!(message);
                message
            })?;

            let savedir = args.path.and_then(|path| Some(PathBuf::from(path))).or_else(|| {
                let path = env::var("DOCUMENT_SAVE_DIRECTORY").with_context(|| {
                    let message = "Documents save directory does not configured. Check and make sure DOCUMENT_SAVE_DIRECTORY environment variable.";
                    tracing::error!(message);
                    message
                }).ok()?;
                Some(PathBuf::from(path))
            }).with_context(|| "Failed to set document save dir.")?;
            tracing::info!("Target documents are in {:?}", savedir);

            let uploader = DocumentUploader::new(&savedir, &core);
            tracing::info!("Start to post documents");
            uploader.upload(args.optimize).await.with_context(|| {
                let message = "Failed to post documents";
                tracing::error!(message);
                message
            })?;
            tracing::info!("Successfully post documents");
        }
    }

    Ok(())
}
