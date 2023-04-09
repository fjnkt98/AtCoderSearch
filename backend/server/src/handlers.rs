use crate::models::*;
use anyhow::Result;
use axum::extract::Extension;
use axum::http::StatusCode;
use axum::response::IntoResponse;
use axum::Json;
use solrust::client::core::{SolrCore, SolrCoreError};
use solrust::types::response::SolrSelectResponse;
use std::sync::Arc;
use tokio::time::Instant;
use uuid::Uuid;

#[tracing::instrument(skip(core))]
pub async fn search_with_qs(
    ValidatedSearchQueryParams(params): ValidatedSearchQueryParams<SearchQueryParams>,
    Extension(core): Extension<Arc<SolrCore>>,
) -> Result<impl IntoResponse, (StatusCode, Json<SearchResultResponse>)> {
    tracing::info!("GET");
    Ok(handle_request(&params.into(), core).await?)
}

#[tracing::instrument(skip(core))]
pub async fn search_with_json(
    Extension(core): Extension<Arc<SolrCore>>,
    ValidatedSearchJson(params): ValidatedSearchJson<SearchParams>,
) -> Result<impl IntoResponse, (StatusCode, Json<SearchResultResponse>)> {
    tracing::info!("POST");
    Ok(handle_request(&params, core).await?)
}

#[tracing::instrument(target = "querylog", skip(core))]
async fn handle_request(
    params: &SearchParams,
    core: Arc<SolrCore>,
) -> Result<impl IntoResponse, (StatusCode, Json<SearchResultResponse>)> {
    let start_process = Instant::now();
    let response: SolrSelectResponse<Document> = match core.select(&params.as_qs()).await {
        Ok(response) => response,
        Err(e) => match e {
            SolrCoreError::RequestError(e) => {
                tracing::error!("{}", e.to_string());
                return Err((
                    StatusCode::INTERNAL_SERVER_ERROR,
                    Json(SearchResultResponse::error(&params, &e.to_string())),
                ));
            }
            SolrCoreError::DeserializeError(e) => {
                tracing::error!("{}", e.to_string());
                return Err((
                    StatusCode::INTERNAL_SERVER_ERROR,
                    Json(SearchResultResponse::error(&params, &e.to_string())),
                ));
            }
            SolrCoreError::UnexpectedError((code, msg)) => {
                tracing::error!("{}:{}", code, msg);
                return Err((
                    StatusCode::BAD_REQUEST,
                    Json(SearchResultResponse::error(&params, &msg)),
                ));
            }
        },
    };

    let total: u32 = response.response.num_found;
    let count: u32 = response.response.docs.len() as u32;
    let rows: u32 = params.limit.unwrap_or(20);
    let index: u32 = (response.response.start / rows) + 1;
    let pages: u32 = (total + rows - 1) / rows;

    // キーワード検索のときのみロギングする
    {
        let params = params.clone();
        let words = params.keyword.clone().and_then(|word| {
            Some(
                word.split_whitespace()
                    .map(|word| word.to_string())
                    .collect::<Vec<String>>(),
            )
        });
        let page = params
            .page
            .and_then(|page| if page == 1 { Some(1) } else { None });
        let sort = params.sort.and_then(|sort| {
            if sort == String::from("-score") {
                Some(sort)
            } else {
                None
            }
        });

        if let (Some(words), Some(_page), Some(_sort), Some(_filter)) =
            (words, page, sort, params.filter)
        {
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
