use crate::models::*;
use anyhow::Result;
use axum::extract::Extension;
use axum::http::StatusCode;
use axum::response::IntoResponse;
use axum::Json;
use solr_client::core::{SolrCore, SolrCoreError};
use solr_client::models::SolrSelectResponse;
use std::collections::HashMap;
use std::sync::Arc;
use tokio::time::Instant;
use tracing::instrument;

#[instrument(skip(core))]
pub async fn search_with_qs(
    ValidatedSearchQueryParams(params): ValidatedSearchQueryParams<SearchParams>,
    Extension(core): Extension<Arc<SolrCore>>,
) -> Result<impl IntoResponse, StatusCode> {
    let start = Instant::now();

    let params = params.as_qs();

    let response = match core.select(&params).await {
        Ok(response) => response,
        Err(e) => match e {
            SolrCoreError::RequestError(e) => {
                tracing::error!("{}", e.to_string());
                return Err(StatusCode::INTERNAL_SERVER_ERROR);
            }
            SolrCoreError::DeserializeError(e) => {
                tracing::error!("{}", e.to_string());
                return Err(StatusCode::INTERNAL_SERVER_ERROR);
            }
            SolrCoreError::UnexpectedError((code, msg)) => {
                tracing::error!("{}:{}", code, msg);
                return Err(StatusCode::BAD_REQUEST);
            }
        },
    };

    let response = generate_response(response, start)
        .await
        .or(Err(StatusCode::INTERNAL_SERVER_ERROR))?;

    Ok((StatusCode::OK, Json(response)))
}

#[instrument(skip(core))]
pub async fn search_with_json(
    Extension(core): Extension<Arc<SolrCore>>,
    ValidatedSearchJson(params): ValidatedSearchJson<SearchParams>,
) -> Result<impl IntoResponse, StatusCode> {
    let start = Instant::now();

    let params = params.as_qs();

    let response = core
        .select(&params)
        .await
        .or(Err(StatusCode::BAD_REQUEST))?;

    let response = generate_response(response, start)
        .await
        .or(Err(StatusCode::INTERNAL_SERVER_ERROR))?;

    Ok((StatusCode::OK, Json(response)))
}

async fn generate_response(
    response: SolrSelectResponse,
    start: Instant,
) -> Result<SearchResultResponse> {
    let docs: Vec<Document> = serde_json::from_value(response.response.docs)?;

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
                FacetResult {
                    counts: value
                        .counts
                        .iter()
                        .map(|(key, count)| FacetCount {
                            key: key.clone(),
                            count: count.clone(),
                        })
                        .collect(),
                    start: Some(value.start.to_string()),
                    end: Some(value.end.to_string()),
                    gap: Some(value.gap.to_string()),
                    before: None,
                    after: None,
                    between: None,
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
        amount: docs.len() as u32,
        facet: if facet.is_empty() { None } else { Some(facet) },
    };

    let items = SearchResultBody { docs: docs };

    Ok(SearchResultResponse {
        stats: stats,
        items: items,
    })
}
