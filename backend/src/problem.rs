use crate::atcodersearch::problem_service_server::ProblemService;
use crate::atcodersearch::{SearchProblemByKeywordRequest, SearchProblemResult};
use meilisearch_sdk::indexes::Index;
use sqlx::{Pool, Postgres};
use std::sync::Arc;
use tonic::{Request, Response, Status};

#[derive(Debug)]
pub struct ProblemSearcher {
    pool: Pool<Postgres>,
    index: Arc<Index>,
}

impl ProblemSearcher {
    pub fn new(pool: Pool<Postgres>, index: Index) -> Self {
        Self {
            pool,
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

        let res = SearchProblemResult {
            time: 0,
            total: 0,
            index: 0,
            count: 0,
            pages: 0,
            items: vec![],
            facet: None,
        };

        Ok(Response::new(res))
    }
}

#[cfg(test)]
mod test {
    use meilisearch_sdk::client::Client;
    use sqlx::postgres::PgPoolOptions;

    use super::*;

    #[tokio::test]
    async fn test_search_problem_by_keyword() -> anyhow::Result<()> {
        let pool = PgPoolOptions::new()
            .connect("postgres://atcodersearch:atcodersearch@localhost:5432/atcodersearch?sslmode=disable")
            .await?;

        let client = Client::new("http://localhost:7700", Some("meili-master-key"))?;
        let index = client.index("problems");

        let searcher = ProblemSearcher::new(pool, index);
        let req = Request::new(SearchProblemByKeywordRequest {
            q: String::from("q"),
            sorts: vec![],
            facets: vec![],
            categories: vec![],
            difficulty_from: None,
            difficulty_to: None,
            colors: vec![],
            experimental: None,
        });
        let res = searcher.search_by_keyword(req).await?;

        let want = SearchProblemResult {
            time: 0,
            total: 0,
            index: 0,
            count: 0,
            pages: 0,
            items: vec![],
            facet: None,
        };
        assert_eq!(res.into_inner(), want);

        Ok(())
    }
}
