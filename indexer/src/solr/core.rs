use crate::solr::models::*;
use reqwest::header::CONTENT_TYPE;
use reqwest::Client;

type Result<T> = std::result::Result<T, SolrError>;

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
            .map_err(|e| SolrError::RequestError(e))?
            .text()
            .await
            .map_err(|e| SolrError::RequestError(e))?;

        let result: SolrCoreList =
            serde_json::from_str(&response).map_err(|e| SolrError::DeserializeError(e))?;

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
            .map_err(|e| SolrError::RequestError(e))?
            .text()
            .await
            .map_err(|e| SolrError::RequestError(e))?;

        let result: SolrSimpleResponse =
            serde_json::from_str(&response).map_err(|e| SolrError::DeserializeError(e))?;

        if result.header.status != 0 {
            return Err(SolrError::CoreReloadError);
        }

        Ok(result.header.status)
    }

    pub async fn select(&self, params: &Vec<(String, String)>) -> Result<()> {
        let _response = self
            .client
            .get(format!("{}/select", self.core_url))
            .query(params)
            .send()
            .await
            .map_err(|e| SolrError::RequestError(e))?
            .text()
            .await
            .map_err(|e| SolrError::RequestError(e))?;

        todo!();
    }

    pub async fn analyze(&self) -> Result<()> {
        todo!();
    }

    pub async fn post(&self, body: Vec<u8>) -> Result<()> {
        let response = self
            .client
            .post(format!("{}/update", self.core_url))
            .header(CONTENT_TYPE, "application/json")
            .body(body)
            .send()
            .await?;

        if response.status().as_u16() != 200 {
            return Err(SolrError::CorePostError);
        }

        Ok(())
    }

    pub async fn commit(&self, optimize: bool) -> Result<()> {
        if optimize {
            self.post(br#"{"optimize": {}}"#.to_vec()).await?;
        } else {
            self.post(br#"{"commit": {}}"#.to_vec()).await?;
        }

        Ok(())
    }

    pub async fn rollback(&self) -> Result<()> {
        self.post(br#"{"rollback": {}}"#.to_vec()).await?;

        Ok(())
    }

    pub async fn truncate(&self) -> Result<()> {
        self.post(br#"{"delete":{"query": "*:*"}}"#.to_vec())
            .await?;

        Ok(())
    }
}
