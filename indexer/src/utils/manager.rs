use crate::solr::client::SolrClient;
use crate::utils::models::Document;
use crate::utils::reader::RecordReader;
use anyhow::{Context, Result};
use futures::TryStreamExt;
use serde_json;
use sqlx::postgres::Postgres;
use sqlx::Pool;
use tokio::fs::File;
use tokio::io::{AsyncReadExt, AsyncWriteExt};

pub struct IndexingManager<'a> {
    reader: RecordReader<'a>,
}

impl<'a> IndexingManager<'a> {
    pub fn new(pool: &'a Pool<Postgres>) -> Self {
        IndexingManager {
            reader: RecordReader::new(pool),
        }
    }

    pub async fn write(&self) -> Result<()> {
        let mut stream = self.reader.read_rows().await?;

        let mut buffer: Vec<Document> = Vec::new();
        while let Some(record) = stream.try_next().await? {
            let document = record
                .to_document()
                .context("Couldn't convert record to document.")?;
            buffer.push(document);
        }

        tracing::info!("{} documents available.", buffer.len());

        let mut file = File::create("/var/tmp/documents.json").await?;
        let contents =
            serde_json::to_string_pretty(&buffer).context("Couldn't serialize documents.")?;
        tracing::info!("Serialized JSON length is: {}", contents.len());

        file.write_all(contents.as_bytes()).await?;

        Ok(())
    }
    pub async fn post(&self) -> Result<()> {
        let client = SolrClient::new("http://localhost", 8983).unwrap();
        let core = client.core("atcoder").await?;

        let mut file = File::open("/var/tmp/documents.json").await?;
        let mut buffer = Vec::new();
        let size = file.read_to_end(&mut buffer).await?;

        tracing::info!("Document size is: {}", size);

        core.post(buffer).await?;
        core.commit(true).await?;

        Ok(())
    }
    pub async fn run(&self) {}
}
