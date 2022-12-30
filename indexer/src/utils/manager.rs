use crate::solr::client::SolrClient;
use crate::utils::extractor::FullTextExtractor;
use crate::utils::models::Document;
use crate::utils::models::*;
use crate::utils::reader::RecordReader;
use futures::future::try_join_all;
use futures::TryStreamExt;
use serde_json;
use sqlx::postgres::Postgres;
use sqlx::Pool;
use tokio::fs;
use tokio::fs::File;
use tokio::io::{AsyncReadExt, AsyncWriteExt};

type Result<T> = std::result::Result<T, IndexingError>;

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
        let mut suffix = 0;
        while let Some(record) = stream.try_next().await? {
            tracing::debug!("Processing {}...", record.problem_id);

            let extractor = FullTextExtractor::new()?;
            let document = record.to_document(&extractor)?;
            buffer.push(document);

            if buffer.len() == 1000 {
                suffix += 1000;
                let filename = format!("doc-{}.json", suffix.to_string());
                tracing::info!("Saving {}", filename);

                let mut file = File::create(format!("/var/tmp/atcoder/{}", filename)).await?;
                let contents = serde_json::to_string_pretty(&buffer)
                    .map_err(|e| IndexingError::SerializeError(e))?;

                file.write_all(contents.as_bytes()).await?;

                buffer.clear();
            }
        }

        if !buffer.is_empty() {
            suffix += buffer.len() as i32;
            let filename = format!("doc-{}.json", suffix.to_string());

            tracing::info!("Saving {}", filename);

            let mut file = File::create(format!("/var/tmp/{}", filename)).await?;
            let contents = serde_json::to_string_pretty(&buffer)
                .map_err(|e| IndexingError::SerializeError(e))?;

            file.write_all(contents.as_bytes()).await?;
        }

        Ok(())
    }
    pub async fn post(&self) -> Result<()> {
        let client = SolrClient::new("http://localhost", 8983).unwrap();
        let core = client.core("atcoder").await?;

        core.reload().await?;
        core.truncate().await?;

        let mut files = fs::read_dir("/var/tmp/atcoder").await?;
        let mut target = Vec::new();
        while let Some(file) = files.next_entry().await? {
            let path = file.path();
            if let Some(extension) = path.extension() {
                if extension == "json" {
                    target.push(path)
                }
            }
        }

        let files: Vec<_> = target.iter().map(|file| File::open(file)).collect();
        let files = try_join_all(files).await?;
        let mut buffers = Vec::new();
        for mut file in files {
            let mut buffer = Vec::new();
            file.read_to_end(&mut buffer).await?;
            buffers.push(buffer);
        }

        let tasks: Vec<_> = buffers
            .into_iter()
            .map(|buffer| core.post(buffer))
            .collect();

        try_join_all(tasks).await?;
        core.commit(true).await?;

        Ok(())
    }
    pub async fn run(&self) {}
}
