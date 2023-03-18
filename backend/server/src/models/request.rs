use crate::models::{FacetResult, SearchResultResponse, SearchResultStats};
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
use serde_qs::Config;
use solrust::querybuilder::{
    common::SolrCommonQueryBuilder,
    dismax::SolrDisMaxQueryBuilder,
    edismax::{EDisMaxQueryBuilder, SolrEDisMaxQueryBuilder},
    facet::{FieldFacetBuilder, RangeFacetBuilder, RangeFacetOtherOptions},
    q::{
        Aggregation, Operator, PhraseQueryOperand, QueryExpression, QueryOperand, RangeQueryOperand,
    },
    sort::SortOrderBuilder,
};
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
    // 検索キーワード
    #[validate(length(max = 200))]
    pub keyword: Option<String>,
    // 1ページ当たり返却数
    #[validate(range(max = 200))]
    pub limit: Option<u32>,
    // 返却ページ番号
    #[validate(range(min = 1))]
    pub page: Option<u32>,
    // フィルタリング条件
    pub filter: Option<FilteringParameters>,
    // ソート順
    #[validate(custom = "validate_sort_option")]
    pub sort: Option<String>,
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
        "start_at" | "-start_at" | "difficulty" | "-difficulty" | "-score" => {}
        _ => return Err(ValidationError::new("Invalid sort field option. Select from start_at, -start_at, difficulty, -difficulty, -score.")),
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
            3600.to_string(),
            400.to_string(),
        )
        .other(RangeFacetOtherOptions::All);

        let rows: u32 = self.limit.unwrap_or(20);
        let page: u32 = self.page.unwrap_or(1);
        let start: u32 = (page - 1) * rows;

        let mut builder = EDisMaxQueryBuilder::new()
            .rows(rows)
            .start(start)
            .qf("text_ja text_en text_1gram")
            .q_alt(&QueryOperand::from("*:*"))
            .op(Operator::AND)
            .sow(true)
            .facet(&category_facet)
            .facet(&difficulty_facet);

        if let Some(q) = &self.keyword {
            let q = RE.replace_all(q, r"\$0");
            if !q.is_empty() {
                builder = builder.q(String::from(q));
            }
        }

        if let Some(s) = &self.sort {
            let sort = if s.starts_with("-") {
                SortOrderBuilder::new().desc(&s[1..])
            } else {
                SortOrderBuilder::new().asc(s)
            };
            builder = builder.sort(&sort);
        }

        if let Some(f) = &self.filter {
            if let Some(category) = &f.category {
                let fq = QueryExpression::sum(
                    category
                        .iter()
                        .map(|c| QueryOperand::from(PhraseQueryOperand::new("category", c)))
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
                    total: 0,
                    index: 0,
                    pages: 0,
                    count: 0,
                    facet: FacetResult::from(None),
                };
                (
                    StatusCode::BAD_REQUEST,
                    Json(SearchResultResponse {
                        stats: stats,
                        items: Vec::new(),
                        message: Some(format!("JSON parse error: [{}]", rejection)),
                    }),
                )
            })?;

        value.validate().map_err(|rejection| {
            tracing::error!("Validation error: {}", rejection);
            let stats = SearchResultStats {
                time: 0,
                total: 0,
                index: 0,
                pages: 0,
                count: 0,
                facet: FacetResult::from(None),
            };
            (
                StatusCode::BAD_REQUEST,
                Json(SearchResultResponse {
                    stats: stats,
                    items: Vec::new(),
                    message: Some(format!("Validation error: [{}]", rejection).replace('\n', ", ")),
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
        let config = Config::new(1, false);
        let query = parts.uri.query().unwrap_or_default();
        let value: T = config.deserialize_str(query).map_err(|rejection| {
            tracing::error!("Parsing error: {}", rejection);
            let stats = SearchResultStats {
                time: 0,
                total: 0,
                index: 0,
                pages: 0,
                count: 0,
                facet: FacetResult::from(None),
            };
            (
                StatusCode::BAD_REQUEST,
                Json(SearchResultResponse {
                    stats: stats,
                    items: Vec::new(),
                    message: Some(format!("Invalid format query string: [{}]", rejection)),
                }),
            )
        })?;

        value.validate().map_err(|rejection| {
            tracing::error!("Validation error: {}", rejection);
            let stats = SearchResultStats {
                time: 0,
                total: 0,
                index: 0,
                pages: 0,
                count: 0,
                facet: FacetResult::from(None),
            };
            (
                StatusCode::BAD_REQUEST,
                Json(SearchResultResponse {
                    stats: stats,
                    items: Vec::new(),
                    message: Some(format!("Validation error: [{}]", rejection).replace('\n', ", ")),
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
            keyword: None,
            limit: None,
            page: None,
            filter: None,
            sort: None,
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
    fn should_do_wildcard_search_when_q_is_empty() {
        let params = SearchParams {
            keyword: Some("".to_string()),
            limit: None,
            page: None,
            filter: None,
            sort: None,
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
    fn rows_should_equal_to_c_parameter() {
        let params = SearchParams {
            keyword: None,
            limit: Some(10),
            page: None,
            filter: None,
            sort: None,
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
    fn start_should_equal_to_0_when_p_is_1() {
        let params = SearchParams {
            keyword: None,
            limit: Some(20),
            page: Some(1),
            filter: None,
            sort: None,
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
    fn start_should_equal_to_20_when_p_is_2() {
        let params = SearchParams {
            keyword: None,
            limit: Some(20),
            page: Some(2),
            filter: None,
            sort: None,
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
            ("start", "20"),
        ]
        .into_iter()
        .map(|(key, value)| (key.to_string(), value.to_string()))
        .collect::<Vec<(String, String)>>();
        expected.sort();

        assert_eq!(qs, expected)
    }
}
