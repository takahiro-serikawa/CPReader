# ComProReader
競プロ用入力ルーチン go言語用

競技プログラミング、以下「競プロ」

go言語の fmt.Scan() は早くないのでその代替

## 使い方

type Reader 宣言以降を提出ファイルにコピペ
main()の前でも後ろでもどちらでもよい

    cr := NewReader(os.Stdin)   // インスタンス生成
    N := cr.Int()               // 整数読み込み
    S := cr.Line()              // 行読み込み (改行区切りstring)
    S := cr.Word()              // 単語読み込み (改行または空白string)
    C := cr.Byte()              // 一文字読み込み (byte)

    N := cr.Int64()
    N := cr.Int32()
    N := cr.Uint()
    N := cr.Uint64()
    N := cr.Uint32()
注: キャストの時、範囲チェックなし
競プロではふつう入力値の妥当性のチェックは省くので問題はないと考える。

    F := cr.Float64()

cr.Int() は改行とスペースを区別しない

    123 456

と、

    123
    456

は、どちらも

    A, B := cr.Int(), cr.Int()  // A->123, B->456

で読める。

数値として解釈出来ない文字にあたったときはpanic

## 不要な関数を削除して提出する場合

```
                             ↓ここから
type Reader
 |-func NewReader(r *os.File) *Reader
 |
 |-func Byte() byte
 | |-func Uint64() uint64
 |   |-func Int64() int64
 |   | |-func Int() int
                             ↑Int()しか使わないのであれば、ここまでをコピペ
 |   | |-func Int32() int32
 |   |-func Uint() uint
 |   |-func Uint32() uint32
 |-func Line() string
 |-func Word() string
 |-func Slice() string
```

 ## メリット
fmt.Scan()より早く、
bufio.Scanner より簡便

cr.Int()は fmt.Scan()よりかなり早く、bufio.ScannerのText()より倍ほど早い。
bufio.Scanner のバッファサイズの設定が不要

 ## 実装
バッファサイズ　デフォルト 64KB

 ## 未実装&既知の問題
- Word()　未テスト

- EOFでpanic()　要検討; 今のところは仕様

- Float64() 最小値テスト通らないので暫定的に ParseFloat()版で差し替え

- この後の改良で互換性はあまり考えない。「競プロ」では単一ファイルにすべてを記述し、提出してそれっきり。
型名 Reader は優先順位高めで変えるかもしれない。ちょっと行儀が悪いような気がしている。

- 複数スペースでの区切り　対応予定なし
- 位取り文字 _ など   対応予定なし

- サンプルコードでインスタンス名を cr としているのは、
競プロでは、問題文中にアルファベット一文字の変数が頻出するのでそれとバッティングする可能性を避けている。

- Win10機で調整しているので、競プロサーバとは性能が違う可能性がある
