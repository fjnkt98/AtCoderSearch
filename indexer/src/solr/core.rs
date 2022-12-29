use crate::solr::models::*;
use reqwest::header::CONTENT_TYPE;
use reqwest::Client;
use serde_json::Value;

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

    pub async fn select(&self, params: &Vec<(String, String)>) -> Result<SolrSelectResponse> {
        let response = self
            .client
            .get(format!("{}/select", self.core_url))
            .query(params)
            .send()
            .await
            .map_err(|e| SolrError::RequestError(e))?
            .text()
            .await
            .map_err(|e| SolrError::RequestError(e))?;

        let result: SolrSelectResponse =
            serde_json::from_str(&response).map_err(|e| SolrError::DeserializeError(e))?;

        Ok(result)
    }

    pub async fn analyze(&self, word: &str, field: &str, analyzer: &str) -> Result<Vec<String>> {
        let params = [("analysis.fieldvalue", word), ("analysis.fieldtype", field)];

        let response = self
            .client
            .get(format!("{}/analysis/field", self.core_url))
            .query(&params)
            .send()
            .await
            .map_err(|e| SolrError::RequestError(e))?
            .text()
            .await
            .map_err(|e| SolrError::RequestError(e))?;

        let result: SolrAnalysisResponse =
            serde_json::from_str(&response).map_err(|e| SolrError::DeserializeError(e))?;

        let result = result.analysis.field_types.get(field).unwrap();
        let result = match analyzer {
            "index" => result.index.as_ref().unwrap(),
            "query" => result.query.as_ref().unwrap(),
            _ => return Err(SolrError::InvalidValueError),
        };
        let result = result.last().unwrap().clone();

        let result = match result {
            Value::Array(array) => array
                .iter()
                .map(|e| e["text"].to_string().trim_matches('"').to_string())
                .collect::<Vec<String>>(),
            _ => Vec::new(),
        };

        Ok(result)
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
