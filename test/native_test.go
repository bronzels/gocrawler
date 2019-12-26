package test

import (
	"testing"
	//"fmt"
	//"encoding/json"

	"gocrawler/web/sinaforex"
)

func Test_Sinaforex_native(t *testing.T) {
	sinaforex.Crawl(sinaforex.URL)
	//t.Error("error")
}
