use crate::models::*;
use anyhow::Result;
use axum::extract::Extension;
use axum::http::StatusCode;
use axum::response::IntoResponse;
use axum::Json;
use solrust::client::core::SolrCore;
use solrust::types::response::SolrSelectResponse;
use std::sync::Arc;
use tokio::time::Instant;
use uuid::Uuid;

type SearchResponse = (StatusCode, Json<SearchResultResponse>);

pub async fn search_with_qs(
    ValidatedSearchQueryParams(params): ValidatedSearchQueryParams<SearchQueryParams>,
    Extension(core): Extension<Arc<SolrCore>>,
) -> SearchResponse {
    let params: SearchParams = params.into();
    match handle_request(&params, core).await {
        Ok(res) => res,
        Err(e) => (
            StatusCode::INTERNAL_SERVER_ERROR,
            Json(SearchResultResponse::error(&params, &e.to_string())),
        ),
    }
}

pub async fn search_with_json(
    Extension(core): Extension<Arc<SolrCore>>,
    ValidatedSearchJson(params): ValidatedSearchJson<SearchParams>,
) -> SearchResponse {
    match handle_request(&params, core).await {
        Ok(res) => res,
        Err(e) => (
            StatusCode::INTERNAL_SERVER_ERROR,
            Json(SearchResultResponse::error(&params, &e.to_string())),
        ),
    }
}

async fn handle_request(params: &SearchParams, core: Arc<SolrCore>) -> Result<SearchResponse> {
    let start_process = Instant::now();
    let response: SolrSelectResponse<Document> = core.select(&params.as_qs()).await?;

    let total: u32 = response.response.num_found;
    let count: u32 = response.response.docs.len() as u32;
    let rows: u32 = params.limit.unwrap_or(20);
    let index: u32 = (response.response.start / rows) + 1;
    let pages: u32 = (total + rows - 1) / rows;

    // クエリログのロギング
    tracing::info!(
        target: "querylog",
        uuid=Uuid::new_v4().to_string(),
        hits=response.response.num_found,
        params=serde_json::to_string(&params).unwrap()
    );

    let stats = SearchResultStats {
        time: Instant::now().duration_since(start_process).as_millis() as u32,
        total,
        index,
        count,
        pages,
        params: params.clone(),
        facet: FacetResult::from(response.facet_counts),
    };

    Ok((
        StatusCode::OK,
        Json(SearchResultResponse {
            stats: stats,
            items: response.response.docs,
            message: None,
        }),
    ))
}

pub async fn healthcheck(Extension(core): Extension<Arc<SolrCore>>) -> impl IntoResponse {
    if let Ok(response) = core.ping().await {
        if response.status == "OK" {
            return (StatusCode::OK, "OK");
        } else {
            return (StatusCode::INTERNAL_SERVER_ERROR, "ERROR");
        }
    } else {
        return (StatusCode::INTERNAL_SERVER_ERROR, "ERROR");
    }
}
