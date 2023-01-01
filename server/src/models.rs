use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize)]
pub struct SearchResultResponse {
    pub stats: SearchResultStats,
    pub items: SearchResultBody,
}

#[derive(Serialize, Deserialize)]
pub struct SearchResultStats {
    pub total: u32,
    pub start: u32,
    pub amount: u32,
    pub facet: FacetResult,
}

#[derive(Serialize, Deserialize)]
pub struct SearchResultBody {
    docs: Vec<Document>,
}

#[derive(Serialize, Deserialize)]
pub struct FacetResult {}

#[derive(Serialize, Deserialize)]
pub struct Document {}
