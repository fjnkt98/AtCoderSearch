mod handlers;
mod models;

use crate::handlers::search;
use anyhow::Result;
use axum::extract::Extension;
use axum::http::StatusCode;
use axum::response::IntoResponse;
use axum::routing::get;
use axum::{Router, Server};
use dotenvy::dotenv;
use solr_client::client::SolrClient;
use solr_client::core::SolrCore;
use std::env;
use std::net::SocketAddr;
use std::sync::Arc;
use tracing_subscriber::fmt;

#[tokio::main]
async fn main() {
    dotenv().ok();

    let log_level = env::var("RUST_LOG").unwrap_or(String::from("debug"));
    env::set_var("RUST_LOG", log_level);
    fmt::init();

    let solr = SolrClient::new("http://localhost", 8983).unwrap();
    let core = solr.core("atcoder").await.unwrap();
    let app = Router::new()
        .route("/", get(index))
        .route("/api/search", get(search))
        .layer(Extension(Arc::new(core)));

    let addr = SocketAddr::from(([0, 0, 0, 0], 8000));
    Server::bind(&addr)
        .serve(app.into_make_service())
        .await
        .unwrap();
}

async fn index(Extension(core): Extension<Arc<SolrCore>>) -> Result<impl IntoResponse, StatusCode> {
    let status = core
        .status()
        .await
        .or(Err(StatusCode::INTERNAL_SERVER_ERROR))?;

    Ok((StatusCode::OK, status.name))
}
