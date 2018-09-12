# 東京アメッシュの画像を Idobata.io に POST する

[Idobata.io](https://idobata.io/home) にテキストや画像を POST するパッケージと、[東京アメッシュ](http://tokyo-ame.jwa.or.jp) から最新の画像を取得するパッケージを作って見たので、合わせてみた。

- Idobata.io にテキストや画像を POST する
    - https://github.com/mattsan/emattsan-go/tree/master/idobata
- 東京アメッシュから最新の画像を取得する
    - https://github.com/mattsan/emattsan-go/tree/master/amesh

## 使い方

Go がインストールされていて `$GOPATH` が設定されている必要があります。

### ビルドとデプロイ

Go の作業ディレクトリに移動してリポジトリを clone します。

```
$ cd $GOPATH/src
$ git clone git@github.com:mattsan/amesh-idobata.git
$ cd amesh-idobata
```

make でビルドします。

```
$ make build
```

POST するルームのエンドポイントを環境変数 `IDOBATA_HOOK_ENDPOINT_URL` に設定して、デプロイします。

```
$ export IDOBATA_HOOK_ENDPOINT_URL=https://idobata.io/hook/custom/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
$ sls deploy
```

### 実行する

コマンドラインから invoke して動作を確認します

```
$ sls invoke -f main
```

### イベントを設定する

`serverless.yml` を編集して定期イベントを設定してデプロイすると定期的に POST するようになる。

```yml
     environment:
       IDOBATA_HOOK_ENDPOINT_URL: ${env:IDOBATA_HOOK_ENDPOINT_URL}
+    events:
+      - schedule: cron(* * * * ? *) # 1 分毎に POST する

```

