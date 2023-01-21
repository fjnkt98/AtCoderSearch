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
use solr_client::query::SortOrderBuilder;
use solr_client::query::{
    QueryBuilder, QueryExpression, QueryOperand, StandardQueryBuilder, StandardQueryOperand,
};
use validator::{Validate, ValidationError};

#[derive(Debug, Serialize, Deserialize, Clone, Validate)]
pub struct SearchParams {
    pub q: Option<String>,
    #[validate(range(max = 200))]
    pub p: Option<u32>,
    pub o: Option<u32>,
    pub category: Option<Vec<String>>,
    #[serde(alias = "difficulty.from")]
    pub difficulty_from: Option<u32>,
    #[serde(alias = "difficulty.to")]
    pub difficulty_to: Option<u32>,
    #[validate(custom = "validate_sort_option")]
    pub s: Option<String>,
}

fn validate_sort_option(value: &str) -> Result<(), ValidationError> {
    match value {
        "start_at" | "-start_at" | "difficulty" | "-difficulty" | "score" | "-score" => {}
        _ => return Err(ValidationError::new("Invalid sort field option.")),
    }
    Ok(())
}

impl SearchParams {
    pub fn as_qs(&self) -> Vec<(String, String)> {
        let mut builder = StandardQueryBuilder::new();
        if let Some(q) = &self.q {
            let op = QueryOperand(format!("text_ja: {}", q));
            builder = builder.q(&op);
        };
        if let Some(p) = self.p {
            builder = builder.rows(p);
        };
        if let Some(o) = self.o {
            builder = builder.start(o);
        }
        if let Some(s) = &self.s {
            let sort = if s.starts_with("-") {
                SortOrderBuilder::new().desc(s)
            } else {
                SortOrderBuilder::new().asc(s)
            };
            builder = builder.sort(&sort);
        }
        if let Some(category) = &self.category {
            let fq = QueryExpression::sum(
                category
                    .iter()
                    .map(|c| QueryOperand::from(StandardQueryOperand::new("category", c)))
                    .collect(),
            );
            builder = builder.fq(&fq);
        }

        // if let Some(difficulty_from) = self.difficulty_from {

        // }

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
        let value: T = serde_qs::from_str(query).map_err(|rejection| {
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
