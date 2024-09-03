pub mod user;

use anyhow::Context;
use futures::stream::FuturesUnordered;
use serde::Serialize;
use std::{fmt::Debug, future::Future, mem, pin::Pin, time::Duration};
use tokio::task::JoinHandle;
use tokio_stream::{wrappers::ReceiverStream, Stream, StreamExt};

pub trait ReadRows<'a> {
    type Row: Debug + ToDocument + Send + Sync + 'static;

    fn read_rows(
        &'a self,
    ) -> impl Future<
        Output = anyhow::Result<
            Pin<Box<dyn Stream<Item = Result<Self::Row, sqlx::Error>> + Send + 'a>>,
        >,
    > + Send;
}

pub trait ToDocument {
    type Document: Debug + Serialize + Send + Sync + 'static;

    fn to_document(self) -> anyhow::Result<Self::Document>;
}

pub async fn update_index<'a, R>(reader: &'a R, chunk_size: usize) -> anyhow::Result<()>
where
    R: ReadRows<'a> + 'a + Send,
{
    let (tx, rx) = tokio::sync::mpsc::channel(chunk_size);

    let mut tasks: FuturesUnordered<JoinHandle<()>> = FuturesUnordered::new();

    let uploader = tokio::task::spawn(async move {
        let rx = ReceiverStream::new(rx);

        let stream = rx.chunks_timeout(chunk_size, Duration::from_secs(1));
        tokio::pin!(stream);
        while let Some(docs) = stream.next().await {
            if let Some(doc) = docs.get(0) {
                todo!();
            }
        }
    });
    tasks.push(uploader);

    let mut stream = reader
        .read_rows()
        .await
        .unwrap_or_else(|e| panic!(" read rows: {:#}", e));

    while let Some(row) = stream
        .try_next()
        .await
        .unwrap_or_else(|e| panic!("stream: {:?}", e))
    {
        let tx = tx.clone();

        let task = tokio::task::spawn(async move {
            let doc = row
                .to_document()
                .unwrap_or_else(|e| panic!("to document: {:#}", e));

            tx.send(doc)
                .await
                .unwrap_or_else(|e| panic!("send document: {}", e));
        });

        tasks.push(task);
    }
    mem::drop(tx);

    while let Some(task) = tasks.next().await {
        task.with_context(|| "task failed")?;
    }

    Ok(())
}
