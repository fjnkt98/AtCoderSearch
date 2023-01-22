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
use solr_client::query::{
    Aggregation, QueryBuilder, QueryExpression, QueryOperand, StandardQueryBuilder,
    StandardQueryOperand,
};
use solr_client::query::{RangeQueryOperand, SortOrderBuilder};
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
            let op = QueryExpression::prod(
                q.split_whitespace()
                    .map(|word| {
                        let text_ja =
                            QueryOperand::from(StandardQueryOperand::new("text_ja", word));
                        let text_en =
                            QueryOperand::from(StandardQueryOperand::new("text_en", word));
                        text_ja + text_en
                    })
                    .collect::<Vec<QueryExpression>>(),
            );
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
                SortOrderBuilder::new().desc(&s[1..])
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
                    .collect::<Vec<QueryOperand>>(),
            );
            builder = builder.fq(&fq);
        }

        if self.difficulty_from.is_some() || self.difficulty_to.is_some() {
            let mut range = RangeQueryOperand::new("difficulty");
            if let Some(from) = self.difficulty_from {
                range = range.ge(from.to_string());
            }
            if let Some(to) = self.difficulty_to {
                range = range.lt(to.to_string());
            }
            let fq = QueryOperand::from(range);
            builder = builder.fq(&fq);
        }

        builder = builder.op("AND");

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
    pub time: u32,
    pub message: Option<String>,
    pub total: u32,
    pub offset: u32,
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

#[cfg(test)]
mod test {
    use super::*;

    #[test]
    fn test_single_q() {
        let params = SearchParams {
            q: Some(String::from("hoge")),
            p: None,
            o: None,
            category: None,
            difficulty_from: None,
            difficulty_to: None,
            s: None,
        };

        let qs = params.as_qs();

        assert_eq!(
            vec![(
                String::from("q"),
                String::from("(text_ja:hoge OR text_en:hoge OR text_phrase:hoge)")
            )],
            qs
        )
    }

    #[test]
    fn test_multiple_q() {
        let params = SearchParams {
            q: Some(String::from("hoge moge")),
            p: None,
            o: None,
            category: None,
            difficulty_from: None,
            difficulty_to: None,
            s: None,
        };

        let qs = params.as_qs();

        assert_eq!(
            vec![(
                String::from("q"),
                String::from("(text_ja:hoge OR text_en:hoge OR text_phrase:hoge) AND (text_ja:moge OR text_en:moge OR text_phrase:moge)")
            )],
            qs
        )
    }
}
