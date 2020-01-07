package sinaforex

const (
	NEWS_urlToCrawl = URL
	NEWS_newsEntry  = "div[class^=main-content]"
)

var NEWS_queriesToExtract = []string{
	"h1[class=main-title]", //NEWS_arrindex_title = 0
	"div[class=top-bar-wrap]>div[class^='top-bar ani']>div[class='top-bar-inner clearfix']>div[class=date-source]>span[class=date]", //NEWS_arrindex_publishedAt    = 1
	"div[class=top-bar-wrap]>div[class^='top-bar ani']>div[class='top-bar-inner clearfix']>div[class=date-source]>a[href]",          //NEWS_arrindex_referredAHref  = 2
	"div[class^='article-content clearfix']>div[class=article-content-left]>div[class=article][id=artibody]>p",                      //NEWS_arrindex_contents       = 3
}

const (
	NEWS_scriptPublishedAt = "ret = strings.replacer(\"年\", \"-\", \"月\", \"-\", \"日\", \"\").replace(\"%s\")"
)
