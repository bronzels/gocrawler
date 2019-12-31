package test

import (
	"testing"
	//"fmt"
	//"encoding/json"

	"github.com/bronzels/gocrawler/src/crawl/sinaforex"
)

func Test_Sinaforex_native(t *testing.T) {
	sinaforex.Crawl(sinaforex.URL)
	//t.Error("error")
}
