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

$ sudo curl -L "https://github.com/docker/compose/releases/download/1.29.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
$ sudo chmod +x /usr/local/bin/docker-compose
$ docker-compose -v
docker-compose version 1.29.0, build 07737305
```

### AWS Copilot CLI のインストール

```bash
$ sudo su -
$ curl -Lo /usr/local/bin/copilot https://github.com/aws/copilot-cli/releases/latest/download/copilot-linux && chmod +x /usr/local/bin/copilot
$ copilot --version
copilot version: v1.8.1
```

Cloud9 以外の各種環境へのインストールについては以下を参照。  
https://aws.github.io/copilot-cli/ja/docs/getting-started/install/

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

### 標準ライブラリのみ

`hello world!` という文字列を出力するプログラム。

#### ビルドして実行可能ファイル(バイナリ)を作成する

```bash
$ pwd
/learn-run-go/01.helloworld/01.simple

$ go build -o ./bin/main main.go
```

#### バイナリを実行する

```bash
$ ./bin/main
hello world!
```

### 外部ライブラリ利用

外部のロガーパッケージを利用して `hello world!` を出力するプログラム。

#### ビルド前に、modules を初期化する

```bash
$ pwd
/learn-run-go/01.helloworld/02.modules

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

`init`の後ろの文字列は何でも良いが、公開したりする場合にはモジュールが存在するリポジトリのパスとするのが良いらしい。

> 例 :  
> `$ go mod init github.com/ww2or3ww/learn-run-go/02.modules`

外部モジュールを利用しないシンプルなプログラムであっても、これらのコマンドは実行しておいた方が良いと思われる。

#### ビルド & 実行

```bash
$ go build -o ./bin/main main.go && ./bin/main
INFO[0000] hello world!
```

## 03.lambda

クエリパラメータとして受け取った JSON に、`"hello": "world!"` という Key-Value を加えて返す Lambda ファンクション。

### modules の初期化

```bash
$ pwd
/learn-run-go/03.lambda/func

$ go mod init func && go mod tidy
go: creating new go.mod: module func
go: to add module requirements and sums:
        go mod tidy
go: finding module for package github.com/aws/aws-lambda-go/events
go: finding module for package github.com/aws/aws-lambda-go/lambda
go: found github.com/aws/aws-lambda-go/events in github.com/aws/aws-lambda-go v1.24.0
go: found github.com/aws/aws-lambda-go/lambda in github.com/aws/aws-lambda-go v1.24.0
```

### 実行(テスト)

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

#### 関数の作成

AWS マネジメントコンソール > [AWS Lambda](https://ap-northeast-1.console.aws.amazon.com/lambda/home?region=ap-northeast-1#/functions) から作成する。  
[関数の作成] > [1 から作成] > [関数名:learn-run-go] > [ランタイム:Go 1.x] > [関数の作成]

#### ビルド & パッケージング(zip)

```bash
$ GOOS=linux GOARCH=amd64 go build -o ../bin/hello main.go
$ (cd ../bin && zip -r ../lambda-package.zip *)
```

Go で作成した場合、ハンドラ名がデフォルトで `hello` となっているため、バイナリのファイル名も `hello` としている。  
Lambda の実行環境に合わせてクロスコンパイルの指定(`GOOS=linux GOARCH=amd64`)をする必要がある。

### zip のアップロード

```bash
$ aws lambda update-function-code --function-name learn-run-go --zip-file fileb://../lambda-package.zip
```

### Lambda のテストイベント

```json
{
  "queryStringParameters": { "hey": "yo!" }
}
```

## 04.webapp

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
/learn-run-go/04.webapp/webapp

$ go build -o ../bin/main main.go && ../bin/main
```

#### Docker(デバッグ用)

```bash
$ pwd
/learn-run-go/04.webapp

$ docker-compose -f docker-compose.debug.yml up --build
```

#### Docker(リリース用)

```bash
$ docker-compose -f docker-compose.release.yml up --build
```

### App Runner によるデプロイ

#### App Runner

```bash
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
