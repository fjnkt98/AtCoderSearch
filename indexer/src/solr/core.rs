use crate::solr::models::*;
use anyhow::{ensure, Context, Result};
use reqwest::Client;

pub struct SolrCore {
    pub name: String,
    pub base_url: String,
    pub core_url: String,
    client: Client,
}

impl SolrCore {
    pub fn new(name: &str, base_url: &str) -> Self {
        let core_url = format!("{}/{}", base_url, name);

        SolrCore {
            name: String::from(name),
            base_url: String::from(base_url),
            core_url: core_url,
            client: reqwest::Client::new(),
        }
    }
    pub async fn status(&self) -> Result<CoreStatus> {
        let path = "solr/admin/cores";

        let response = self
            .client
            .get(format!("{}/{}", self.base_url, path))
            .query(&[("action", "status"), ("core", &self.name)])
            .send()
            .await
            .context("Failed to get response.")?
            .text()
            .await
            .context("Failed to get JSON body from response.")?;

        let result: SolrCoreList =
            serde_json::from_str(&response).context("Failed to parse JSON body")?;

        let status = result.status.unwrap().get(&self.name).unwrap().clone();

        Ok(status)
    }

    pub async fn reload(&self) -> Result<u32> {
        let path = "solr/admin/cores";

        let response = self
            .client
            .get(format!("{}/{}", self.base_url, path))
            .query(&[("action", "status"), ("core", &self.name)])
            .send()
            .await
            .context("Failed to get response.")?
            .text()
            .await
            .context("Failed to get JSON body from response.")?;

        let result: SolrSimpleResponse =
            serde_json::from_str(&response).context("Failed to parse JSON body")?;

        ensure!(result.header.status == 0);

        Ok(result.header.status)
    }

    pub async fn select(&self) -> Result<()> {
        todo!();
    }

    pub async fn analyze(&self) -> Result<()> {
        todo!();
    }
}
