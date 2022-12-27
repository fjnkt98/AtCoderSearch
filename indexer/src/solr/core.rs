use crate::solr::models::*;
use anyhow::{ensure, Context, Result};
use reqwest::header::CONTENT_TYPE;
use reqwest::Client;

pub struct SolrCore {
    pub name: String,
    pub base_url: String,
    pub core_url: String,
    client: Client,
}

impl SolrCore {
    pub fn new(name: &str, base_url: &str) -> Self {
        let core_url = format!("{}/solr/{}", base_url, name);

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

    pub async fn select(&self, params: &Vec<(String, String)>) -> Result<()> {
        let _response = self
            .client
            .get(format!("{}/select", self.core_url))
            .query(params)
            .send()
            .await
            .context("Failed to get response.")?
            .text()
            .await
            .context("Failed to get JSON body from response.")?;

        todo!();
    }

    pub async fn analyze(&self) -> Result<()> {
        todo!();
    }

    pub async fn post(&self, body: String) -> Result<()> {
        let response = self
            .client
            .post(format!("{}/update", self.core_url))
            .header(CONTENT_TYPE, "application/json")
            .body(body)
            .send()
            .await?;

        ensure!(response.status().as_u16() == 200);

        Ok(())
    }

    pub async fn commit(&self, optimize: bool) -> Result<()> {
        if optimize {
            self.post(String::from(r#"{"optimize": {}}"#)).await?;
        } else {
            self.post(String::from(r#"{"commit": {}}"#)).await?;
        }

        Ok(())
    }

    pub async fn rollback(&self) -> Result<()> {
        self.post(String::from(r#"{"rollback": {}}"#)).await?;

        Ok(())
    }
}

// #[derive(Default)]
// pub struct ParameterBuilder {
//     pub q: Option<String>,
//     pub start: Option<u32>,
//     pub rows: Option<u32>,
//     pub fq: Option<String>,
//     pub fl: Option<String>,
//     pub sort: Option<String>,
// }

// impl ParameterBuilder {
//     pub fn new() -> Self {
//         ParameterBuilder {
//             q: None,
//             start: None,
//             rows: None,
//             fq: None,
//             fl: None,
//             sort: None,
//         }
//     }

//     pub fn q(self) -> Self {
//         self
//     }

//     pub fn build(&self) -> Vec<(String, String)> {
//         todo!();
//     }
// }
