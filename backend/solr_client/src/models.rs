use itertools::Itertools;
use serde::{Deserialize, Deserializer, Serialize};
use serde_json::Value;
use std::collections::HashMap;

#[derive(Serialize, Deserialize, Debug)]
pub struct ResponseHeader {
    pub status: u32,
    #[serde(alias = "QTime")]
    pub qtime: u32,
    pub params: Option<HashMap<String, Value>>,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct SolrErrorInfo {
    pub metadata: Vec<String>,
    pub msg: String,
    pub code: u32,
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
    pub error: Option<SolrErrorInfo>,
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
    pub error: Option<SolrErrorInfo>,
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
    pub error: Option<SolrErrorInfo>,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct SolrSelectResponse {
    #[serde(alias = "responseHeader")]
    pub header: ResponseHeader,
    pub response: SolrSelectResponseBody,
    pub facet_counts: Option<FacetResult>,
    pub error: Option<SolrErrorInfo>,
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
pub struct FacetResult {
    pub facet_queries: Value,
    #[serde(deserialize_with = "deserialize_facet_fields")]
    pub facet_fields: HashMap<String, Vec<(String, u32)>>,
    pub facet_ranges: HashMap<String, RangeFacet>,
    pub facet_intervals: Value,
    pub facet_heatmaps: Value,
}

fn deserialize_facet_fields<'de, D>(
    deserializer: D,
) -> Result<HashMap<String, Vec<(String, u32)>>, D::Error>
where
    D: Deserializer<'de>,
{
    let value: HashMap<String, Vec<Value>> = Deserialize::deserialize(deserializer)?;
    let value: HashMap<String, Vec<(String, u32)>> = value
        .iter()
        .map(|(k, v)| {
            (
                k.to_string(),
                v.iter()
                    .tuples()
                    .map(|(v1, v2)| {
                        (
                            v1.as_str().unwrap_or("").to_string(),
                            v2.as_u64().unwrap_or(0) as u32,
                        )
                    })
                    .collect::<Vec<(String, u32)>>(),
            )
        })
        .collect();

    Ok(value)
}

#[derive(Serialize, Deserialize, Debug)]
pub struct RangeFacet {
    #[serde(deserialize_with = "deserialize_range_facet_counts")]
    pub counts: Vec<(String, u32)>,
    pub gap: i32,
    pub start: i32,
    pub end: i32,
    pub before: Option<i32>,
    pub after: Option<i32>,
    pub between: Option<i32>,
}

fn deserialize_range_facet_counts<'de, D>(deserializer: D) -> Result<Vec<(String, u32)>, D::Error>
where
    D: Deserializer<'de>,
{
    let value: Vec<Value> = Deserialize::deserialize(deserializer)?;
    let value: Vec<(String, u32)> = value
        .iter()
        .tuples()
        .map(|(v1, v2)| {
            (
                v1.as_str().unwrap_or("").to_string(),
                v2.as_u64().unwrap_or(0) as u32,
            )
        })
        .collect();

    Ok(value)
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
    pub error: Option<SolrErrorInfo>,
}
