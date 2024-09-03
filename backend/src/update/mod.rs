pub mod user;

use anyhow::Context;
use futures::stream::FuturesUnordered;
// use meilisearch_sdk::client::Client;
use serde::Serialize;
use std::{fmt::Debug, future::Future, mem, time::Duration};
use tokio::sync::mpsc::Sender;
use tokio::task::JoinHandle;
use tokio_stream::{wrappers::ReceiverStream, StreamExt};

pub trait ReadRows {
    type Row: Debug + ToDocument + Send + Sync + 'static;

    fn read_rows(&self, tx: Sender<Self::Row>) -> impl Future<Output = anyhow::Result<()>> + Send;
}

pub trait ToDocument {
    type Document: Debug + Serialize + Send + Sync + 'static;

    fn to_document(self) -> impl Future<Output = anyhow::Result<Self::Document>> + Send;
}

pub async fn update_index<'a, R>(reader: R, chunk_size: usize) -> anyhow::Result<()>
where
    R: ReadRows + Send + 'static,
{
    let (row_tx, mut row_rx) = tokio::sync::mpsc::channel::<R::Row>(chunk_size);
    let (doc_tx, doc_rx) =
        tokio::sync::mpsc::channel::<<R::Row as ToDocument>::Document>(chunk_size);

    let mut tasks = FuturesUnordered::<JoinHandle<anyhow::Result<()>>>::new();

    // spawn indexer as background task
    tasks.push(tokio::task::spawn(async move {
        let rx = ReceiverStream::new(doc_rx);

        let stream = rx.chunks_timeout(chunk_size, Duration::from_secs(1));
        tokio::pin!(stream);
        while let Some(docs) = stream.next().await {
            if let Some(doc) = docs.get(0) {
                println!("{:?}", doc);
            }
        }
        Ok(())
    }));

    // spawn reader as background task
    tasks.push(tokio::task::spawn(async move {
        reader
            .read_rows(row_tx)
            .await
            .with_context(|| "read rows")?;
        Ok(())
    }));

    // receive rows and spawn task which converts rows to documents
    tasks.push(tokio::task::spawn(async move {
        let mut tasks = FuturesUnordered::<JoinHandle<anyhow::Result<()>>>::new();

        while let Some(row) = row_rx.recv().await {
            let doc_tx = doc_tx.clone();

            tasks.push(tokio::task::spawn(async move {
                let doc = row.to_document().await.with_context(|| "to document")?;
                doc_tx.send(doc).await.with_context(|| "send document")?;
                Ok(())
            }));
        }

        while let Some(task) = tasks.next().await {
            task.with_context(|| "task failed")?
                .with_context(|| "task failed")?;
        }

        mem::drop(doc_tx);

        Ok(())
    }));

    // error handling
    while let Some(task) = tasks.next().await {
        task.with_context(|| "task failed")?
            .with_context(|| "task failed")?;
    }
    Ok(())
}
