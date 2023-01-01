# Crawler

AtCoderおよびAtCoder Problemsからコンテストおよび問題の情報を取得するオブジェクトを提供するクレート。

## `.env`ファイル

ディレクトリのルートに、以下の環境変数を定義した`.env`ファイルを準備する必要がある。

- `PGUSER`: PostgreSQLのデフォルトユーザ
- `PGPASSWORD`: PostgreSQLのデフォルトパスワード
- `PGHOST`: PostgreSQLのデフォルトホスト
- `PGPORT`: PostgreSQLのデフォルトポート
- `PGDATABASE`: PostgreSQLのデフォルトデータベース
- `DATABASE_URL`: データベースURL
