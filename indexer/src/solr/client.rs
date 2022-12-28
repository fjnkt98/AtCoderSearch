use crate::solr::core::SolrCore;
use crate::solr::models::*;
use reqwest::Client;
use url::Url;

type Result<T> = std::result::Result<T, SolrError>;

#[derive(Debug)]
pub struct SolrClient {
    url: String,
    client: Client,
}

impl SolrClient {
    pub fn new(url: &str, port: u32) -> Result<Self> {
        let url = Url::parse(url).map_err(|e| SolrError::UrlParseError(e))?;

        let scheme = url.scheme();
        let host = url.host_str().ok_or_else(|| SolrError::InvalidHostError)?;

        Ok(SolrClient {
            url: format!("{}://{}:{}", scheme, host, port),
            client: reqwest::Client::new(),
        })
    }

    pub async fn status(&self) -> Result<SolrSystemInfo> {
        let path = "solr/admin/info/system";

        let response = self
            .client
            .get(format!("{}/{}", self.url, path))
            .send()
            .await
            .map_err(|e| SolrError::RequestError(e))?
            .text()
            .await
            .map_err(|e| SolrError::RequestError(e))?;

        let result: SolrSystemInfo =
            serde_json::from_str(&response).map_err(|e| SolrError::DeserializeError(e))?;

        Ok(result)
    }

    pub async fn cores(&self) -> Result<SolrCoreList> {
        let path = "solr/admin/cores";

        let response = self
            .client
            .get(format!("{}/{}", self.url, path))
            .send()
            .await
            .map_err(|e| SolrError::RequestError(e))?
            .text()
            .await
            .map_err(|e| SolrError::RequestError(e))?;

        let result: SolrCoreList =
            serde_json::from_str(&response).map_err(|e| SolrError::DeserializeError(e))?;

        Ok(result)
    }

    pub async fn core(&self, name: &str) -> Result<SolrCore> {
        let cores = self
            .cores()
            .await?
            .status
            .ok_or_else(|| SolrError::SpecifiedCoreNotFoundError)?;

        if !cores.contains_key(name) {
            return Err(SolrError::SpecifiedCoreNotFoundError);
        }

        Ok(SolrCore::new(name, &self.url))
    }
}
