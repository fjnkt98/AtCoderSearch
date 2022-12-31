use crate::solr::core::SolrCore;
use crate::solr::models::*;
use reqwest::Client;
use url::Url;

type Result<T> = std::result::Result<T, SolrError>;

#[derive(Debug)]
pub struct SolrClient {
    /// SolrインスタンスのルートURL。e.g.) http://localhost:8983
    url: String,
    /// reqwest HTTPクライアント
    client: Client,
}

impl SolrClient {
    /// コンストラクタ
    /// 引数で与えられたURLはスキーマ(http)とホスト名しか使わない。
    /// 冗長なURL(e.g.) http://localhost:8983/solr)が与えられても、ポート番号やパスはすべて無視される
    pub fn new(url: &str, port: u32) -> Result<Self> {
        let url = Url::parse(url).map_err(|e| SolrError::UrlParseError(e))?;

        let scheme = url.scheme();
        let host = url.host_str().ok_or_else(|| SolrError::InvalidHostError)?;

        Ok(SolrClient {
            url: format!("{}://{}:{}", scheme, host, port),
            client: reqwest::Client::new(),
        })
    }

    /// Solrインスタンスのステータスを取得するメソッド
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

    /// Solrインスタンスに存在するコアの一覧を取得するメソッド
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

    /// Solrコアオブジェクトを取得するメソッド
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

#[cfg(test)]
mod tests {
    use super::*;

    /// Solrクライアント作成の通常系テスト
    #[test]
    fn test_create_solr_client() {
        let client = SolrClient::new("http://localhost", 8983).unwrap();
        assert_eq!(client.url, "http://localhost:8983");
    }

    /// Solrクライアント作成の通常系テスト
    ///
    /// 冗長なURLを与えられたときの動作を確認する。
    /// 与えられたURLのスキーマとホストのみを読み取る仕様なので、冗長なURLを与えられても
    /// スキーマとホスト以外の情報は無視される。
    #[test]
    fn test_create_solr_client_with_redundant_url() {
        let client = SolrClient::new("http://localhost:8983/solr", 8983).unwrap();
        assert_eq!(client.url, "http://localhost:8983");
    }

    /// Solrクライアント作成の異常系テスト
    #[test]
    fn test_create_solr_client_with_invalid_url() {
        let client = SolrClient::new("hogehoge", 3000);
        assert!(client.is_err());
    }

    /// Solrクライアントのステータス取得の正常系テスト
    ///
    /// 以下のコマンドでDockerコンテナを起動してからテストを実行すること。
    ///
    /// ```ignore
    /// docker run --rm -d -p 8983:8983 solr:9.1.0 solr-precreate example
    /// ```
    #[tokio::test]
    #[ignore]
    async fn test_get_status() {
        let client = SolrClient::new("http://localhost", 8983).unwrap();

        let response = client.status().await.unwrap();
        assert_eq!(response.header.status, 0);
    }

    /// コア一覧取得の正常系テスト
    ///
    /// 以下のコマンドでDockerコンテナを起動してからテストを実行すること。
    ///
    /// ```ignore
    /// docker run --rm -d -p 8983:8983 solr:9.1.0 solr-precreate example
    /// ```
    #[tokio::test]
    #[ignore]
    async fn test_get_cores() {
        let client = SolrClient::new("http://localhost", 8983).unwrap();

        let response = client.cores().await.unwrap();
        assert!(response.status.unwrap().contains_key("example"));
    }

    /// コア一覧を文字列のベクタとして取得する機能の正常系テスト
    ///
    /// 以下のコマンドでDockerコンテナを起動してからテストを実行すること。
    ///
    /// ```ignore
    /// docker run --rm -d -p 8983:8983 solr:9.1.0 solr-precreate example
    /// ```
    #[tokio::test]
    #[ignore]
    async fn test_get_cores_as_vec() {
        let client = SolrClient::new("http://localhost", 8983).unwrap();

        let response = client.cores().await.unwrap();
        let cores = response.as_vec().unwrap();
        assert_eq!(cores, vec![String::from("example")]);
    }

    /// コアオブジェクト取得メソッドの正常系テスト
    ///
    /// 以下のコマンドでDockerコンテナを起動してからテストを実行すること。
    ///
    /// ```ignore
    /// docker run --rm -d -p 8983:8983 solr:9.1.0 solr-precreate example
    /// ```
    #[tokio::test]
    #[ignore]
    async fn test_get_core() {
        let client = SolrClient::new("http://localhost", 8983).unwrap();

        let core = client.core("example").await.unwrap();
        assert_eq!(core.name, String::from("example"));
    }

    /// 存在しないコアを指定したときのテスト
    ///
    /// 以下のコマンドでDockerコンテナを起動してからテストを実行すること。
    ///
    /// ```ignore
    /// docker run --rm -d -p 8983:8983 solr:9.1.0 solr-precreate example
    /// ```
    #[tokio::test]
    #[ignore]
    async fn test_get_non_existent_core() {
        let client = SolrClient::new("http://localhost", 8983).unwrap();

        let core = client.core("hoge").await;
        assert!(core.is_err());
    }
}
