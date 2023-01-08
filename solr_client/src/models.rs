use serde::{Deserialize, Serialize};
use serde_json::Value;
use std::collections::HashMap;
use thiserror::Error;

#[derive(Debug, Error)]
pub enum SolrError {
    #[error("Failed to request to solr")]
    RequestError(#[from] reqwest::Error),
    #[error("Failed to parse URL")]
    UrlParseError(#[from] url::ParseError),
    #[error("Given URL host is invalid")]
    InvalidHostError,
    #[error("Failed to deserialize JSON data")]
    DeserializeError(#[from] serde_json::Error),
    #[error("Specified core name does not exist")]
    SpecifiedCoreNotFoundError,
    #[error("Failed to reload core")]
    CoreReloadError,
    #[error("Failed to post data")]
    CorePostError,
    #[error("Invalid argument has given")]
    InvalidValueError,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct ResponseHeader {
    pub status: u32,
    #[serde(alias = "QTime")]
    pub qtime: u32,
    pub params: Option<HashMap<String, String>>,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct LuceneInfo {
    #[serde(alias = "solr-spec-version")]
    pub solr_spec_version: String,
    #[serde(alias = "solr-impl-version")]
    pub solr_impl_version: String,
    #[serde(alias = "lucene-spec-version")]
    pub lucene_spec_version: String,
    #[serde(alias = "lucene-impl-version")]
    pub lucene_impl_version: String,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct SolrSystemInfo {
    #[serde(alias = "responseHeader")]
    pub header: ResponseHeader,
    pub mode: String,
    pub solr_home: String,
    pub core_root: String,
    pub lucene: LuceneInfo,
    pub jvm: Value,
    pub security: Value,
    pub system: Value,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct IndexInfo {
    #[serde(alias = "numDocs")]
    pub num_docs: u64,
    #[serde(alias = "maxDoc")]
    pub max_doc: u64,
    #[serde(alias = "deletedDocs")]
    pub deleted_docs: u64,
    pub version: u64,
    #[serde(alias = "segmentCount")]
    pub segment_count: u64,
    pub current: bool,
    #[serde(alias = "hasDeletions")]
    pub has_deletions: bool,
    pub directory: String,
    #[serde(alias = "segmentsFile")]
    pub segments_file: String,
    #[serde(alias = "segmentsFileSizeInBytes")]
    pub segments_file_size_in_bytes: u64,
    #[serde(alias = "userData")]
    pub user_data: Value,
    #[serde(alias = "sizeInBytes")]
    pub size_in_bytes: u64,
    pub size: String,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct CoreStatus {
    pub name: String,
    #[serde(alias = "instanceDir")]
    pub instance_dir: String,
    #[serde(alias = "dataDir")]
    pub data_dir: String,
    pub config: String,
    pub schema: String,
    #[serde(alias = "startTime")]
    pub start_time: String,
    pub uptime: u64,
    pub index: IndexInfo,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct SolrCoreList {
    #[serde(alias = "responseHeader")]
    pub header: ResponseHeader,
    #[serde(alias = "initFailures")]
    pub init_failures: Value,
    pub status: Option<HashMap<String, CoreStatus>>,
}

impl SolrCoreList {
    pub fn as_vec(&self) -> Option<Vec<String>> {
        if let Some(cores) = &self.status {
            Some(cores.keys().cloned().collect())
        } else {
            return None;
        }
    }
}

#[derive(Serialize, Deserialize, Debug)]
pub struct SolrSimpleResponse {
    #[serde(alias = "responseHeader")]
    pub header: ResponseHeader,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct SolrSelectResponse {
    #[serde(alias = "responseHeader")]
    pub header: ResponseHeader,
    pub response: SolrSelectResponseBody,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct SolrSelectResponseBody {
    #[serde(alias = "numFound")]
    pub num_found: u32,
    pub start: u32,
    #[serde(alias = "numFoundExact")]
    pub num_found_exact: bool,
    pub docs: Value,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct SolrAnalysisBody {
    pub field_types: HashMap<String, SolrAnalysisField>,
    pub field_names: HashMap<String, SolrAnalysisField>,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct SolrAnalysisField {
    pub index: Option<Vec<Value>>,
    pub query: Option<Vec<Value>>,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct SolrAnalysisResponse {
    #[serde(alias = "responseHeader")]
    pub header: ResponseHeader,
    pub analysis: SolrAnalysisBody,
}
