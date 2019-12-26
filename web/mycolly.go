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
		fmt.Println(r.StatusCode)
	})

	return c, detailCollector
}

func otherSetCollector(c *colly.Collector, detailCollector *colly.Collector, urlToCrawl string) {
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		url := e.Attr("href")
		detailCollector.Visit(url)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
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
	ql := qlang.New()

	detailCollector.OnHTML(newsEntry, func(e *colly.HTMLElement) {
		var extracteds [len(queriesToExtract)]string
		for index, query := range queriesToExtract {
			if query != "" {
				extracteds[index] = e.ChildText(query)
			} else {
				extracteds[index] = ""
			}
		}

		n := News{}

		url := e.Request.URL.String()
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
		tPublishedAtStr := extracteds[NEWS_arrindex_publishedAt]
		err := ql.SafeEval(fmt.Sprintf(scriptPublishedAt, tPublishedAtStr))
		if err != nil {
			log.Fatal(err)
			return
		}

		scriptedPublishedAtStr := fmt.Sprintf("%s", ql.Var(SCRIPT_VAR_NAME_RET))
		fmt.Println("scriptedPublishedAtStr:", scriptedPublishedAtStr)
		timeTemplate := "2006-01-02 15:04"
		n.PublishedAt, _ = time.ParseInLocation(timeTemplate, scriptedPublishedAtStr, time.Local)

		n.Contents = extracteds[NEWS_arrindex_contents]
		n.CrawledAt = time.Now()
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
