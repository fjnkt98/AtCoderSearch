use crate::solr::core::SolrCore;
use crate::solr::models::*;
use anyhow::{anyhow, ensure, Context, Result};
use reqwest::Client;
// use serde_json::Value;
use url::Url;

#[derive(Debug)]
pub struct SolrClient {
    url: String,
    client: Client,
}

impl SolrClient {
    pub fn new(url: &str, port: u32) -> Result<Self> {
        let url = Url::parse(url)?;

        let scheme = url.scheme();
        let host = url
            .host_str()
            .ok_or_else(|| anyhow!("Failed to parse URL host."))?;

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
            .context("Failed to get response.")?
            .text()
            .await
            .context("Failed to get JSON body from response.")?;

        let result: SolrSystemInfo =
            serde_json::from_str(&response).context("Failed to parse JSON body")?;

        Ok(result)
    }

    pub async fn cores(&self) -> Result<SolrCoreList> {
        let path = "solr/admin/cores";

        let response = self
            .client
            .get(format!("{}/{}", self.url, path))
            .send()
            .await
            .context("Failed to get response.")?
            .text()
            .await
            .context("Failed to get JSON body from response.")?;

        let result: SolrCoreList =
            serde_json::from_str(&response).context("Failed to parse JSON body")?;

        Ok(result)
    }

    pub async fn core(&self, name: &str) -> Result<SolrCore> {
        let cores = self
            .cores()
            .await?
            .status
            .ok_or_else(|| anyhow!("Any cores does not exists."))?;

        ensure!(cores.contains_key(name), "Specified core does not exists.");

        Ok(SolrCore::new(name, &self.url))
    }
}
