package web

import (
	"fmt"
	"log"
	"strings"
	"time"

	colly "github.com/gocolly/colly"
	//lua "github.com/yuin/gopher-lua"
	"github.com/qiniu/qlang"
	_ "github.com/qiniu/qlang/lib/builtin" // 导入 builtin 包
)

var strings_Exports = map[string]interface{}{
	"replacer": strings.NewReplacer,
	"reader":   strings.NewReader,
}

func getCollector(myAsync bool, userAgent string, parallelism int, randomDelay int) (*colly.Collector, *colly.Collector) {
	c := colly.NewCollector(
		colly.Async(myAsync),
		colly.UserAgent(userAgent),
	)

	c.Limit(&colly.LimitRule{
		Parallelism: parallelism,
		RandomDelay: time.Duration(randomDelay) * time.Second,
	})

	detailCollector := c.Clone()

	detailCollector.OnResponse(func(r *colly.Response) {
		log.Println(r.StatusCode)
	})

	return c, detailCollector
}

func otherSetCollector(c *colly.Collector, detailCollector *colly.Collector, urlToCrawl string) {
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		url := e.Attr("href")
		if strings.HasPrefix(url, "http") {
			log.Println("Crawling:", url)
			detailCollector.Visit(url)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting:", r.URL.String())
	})

	c.Visit(urlToCrawl)

	c.Wait()
	detailCollector.Wait()
}

func crawlNews(myAsync bool, userAgent string, parallelism int, randomDelay int,

	urlToCrawl string,

	newsEntry string,
	queriesToExtract [4]string,
	scriptPublishedAt string,
) {

	c, detailCollector := getCollector(myAsync, userAgent, parallelism, randomDelay)
	/*
		// 创建一个lua解释器实例
		l := lua.NewState()
		defer l.Close()
	*/

	qlang.Import("strings", strings_Exports) // 导入一个自定义的包，叫 strings（和标准库同名）

	detailCollector.OnHTML(newsEntry, func(e *colly.HTMLElement) {
		url := e.Request.URL.String()

		var extracteds [len(queriesToExtract)]string
		for index, query := range queriesToExtract {
			var extracted = ""
			if query != "" {
				extracted = e.ChildText(query)
			} else {
				extracted = ""
			}
			if (index != NEWS_arrindex_publishedAt && index != NEWS_arrindex_referredAHref) && extracted == "" {
				log.Println("empty index:", index)
				return
			}
			extracteds[index] = extracted
		}

		n := News{}

		n.EntryURL = urlToCrawl

		n.URL = url
		n.URLid = MyMd5(url)

		//跳转到微信qq页面
		//https://finance.sina.com.cn/money/forex/forexroll/2019-12-26/doc-iihnzhfz8326034.shtml
		//微信qq页面再次跳转页面不可用
		referredAHref := extracteds[NEWS_arrindex_referredAHref]
		if referredAHref != "" {
			n.QuoteeURL = referredAHref
			n.QuoteeURLid = MyMd5(n.QuoteeURL)
		}

		n.Title = extracteds[NEWS_arrindex_title]

		nowTime := time.Now()

		tPublishedAtStr := extracteds[NEWS_arrindex_publishedAt]
		log.Println("tPublishedAtStr:", tPublishedAtStr)
		if tPublishedAtStr != "" {
			ql := qlang.New()

			expr := fmt.Sprintf(scriptPublishedAt, tPublishedAtStr)
			err := ql.SafeEval(expr)
			if err != nil {
				log.Fatal(err)
				return
			}

			scriptedPublishedAtStr := fmt.Sprintf("%s", ql.Var(SCRIPT_VAR_NAME_RET))
			log.Println("scriptedPublishedAtStr:", scriptedPublishedAtStr)
			timeTemplate := "2006-01-02 15:04"
			n.PublishedAt, _ = time.ParseInLocation(timeTemplate, scriptedPublishedAtStr, time.Local)
		} else {
			n.PublishedAt = nowTime
		}

		n.Contents = extracteds[NEWS_arrindex_contents]
		n.CrawledAt = nowTime
		log.Println(n)
		//}
	})

	otherSetCollector(c, detailCollector, urlToCrawl)
}

func CrawlNews(myAsync bool, userAgent string, parallelism int, randomDelay int) func(string, string, [4]string, string) {
	return func(
		urlToCrawl string,
		newsEntry string,
		queriesToExtract [4]string,
		scriptPublishedAt string,
	) {
		crawlNews(myAsync, userAgent, parallelism, randomDelay,

			urlToCrawl,

			newsEntry,
			queriesToExtract,
			scriptPublishedAt)
	}
}
