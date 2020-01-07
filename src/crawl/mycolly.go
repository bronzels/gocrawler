package crawl

import (
	"fmt"
	"log"
	"strings"
	"time"

	colly "github.com/gocolly/colly"
	//lua "github.com/yuin/gopher-lua"
	"github.com/qiniu/qlang"
	_ "github.com/qiniu/qlang/lib/builtin" // 导入 builtin 包

	gocommon "github.com/bronzels/gocommon"
	"github.com/bronzels/gocrawler/pb"
)

var strings_Exports = map[string]interface{}{
	"replacer": strings.NewReplacer,
	"reader":   strings.NewReader,
}

func init() {
	qlang.Import("strings", strings_Exports) // 导入一个自定义的包，叫 strings（和标准库同名）
}

func logPrint(level gocommon.LogLevel, logFlag pb.Logflag, pLogged *strings.Builder, a ...interface{}) {
	if logFlag != pb.Logflag_NO {
		logPrefix := gocommon.LOGprefix(level, 3)
		if logPrefix != "" {
			if pLogged != nil {
				fmt.Fprint(pLogged, a, "\n")
			} else {
				log.Println(a)
			}
		}
	}
}

func getCollector(myAsync bool, userAgent string, parallelism int, randomDelay int, logFlag pb.Logflag, pLogged *strings.Builder) (*colly.Collector, *colly.Collector) {
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
		logPrint(gocommon.LogLevel_INFO, logFlag, pLogged, "r.StatusCode:", r.StatusCode)
	})

	return c, detailCollector
}

func otherSetCollector(c *colly.Collector, detailCollector *colly.Collector, urlToCrawl string, logFlag pb.Logflag, pLogged *strings.Builder) {
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		url := e.Attr("href")
		if strings.HasPrefix(url, "http") {
			logPrint(gocommon.LogLevel_INFO, logFlag, pLogged, "Crawling:", url)
			detailCollector.Visit(url)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		logPrint(gocommon.LogLevel_INFO, logFlag, pLogged, "Visiting:", r.URL.String())
	})

	c.Visit(urlToCrawl)

	c.Wait()
	detailCollector.Wait()
}
