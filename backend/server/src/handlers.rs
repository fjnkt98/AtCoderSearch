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

    let stats = SearchResultStats {
        time: now.duration_since(start).as_millis() as u32,
        message: None,
        total: response.response.num_found,
        offset: response.response.start,
        amount: response.response.docs.len() as u32,
        facet: facet,
    };

    let items = SearchResultBody {
        docs: response.response.docs,
    };

    // キーワード検索、かつページング無しのときのみクエリログをロギングする。
    if let Some(params) = response.header.params {
        if let Some(q) = params.get("q") {
            if q.is_string() && response.response.start == 0 {
                tracing::info!(target: "querylog", "{} {}", q, response.response.num_found)
            }
        }
    }

    Ok(SearchResultResponse {
        stats: stats,
        items: items,
    })
}

fn generate_error_response(message: &str) -> SearchResultResponse {
    let stats = SearchResultStats {
        time: 0,
        message: Some(message.to_string()),
        total: 0,
        offset: 0,
        amount: 0,
        facet: HashMap::new(),
    };
    let items = SearchResultBody { docs: Vec::new() };
    let response = SearchResultResponse { stats, items };

    response
}
