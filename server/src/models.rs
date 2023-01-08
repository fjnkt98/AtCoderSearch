use axum::async_trait;
use axum::extract::FromRequest;
use axum::http::StatusCode;
use axum::{BoxError, Json};
use http_body::Body;
use hyper::Request;
use serde::de::DeserializeOwned;
use serde::{Deserialize, Serialize};
use validator::Validate;

pub struct ValidatedSearchParams<T>(pub T);

#[async_trait]
impl<T, S, B> FromRequest<S, B> for ValidatedSearchParams<T>
where
    T: DeserializeOwned + Validate,
    B: Body + Send + 'static,
    B::Data: Send,
    B::Error: Into<BoxError>,
    S: Send + Sync,
{
    type Rejection = (StatusCode, String);

    async fn from_request(req: Request<B>, state: &S) -> Result<Self, Self::Rejection> {
        let Json(value) = Json::<T>::from_request(req, state)
            .await
            .map_err(|rejection| {
                (
                    StatusCode::BAD_REQUEST,
                    format!("JSON parse error: [{}]", rejection),
                )
            })?;

        value.validate().map_err(|rejection| {
            (
                StatusCode::BAD_REQUEST,
                format!("Validation error: [{}]", rejection).replace('\n', ", "),
            )
        })?;

        Ok(ValidatedSearchParams(value))
    }
}

#[derive(Debug, Serialize, Deserialize, Clone, Validate)]
pub struct SearchParams {
    pub q: Option<String>,
    #[validate(range(max = 1000))]
    pub o: Option<u32>,
    #[validate(range(max = 200))]
    pub l: Option<u32>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct SearchResultResponse {
    pub stats: SearchResultStats,
    pub items: SearchResultBody,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct SearchResultStats {
    pub total: u32,
    pub start: u32,
    pub amount: u32,
    pub facet: FacetResult,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct SearchResultBody {
    docs: Vec<Document>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct FacetResult {}

#[derive(Debug, Serialize, Deserialize)]
pub struct Document {}
