use crate::models::{SearchParams, ValidatedSearchParams};
use anyhow::Result;
use axum::extract::{Extension, Form, Query};
use axum::http::StatusCode;
use axum::response::IntoResponse;
use axum::Json;
use solr_client::core::SolrCore;
use std::sync::Arc;
use tracing::instrument;

#[instrument(skip(core))]
pub async fn search_with_qs(
    Extension(core): Extension<Arc<SolrCore>>,
    Query(params): Query<SearchParams>,
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
    Form(params): Form<SearchParams>,
) -> Result<impl IntoResponse, StatusCode> {
    let q = params.q.unwrap_or(String::from("*"));
    let params = vec![(String::from("q"), format!("text_ja:{}", q))];

    let response = core
        .select(&params)
        .await
        .or(Err(StatusCode::BAD_REQUEST))?;

    Ok((StatusCode::OK, Json(response)))
}
