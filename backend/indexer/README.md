# Indexer

Solrコアのインデクシング処理を行うクレート。
以下の3つの機能が実装されている。

- AtCoderおよびAtCoder Problemsからコンテストおよび問題の情報を取得し、データベースに格納する
- データベースからコンテストおよび問題の情報を読み出し、ドキュメントに加工してJSONファイルに保存する
- JSONファイルをSolrコアにポストする

## `.env`ファイル

ディレクトリのルートに、以下の環境変数を定義した`.env`ファイルを準備する必要がある。

- `PGUSER`: PostgreSQLのデフォルトユーザ
- `PGPASSWORD`: PostgreSQLのデフォルトパスワード
- `PGHOST`: PostgreSQLのデフォルトホスト
- `PGPORT`: PostgreSQLのデフォルトポート
- `PGDATABASE`: PostgreSQLのデフォルトデータベース
- `DATABASE_URL`: データベースURL
- `SOLR_HOST`: Solrインスタンスのホスト(e.g. `http://localhost`)
- `SOLR_PORT`: Solrインスタンスのポート番号
- `CORE_NAME`: 対象コア名。現状1つのコアについてのみ対応している
- `DOCUMENT_SAVE_DIRECTORY`: ドキュメントJSONのデフォルトの保存先。ドキュメント生成時にパスを指定しない場合、このディレクトリにJSONが保存される。
