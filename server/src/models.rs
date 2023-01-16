use axum::async_trait;
use axum::extract::{FromRequest, FromRequestParts};
use axum::http::StatusCode;
use axum::{BoxError, Form, Json};
use chrono::{DateTime, Local};
use http::request::Parts;
use http_body::Body;
use hyper::Request;
use serde::de::{DeserializeOwned, Error, Unexpected};
use serde::{Deserialize, Deserializer, Serialize, Serializer};
use solr_client::query::{QueryBuilder, StandardQueryBuilder, StandardQueryOperand};
use validator::Validate;

#[derive(Debug, Serialize, Deserialize, Clone, Validate)]
pub struct SearchParams {
    pub q: Option<String>,
    #[validate(range(max = 1000))]
    pub o: Option<u32>,
    #[validate(range(max = 200))]
    pub l: Option<u32>,
}

impl SearchParams {
    pub fn as_qs(&self) -> Vec<(String, String)> {
        let mut builder = StandardQueryBuilder::new();
        if let Some(q) = &self.q {
            let op = StandardQueryOperand::new("text_ja", q);
            builder = builder.q(&op);
        };
        if let Some(l) = self.l {
            builder = builder.rows(l);
        };
        if let Some(o) = self.o {
            builder = builder.start(o);
        }

        builder.build()
    }
}

pub struct ValidatedSearchJson<T>(pub T);

#[async_trait]
impl<T, S, B> FromRequest<S, B> for ValidatedSearchJson<T>
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

        Ok(ValidatedSearchJson(value))
    }
}

pub struct ValidatedSearchForm<T>(pub T);

#[async_trait]
impl<T, S, B> FromRequest<S, B> for ValidatedSearchForm<T>
where
    T: Validate,
    Form<T>: FromRequest<S, B>,
    B: Send + 'static,
    S: Send + Sync,
{
    type Rejection = (StatusCode, String);

    async fn from_request(req: Request<B>, state: &S) -> Result<Self, Self::Rejection> {
        let Form(value) = Form::<T>::from_request(req, state)
            .await
            .map_err(|_| (StatusCode::BAD_REQUEST, format!("Invalid format forms")))?;

        value.validate().map_err(|rejection| {
            (
                StatusCode::BAD_REQUEST,
                format!("Validation error: [{}]", rejection).replace('\n', ", "),
            )
        })?;

        Ok(ValidatedSearchForm(value))
    }
}

pub struct ValidatedSearchQueryParams<T>(pub T);

#[async_trait]
impl<T, S> FromRequestParts<S> for ValidatedSearchQueryParams<T>
where
    T: DeserializeOwned + Validate,
    S: Send + Sync,
{
    type Rejection = (StatusCode, String);

    async fn from_request_parts(parts: &mut Parts, _state: &S) -> Result<Self, Self::Rejection> {
        let query = parts.uri.query().unwrap_or_default();
        let value: T = serde_urlencoded::from_str(query).map_err(|rejection| {
            (
                StatusCode::BAD_REQUEST,
                format!("Invalid format query string: [{}]", rejection),
            )
        })?;

        value.validate().map_err(|rejection| {
            (
                StatusCode::BAD_REQUEST,
                format!("Validation error: [{}]", rejection).replace('\n', ", "),
            )
        })?;

        Ok(ValidatedSearchQueryParams(value))
    }
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
    pub facet: Option<FacetResult>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct SearchResultBody {
    pub docs: Vec<Document>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct FacetResult {}

#[derive(Debug, Serialize, Deserialize)]
pub struct Document {
    pub problem_id: String,
    pub problem_title: String,
    pub problem_url: String,
    pub contest_id: String,
    pub contest_title: String,
    pub contest_url: String,
    pub difficulty: i32,
    #[serde(serialize_with = "serialize", deserialize_with = "deserialize")]
    pub start_at: DateTime<Local>,
    pub duration: i64,
    pub rate_change: String,
    pub category: String,
}

fn serialize<S>(value: &DateTime<Local>, serializer: S) -> Result<S::Ok, S::Error>
where
    S: Serializer,
{
    serializer.serialize_str(&value.to_rfc3339())
}

fn deserialize<'de, D>(deserializer: D) -> Result<DateTime<Local>, D::Error>
where
    D: Deserializer<'de>,
{
    let value = String::deserialize(deserializer)?;
    if let Ok(timestamp) = DateTime::parse_from_rfc3339(&value) {
        return Ok(timestamp.with_timezone(&Local));
    } else {
        return Err(Error::invalid_value(
            Unexpected::Str(&value),
            &"Invalid timestamp string",
        ));
    }
}
