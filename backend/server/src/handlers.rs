use crate::models::*;
use anyhow::Result;
use axum::extract::Extension;
use axum::http::StatusCode;
use axum::response::IntoResponse;
use axum::Json;
use solrust::client::core::{SolrCore, SolrCoreError};
use solrust::types::response::{SolrRangeFacetKind, SolrSelectResponse};
use std::collections::HashMap;
use std::sync::Arc;
use tokio::time::Instant;
use uuid::Uuid;

#[tracing::instrument(target = "querylog", skip(core))]
pub async fn search_with_qs(
    ValidatedSearchQueryParams(params): ValidatedSearchQueryParams<SearchParams>,
    Extension(core): Extension<Arc<SolrCore>>,
) -> Result<impl IntoResponse, (StatusCode, Json<SearchResultResponse>)> {
    let start = Instant::now();

    let params = params.as_qs();

    let response = match core.select(&params).await {
        Ok(response) => response,
        Err(e) => match e {
            SolrCoreError::RequestError(e) => {
                tracing::error!("{}", e.to_string());
                return Err((
                    StatusCode::INTERNAL_SERVER_ERROR,
                    Json(generate_error_response(&e.to_string())),
                ));
            }
            SolrCoreError::DeserializeError(e) => {
                tracing::error!("{}", e.to_string());
                return Err((
                    StatusCode::INTERNAL_SERVER_ERROR,
                    Json(generate_error_response(&e.to_string())),
                ));
            }
            SolrCoreError::UnexpectedError((code, msg)) => {
                tracing::error!("{}:{}", code, msg);
                return Err((StatusCode::BAD_REQUEST, Json(generate_error_response(&msg))));
            }
        },
    };

    let response = generate_response(response, start).await.or_else(|e| {
        tracing::error!("{}", e.to_string());
        return Err((
            StatusCode::INTERNAL_SERVER_ERROR,
            Json(generate_error_response("Failed to generate response.")),
        ));
    })?;

    Ok((StatusCode::OK, Json(response)))
}

#[tracing::instrument(target = "querylog", skip(core))]
pub async fn search_with_json(
    Extension(core): Extension<Arc<SolrCore>>,
    ValidatedSearchJson(params): ValidatedSearchJson<SearchParams>,
) -> Result<impl IntoResponse, (StatusCode, Json<SearchResultResponse>)> {
    let start = Instant::now();

    let params = params.as_qs();

    let response = match core.select(&params).await {
        Ok(response) => response,
        Err(e) => match e {
            SolrCoreError::RequestError(e) => {
                tracing::error!("{}", e.to_string());
                return Err((
                    StatusCode::INTERNAL_SERVER_ERROR,
                    Json(generate_error_response(&e.to_string())),
                ));
            }
            SolrCoreError::DeserializeError(e) => {
                tracing::error!("{}", e.to_string());
                return Err((
                    StatusCode::INTERNAL_SERVER_ERROR,
                    Json(generate_error_response(&e.to_string())),
                ));
            }
            SolrCoreError::UnexpectedError((code, msg)) => {
                tracing::error!("{}:{}", code, msg);
                return Err((StatusCode::BAD_REQUEST, Json(generate_error_response(&msg))));
            }
        },
    };

    let response = generate_response(response, start).await.or_else(|e| {
        tracing::error!("{}", e.to_string());
        return Err((
            StatusCode::INTERNAL_SERVER_ERROR,
            Json(generate_error_response("Failed to generate response.")),
        ));
    })?;

    Ok((StatusCode::OK, Json(response)))
}

