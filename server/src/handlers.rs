use crate::models::*;
use anyhow::Result;
use axum::extract::Extension;
use axum::http::StatusCode;
use axum::response::IntoResponse;
use axum::Json;
use solr_client::core::SolrCore;
use solr_client::models::*;
use std::sync::Arc;
use tracing::instrument;

#[instrument(skip(core))]
pub async fn search_with_qs(
    ValidatedSearchQueryParams(params): ValidatedSearchQueryParams<SearchParams>,
    Extension(core): Extension<Arc<SolrCore>>,
) -> Result<impl IntoResponse, StatusCode> {
    let q = params.q.unwrap_or(String::from("*"));
    let params = vec![(String::from("q"), format!("text_ja:{}", q))];

    let response = core
        .select(&params)
        .await
        .or(Err(StatusCode::BAD_REQUEST))?;

    let response = generate_response(response)
        .await
        .or(Err(StatusCode::INTERNAL_SERVER_ERROR))?;

    Ok((StatusCode::OK, Json(response)))
}

#[instrument(skip(core))]
pub async fn search_with_form(
    Extension(core): Extension<Arc<SolrCore>>,
    ValidatedSearchForm(params): ValidatedSearchForm<SearchParams>,
) -> Result<impl IntoResponse, StatusCode> {
    let q = params.q.unwrap_or(String::from("*"));
    let params = vec![(String::from("q"), format!("text_ja:{}", q))];

    let response = core
        .select(&params)
        .await
        .or(Err(StatusCode::BAD_REQUEST))?;

    let response = generate_response(response)
        .await
        .or(Err(StatusCode::INTERNAL_SERVER_ERROR))?;

    Ok((StatusCode::OK, Json(response)))
}

async fn generate_response(response: SolrSelectResponse) -> Result<SearchResultResponse> {
    let docs: Vec<Document> = serde_json::from_value(response.response.docs)?;

    let stats = SearchResultStats {
        total: response.response.num_found,
        start: response.response.start,
        amount: docs.len() as u32,
        facet: None,
    };

    let items = SearchResultBody { docs: docs };

    Ok(SearchResultResponse {
        stats: stats,
        items: items,
    })
}
