package target

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"
)

type Quoter struct {
	URL         string
	URLid       string
	QuoteeURL   string
	QuoteeURLid string
}

type News struct {
	EntryURL    string
	URL         string
	URLid       string
	Title       string
	PublishedAt time.Time
	Contents    string
	CrawledAt   time.Time
}

func MyMd5(s string) string {
	m := md5.New()
	m.Write([]byte(s))
	return fmt.Sprintf(hex.EncodeToString(m.Sum(nil)))
	/*
	   m := md5.Sum([]byte (s))
	   return fmt.Sprintf(hex.EncodeToString(m[:]))
	*/
	/*
	   m := md5.Sum([]byte(s))
	   fmt.Printf("%x", m)
	   return fmt.Sprintf()
	*/
	/*
	   m := md5.New()
	   io.WriteString(m, s)
	   return fmt.Sprintf(hex.EncodeToString(m.Sum(nil)))
	*/
}
