# oaro (Out AWS Status RSS OPML)
AWSのステータス情報RSSを取得し、OPMLファイルに出力する。

## 目的
AWSの障害情報RSSのうち、Tokyoのものだけ取り出したい。

## 使い方
oaro.exe -cn Tokyo

## 成果物
`./dist/aws_status_rss_YYYYMMDD.opml`
上記OPMLをFeedlyなりInoreaderなりに読み込ませればOK
