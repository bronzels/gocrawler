package test

import (
	"testing"
	//"fmt"
	//"encoding/json"

	"gocrawler/web"
	"gocrawler/web/sinaforex"
)

var envFunc = web.CrawlNews(web.COLLY_ENV_myAsync, web.COLLY_ENV_userAgent, web.COLLY_ENV_parallelism, web.COLLY_ENV_randomDelay)

func Test_Sinaforex_config(t *testing.T) {
	envFunc(sinaforex.NEWS_urlToCrawl, sinaforex.NEWS_newsEntry, sinaforex.NEWS_queriesToExtract, sinaforex.NEWS_scriptPublishedAt)
	//t.Error("error")
}
