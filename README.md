# 学園祭用Webアプリ

## 機能一覧

### スタンプラリー:
各店舗に置いてあるNFCタグ(またはQRコード)を読み込んでスタンプを集め、ビンゴ(3*3)を行う。  
ビンゴが揃ったらサークルの教室まで来て、景品と交換(再度、遊べるようにリセット)  

### マップ:
どの階のどの場所でどのクラス/サークルが何を運営しているかが視覚的に分かる  

---

## 開発環境

| 項目 | 内容 |
|---|---|  
|IDE|VSCode|  
|フロントエンド|Nodejs(25.5.0)|  
|バックエンド|Go(1.26.3)|  


## ディレクトリ構成

backendのディレクトリ構成は以下の通りです。
``` cmd
└───go_back
    ├───Controller
    └───model
```

フロントエンドのディレクトリ構成は以下の通りです。
``` cmd
└───react_front
    ├───public
    └───src
        ├───assets
```

## 実行方法  
### バックエンド
go をインストール 
``` cmd
scoop install go
```

go_back ディレクトリに移動して、以下のコマンドを実行
``` cmd
go run main.go
```

### フロントエンド
nodejs をインストール
``` cmd
scoop install nodejs
```

stamp_rally_front ディレクトリに移動して、以下のコマンドを実行
``` cmd
npm install
npm run dev
```