async fn generate_response(
    response: SolrSelectResponse<Document>,
    start: Instant,
) -> Result<SearchResultResponse> {
    let mut facet: HashMap<String, FacetResult> = HashMap::new();
    if let Some(facet_counts) = response.facet_counts {
        for (key, value) in facet_counts.facet_fields.iter() {
            facet.insert(
                key.clone(),
                FacetResult {
                    counts: value
                        .iter()
                        .map(|(key, count)| FacetCount {
                            key: key.clone(),
                            count: count.clone(),
                        })
                        .collect(),
                    start: None,
                    end: None,
                    gap: None,
                    before: None,
                    after: None,
                    between: None,
                },
            );
        }
        for (key, value) in facet_counts.facet_ranges.iter() {
            facet.insert(
                key.clone(),
                match value {
                    SolrRangeFacetKind::Integer(count) => FacetResult {
                        counts: count
                            .counts
                            .iter()
                            .map(|(key, count)| FacetCount {
                                key: key.clone(),
                                count: count.clone(),
                            })
                            .collect(),
                        start: Some(count.start.to_string()),
                        end: Some(count.end.to_string()),
                        gap: Some(count.gap.to_string()),
                        before: count.before.and_then(|before| Some(before.to_string())),
                        after: count.after.and_then(|after| Some(after.to_string())),
                        between: count.after.and_then(|between| Some(between.to_string())),
                    },
                    SolrRangeFacetKind::Float(count) => FacetResult {
                        counts: count
                            .counts
                            .iter()
                            .map(|(key, count)| FacetCount {
                                key: key.clone(),
                                count: count.clone(),
                            })
                            .collect(),
                        start: Some(count.start.to_string()),
                        end: Some(count.end.to_string()),
                        gap: Some(count.gap.to_string()),
                        before: count.before.and_then(|before| Some(before.to_string())),
                        after: count.after.and_then(|after| Some(after.to_string())),
                        between: count.after.and_then(|between| Some(between.to_string())),
                    },
                    SolrRangeFacetKind::DateTime(count) => FacetResult {
                        counts: count
                            .counts
                            .iter()
                            .map(|(key, count)| FacetCount {
                                key: key.clone(),
                                count: count.clone(),
                            })
                            .collect(),
                        start: Some(count.start.format("%Y-%m-%dT%H:%M:%S%:z").to_string()),
                        end: Some(count.end.format("%Y-%m-%dT%H:%M:%S%:z").to_string()),
                        gap: Some(count.gap.clone()),
                        before: count.before.and_then(|before| {
                            Some(before.format("%Y-%m-%dT%H:%M:%S%:z").to_string())
                        }),
                        after: count.after.and_then(|after| {
                            Some(after.format("%Y-%m-%dT%H:%M:%S%:z").to_string())
                        }),
                        between: count.after.and_then(|between| {
                            Some(between.format("%Y-%m-%dT%H:%M:%S%:z").to_string())
                        }),
                    },
                },
            );
        }
    }

    let now = Instant::now();

    let total: u32 = response.response.num_found;
    let count: u32 = response.response.docs.len() as u32;
    let rows: u32 = response
        .header
        .params
        .as_ref()
        .and_then(|params| {
            params
                .get("rows")
                .and_then(|rows| rows.as_str())
                .and_then(|rows| rows.parse::<u32>().ok())
        })
        .unwrap_or(20);
    let index: u32 = (response.response.start / rows) + 1;
    let pages: u32 = (total + rows - 1) / rows;
    let stats = SearchResultStats {
        time: now.duration_since(start).as_millis() as u32,
        total,
        index,
        count,
        pages,
        facet: facet,
    };

    // キーワード検索のときのみロギングする
    if let Some(params) = response.header.params {
        if let Some(q) = params.get("q").and_then(|q| {
            let q = q.to_string().trim_matches('"').to_string();
            if q.is_empty() {
                None
            } else {
                Some(q)
            }
        }) {
            let words = q.split_whitespace().collect::<Vec<&str>>();
            // 空文字列でないキーワード検索が実行されたときのみロギングする
            // オフセットが指定された場合はロギングしない(多重カウントになるので)
            if response.response.start == 0 && words.len() > 0 {
                // 複数キーワード検索か否かを判断する
                let word_type = if words.len() > 1 {
                    "multiple"
                } else {
                    "single"
                };
                let uuid = Uuid::new_v4();
                for (position, word) in words.iter().enumerate() {
                    tracing::info!(target: "querylog", "{} {} {} {} {}", uuid, word, response.response.num_found, position, word_type)
                }
            }
        }
    }

    Ok(SearchResultResponse {
        stats: stats,
        items: response.response.docs,
        message: None,
    })
}

fn generate_error_response(message: &str) -> SearchResultResponse {
    let stats = SearchResultStats {
        time: 0,
        total: 0,
        index: 0,
        pages: 0,
        count: 0,
        facet: HashMap::new(),
    };
    let response = SearchResultResponse {
        stats,
        items: Vec::new(),
        message: Some(message.to_string()),
    };

    response
}
