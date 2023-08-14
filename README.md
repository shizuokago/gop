# gop

gop は Go のモジュールの絶対パスを扱うコマンドです。
Modules内にあるパスを出力します。

# install

```
$ go install github.com/shizuokago/gop/_cmd/gop
```

GOBINにコマンドが生成されます

# 実行

```
$ gop {パッケージ名}
```

指定したパッケージの最新バージョンのパスを返します。

```
$ gop golang.org/x/sys
D:\Go\PATH\pkg\mod\golang.org\x\sys@v0.11.0
```

あくまで手元にある最新で、
オンライン上にある最新のバージョンを表示するわけではありません。

## バージョンリストを表示

```
$ gop -list {パッケージ名}
```

```
$ gop -list golang.org/x/sys
golang.org/x/sys:
    0.8.0 => D:\Go\PATH\pkg\mod\golang.org\x\sys@v0.8.0
    0.10.0 => D:\Go\PATH\pkg\mod\golang.org\x\sys@v0.10.0
    0.11.0 => D:\Go\PATH\pkg\mod\golang.org\x\sys@v0.11.0
```

指定したパッケージのバージョンのリスト表示を行います

## go.mod対象のパッケージ(experimental)

gop -p {packagename}

利用しているgo.mod内のパッケージを調べます
パッケージ名を省略すると全リストを表示します

```
$ gop -p golang.org/x/sys
D:\Go\PATH\pkg\mod\golang.org\x\sys@v0.11.0
```

GOMOD が指す go.mod で指定しているパッケージから検索します。

-p に -list を追加した場合

```
$ gop -p -list golang.org/x/sys
golang.org/x/sys:
    0.8.0 => D:\Go\PATH\pkg\mod\golang.org\x\sys@v0.8.0
    0.10.0 => D:\Go\PATH\pkg\mod\golang.org\x\sys@v0.10.0
  * 0.11.0 => D:\Go\PATH\pkg\mod\golang.org\x\sys@v0.11.0
```

存在するバージョンの一覧と利用しているバージョンに*がつきます。

## 全パッケージ

```
$ gop -p
$ gop
```

パッケージを指定しないで処理すると一覧を表示します。
特に -p を利用するとプロジェクトが利用しているパッケージの位置を返します。
※ -p も指定しない場合、存在する全パッケージになります。

```
$ gop -p
D:\Go\PATH\pkg\mod\github.com\fsnotify\fsnotify@v1.6.0
D:\Go\PATH\pkg\mod\golang.org\x\sys@v0.11.0
D:\Go\PATH\pkg\mod\golang.org\x\xerrors@v0.0.0-20220907171357-04be3eba64a2
```

-list を行うとそれぞれのパッケージのバージョンの位置を返します。

# Issue

バージョンが複数ある場合に
指定なしを 0 or 1 のバージョンに固定する
