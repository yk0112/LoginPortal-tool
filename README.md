# LoginPortal-tool
東工大ポータルに自動ログインするCLIツール

# Installation
```
$ git clone https://github.com/yk0112/LoginPortal-tool.git
```

# Usage

## WebDriverのインストール
使用したいブラウザのドライバーをインストールし、パスを通しておく

## 認証情報の設定
アカウント名, パスワード, マトリクス表などを以下のコマンドで設定する.
```
$ go run main.go init
```

## ログイン開始
以下のコマンドを実行することで、ブラウザが開き自動でポータルサイトにログインする.
```
$ go run main.go login
```
# Command 
```
`help`     - ツールの使用方法や機能に関する情報を表示する
`init`     - 東工大ポータルの認証情報を設定する
`login`     - 自動的にポータルサイトを開いてログインする
```
# Dependencies
- Go 1.22.0 or greater
- godotenv v1.5.1
- cobra v1.8.0
- agouti v3.0.0
