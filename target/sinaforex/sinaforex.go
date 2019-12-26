package sinaforex

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly"

	"gocrawler/target"
)

func Crawl(url2crawl string) {
	c := colly.NewCollector(
		colly.Async(true),
		colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64)    AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36"),
	)

	c.Limit(&colly.LimitRule{
		Parallelism: 2,
		RandomDelay: 5 * time.Second,
	})

	detailCollector := c.Clone()

	detailCollector.OnResponse(func(r *colly.Response) {
		fmt.Println(r.StatusCode)
	})

	detailCollector.OnHTML("div[class^=main-content]", func(e *colly.HTMLElement) {
		//detailCollector.OnHTML("div", func(e *colly.HTMLElement) {
		//topclass := e.Attr("class")
		//if strings.HasSuffix(topclass, "main-content") {
		/*
			dataSource1 := e.DOM.Find("div[class=top-bar-wrap]")
			dataSource2 := dataSource1.Find("div[class^='top-bar ani']")
			dataSource3 := dataSource2.Find("div[class='top-bar-inner clearfix']")
			dataSource4 := dataSource3.Find("div[class=date-source]")
				dataSource := dataSource4//e.DOM.Find("div[class=top-bar-wrap]>div[class=top-bar ani]>div[class=top-bar-inner clearfix]>div[class=date-source]")
		*/
		dataSource := e.DOM.Find("div[class=top-bar-wrap]>div[class^='top-bar ani']>div[class='top-bar-inner clearfix']>div[class=date-source]")
		referredAHref := dataSource.Find("a[href]")
		url := e.Request.URL.String()
		//跳转到微信qq页面
		//https://finance.sina.com.cn/money/forex/forexroll/2019-12-26/doc-iihnzhfz8326034.shtml
		//微信qq页面再次跳转页面不可用
		if referredAHref.Nodes != nil {
			q := target.Quoter{}
			q.URL = url
			q.URLid = target.MyMd5(url)
			q.QuoteeURL, _ = referredAHref.Attr("href")
			q.QuoteeURLid = target.MyMd5(q.QuoteeURL)
			log.Println(q)
		} else {
			n := target.News{}
			n.URL = url
			n.URLid = target.MyMd5(url)
			n.Title = e.ChildText("h1[class=main-title]")
			tPublishedAtStr := dataSource.Find("span[class=date]").Text()
			tPublishedAtStr = strings.ReplaceAll(tPublishedAtStr, "年", "-")
			tPublishedAtStr = strings.ReplaceAll(tPublishedAtStr, "月", "-")
			tPublishedAtStr = strings.ReplaceAll(tPublishedAtStr, "日", "")
			//timeTemplate := "2019年12月25日 15:13"
			//timeTemplate := "2019-12-25 15:13"
			timeTemplate := "2006-01-02 15:04"
			n.PublishedAt, _ = time.ParseInLocation(timeTemplate, tPublishedAtStr, time.Local)
			/*
				dataSource1 := e.DOM.Find("div[class^='article-content clearfix']")
				dataSource2 := dataSource1.Find("div[class=article-content-left]")
				dataSource3 := dataSource2.Find("div[class=article][id=artibody]")
				n.Contents = dataSource3.Find("p").Text()//e.ChildText("div[class^='article-content clearfix']>div[class=article-content-left]>div[class=article,id=artibody]>p")
			*/
			n.Contents = e.ChildText("div[class^='article-content clearfix']>div[class=article-content-left]>div[class=article][id=artibody]>p")
			n.CrawledAt = time.Now()
			log.Println(n)
		}
		//}
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		url := e.Attr("href")
		/*
			if strings.Compare(e.Attr("target"), "_blank") == 0 && strings.HasSuffix(url, "shtml") {
				detailCollector.Visit(url)
			}
		*/
		detailCollector.Visit(url)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.Visit(url2crawl)

	c.Wait()
	detailCollector.Wait()
}
