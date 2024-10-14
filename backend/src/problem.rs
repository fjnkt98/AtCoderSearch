use crate::atcodersearch::problem_service_server::ProblemService;
use crate::atcodersearch::{Problem, SearchProblemByKeywordRequest, SearchProblemResult};
use itertools::Itertools;
use meilisearch_sdk::indexes::Index;
use meilisearch_sdk::search::Selectors;
use serde::Deserialize;
use sqlx::{Pool, Postgres};
use std::sync::Arc;
use tonic::{Request, Response, Status};

#[derive(Debug)]
pub struct ProblemSearcher {
    _pool: Pool<Postgres>,
    index: Arc<Index>,
}

impl ProblemSearcher {
    pub fn new(pool: Pool<Postgres>, index: Index) -> Self {
        Self {
            _pool: pool,
            index: Arc::new(index),
        }
    }
}

#[tonic::async_trait]
impl ProblemService for ProblemSearcher {
    async fn search_by_keyword(
        &self,
        request: Request<SearchProblemByKeywordRequest>,
    ) -> Result<Response<SearchProblemResult>, Status> {
        let req = request.into_inner();

        let mut query = self.index.search();
        query
            .with_query(&req.q)
            .with_attributes_to_retrieve(Selectors::Some(&[
                "problemId",
                "problemTitle",
                "problemUrl",
                "contestId",
                "contestTitle",
                "contestUrl",
                "difficulty",
                "color",
                "startAt",
                "duration",
                "rateChange",
                "category",
            ]));

        let limit = req.limit.unwrap_or(20) as usize;
        let offset = if let Some(page) = req.page {
            if page <= 0 {
                1
            } else {
                (page as usize - 1) * limit
            }
        } else {
            0
        };
        query.with_limit(limit).with_offset(offset);

        let result = match query.execute::<Item>().await {
            Ok(r) => r,
            Err(e) => {
                // tracing::error!("failed to execute search query: {:?}", e);
                return Err(Status::unknown(format!(
                    "failed to execute search query: {:?}",
                    e
                )));
            }
        };

        let res = SearchProblemResult {
            time: result.processing_time_ms as i64,
            total: result.total_hits.unwrap_or(0) as i64,
            index: result
                .offset
                .and_then(|offset| Some((offset / limit) + 1))
                .unwrap_or(0) as i64,
            count: result.hits.len() as i64,
            pages: result
                .total_hits
                .and_then(|hits| Some((hits + limit - 1) / limit))
                .unwrap_or(0) as i64,
            items: result
                .hits
                .iter()
                .cloned()
                .map(|r| r.result.into())
                .collect_vec(),
            facet: None,
        };

        Ok(Response::new(res))
    }
}

#[derive(Debug, Clone, PartialEq, Deserialize)]
#[serde(rename_all = "camelCase")]
struct Item {
    pub problem_id: String,
    pub problem_title: String,
    pub problem_url: String,
    pub contest_id: String,
    pub contest_title: String,
    pub contest_url: String,
    pub difficulty: Option<i64>,
    pub color: Option<String>,
    pub start_at: i64,
    pub duration: i64,
    pub rate_change: String,
    pub category: String,
}

impl Into<Problem> for Item {
    fn into(self) -> Problem {
        Problem {
            problem_id: self.problem_id,
            problem_title: self.problem_title,
            problem_url: self.problem_url,
            contest_id: self.contest_id,
            contest_title: self.contest_title,
            contest_url: self.contest_url,
            difficulty: self.difficulty,
            color: self.color,
            start_at: self.start_at,
            duration: self.duration,
            rate_change: self.rate_change,
            category: self.category,
        }
    }
}

#[cfg(test)]
mod test {
    use crate::testutil::{
        // create_client_from_container,
        create_db_container,
        // create_engine_container,
        create_pool_from_container,
    };
    use meilisearch_sdk::client::Client;
    use testcontainers::runners::AsyncRunner;

    use super::*;

    #[tokio::test]
    async fn test_search_problem_by_keyword() -> anyhow::Result<()> {
        let db = create_db_container()?.start().await?;
        let pool = create_pool_from_container(&db).await?;

        // let engine = create_engine_container()?.start().await?;
        // let client = create_client_from_container(&engine).await?;
        let client = Client::new("http://localhost:7700", Some("meili-master-key"))?;
        let index = client.index("problems");

        let searcher = ProblemSearcher::new(pool, index);
        let req = Request::new(SearchProblemByKeywordRequest {
            limit: Some(50),
            page: None,
            q: String::from("ABC"),
            sorts: vec![],
            facets: vec![],
            categories: vec![],
            difficulty_from: None,
            difficulty_to: None,
            colors: vec![],
            experimental: None,
        });
        let _res = searcher.search_by_keyword(req).await?;

        Ok(())
    }
}
