use crate::models::{SearchParams, ValidatedSearchForm, ValidatedSearchQueryParams};
use anyhow::Result;
use axum::extract::Extension;
use axum::http::StatusCode;
use axum::response::IntoResponse;
use axum::Json;
use solr_client::core::SolrCore;
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

    Ok((StatusCode::OK, Json(response)))
}
