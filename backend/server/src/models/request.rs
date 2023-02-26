use crate::models::{SearchResultBody, SearchResultResponse, SearchResultStats};
use axum::async_trait;
use axum::extract::{FromRequest, FromRequestParts};
use axum::http::StatusCode;
use axum::{BoxError, Form, Json};
use http::request::Parts;
use http_body::Body;
use hyper::Request;
use once_cell::sync::Lazy;
use regex::Regex;
use serde::de::DeserializeOwned;
use serde::{Deserialize, Serialize};
use solrust::querybuilder::{
    common::SolrCommonQueryBuilder,
    dismax::SolrDisMaxQueryBuilder,
    edismax::{EDisMaxQueryBuilder, SolrEDisMaxQueryBuilder},
    facet::{FieldFacetBuilder, RangeFacetBuilder, RangeFacetOtherOptions},
    q::{
        Aggregation, Operator, QueryExpression, QueryOperand, RangeQueryOperand,
        StandardQueryOperand,
    },
    sort::SortOrderBuilder,
};
use std::collections::HashMap;
use validator::{Validate, ValidationError};

// Solrの特殊文字をエスケープする正規表現
// 本当はsolrustの方で実装したかったけど間に合わなかったので自力実装
static RE: Lazy<Regex> = Lazy::new(|| {
    Regex::new(r#"(\+|\-|&&|\|\||!|\(|\)|\{|\}|\[|\]|\^|"|\~|\*|\?|:|/|AND|OR)"#).unwrap()
});

/// 検索APIのクエリパラメータ
/// クエリパラメータとJSONパラメータ兼用
#[derive(Debug, Serialize, Deserialize, Clone, Validate)]
pub struct SearchParams {
    pub q: Option<String>,
    #[validate(range(max = 200))]
    pub p: Option<u32>,
    pub o: Option<u32>,
    pub f: Option<FilteringParameters>,
    #[validate(custom = "validate_sort_option")]
    pub s: Option<String>,
}

/// fパラメータに指定できる値
#[derive(Debug, Serialize, Deserialize, Clone, Validate)]
pub struct FilteringParameters {
    pub category: Option<Vec<String>>,
    pub difficulty: Option<RangeFilteringParameter<u32>>,
}

/// 範囲フィルタリングに指定できる値
#[derive(Debug, Serialize, Deserialize, Clone, Validate)]
pub struct RangeFilteringParameter<T> {
    pub from: Option<T>,
    pub to: Option<T>,
}

/// ソートフィールドを限定するためのバリデーション関数
fn validate_sort_option(value: &str) -> Result<(), ValidationError> {
    match value {
        "start_at" | "-start_at" | "difficulty" | "-difficulty" | "score" | "-score" => {}
        _ => return Err(ValidationError::new("Invalid sort field option.")),
    }
    Ok(())
}

impl SearchParams {
    /// リクエストパラメータからSolrへのリクエストパラメータを生成するメソッド
    /// Solrへ送るパラメータはすべてここで生成する
    pub fn as_qs(&self) -> Vec<(String, String)> {
        let category_facet = FieldFacetBuilder::new("category").min_count(1);
        let difficulty_facet = RangeFacetBuilder::new(
            "difficulty",
            0.to_string(),
            2000.to_string(),
            400.to_string(),
        )
        .other(RangeFacetOtherOptions::All);

        let mut builder = EDisMaxQueryBuilder::new()
            .rows(self.p.unwrap_or(20))
            .start(self.o.unwrap_or(0))
            .qf("text_ja text_en text_1gram")
            .q_alt(&QueryOperand::from("*:*"))
            .op(Operator::AND)
            .sow(true)
            .facet(&category_facet)
            .facet(&difficulty_facet);

        if let Some(q) = &self.q {
            let q = RE.replace_all(q, r"\$0");
            if !q.is_empty() {
                builder = builder.q(String::from(q));
            }

        if let Some(s) = &self.s {
            let sort = if s.starts_with("-") {
                SortOrderBuilder::new().desc(&s[1..])
            } else {
                SortOrderBuilder::new().asc(s)
            };
            builder = builder.sort(&sort);
        }

        if let Some(f) = &self.f {
            if let Some(category) = &f.category {
                let fq = QueryExpression::sum(
                    category
                        .iter()
                        .map(|c| QueryOperand::from(StandardQueryOperand::new("category", c)))
                        .collect::<Vec<QueryOperand>>(),
                );
                builder = builder.fq(&fq);
            }

            if let Some(difficulty) = &f.difficulty {
                if difficulty.from.is_some() || difficulty.to.is_some() {
                    let mut range = RangeQueryOperand::new("difficulty");
                    if let Some(from) = difficulty.from {
                        range = range.ge(from.to_string());
                    }
                    if let Some(to) = difficulty.to {
                        range = range.lt(to.to_string());
                    }
                    let fq = QueryOperand::from(range);
                    builder = builder.fq(&fq);
                }
            }
        }

        builder.build()
    }
}

/// JSON形式で入力されたリクエストパラメータをバリデーションする構造体
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
    type Rejection = (StatusCode, Json<SearchResultResponse>);

    async fn from_request(req: Request<B>, state: &S) -> Result<Self, Self::Rejection> {
        let Json(value) = Json::<T>::from_request(req, state)
            .await
            .map_err(|rejection| {
                tracing::error!("Parsing error: {}", rejection);
                let stats = SearchResultStats {
                    time: 0,
                    message: Some(format!("JSON parse error: [{}]", rejection)),
                    total: 0,
                    offset: 0,
                    amount: 0,
                    facet: HashMap::new(),
                };
                let body = SearchResultBody { docs: Vec::new() };
                (
                    StatusCode::BAD_REQUEST,
                    Json(SearchResultResponse {
                        stats: stats,
                        items: body,
                    }),
                )
            })?;

        value.validate().map_err(|rejection| {
            tracing::error!("Validation error: {}", rejection);
            let stats = SearchResultStats {
                time: 0,
                message: Some(format!("Validation error: [{}]", rejection).replace('\n', ", ")),
                total: 0,
                offset: 0,
                amount: 0,
                facet: HashMap::new(),
            };
            let body = SearchResultBody { docs: Vec::new() };
            (
                StatusCode::BAD_REQUEST,
                Json(SearchResultResponse {
                    stats: stats,
                    items: body,
                }),
            )
        })?;

        Ok(ValidatedSearchJson(value))
    }
}

/// Form形式で入力されたリクエストパラメータ
/// Form形式は実装していないので、現在は使用しない(がせっかく実装したので残している。どうせ忘れるし)
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

/// クエリパラメータをバリデーションする構造体
pub struct ValidatedSearchQueryParams<T>(pub T);

#[async_trait]
impl<T, S> FromRequestParts<S> for ValidatedSearchQueryParams<T>
where
    T: DeserializeOwned + Validate,
    S: Send + Sync,
{
    type Rejection = (StatusCode, Json<SearchResultResponse>);

    async fn from_request_parts(parts: &mut Parts, _state: &S) -> Result<Self, Self::Rejection> {
        let query = parts.uri.query().unwrap_or_default();
        let value: T = serde_qs::from_str(query).map_err(|rejection| {
            tracing::error!("Parsing error: {}", rejection);
            let stats = SearchResultStats {
                time: 0,
                message: Some(format!("Invalid format query string: [{}]", rejection)),
                total: 0,
                offset: 0,
                amount: 0,
                facet: HashMap::new(),
            };
            let body = SearchResultBody { docs: Vec::new() };
            (
                StatusCode::BAD_REQUEST,
                Json(SearchResultResponse {
                    stats: stats,
                    items: body,
                }),
            )
        })?;

        value.validate().map_err(|rejection| {
            tracing::error!("Validation error: {}", rejection);
            let stats = SearchResultStats {
                time: 0,
                message: Some(format!("Validation error: [{}]", rejection).replace('\n', ", ")),
                total: 0,
                offset: 0,
                amount: 0,
                facet: HashMap::new(),
            };
            let body = SearchResultBody { docs: Vec::new() };

            (
                StatusCode::BAD_REQUEST,
                Json(SearchResultResponse {
                    stats: stats,
                    items: body,
                }),
            )
        })?;

        Ok(ValidatedSearchQueryParams(value))
    }
}

#[cfg(test)]
mod test {
    use super::*;

    #[test]
    fn test_default() {
        let params = SearchParams {
            q: None,
            p: None,
            o: None,
            f: None,
            s: None,
        };

        let mut qs = params.as_qs();
        qs.sort();
        let mut expected = vec![
            ("defType", "edismax"),
            ("f.category.facet.mincount", "1"),
            ("f.difficulty.facet.range.end", "2000"),
            ("f.difficulty.facet.range.gap", "400"),
            ("f.difficulty.facet.range.other", "all"),
            ("f.difficulty.facet.range.start", "0"),
            ("facet", "true"),
            ("facet.field", "category"),
            ("facet.range", "difficulty"),
            ("q.alt", "*:*"),
            ("q.op", "AND"),
            ("qf", "text_ja text_en text_1gram"),
            ("sow", "true"),
            ("rows", "20"),
            ("start", "0"),
        ]
        .into_iter()
        .map(|(key, value)| (key.to_string(), value.to_string()))
        .collect::<Vec<(String, String)>>();
        expected.sort();

        assert_eq!(qs, expected)
    }

    #[test]
    fn should_wildcard_search_when_q_is_empty() {
        let params = SearchParams {
            q: Some("".to_string()),
            p: None,
            o: None,
            f: None,
            s: None,
        };

        let mut qs = params.as_qs();
        qs.sort();
        let mut expected = vec![
            ("defType", "edismax"),
            ("f.category.facet.mincount", "1"),
            ("f.difficulty.facet.range.end", "2000"),
            ("f.difficulty.facet.range.gap", "400"),
            ("f.difficulty.facet.range.other", "all"),
            ("f.difficulty.facet.range.start", "0"),
            ("facet", "true"),
            ("facet.field", "category"),
            ("facet.range", "difficulty"),
            ("q.alt", "*:*"),
            ("q.op", "AND"),
            ("qf", "text_ja text_en text_1gram"),
            ("sow", "true"),
            ("rows", "20"),
            ("start", "0"),
        ]
        .into_iter()
        .map(|(key, value)| (key.to_string(), value.to_string()))
        .collect::<Vec<(String, String)>>();
        expected.sort();

        assert_eq!(qs, expected)
    }

    #[test]
    fn rows_should_equal_to_p_parameter() {
        let params = SearchParams {
            q: None,
            p: Some(10),
            o: None,
            f: None,
            s: None,
        };

        let mut qs = params.as_qs();
        qs.sort();
        let mut expected = vec![
            ("defType", "edismax"),
            ("f.category.facet.mincount", "1"),
            ("f.difficulty.facet.range.end", "2000"),
            ("f.difficulty.facet.range.gap", "400"),
            ("f.difficulty.facet.range.other", "all"),
            ("f.difficulty.facet.range.start", "0"),
            ("facet", "true"),
            ("facet.field", "category"),
            ("facet.range", "difficulty"),
            ("q.alt", "*:*"),
            ("q.op", "AND"),
            ("qf", "text_ja text_en text_1gram"),
            ("sow", "true"),
            ("rows", "10"),
            ("start", "0"),
        ]
        .into_iter()
        .map(|(key, value)| (key.to_string(), value.to_string()))
        .collect::<Vec<(String, String)>>();
        expected.sort();

        assert_eq!(qs, expected)
    }

    #[test]
    fn start_should_equal_to_o() {
        let params = SearchParams {
            q: None,
            p: None,
            o: Some(10),
            f: None,
            s: None,
        };

        let mut qs = params.as_qs();
        qs.sort();
        let mut expected = vec![
            ("defType", "edismax"),
            ("f.category.facet.mincount", "1"),
            ("f.difficulty.facet.range.end", "2000"),
            ("f.difficulty.facet.range.gap", "400"),
            ("f.difficulty.facet.range.other", "all"),
            ("f.difficulty.facet.range.start", "0"),
            ("facet", "true"),
            ("facet.field", "category"),
            ("facet.range", "difficulty"),
            ("q.alt", "*:*"),
            ("q.op", "AND"),
            ("qf", "text_ja text_en text_1gram"),
            ("sow", "true"),
            ("rows", "20"),
            ("start", "10"),
        ]
        .into_iter()
        .map(|(key, value)| (key.to_string(), value.to_string()))
        .collect::<Vec<(String, String)>>();
        expected.sort();

        assert_eq!(qs, expected)
    }
}
