# learn-run-go

Go の実行やデプロイについて学習するためのリポジトリ。

## 環境情報

(Cloud9 2021/06/25)

```bash
$ go version
go version go1.15.12 linux/amd64
$ docker -v
Docker version 20.10.4, build d3cb89e
```

### docker-compose のインストール

```bash
$ docker-compose -v
bash: docker-compose: command not found

$ sudo curl -L "https://github.com/docker/compose/releases/download/1.29.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose && chmod +x /usr/local/bin/docker-compose
$ sudo chmod +x /usr/local/bin/docker-compose
$ docker-compose -v
docker-compose version 1.29.0, build 07737305
```

### AWS Copilot CLI のインストール

```bash
$ sudo su -
[root@ ~]# curl -Lo /usr/local/bin/copilot https://github.com/aws/copilot-cli/releases/latest/download/copilot-linux && chmod +x /usr/local/bin/copilot
[root@ ~]# exit

$ copilot --version
copilot version: v1.8.1
```

Cloud9 以外の各種環境へのインストールについては以下を参照。  
`https://aws.github.io/copilot-cli/ja/docs/getting-started/install/`

### AWS プロファイル設定

(Cloud9 で AWS Copilot CLI を実行するために必要)

Preferences > AWS SETTINGS > AWS managed temporary credentials : OFF

```bash
$ aws configure
AWS Access Key ID [None]: XXXXXXXXXXXXXXXXXXXX
AWS Secret Access Key [None]: XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
Default region name [None]: ap-northeast-1
Default output format [None]:
```

## 01.helloworld

### a) 標準ライブラリのみ

`hello world!` という文字列を出力するプログラム。

#### ビルドして実行可能ファイル(バイナリ)を作成し、それを実行する。

```bash
$ pwd
/learn-run-go/01.helloworld/a.simple

# ビルドしてバイナリを作成する。
$ go build -o ./bin/main main.go

# バイナリを実行する。
$ ./bin/main
hello world!
```

#### ビルドして実行する(バイナリは出力されない)

```bash
$ pwd
/learn-run-go/01.helloworld/a.simple

$ hello world!
hello world!
```


### b) 外部ライブラリ利用

外部のロガーパッケージを利用して `hello world!` を出力するプログラム。

#### ビルド前に、modules を初期化する

```bash
$ pwd
/learn-run-go/01.helloworld/b.modules

# 既に存在する mod, sum ファイルを削除しておく。
$ rm -f go.mod go.sum

$ go mod init main
go: creating new go.mod: module main
go: to add module requirements and sums:
        go mod tidy

$ go mod tidy
go: finding module for package github.com/sirupsen/logrus
go: downloading github.com/sirupsen/logrus v1.8.1
go: found github.com/sirupsen/logrus in github.com/sirupsen/logrus v1.8.1
go: downloading golang.org/x/sys v0.0.0-20191026070338-33540a1f6037
go: downloading github.com/stretchr/testify v1.2.2
```

`go mod init main` により `go.mod` ファイルが作成される。  
`go mod tidy` により `go.sum` ファイルが作成される。  
これらを行わないと、利用しているライブラリの情報不足により、ビルドできない。

これらのファイルはGitなどのバージョン管理システムの管理に含めた方が良い。
(毎回最新バージョンを利用する場合は含めなくても良い。本項ではあえて含めていない。)

#### ビルドして実行可能ファイル(バイナリ)を作成し、それを実行する。(ワンライナー)

```bash
$ go build -o ./bin/main main.go && ./bin/main
INFO[0000] hello world!
```

## 02.lambda

クエリパラメータとして受け取った JSON に `"hello": "world!"` という Key-Value を加えて返す Lambda ファンクション。

### modules の初期化

```bash
$ pwd
/learn-run-go/02.lambda/func

$ go mod init func && go mod tidy
go: creating new go.mod: module func
go: to add module requirements and sums:
        go mod tidy
go: finding module for package github.com/aws/aws-lambda-go/events
go: finding module for package github.com/aws/aws-lambda-go/lambda
go: found github.com/aws/aws-lambda-go/events in github.com/aws/aws-lambda-go v1.24.0
go: found github.com/aws/aws-lambda-go/lambda in github.com/aws/aws-lambda-go v1.24.0
```

### ローカルで実行

#### main メソッドから実行する

