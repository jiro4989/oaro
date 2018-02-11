package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "oaro"
	app.Usage = "AWSの障害情報RSSをOPML形式で出力します。"
	app.Version = "1.0.0"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "country-name,cn",
			Usage: "絞り込む国名",
		},
		cli.StringFlag{
			Name:  "url,u",
			Value: "http://status.aws.amazon.com/",
			Usage: "HTML取得元のAWSの障害情報のURL",
		},
	}

	app.Action = func(c *cli.Context) error {
		log.Println("start:", app.Name)

		// コマンドライン引数の取得
		countryName := c.String("cn")
		url := c.String("u")

		// goqueryで取得してきたopml作成に必要な情報のリスト
		records := make([]map[string]string, 0)

		// Webサイトからデータ取得開始
		doc, err := goquery.NewDocument(url)
		if err != nil {
			log.Println(err)
			return err
		}
		doc.Find("tbody").Eq(7).Find("tr").Each(func(_ int, s *goquery.Selection) {
			record := make(map[string]string)

			// サービス名の取得
			td := s.Find("td.bb.top.pad8")
			td.Each(func(_ int, s *goquery.Selection) {
				text := s.Text()

				if countryName == "" {
					record["service"] = text
					return
				}

				if strings.Contains(text, "("+countryName+")") {
					record["service"] = text
				}
			})

			// RSSのリンクの取得
			a := s.Find("a")
			a.Each(func(_ int, s *goquery.Selection) {
				t, _ := s.Attr("href")
				record["rss"] = t
			})

			if record["service"] == "" {
				return
			}

			if record["rss"] == "" {
				return
			}

			records = append(records, record)
		})

		opmlFormat := `<?xml version="1.0" encoding="UTF-8"?>
<opml version="1.0">
	<head>
	<title>AWS Service Status(Asia Pacific)</title>
	</head>
	<body>
		$body
	</body>
</opml>
`
		// ファイル出力先ディレクトリ
		os.Mkdir("dist", os.ModeDir)
		today := time.Now().Format("20060102")

		fileName := "dist/aws_status_rss_" + today + ".opml"
		outline := convertOutline(records)
		opml := strings.Replace(opmlFormat, "$body", outline, -1)
		b := []byte(opml)
		if err := ioutil.WriteFile(fileName, b, os.ModePerm); err != nil {
			log.Println(err)
			return err
		}

		log.Println(fmt.Sprintf("[success]create: filename=<%s>, filesize=<%d>", fileName, len(b)))

		log.Println("complete:" + app.Name)
		return nil
	}

	app.Run(os.Args)
}

// マップデータをOPMLに埋め込み要素文字列に変換する
func convertOutline(records []map[string]string) string {
	outlineFormat := `<outline type="rss" text="$title" title="$title" xmlUrl="http://status.aws.amazon.com$rss" htmlUrl="http://status.aws.amazon.com/" />`

	outlines := make([]string, len(records))
	for i, record := range records {
		outline := strings.Replace(outlineFormat, "$title", record["service"], -1)
		outline = strings.Replace(outline, "$rss", record["rss"], -1)
		outlines[i] = outline
	}
	outline := strings.Join(outlines, "\n")
	return outline
}
