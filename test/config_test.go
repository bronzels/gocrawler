package test

import (
	"testing"
	//"fmt"
	//"encoding/json"

	"github.com/bronzels/gocrawler/src/crawl"
	"github.com/bronzels/gocrawler/src/crawl/sinaforex"
)

var envFunc = crawl.CrawlNews(crawl.COLLY_ENV_myAsync, crawl.COLLY_ENV_userAgent, crawl.COLLY_ENV_parallelism, crawl.COLLY_ENV_randomDelay)

func Test_Sinaforex_config(t *testing.T) {
	envFunc(sinaforex.NEWS_urlToCrawl, sinaforex.NEWS_newsEntry, sinaforex.NEWS_queriesToExtract, sinaforex.NEWS_scriptPublishedAt)
	//t.Error("error")
}