```bash
$ pwd
/learn-run-go/02.lambda/func

$ go build -o ../bin/main main.go && ../bin/main
=== start main ===
2021/06/25 06:42:08 expected AWS Lambda environment variables [_LAMBDA_SERVER_PORT AWS_LAMBDA_RUNTIME_API] are not defined
```

main メソッドから Lambda のエントリーポイントを呼べずに終了している。

#### 単体テストで Lambda のエントリーポイントを直接実行する

```bash
$ go test
query = map[hey:yo!]
StatusCode=200, Body={
   "hello": "world!",
   "hey": "yo!"
}
PASS
ok      func    0.131s
```

### デプロイ

#### Lambda 関数の作成

AWS マネジメントコンソール > [AWS Lambda](https://ap-northeast-1.console.aws.amazon.com/lambda/home?region=ap-northeast-1#/functions) から作成する。  
[関数の作成] > [1 から作成] > [関数名:learn-run-go] > [ランタイム:Go 1.x] > [関数の作成]

#### ビルド & パッケージング(zip)

```bash
command
$ GOOS=linux GOARCH=amd64 go build -o ../bin/hello main.go
$ (cd ../bin && zip -r ../lambda-package.zip *)
```

Go で作成した場合、ハンドラ名がデフォルトで `hello` となっているため、バイナリのファイル名も `hello` としている。  
Lambda の実行環境に合わせてクロスコンパイルの指定(`GOOS=linux GOARCH=amd64`)をしておく。  
(Cloud9 の場合クロスコンパイルの指定は必要ないが、Lambda にデプロイする際のおまじないと思って付けておく。)

### zip のアップロード

```bash
command
$ aws lambda update-function-code --function-name learn-run-go --zip-file fileb://../lambda-package.zip
```

### Lambda のテストイベント

```json
{
  "queryStringParameters": { "hey": "yo!" }
}
```

## 03.webapp

`hello` と `world` のページをもつ Web アプリケーション

```bash
03.webapp
|--docker-compose.debug.yml
|--docker-compose.release.yml
|--webapp
|  |--Dockerfile.debug
|  |--Dockerfile.release
|  |--main.go
|  |--controllers
|  |  |--server.go
|  |  |--route_main.go
|  |--views
|  |  |--css
|  |  |  |--bootstrap.min.css
|  |  |--js
|  |  |  |--bootstrap.bundle.min.js
|  |  |  |--jquery-3.6.0.min.js
|  |  |--templates
|  |  |  |--layout.html
|  |  |  |--public_navbar.html
|  |  |  |--hello.html
|  |  |  |--world.html
```

### ビルドして実行

#### ローカル

```bash
$ pwd
/learn-run-go/03.webapp/main

$ go mod init main && go mod tidy
$ go build -o ../bin/main main.go && ../bin/main
```

#### Docker(デバッグ用)

```bash
command
$ pwd
/learn-run-go/03.webapp

$ docker-compose -f docker-compose.debug.yml up --build
```

#### Docker(リリース用)

```bash
command
$ docker-compose -f docker-compose.release.yml up --build
```

#### イメージサイズの確認

```bash
command
$ docker images
REPOSITORY                TAG       IMAGE ID       CREATED          SIZE
go_webapp_image_debug     latest    aab1c76b6b7f   52 seconds ago   338MB
go_webapp_image_release   latest    7d8ca74c9b04   16 seconds ago   12.8MB
```

### App Runner によるデプロイ

#### App Runner

```bash
command
$ copilot init
? Application name: appname
? Workload type: Request-Driven Web Service
? Service name: servname
? Dockerfile: webapp/Dockerfile.release
? Port: 80
Ok great, we'll set up a Request-Driven Web Service named servname in application appname listening on port 80.
:
All right, you're all set for local development.
? Deploy: Yes
:
✔ Proposing infrastructure changes for stack appname-test-servname
- Creating the infrastructure for stack appname-test-servname                     [create complete]  [255.6s]
  - An IAM Role for App Runner to use on your behalf to pull your image from ECR  [create complete]  [20.7s]
  - An IAM role to control permissions for the containers in your service         [create complete]  [23.4s]
  - An App Runner service to run and manage your containers                       [create complete]  [225.4s]
✔ Deployed servname, you can access it at https://xxxxxxxxxx.ap-northeast-1.awsapprunner.com.
```

## その他

### VSCode でデバッグ

`.vscode/launch.json` の `connfigurachionns[0].program` を目的の main パッケージがあるフォルダに切り替える。  
F5 押下により、指定した main パッケージのプログラムが開始されます。
