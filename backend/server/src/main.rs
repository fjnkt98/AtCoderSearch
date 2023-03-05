mod handlers;
mod models;

use crate::handlers::{search_with_json, search_with_qs};
use axum::extract::Extension;
use axum::routing::get;
use axum::{Router, Server};
use dotenvy::dotenv;
use hyper::header::CONTENT_TYPE;
use solrust::client::core::SolrCore;
use solrust::client::solr::SolrClient;
use std::env;
use std::net::SocketAddr;
use std::sync::Arc;
use tower_http::cors::{AllowOrigin, Any, CorsLayer};
use tracing_appender::rolling::{RollingFileAppender, Rotation};
use tracing_subscriber::layer::SubscriberExt;
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
            .add_directive("querylog=off".parse().unwrap())
    };

    let log_dir = env::var("LOG_DIRECTORY").unwrap_or(String::from("/var/tmp/atcoder/log"));

    // システムログ(コンソールへ出力)
    let layer1 = fmt::Layer::new().with_filter(create_filter());

    // システムログ(ファイルへ出力)
    let (file, _guard) = tracing_appender::non_blocking(RollingFileAppender::new(
        Rotation::DAILY,
        log_dir.clone(),
        "atcoder.log",
    ));
    let layer2 = fmt::Layer::new()
        .with_writer(file)
        .with_filter(create_filter());
    // クエリログ(ファイルへ出力)
    let (file, _guard) = tracing_appender::non_blocking(RollingFileAppender::new(
        Rotation::DAILY,
        log_dir.clone(),
        "query.log",
    ));
    let layer3 = fmt::Layer::new().with_writer(file).with_filter(
        EnvFilter::builder()
            .with_default_directive(LevelFilter::INFO.into())
            .from_env_lossy()
            .add_directive("querylog=info".parse().unwrap())
            .add_directive("server=off".parse().unwrap()),
    );

    let subscriber = Registry::default().with(layer1).with(layer2).with(layer3);
    tracing::subscriber::set_global_default(subscriber).expect("Failed to set subscriber.");

    let solr_host = env::var("SOLR_HOST").unwrap_or(String::from("http://localhost"));
    let solr_port: u32 = (env::var("SOLR_PORT").unwrap_or(String::from("8983")))
        .parse()
        .unwrap_or(8983);
    let solr = SolrClient::new(&solr_host, solr_port).expect("Failed to create SolrClient instance. check that your Solr instance is running properly and that the SOLR_HOST and SOLR_PORT environment variable is set correctly.");

    let core_name = env::var("CORE_NAME").unwrap_or(String::from("atcoder"));
    let core = solr
        .core(&core_name)
        .await
        .expect("Failed to create SolrCore instance. Check that your Solr instance is running properly, that the CORE_NAME environment variable is set correctly, and that the specified core exists.");
    let app = create_router(core);

    let port: u16 = (env::var("API_SERVER_LISTEN_PORT").unwrap_or(String::from("8000")))
        .parse()
        .unwrap();
    let addr = SocketAddr::from(([0, 0, 0, 0], port));
    tracing::debug!("Server start at port {}", port);
    Server::bind(&addr)
        .serve(app.into_make_service())
        .await
        .unwrap();
}

fn create_router(core: SolrCore) -> Router {
    Router::new()
        .route("/api/search", get(search_with_qs).post(search_with_json))
        .layer(Extension(Arc::new(core)))
        .layer(
            CorsLayer::new()
                .allow_origin(AllowOrigin::exact("http://localhost:8080".parse().unwrap()))
                .allow_methods(Any)
                .allow_headers(vec![CONTENT_TYPE]),
        )
}

#[cfg(test)]
mod test {
    use super::*;
    use axum::{
        body::Body,
        http::{header, Method, Request},
    };
    use tower::ServiceExt;

    async fn create_app() -> Router {
        let solr_host = env::var("SOLR_HOST").unwrap_or(String::from("http://localhost"));
        let solr_port: u32 = (env::var("SOLR_PORT").unwrap_or(String::from("8983")))
            .parse()
            .unwrap_or(8983);
        let solr = SolrClient::new(&solr_host, solr_port).expect("Failed to create SolrClient instance. check that your Solr instance is running properly and that the SOLR_HOST and SOLR_PORT environment variable is set correctly.");

        let core_name = env::var("CORE_NAME").unwrap_or(String::from("atcoder"));
        let core = solr
        .core(&core_name)
        .await
        .expect("Failed to create SolrCore instance. Check that your Solr instance is running properly, that the CORE_NAME environment variable is set correctly, and that the specified core exists.");

        create_router(core)
    }

    #[ignore]
    #[tokio::test]
    async fn get_default() {
        let req = Request::builder()
            .uri("/api/search")
            .method(Method::GET)
            .body(Body::from(""))
            .unwrap();
        let res = create_app().await.oneshot(req).await.unwrap();

        assert_eq!(res.status(), 200);
    }

    #[ignore]
    #[tokio::test]
    async fn post_default() {
        let req = Request::builder()
            .uri("/api/search")
            .method(Method::POST)
            .header(header::CONTENT_TYPE, mime::APPLICATION_JSON.as_ref())
            .body(Body::from(r#"{}"#))
            .unwrap();
        let res = create_app().await.oneshot(req).await.unwrap();

        assert_eq!(res.status(), 200);
    }

    #[ignore]
    #[tokio::test]
    async fn post_with_q() {
        let req = Request::builder()
            .uri("/api/search")
            .method(Method::POST)
            .header(header::CONTENT_TYPE, mime::APPLICATION_JSON.as_ref())
            .body(Body::from(r#"{"q": "高橋"}"#))
            .unwrap();
        let res = create_app().await.oneshot(req).await.unwrap();

        assert_eq!(res.status(), 200);
    }

    #[ignore]
    #[tokio::test]
    async fn post_with_p() {
        let req = Request::builder()
            .uri("/api/search")
            .method(Method::POST)
            .header(header::CONTENT_TYPE, mime::APPLICATION_JSON.as_ref())
            .body(Body::from(r#"{"p": 200}"#))
            .unwrap();
        let res = create_app().await.oneshot(req).await.unwrap();

        assert_eq!(res.status(), 200);
    }

    #[ignore]
    #[tokio::test]
    async fn should_return_error_when_p_is_greater_than_200() {
        let req = Request::builder()
            .uri("/api/search")
            .method(Method::POST)
            .header(header::CONTENT_TYPE, mime::APPLICATION_JSON.as_ref())
            .body(Body::from(r#"{"p": 201}"#))
            .unwrap();
        let res = create_app().await.oneshot(req).await.unwrap();

        assert_eq!(res.status(), 400);
    }

    #[ignore]
    #[tokio::test]
    async fn should_return_error_when_p_is_minus() {
        let req = Request::builder()
            .uri("/api/search")
            .method(Method::POST)
            .header(header::CONTENT_TYPE, mime::APPLICATION_JSON.as_ref())
            .body(Body::from(r#"{"p": -1}"#))
            .unwrap();
        let res = create_app().await.oneshot(req).await.unwrap();

        assert_eq!(res.status(), 400);
    }
}
