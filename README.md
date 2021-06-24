# learn-run-go
Goの実行やデプロイについて学習するためのリポジトリ。  

## 環境情報
```
Mac (2021/06/24)
$ go version
go version go1.16.5 darwin/amd64

$ docker -v      
Docker version 20.10.7, build f0df350

$ docker-compose -v
docker-compose version 1.29.2, build 5becea4c

```

```
Cloud9 (2021/06/24)
$ go version
go version go1.15.12 linux/amd64
$ docker -v
Docker version 20.10.4, build d3cb89e
$ docker-compose -v
bash: docker-compose: command not found

$ sudo curl -L "https://github.com/docker/compose/releases/download/1.29.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
$ sudo chmod +x /usr/local/bin/docker-compose
$ docker-compose -v
docker-compose version 1.29.0, build 07737305
```

## 01.helloworld
`hello world!` という文字列を出力するプログラム。  
ライブラリを使用しない。

### ビルドして実行可能ファイル(バイナリ)を作成する
```
$ pwd
/learn-run-go/01.helloworld

$ go build -o ./bin/main main.go
```

### バイナリを実行する
```
$ ./bin/main
hello world!
```

## 02.modules
現在時刻(GMT)を出力するプログラム。  
ライブラリを利用する。

### ビルド前に、modulesを初期化する
```
$ pwd
/learn-run-go/02.modules

$ go mod init main
go: creating new go.mod: module main
```
このModules初期化コマンドにより `go.mod` ファイルが作成され、ビルドできるようになる。  
`init`の後ろの文字列は何でも良いが、公開したりする場合にはモジュールが存在するリポジトリのパスとするのが良いらしい。  
例 :   
`$ go mod init github.com/ww2or3ww/learn-run-go/02.modules`

### ビルド
```
$ go build -o ./bin/main main.go
go: finding module for package github.com/Songmu/go-httpdate
go: found github.com/Songmu/go-httpdate in github.com/Songmu/go-httpdate v1.0.0
```
ビルド時に依存パッケージがあった場合、それらをインストールしてからビルドすることになる。  
`go.mod`ファイルは更新され、また、`go.sum`ファイルが新規作成される。


### 実行
```
$ ./bin/main
Tue, 22 Jun 2021 09:44:15 GMT
```

## 03.lambda
クエリパラメータとして受け取ったJSONに、`"hello": "world!"` というKey-Valueを加えて返すLambdaファンクション。

### modulesの初期化
```
$ pwd
/learn-run-go/03.lambda/func

$ go mod init func
go: creating new go.mod: module func
```

### 実行(テスト)
```
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
AWSマネジメントコンソール > [AWS Lambda](https://ap-northeast-1.console.aws.amazon.com/lambda/home?region=ap-northeast-1#/functions) から作成する。  
[関数の作成] > [1から作成] > [関数名:learn-run-go] > [ランタイム:Go 1.x] > [関数の作成]

#### ビルド & パッケージング(zip)
```
$ GOOS=linux GOARCH=amd64 go build -o ../bin/hello main.go
$ (cd ../bin && zip -r ../lambda-package.zip *)
```
Goで作成した場合、ハンドラ名がデフォルトで `hello` となっているため、バイナリのファイル名も `hello` としている。  
Lambdaの実行環境に合わせてクロスコンパイルの指定(`GOOS=linux GOARCH=amd64`)をする必要がある。

### zipのアップロード
```
$ aws lambda update-function-code --function-name learn-run-go --zip-file fileb://../lambda-package.zip
```

### Lambdaのテストイベント
```
{
  "queryStringParameters": {"hey": "yo!"}
}
```


## 04.webapp
`hello` と `world` のページをもつWebアプリケーション

```
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
```
$ pwd
/learn-run-go/04.webapp/webapp

$ go build -o ../bin/main main.go && ../bin/main
```

#### Docker(デバッグ用)
```
$ pwd
/learn-run-go/04.webapp

$ docker-compose -f docker-compose.debug.yml up --build
```

#### Docker(リリース用)
```
$ docker-compose -f docker-compose.release.yml up --build
```

### デプロイ
#### AWS Copilotのインストール(Cloud9)
```
$ curl -Lo /usr/local/bin/copilot https://github.com/aws/copilot-cli/releases/latest/download/copilot-linux && chmod +x /usr/local/bin/copilot
$ copilot --version
copilot version: v1.8.1
```
各種環境へのインストールについては以下を参照。  
https://aws.github.io/copilot-cli/ja/docs/getting-started/install/

#### App Runner 
```
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
### VSCodeでデバッグ
`.vscode/launch.json` の `connfigurachionns[0].program` を目的のmainパッケージがあるフォルダに切り替える。  
F5押下により、指定したmainパッケージのプログラムが開始されます。  
