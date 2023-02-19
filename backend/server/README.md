# Server

## 環境変数

ディレクトリのルートに、以下の環境変数を定義した`.env`ファイルを準備する必要がある。

- `API_SERVER_LISTEN_PORT`: APIサーバがリクエストを待ち受けるポート番号。デフォルトは8000。
- `CORE_NAME`: Solrコア名。
- `LOG_DIRECTORY`: ログを出力するファイル名。デフォルトは`/var/tmp/atcoder/log`。
- `SOLR_HOST`: Solrのホスト。デフォルトは`http://localhost`。
- `SOLR_PORT`: Solrのポート番号。デフォルトは`8983`。
- `RUST_LOG`: ログレベル。
