/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package main

import (
	"flag"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"

	pb "github.com/bronzels/gocrawler/pb"
	"github.com/bronzels/gocrawler/src/crawl"
	"github.com/bronzels/gocrawler/src/crawl/sinaforex"
)

var (
	pHost = flag.String("host", "localhost", "host of grpc server")
	pPort = flag.Int("port", 30092, "port of grpc server")

	pStreamMode = flag.Bool("stream", false, "request in stream mode or not")
)

func main() {
	flag.Parse()

	// Set up a connection to the crawlnewsserver.
	//*pPort = 50058
	//*pStreamMode = true
	conn, err := grpc.Dial(fmt.Sprintf("%v:%v", *pHost, *pPort), grpc.WithInsecure(), grpc.WithMaxMsgSize(1024*1024*8))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewCrawlerClient(conn)

	// Contact the crawlnewsserver and print out its response.
	rq := pb.CrawlNewsRequest{
		Myasync:           crawl.COLLY_ENV_myAsync,
		UserAgent:         crawl.COLLY_ENV_userAgent,
		Parallelism:       crawl.COLLY_ENV_parallelism,
		RandomDelay:       crawl.COLLY_ENV_randomDelay,
		Url2Crawl:         sinaforex.URL,
		NewsEntry:         sinaforex.NEWS_newsEntry,
		Queries2Extract:   sinaforex.NEWS_queriesToExtract,
		ScriptPublishedAt: sinaforex.NEWS_scriptPublishedAt,
		Logflag:           pb.Logflag_ACCUMULATED,
		//Logflag: pb.Logflag_ORDINARY,
	}

	if !*pStreamMode {
		log.Println("call CrawlNews")
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		r, err := c.CrawlNews(ctx, &rq)
		if err != nil {
			log.Fatalf("could not crawl: %v", err)
		}
		log.Printf("crawled: %d,\nlogged: \n%s\n", r.Crawled, r.Logged)
	} else {
		log.Println("call CrawlNewsStreamServer")
		streamCtx, streamCancel := context.WithTimeout(context.Background(), 200*time.Second)
		defer streamCancel()
		stream, err := c.CrawlNewsStreamServer(streamCtx, &rq)
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			//log.Println("res:", res)
			log.Println(
				"res.EntryUrl:", res.EntryUrl,
				"res.Url:", res.Url,
				"res.UrlId:", res.UrlId,
				"res.QuoteeUrl:", res.QuoteeUrl,
				"res.QuoteeUrlId:", res.QuoteeUrlId,
				"res.PublishedAt:", res.PublishedAt,
				"res.CrawledAt:", res.CrawledAt,
				"res.Title:", res.Title,
				"res.Contents:", res.Contents,
			)
		}
	}

}
