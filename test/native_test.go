package test

import (
	"testing"
	//"fmt"
	//"encoding/json"

	"gocrawler/target/sinaforex"
)

func Test_Sinaforex(t *testing.T) {
	sinaforex.Crawl(URL_SINAFOREX)
	//t.Error("error")
}
