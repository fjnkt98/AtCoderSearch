mod handlers;
mod models;

use crate::handlers::{search_with_json, search_with_qs};
use axum::extract::Extension;
use axum::routing::get;
use axum::{Router, Server};
use dotenvy::dotenv;
use solrust::client::core::SolrCore;
use solrust::client::solr::SolrClient;
use std::env;
use std::net::SocketAddr;
use std::sync::Arc;
use tracing_appender::rolling::{RollingFileAppender, Rotation};
use tracing_subscriber::prelude::__tracing_subscriber_SubscriberExt;
use tracing_subscriber::Layer;
use tracing_subscriber::{
    filter::{EnvFilter, LevelFilter},
    fmt, Registry,
};

#[tokio::main]
async fn main() {
    dotenv().ok();

    let log_level = env::var("RUST_LOG").unwrap_or(String::from("info"));
    env::set_var("RUST_LOG", "info");
    let create_filter = || {
        EnvFilter::builder()
            .with_default_directive(LevelFilter::INFO.into())
            .from_env_lossy()
            .add_directive(format!("solrust={}", log_level).parse().unwrap())
            .add_directive(format!("server={}", log_level).parse().unwrap())
    };
    let layer1 = fmt::Layer::new().with_filter(create_filter());
    let (file, _guard) = tracing_appender::non_blocking(RollingFileAppender::new(
        Rotation::DAILY,
        "/var/tmp/atcoder/log",
        "query.log",
    ));
    let layer2 = fmt::Layer::new()
        .with_writer(file)
        .with_filter(create_filter());
    let subscriber = Registry::default().with(layer1).with(layer2);
    tracing::subscriber::set_global_default(subscriber).expect("Failed to set subscriber.");

    let solr = SolrClient::new("http://localhost", 8983).unwrap();
    let core = solr.core("atcoder").await.unwrap();
    let app = create_router(core);

    let addr = SocketAddr::from(([0, 0, 0, 0], 8000));
    tracing::info!("Server start");
    tracing::debug!("Server start at port 8000");
    Server::bind(&addr)
        .serve(app.into_make_service())
        .await
        .unwrap();
}

fn create_router(core: SolrCore) -> Router {
    Router::new()
        .route("/api/search", get(search_with_qs).post(search_with_json))
        .layer(Extension(Arc::new(core)))
}
