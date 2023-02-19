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

## インデクシング実行コマンド

### クローリング

```bash
indexer crawl [OPTIONS]
```

`--all`オプションをつけると、AtCoderに存在するすべての問題情報をクロールする。
クローリング間隔は1000msにしているので、非常に長い時間がかかるので注意。

デフォルトでは差分クローリングを行う。データベースに存在しない問題のみをクロールする。

### ドキュメントJSONファイル生成

```bash
indexer generate [PATH]
```

パスを指定しない場合、環境変数`DOCUMENT_SAVE_DIRECTORY`に指定したパスにJSONファイルが保存される。

### ドキュメントJSONファイルポスト

```indexer
indexer post [OPTIONS] [PATH]
```

パスを指定しない場合、環境変数`DOCUMENT_SAVE_DIRECTORY`に指定したパスに保存されたJSONファイルを使用する。

`--optimize`オプションをつけると、ポスト実行後のコミット時に最適化をリクエストする。これにより、インデックスファイルが1つにまとめられる。
デフォルト設定では最適化オプションは適用されない。
