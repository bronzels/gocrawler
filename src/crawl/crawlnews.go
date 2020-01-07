package crawl

import (
	"fmt"
	"strings"
	"time"

	colly "github.com/gocolly/colly"
	//lua "github.com/yuin/gopher-lua"
	"github.com/qiniu/qlang"
	_ "github.com/qiniu/qlang/lib/builtin" // 导入 builtin 包

	gocommon "github.com/bronzels/gocommon"
	"github.com/bronzels/gocrawler/pb"
)

func crawlNews(myAsync bool, userAgent string, parallelism int, randomDelay int,

	urlToCrawl string,

	newsEntry string,
	queriesToExtract []string,
	scriptPublishedAt string,
	logFlag pb.Logflag,
	stream pb.Crawler_CrawlNewsStreamServerServer,
) (string, int) {
	var logged strings.Builder
	var pLogged *strings.Builder
	if logFlag == pb.Logflag_ACCUMULATED {
		pLogged = &logged
	} else {
		pLogged = nil
	}

	c, detailCollector := getCollector(myAsync, userAgent, parallelism, randomDelay, logFlag, pLogged)
	/*
		// 创建一个lua解释器实例
		l := lua.NewState()
		defer l.Close()
	*/

	var crawledPage = 0

	detailCollector.OnHTML(newsEntry, func(e *colly.HTMLElement) {
		url := e.Request.URL.String()
		var extracteds = make(map[int]string)
		for index, query := range queriesToExtract {
			var extracted = ""
			if query != "" {
				extracted = e.ChildText(query)
			} else {
				extracted = ""
			}
			if (index != NEWS_arrindex_publishedAt && index != NEWS_arrindex_referredAHref) && extracted == "" {
				logPrint(gocommon.LogLevel_INFO, logFlag, pLogged, "empty index:", index)
				return
			}
			extracteds[index] = extracted
		}

		//n := News{}
		n := pb.CrawlNewsStreamServerReply{}

		//n.EntryURL = urlToCrawl
		n.EntryUrl = urlToCrawl

		//n.URL = url
		n.Url = url
		//n.URLid = MyMd5(url)
		n.UrlId = MyMd5(url)

		//跳转到微信qq页面
		//https://finance.sina.com.cn/money/forex/forexroll/2019-12-26/doc-iihnzhfz8326034.shtml
		//微信qq页面再次跳转页面不可用
		referredAHref := extracteds[NEWS_arrindex_referredAHref]
		if referredAHref != "" {
			//n.QuoteeURL = referredAHref
			n.QuoteeUrl = referredAHref
			//n.QuoteeURLid = MyMd5(n.QuoteeURL)
			n.QuoteeUrlId = MyMd5(n.QuoteeUrl)
		}

		n.Title = extracteds[NEWS_arrindex_title]

		nowTime := time.Now().Unix()

		tPublishedAtStr := extracteds[NEWS_arrindex_publishedAt]
		logPrint(gocommon.LogLevel_INFO, logFlag, pLogged, "tPublishedAtStr:", tPublishedAtStr)
		if tPublishedAtStr != "" {
			ql := qlang.New()

			expr := fmt.Sprintf(scriptPublishedAt, tPublishedAtStr)
			err := ql.SafeEval(expr)
			if err != nil {
				logPrint(gocommon.LogLevel_ERROR, logFlag, pLogged, "qlang eval err:", err.Error())
				return
			}

			scriptedPublishedAtStr := fmt.Sprintf("%s", ql.Var(SCRIPT_VAR_NAME_RET))
			logPrint(gocommon.LogLevel_INFO, logFlag, pLogged, "scriptedPublishedAtStr:", scriptedPublishedAtStr)
			timeTemplate := "2006-01-02 15:04"
			//n.PublishedAt, _ = time.ParseInLocation(timeTemplate, scriptedPublishedAtStr, time.Local)
			publishedAt, _ := time.ParseInLocation(timeTemplate, scriptedPublishedAtStr, time.Local)
			n.PublishedAt = publishedAt.Unix()
		} else {
			n.PublishedAt = nowTime
		}

		n.Contents = extracteds[NEWS_arrindex_contents]
		n.CrawledAt = nowTime
		crawledPage += 1
		if stream != nil {
			stream.Send(&n)
		} else {
			//logPrint(gocommon.LogLevel_INFO, logFlag, pLogged, "n:", n)
			logPrint(gocommon.LogLevel_INFO, logFlag, pLogged,
				"n.EntryUrl:", n.EntryUrl,
				"n.Url:", n.Url,
				"n.UrlId:", n.UrlId,
				"n.QuoteeUrl:", n.QuoteeUrl,
				"n.QuoteeUrlId:", n.QuoteeUrlId,
				"n.PublishedAt:", n.PublishedAt,
				"n.CrawledAt:", n.CrawledAt,
				"n.Title:", n.Title,
				"n.Contents:", n.Contents,
			)
		}
		//fmt.Fprint(&logged, "n:", n)//log.Println(n)
		//}
	})

	otherSetCollector(c, detailCollector, urlToCrawl, logFlag, pLogged)

	return logged.String(), crawledPage
}

func CrawlNews(myAsync bool, userAgent string, parallelism int, randomDelay int) func(string, string, []string, string, pb.Logflag, pb.Crawler_CrawlNewsStreamServerServer) (string, int) {
	return func(
		urlToCrawl string,
		newsEntry string,
		queriesToExtract []string,
		scriptPublishedAt string,
		logFlag pb.Logflag,
		stream pb.Crawler_CrawlNewsStreamServerServer,
	) (string, int) {
		return crawlNews(myAsync, userAgent, parallelism, randomDelay,

			urlToCrawl,

			newsEntry,
			queriesToExtract,
			scriptPublishedAt,
			logFlag,
			stream)
	}
}
