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
use tracing_subscriber::fmt;

#[tokio::main]
async fn main() {
    dotenv().ok();

    let log_level = env::var("RUST_LOG").unwrap_or(String::from("debug"));
    env::set_var("RUST_LOG", log_level);
    fmt::init();

    let solr = SolrClient::new("http://localhost", 8983).unwrap();
    let core = solr.core("atcoder").await.unwrap();
    let app = create_router(core);

    let addr = SocketAddr::from(([0, 0, 0, 0], 8000));
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
