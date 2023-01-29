use crate::models::document::{Document, Record};
use crate::models::errors::GeneratingError;
use crate::utils::reader::RecordReader;
use futures::stream::FuturesUnordered;
use futures::{StreamExt, TryStreamExt};
use sqlx::postgres::Postgres;
use sqlx::Pool;
use std::io::Write;
use std::path::{Path, PathBuf};
// use tokio::fs::File;
// use tokio::io::{AsyncReadExt, AsyncWriteExt};
use std::fs::File;

use super::extractor::FullTextExtractor;

type Result<T> = std::result::Result<T, GeneratingError>;

pub struct DocumentGenerator<'a> {
    reader: RecordReader<'a>,
    savedir: PathBuf,
}

impl<'a> DocumentGenerator<'a> {
    pub fn new(pool: &'a Pool<Postgres>, savedir: &Path) -> Self {
        DocumentGenerator {
            reader: RecordReader::new(pool),
            savedir: savedir.to_path_buf(),
        }
    }

    pub async fn generate(&self, chunk_size: usize) -> Result<()> {
        let mut stream = self.reader.read_rows().await?;
        let mut buffer: Vec<Record> = Vec::new();
        let mut tasks = FuturesUnordered::new();

        let mut suffix = 0;
        while let Some(record) = stream.try_next().await? {
            buffer.push(record);

            if buffer.len() >= chunk_size {
                suffix += chunk_size;
                let mut savedir = self.savedir.clone();
                let task = tokio::spawn(async move {
                    let extractor = FullTextExtractor::new().unwrap();
                    let documents: Result<Vec<Document>> = buffer
                        .into_iter()
                        .map(|record| record.to_document(&extractor))
                        .collect();

                    if let Ok(documents) = documents {
                        let filename = format!("doc-{}.json", suffix.to_string());
                        savedir.push(filename);

                        let contents = serde_json::to_string_pretty(&documents)
                            .map_err(|e| GeneratingError::SerializeError(e))
                            .unwrap();
                        let mut file = File::create(savedir).unwrap();
                        file.write_all(contents.as_bytes()).unwrap();
                    }
                });
                tasks.push(task);

                buffer = Vec::new();
            }
        }
        if !buffer.is_empty() {
            suffix += buffer.len();
            let mut savedir = self.savedir.clone();
            let task = tokio::spawn(async move {
                let extractor = FullTextExtractor::new().unwrap();
                let documents: Result<Vec<Document>> = buffer
                    .into_iter()
                    .map(|record| record.to_document(&extractor))
                    .collect();

                if let Ok(documents) = documents {
                    let filename = format!("doc-{}.json", suffix.to_string());
                    savedir.push(filename);

                    let contents = serde_json::to_string_pretty(&documents)
                        .map_err(|e| GeneratingError::SerializeError(e))
                        .unwrap();
                    let mut file = File::create(savedir).unwrap();
                    file.write_all(contents.as_bytes()).unwrap();
                }
            });
            tasks.push(task);
        }

        while let Some(task) = tasks.next().await {
            let () = task.unwrap();
        }

        Ok(())
    }
}
