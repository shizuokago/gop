# gop

gop は Go のモジュールの絶対パスを扱うコマンドです。
Modules内にあるパスを出力します。

# 似たような所作

```
$ go list -m -f {{.Dir}} all
```

でほとんど同じことができます。

# install

```
$ go install github.com/shizuokago/gop/_cmd/gop
```

GOBINにコマンドが生成されます

# 実行

```
$ gop
D:\Go\PATH\pkg\mod\golang.org\x\mod@v0.12.0
D:\Go\PATH\pkg\mod\golang.org\x\net@v0.14.0
D:\Go\PATH\pkg\mod\golang.org\x\sync@v0.3.0
---省略
```

一覧はgo.modの指定に依存するため、go.modがない位置ではエラーになります。

## パッケージ指定

```
$ gop {パッケージ名}
```

指定したパッケージのバージョンのパスを返します。

```
$ gop golang.org/x/sys
D:\Go\PATH\pkg\mod\golang.org\x\sys@v0.11.0
```

## バージョンリストを表示

```
$ gop -list {パッケージ名}
```

```
$ gop -list golang.org/x/sys
golang.org/x/sys:
    0.8.0 => D:\Go\PATH\pkg\mod\golang.org\x\sys@v0.8.0
    0.10.0 => D:\Go\PATH\pkg\mod\golang.org\x\sys@v0.10.0
  * 0.11.0 => D:\Go\PATH\pkg\mod\golang.org\x\sys@v0.11.0
```

指定したパッケージのバージョンのリスト表示を行います
パッケージを指定せずに、-listを行うとすべてのパッケージのバージョンを返します。

## 全指定(-all)

```
$ gop -all {packagename}
```

端末にあるパッケージが表示できます。
あまり利用用途はないと思いますが

```
$ gop -all -list
```

で端末内にあるすべてのパッケージのすべてのバージョンを表示できます。

# Issue

- バージョン指定周り
- ブランチ指定周り
- replace

