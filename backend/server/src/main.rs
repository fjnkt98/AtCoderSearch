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
    let _ = Server::bind(&addr).serve(app.into_make_service()).await;
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
        http::{Method, Request},
    };
    use reqwest::Url;
    use rstest::*;
    use tower::ServiceExt;

    /// テスト用のヘルパーメソッド
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
    #[rstest]
    #[case("")]
    #[case("高橋")]
    #[tokio::test]
    async fn test_keyword(#[case] keyword: &str) {
        let uri =
            Url::parse_with_params("https://localhost:8000/api/search", &[("keyword", keyword)])
                .unwrap();

        let req = Request::builder()
            .uri(format!("/api/search?{}", uri.query().unwrap()))
            .method(Method::GET)
            .body(Body::empty())
            .unwrap();
        let res = create_app().await.oneshot(req).await.unwrap();

        assert_eq!(res.status(), 200);
    }

    #[ignore]
    #[rstest]
    #[case("")]
    #[case("1")]
    #[case("20")]
    #[case("200")]
    #[tokio::test]
    async fn test_limit(#[case] limit: &str) {
        let uri = Url::parse_with_params("https://localhost:8000/api/search", &[("limit", limit)])
            .unwrap();

        let req = Request::builder()
            .uri(format!("/api/search?{}", uri.query().unwrap()))
            .method(Method::GET)
            .body(Body::empty())
            .unwrap();
        let res = create_app().await.oneshot(req).await.unwrap();

        assert_eq!(res.status(), 200);
    }

    #[ignore]
    #[rstest]
    #[case("")]
    #[case("1")]
    #[case("2")]
    #[tokio::test]
    async fn test_page(#[case] page: &str) {
        let uri =
            Url::parse_with_params("https://localhost:8000/api/search", &[("page", page)]).unwrap();

        let req = Request::builder()
            .uri(format!("/api/search?{}", uri.query().unwrap()))
            .method(Method::GET)
            .body(Body::empty())
            .unwrap();
        let res = create_app().await.oneshot(req).await.unwrap();

        assert_eq!(res.status(), 200);
    }

    #[ignore]
    #[rstest]
    #[case("")]
    #[case("ABC")]
    #[case("Other Contests")]
    #[case("ABC-Like")]
    #[tokio::test]
    async fn test_category(#[case] category: &str) {
        let uri = Url::parse_with_params(
            "https://localhost:8000/api/search",
            &[("filter.category[]", category)],
        )
        .unwrap();

        let req = Request::builder()
            .uri(format!("/api/search?{}", uri.query().unwrap()))
            .method(Method::GET)
            .body(Body::empty())
            .unwrap();
        let res = create_app().await.oneshot(req).await.unwrap();

        assert_eq!(res.status(), 200);
    }

    #[ignore]
    #[rstest]
    #[case("")]
    #[case("-score")]
    #[case("difficulty")]
    #[case("-difficulty")]
    #[case("start_at")]
    #[case("-start_at")]
    #[tokio::test]
    async fn test_sort(#[case] sort: &str) {
        let uri =
            Url::parse_with_params("https://localhost:8000/api/search", &[("sort", sort)]).unwrap();

        let req = Request::builder()
            .uri(format!("/api/search?{}", uri.query().unwrap()))
            .method(Method::GET)
            .body(Body::empty())
            .unwrap();
        let res = create_app().await.oneshot(req).await.unwrap();

        assert_eq!(res.status(), 200);
    }

    #[ignore]
    #[rstest(
        difficulty_from => ["", "400"],
        difficulty_to => ["", "1200"],
    )]
    #[tokio::test]
    async fn test_difficulty(difficulty_from: &str, difficulty_to: &str) {
        let uri = Url::parse_with_params(
            "https://localhost:8000/api/search",
            &[
                ("filter.difficulty.from", difficulty_from),
                ("filter.difficulty.to", difficulty_to),
            ],
        )
        .unwrap();

        let req = Request::builder()
            .uri(format!("/api/search?{}", uri.query().unwrap()))
            .method(Method::GET)
            .body(Body::empty())
            .unwrap();
        let res = create_app().await.oneshot(req).await.unwrap();

        assert_eq!(res.status(), 200);
    }

    #[ignore]
    #[rstest]
    #[case("pZPcAlHuRZYkroZaOotJqmEsChRMDsQeycDQdqwjCWxQkClrzcvZdeggFaNSslXlktjUozDklVBEJiHYYzpkNhwmGBGbWieIdtIsiANyyCyicKejnFlgSkiWIbWkrfhFKkVqNaLceqZpdFBWREDiIeWRhLuloiXbanQHYdxSqvrYizuJMenLKMmPutwqNRSSNlijxUYfa")]
    #[tokio::test]
    async fn should_400_when_keyword_length_is_greater_than_200(#[case] keyword: &str) {
        let uri =
            Url::parse_with_params("https://localhost:8000/api/search", &[("keyword", keyword)])
                .unwrap();

        let req = Request::builder()
            .uri(format!("/api/search?{}", uri.query().unwrap()))
            .method(Method::GET)
            .body(Body::empty())
            .unwrap();
        let res = create_app().await.oneshot(req).await.unwrap();

        assert_eq!(res.status(), 400);
    }

    #[ignore]
    #[rstest]
    #[case(-1)]
    #[case(0)]
    #[case(201)]
    #[tokio::test]
    async fn should_400_when_limit_param_is_out_of_limitation(#[case] limit: i32) {
        let uri = Url::parse_with_params(
            "https://localhost:8000/api/search",
            &[("limit", &limit.to_string())],
        )
        .unwrap();

        let req = Request::builder()
            .uri(format!("/api/search?{}", uri.query().unwrap()))
            .method(Method::GET)
            .body(Body::empty())
            .unwrap();
        let res = create_app().await.oneshot(req).await.unwrap();

        assert_eq!(res.status(), 400);
    }

    #[ignore]
    #[rstest]
    #[case(-1)]
    #[case(0)]
    #[tokio::test]
    async fn should_400_when_page_param_is_out_of_limitation(#[case] page: i32) {
        let uri = Url::parse_with_params(
            "https://localhost:8000/api/search",
            &[("page", &page.to_string())],
        )
        .unwrap();

        let req = Request::builder()
            .uri(format!("/api/search?{}", uri.query().unwrap()))
            .method(Method::GET)
            .body(Body::empty())
            .unwrap();
        let res = create_app().await.oneshot(req).await.unwrap();

        assert_eq!(res.status(), 400);
    }
}
