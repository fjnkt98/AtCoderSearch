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
use solrust::querybuilder::{
    common::SolrCommonQueryBuilder,
    dismax::SolrDisMaxQueryBuilder,
    edismax::{EDisMaxQueryBuilder, SolrEDisMaxQueryBuilder},
    facet::{FieldFacetBuilder, FieldFacetSortOrder, RangeFacetBuilder, RangeFacetOtherOptions},
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
#[derive(Debug, Serialize, Deserialize, Clone, Validate, PartialEq, Eq)]
pub struct SearchParams {
    // 検索キーワード
    #[validate(length(max = 200))]
    pub keyword: Option<String>,
    // 1ページ当たり返却数
    #[validate(range(min = 1, max = 200))]
    pub limit: Option<u32>,
    // 返却ページ番号
    #[validate(range(min = 1))]
    pub page: Option<u32>,
    // フィルタリング条件
    pub filter: Option<FilteringParameters>,
    // ソート順
    #[validate(custom = "validate_sort_option")]
    pub sort: Option<String>,
    // ファセットカウントのフィールド
    #[validate(custom = "validate_facet_field_vec")]
    pub facet: Option<Vec<String>>,
}

impl SearchParams {
    pub fn empty() -> Self {
        Self {
            keyword: None,
            limit: None,
            page: None,
            filter: None,
            sort: None,
            facet: None,
        }
    }
}

/// fパラメータに指定できる値
#[derive(Debug, Serialize, Deserialize, Clone, Validate, PartialEq, Eq)]
pub struct FilteringParameters {
    pub category: Option<Vec<String>>,
    pub difficulty: Option<RangeFilteringParameter<u32>>,
}

/// 範囲フィルタリングに指定できる値
#[derive(Debug, Serialize, Deserialize, Clone, Validate, PartialEq, Eq)]
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

/// ファセットカウントのフィールドを限定するためのバリデーション関数
fn validate_facet_field_vec(values: &[String]) -> Result<(), ValidationError> {
    for value in values.iter() {
        if value != "category" || value != "difficulty" {
            return Err(ValidationError::new(
                "Invalid facet field option. Select from category or difficulty.",
            ));
        }
    }

    Ok(())
}

impl SearchParams {
    /// リクエストパラメータからSolrへのリクエストパラメータを生成するメソッド
    /// Solrへ送るパラメータはすべてここで生成する
    pub fn as_qs(&self) -> Vec<(String, String)> {
        let rows: u32 = self.limit.unwrap_or(20);
        let page: u32 = self.page.unwrap_or(1);
        let start: u32 = (page - 1) * rows;

        let mut builder = EDisMaxQueryBuilder::new()
            .rows(rows)
            .start(start)
            .qf("text_ja text_en text_1gram")
            .q_alt(&QueryOperand::from("*:*"))
            .op(Operator::AND)
            .sow(true);

        if let Some(facet) = &self.facet {
            for field in facet.iter() {
                if field == "category" {
                    let category_facet = FieldFacetBuilder::new("category")
                        .min_count(0)
                        .sort(FieldFacetSortOrder::Index);
                    builder = builder.facet(&category_facet);
                }
                if field == "difficulty" {
                    let difficulty_facet = RangeFacetBuilder::new(
                        "difficulty",
                        0.to_string(),
                        4000.to_string(),
                        400.to_string(),
                    )
                    .other(RangeFacetOtherOptions::All);

                    builder = builder.facet(&difficulty_facet);
                }
            }
        }

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
                    params: SearchParams::empty(),
                    facet: FacetResult::empty(),
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
                params: SearchParams::empty(),
                facet: FacetResult::empty(),
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

/// クエリパラメータをパースした結果を格納し、SearchParams構造体へ変換するための中間構造体
#[derive(Debug, Serialize, Deserialize, Validate, PartialEq, Eq)]
pub struct SearchQueryParams {
    #[validate(length(max = 200))]
    pub keyword: Option<String>,
    #[validate(range(min = 1, max = 200))]
    pub limit: Option<u32>,
    #[validate(range(min = 1))]
    pub page: Option<u32>,
    #[serde(alias = "filter.category")]
    pub filter_category: Option<String>,
    #[serde(alias = "filter.difficulty.from")]
    pub filter_difficulty_from: Option<u32>,
    #[serde(alias = "filter.difficulty.to")]
    pub filter_difficulty_to: Option<u32>,
    #[validate(custom = "validate_sort_option")]
    pub sort: Option<String>,
    #[validate(custom = "validate_facet_field_str")]
    pub facet: Option<String>,
}

fn validate_facet_field_str(value: &str) -> Result<(), ValidationError> {
    for value in value.split(',').into_iter() {
        if value != "category" && value != "difficulty" {
            return Err(ValidationError::new(
                "Invalid facet field option. Select from category or difficulty.",
            ));
        }
    }

    Ok(())
}

impl Into<SearchParams> for SearchQueryParams {
    fn into(self) -> SearchParams {
        let filter = Some(FilteringParameters {
            category: self.filter_category.and_then(|category| {
                Some(
                    category
                        .split(',')
                        .map(String::from)
                        .collect::<Vec<String>>(),
                )
            }),
            difficulty: Some(RangeFilteringParameter {
                from: self.filter_difficulty_from,
                to: self.filter_difficulty_to,
            }),
        });
        SearchParams {
            keyword: self.keyword,
            limit: self.limit,
            page: self.page,
            filter: filter,
            sort: self.sort,
            facet: self.facet.and_then(|facet| {
                Some(facet.split(',').map(String::from).collect::<Vec<String>>())
            }),
        }
    }
}

#[async_trait]
impl<T, S> FromRequestParts<S> for ValidatedSearchQueryParams<T>
where
    T: DeserializeOwned + Validate,
    S: Send + Sync,
{
    type Rejection = (StatusCode, Json<SearchResultResponse>);

    async fn from_request_parts(parts: &mut Parts, _state: &S) -> Result<Self, Self::Rejection> {
        let query = parts.uri.query().unwrap_or_default();
        let value: T = serde_urlencoded::from_str(query).map_err(|rejection| {
            tracing::error!("Parsing error: {}", rejection);
            let stats = SearchResultStats {
                time: 0,
                total: 0,
                index: 0,
                pages: 0,
                count: 0,
                params: SearchParams::empty(),
                facet: FacetResult::empty(),
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
                params: SearchParams::empty(),
                facet: FacetResult::empty(),
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
    use itertools::{sorted, Itertools};
    use rstest::*;

    /// すべてのパラメータがデフォルト値のときのテスト
    #[test]
    fn test_default() {
        let params = SearchParams {
            keyword: None,
            limit: None,
            page: None,
            filter: None,
            sort: None,
            facet: None,
        };

        let qs = sorted(params.as_qs()).collect_vec();
        let expected = sorted(
            [
                ("defType", "edismax"),
                ("q.alt", "*:*"),
                ("q.op", "AND"),
                ("qf", "text_ja text_en text_1gram"),
                ("sow", "true"),
                ("rows", "20"),
                ("start", "0"),
            ]
            .map(|(key, value)| (key.to_string(), value.to_string())),
        )
        .collect_vec();

        assert_eq!(qs, expected)
    }

    /// 検索キーワードが空文字列のときはすべてのドキュメントを取得する
    #[test]
    fn should_do_wildcard_search_when_keyword_is_empty() {
        let params = SearchParams {
            keyword: Some("".to_string()),
            limit: None,
            page: None,
            filter: None,
            sort: None,
            facet: None,
        };

        let qs = sorted(params.as_qs()).collect_vec();
        let expected = sorted(
            [
                ("defType", "edismax"),
                ("q.alt", "*:*"),
                ("q.op", "AND"),
                ("qf", "text_ja text_en text_1gram"),
                ("sow", "true"),
                ("rows", "20"),
                ("start", "0"),
            ]
            .into_iter()
            .map(|(key, value)| (key.to_string(), value.to_string())),
        )
        .collect_vec();

        assert_eq!(qs, expected)
    }

    /// limitパラメータに指定した値がrowsパラメータになることを確かめるテスト
    #[test]
    fn rows_should_equal_to_limit_parameter() {
        let params = SearchParams {
            keyword: None,
            limit: Some(10),
            page: None,
            filter: None,
            sort: None,
            facet: None,
        };

        let qs = params
            .as_qs()
            .into_iter()
            .filter(|(key, _)| key == "rows")
            .collect_vec();

        let expected = vec![("rows".to_string(), "10".to_string())];

        assert_eq!(qs, expected)
    }

    /// pageパラメータに1が指定されたらstartパラメータが0になることを確かめる
    #[test]
    fn start_should_equal_to_0_when_page_is_1() {
        let params = SearchParams {
            keyword: None,
            limit: Some(20),
            page: Some(1),
            filter: None,
            sort: None,
            facet: None,
        };

        let qs = params
            .as_qs()
            .into_iter()
            .filter(|(key, _)| key == "start")
            .collect_vec();
        let expected = vec![("start".to_string(), "0".to_string())];

        assert_eq!(qs, expected)
    }

    /// 2ページ目を指定したときstartパラメータに20が指定されることを確かめる
    #[test]
    fn start_should_equal_to_20_when_page_is_2() {
        let params = SearchParams {
            keyword: None,
            limit: Some(20),
            page: Some(2),
            filter: None,
            sort: None,
            facet: None,
        };

        let qs = params
            .as_qs()
            .into_iter()
            .filter(|(key, _)| key == "start")
            .collect_vec();
        let expected = vec![("start".to_string(), "20".to_string())];

        assert_eq!(qs, expected)
    }

    /// 単一カテゴリ絞り込みテスト
    #[test]
    fn filter_by_single_category() {
        let params = SearchParams {
            keyword: None,
            limit: None,
            page: None,
            filter: Some(FilteringParameters {
                category: Some(vec![String::from("ABC")]),
                difficulty: None,
            }),
            sort: None,
            facet: None,
        };

        let qs = params
            .as_qs()
            .into_iter()
            .filter(|(key, _)| key == "fq")
            .collect_vec();
        let expected = vec![("fq".to_string(), r#"category:"ABC""#.to_string())];

        assert_eq!(qs, expected);
    }

    /// 複数カテゴリ絞り込みテスト
    #[test]
    fn filter_by_multiple_category() {
        let params = SearchParams {
            keyword: None,
            limit: None,
            page: None,
            filter: Some(FilteringParameters {
                category: Some(vec![
                    String::from("ABC"),
                    String::from("Other Contests"),
                    String::from("ABC-Like"),
                ]),
                difficulty: None,
            }),
            sort: None,
            facet: None,
        };

        let qs = params
            .as_qs()
            .into_iter()
            .filter(|(key, _)| key == "fq")
            .collect_vec();
        let expected = vec![(
            "fq".to_string(),
            r#"category:"ABC" OR category:"Other Contests" OR category:"ABC\-Like""#.to_string(),
        )];

        assert_eq!(qs, expected);
    }

    /// 難易度絞り込み
    #[test]
    fn filter_by_difficulty() {
        let params = SearchParams {
            keyword: None,
            limit: None,
            page: None,
            filter: Some(FilteringParameters {
                category: None,
                difficulty: Some(RangeFilteringParameter {
                    from: Some(800),
                    to: Some(1200),
                }),
            }),
            sort: None,
            facet: None,
        };

        let qs = params
            .as_qs()
            .into_iter()
            .filter(|(key, _)| key == "fq")
            .collect_vec();
        let expected = vec![("fq".to_string(), r#"difficulty:[800 TO 1200}"#.to_string())];

        assert_eq!(qs, expected);
    }

    /// 昇順ソート指定
    #[test]
    fn sort_by_ascending() {
        let params = SearchParams {
            keyword: None,
            limit: None,
            page: None,
            filter: None,
            sort: Some(String::from("start_at")),
            facet: None,
        };

        let qs = params
            .as_qs()
            .into_iter()
            .filter(|(key, _)| key == "sort")
            .collect_vec();
        let expected = vec![("sort".to_string(), "start_at asc".to_string())];
        assert_eq!(qs, expected);
    }

    /// 降順ソート指定
    #[test]
    fn sort_by_descending() {
        let params = SearchParams {
            keyword: None,
            limit: None,
            page: None,
            filter: None,
            sort: Some(String::from("-score")),
            facet: None,
        };

        let qs = params
            .as_qs()
            .into_iter()
            .filter(|(key, _)| key == "sort")
            .collect_vec();
        let expected = vec![("sort".to_string(), "score desc".to_string())];
        assert_eq!(qs, expected);
    }

    #[test]
    fn category_facet() {
        let params = SearchParams {
            keyword: None,
            limit: None,
            page: None,
            filter: None,
            sort: None,
            facet: Some(vec![String::from("category")]),
        };

        let expected = sorted(
            [
                ("defType", "edismax"),
                ("f.category.facet.mincount", "0"),
                ("f.category.facet.sort", "index"),
                ("facet", "true"),
                ("facet.field", "category"),
                ("q.alt", "*:*"),
                ("q.op", "AND"),
                ("qf", "text_ja text_en text_1gram"),
                ("sow", "true"),
                ("rows", "20"),
                ("start", "0"),
            ]
            .into_iter()
            .map(|(key, value)| (key.to_string(), value.to_string())),
        )
        .collect_vec();

        let qs = sorted(params.as_qs().into_iter()).collect_vec();
        assert_eq!(qs, expected);
    }

    #[test]
    fn difficulty_facet() {
        let params = SearchParams {
            keyword: None,
            limit: None,
            page: None,
            filter: None,
            sort: None,
            facet: Some(vec![String::from("difficulty")]),
        };

        let expected = sorted(
            [
                ("defType", "edismax"),
                ("f.difficulty.facet.range.end", "4000"),
                ("f.difficulty.facet.range.gap", "400"),
                ("f.difficulty.facet.range.other", "all"),
                ("f.difficulty.facet.range.start", "0"),
                ("facet", "true"),
                ("facet.range", "difficulty"),
                ("q.alt", "*:*"),
                ("q.op", "AND"),
                ("qf", "text_ja text_en text_1gram"),
                ("sow", "true"),
                ("rows", "20"),
                ("start", "0"),
            ]
            .into_iter()
            .map(|(key, value)| (key.to_string(), value.to_string())),
        )
        .collect_vec();

        let qs = sorted(params.as_qs().into_iter()).collect_vec();
        assert_eq!(qs, expected);
    }

    #[rstest]
    #[case("keyword=", SearchQueryParams {keyword: Some(String::from("")), limit: None, page: None, filter_category: None, filter_difficulty_from: None, filter_difficulty_to: None, sort: None, facet: None})]
    #[case("keyword=hoge", SearchQueryParams {keyword: Some(String::from("hoge")), limit: None, page: None, filter_category: None, filter_difficulty_from: None, filter_difficulty_to: None, sort: None, facet: None })]
    #[case("limit=20", SearchQueryParams {keyword: None, limit: Some(20), page: None, filter_category: None, filter_difficulty_from: None, filter_difficulty_to: None, sort: None, facet: None })]
    #[case("page=2", SearchQueryParams {keyword: None, limit: None, page: Some(2), filter_category: None, filter_difficulty_from: None, filter_difficulty_to: None, sort: None, facet: None })]
    #[case("filter_category=ABC,ARC", SearchQueryParams {keyword: None, limit: None, page: None, filter_category: Some(String::from("ABC,ARC")), filter_difficulty_from: None, filter_difficulty_to: None, sort: None, facet: None})]
    #[case("filter_difficulty_from=800", SearchQueryParams {keyword: None, limit: None, page: None, filter_category: None, filter_difficulty_from: Some(800), filter_difficulty_to: None, sort: None, facet: None })]
    #[case("filter_difficulty_to=1000", SearchQueryParams {keyword: None, limit: None, page: None, filter_category: None, filter_difficulty_from: None, filter_difficulty_to: Some(1000), sort: None, facet: None })]
    #[case("sort=-score", SearchQueryParams {keyword: None, limit: None, page: None, filter_category: None, filter_difficulty_from: None, filter_difficulty_to: None, sort: Some(String::from("-score")), facet: None})]
    fn deserialize_query_parameters(#[case] input: &str, #[case] expected: SearchQueryParams) {
        let value = serde_urlencoded::from_str::<SearchQueryParams>(input);

        assert!(value.is_ok());
        assert_eq!(value.unwrap(), expected);
    }

    #[test]
    fn convert_query_params_to_search_params() {
        let query_params = SearchQueryParams {
            keyword: Some(String::from("hoge")),
            limit: Some(20),
            page: None,
            filter_category: Some(String::from("ABC,ARC")),
            filter_difficulty_from: Some(800),
            filter_difficulty_to: None,
            sort: None,
            facet: None,
        };

        let params: SearchParams = query_params.into();

        let expected = SearchParams {
            keyword: Some(String::from("hoge")),
            limit: Some(20),
            page: None,
            filter: Some(FilteringParameters {
                category: Some(vec![String::from("ABC"), String::from("ARC")]),
                difficulty: Some(RangeFilteringParameter {
                    from: Some(800),
                    to: None,
                }),
            }),
            sort: None,
            facet: None,
        };

        assert_eq!(params, expected);
    }
}
