mod handlers;
mod models;

use crate::handlers::{healthcheck, search_with_json, search_with_qs};
use axum::{extract::Extension, routing::get, Router, Server};
use dotenvy::dotenv;
use hyper::header::CONTENT_TYPE;
use solrust::client::{core::SolrCore, solr::SolrClient};
use std::net::SocketAddr;
use std::sync::Arc;
use std::{env, str::FromStr};
use tower_http::cors::{AllowOrigin, Any, CorsLayer};
use tracing_subscriber::{
    filter::{EnvFilter, LevelFilter},
    fmt,
};

#[tokio::main]
async fn main() {
    dotenv().ok();

    let log_level = env::var("RUST_LOG").unwrap_or(String::from("info"));
    let filter = EnvFilter::builder()
        .with_default_directive(LevelFilter::from_str(&log_level).unwrap().into())
        .from_env_lossy()
        .add_directive("hyper=info".parse().unwrap());

    let format = fmt::format()
        .with_level(true)
        .with_target(true)
        .with_ansi(false);

    let subscriber = tracing_subscriber::fmt()
        .with_env_filter(filter)
        .event_format(format)
        .finish();
    tracing::subscriber::set_global_default(subscriber).expect("Failed to set subscriber.");

    let solr_host = env::var("SOLR_HOST").unwrap_or_else(|_| {
        tracing::info!("SOLR_HOST environment variable is not set. Default value `http://localhost` will be used.");
        String::from("http://localhost")}
    );
    let solr_port = match env::var("SOLR_PORT") {
        Ok(v) => match v.parse::<u32>() {
            Ok(port) => port,
            Err(e) => {
                let message = format!(
                    "Failed to parse SOLR_PORT value into u32. [{}]",
                    e.to_string()
                );
                tracing::error!(message);
                panic!("{}", message);
            }
        },
        Err(_) => {
            tracing::info!(
                "SOLR_PORT environment variable is not set. Default value `8983` will be used."
            );
            8983u32
        }
    };
    tracing::info!("Connect to Solr instance at {}:{}", solr_host, solr_port);
    let solr = SolrClient::new(&solr_host, solr_port).unwrap_or_else(|e| {
        let message = format!("{} Failed to create SolrClient instance. check that your Solr instance is running properly and that the SOLR_HOST and SOLR_PORT environment variable is set correctly.", e.to_string());
        tracing::error!(message);
        panic!("{}", message);
    });

    let core_name = env::var("CORE_NAME").unwrap_or_else(|_| {
        let message = "SOLR_HOST environment variable must be set.";
        tracing::info!(message);
        panic!("{}", message);
    });
    tracing::info!("Connect to Solr core {}", core_name);
    let core = solr
        .core(&core_name)
        .await
        .unwrap_or_else(|e| {
            let message = format!("[{}] Failed to create SolrCore instance. Check that your Solr instance is running properly, that the CORE_NAME environment variable is set correctly, and that the specified core exists.", e.to_string());
            tracing::error!(message);
            panic!("{}", message);
        });

    let app = create_router(core);

    let port: u16 = match env::var("API_SERVER_LISTEN_PORT") {
        Ok(v) => match v.parse::<u16>() {
            Ok(port) => port,
            Err(e) => {
                let message = format!(
                    "Failed to parse API_SERVER_LISTEN_PORT into u16. [{}]",
                    e.to_string()
                );
                tracing::error!(message);
                panic!("{}", message);
            }
        },
        Err(_) => {
            tracing::info!(
                "API_SERVER_LISTEN_PORT environment variable is not set. Default value `8000` will be used."
            );
            8000u16
        }
    };

    let addr = SocketAddr::from(([0, 0, 0, 0], port));
    tracing::info!("Server start at port {}", port);
    Server::bind(&addr)
        .serve(app.into_make_service())
        .await
        .expect("Failed to bind server.");
}

fn create_router(core: SolrCore) -> Router {
    let origin = env::var("FRONTEND_ORIGIN_URL").unwrap_or(String::from("http://localhost:8080"));
    Router::new()
        .route("/api/search", get(search_with_qs).post(search_with_json))
        .route("/api/healthcheck", get(healthcheck))
        .layer(Extension(Arc::new(core)))
        .layer(
            CorsLayer::new()
                .allow_origin(AllowOrigin::exact(origin.parse().unwrap()))
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
    #[case("ABC,ARC")]
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
        difficulty_from => ["400"],
        difficulty_to => ["1200"],
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

    #[ignore]
    #[rstest]
    #[case("facet=category")]
    #[case("facet=difficulty")]
    #[case("facet=category,difficulty")]
    #[tokio::test]
    async fn specify_facet_count_fields(#[case] params: &str) {
        let req = Request::builder()
            .uri(format!("/api/search?{}", params))
            .method(Method::GET)
            .body(Body::empty())
            .unwrap();
        let res = create_app().await.oneshot(req).await.unwrap();

        assert_eq!(res.status(), 200);
    }

    #[ignore]
    #[tokio::test]
    async fn test_healthcheck_api() {
        let req = Request::builder()
            .uri("/api/healthcheck")
            .method(Method::GET)
            .body(Body::empty())
            .unwrap();
        let res = create_app().await.oneshot(req).await.unwrap();
        assert_eq!(res.status(), 200);
    }
}
