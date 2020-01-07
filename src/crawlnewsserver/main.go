package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"

	pb "github.com/bronzels/gocrawler/pb"
	"github.com/bronzels/gocrawler/src/crawl"
)

const (
	port = ":50058"
)

// crawlnewsserver is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) CrawlNews(ctx context.Context, in *pb.CrawlNewsRequest) (*pb.CrawlNewsReply, error) {
	log.Println("Start to CrawlNews")
	//var envFunc = crawl.CrawlNews(crawl.COLLY_ENV_myAsync, crawl.COLLY_ENV_userAgent, crawl.COLLY_ENV_parallelism, crawl.COLLY_ENV_randomDelay)
	var envFunc = crawl.CrawlNews(in.Myasync, in.UserAgent, int(in.Parallelism), int(in.RandomDelay))
	//envFunc(sinaforex.NEWS_urlToCrawl, sinaforex.NEWS_newsEntry, sinaforex.NEWS_queriesToExtract, sinaforex.NEWS_scriptPublishedAt)
	logged, crawled := envFunc(in.Url2Crawl, in.NewsEntry, in.Queries2Extract, in.ScriptPublishedAt, in.Logflag, nil)
	newsReply := pb.CrawlNewsReply{Crawled: uint32(crawled), Logged: logged}
	return &newsReply, nil
}

func (s *server) CrawlNewsStreamServer(in *pb.CrawlNewsRequest, stream pb.Crawler_CrawlNewsStreamServerServer) error {
	log.Println("Start to CrawlNewsStreamServer")
	var envFunc = crawl.CrawlNews(in.Myasync, in.UserAgent, int(in.Parallelism), int(in.RandomDelay))
	envFunc(in.Url2Crawl, in.NewsEntry, in.Queries2Extract, in.ScriptPublishedAt, in.Logflag, stream)
	return nil
}

func main() {
	srv := new(server)
	c, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{
		grpc.MaxMsgSize(1024 * 1024 * 8),
	}
	s := grpc.NewServer(opts...)
	pb.RegisterCrawlerServer(s, srv)
	log.Println("Server starting...")
	s.Serve(c)
}
