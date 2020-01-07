go get -u github.com/kardianos/govendor

cd $GOPATH/src/github.com/bronzels/gocrawler
rm -rf vendor
govendor init

govendor fetch github.com/gocolly/colly@=v1.2.0
#v2.0.0
cd vendor/github.com/gocolly
#mv colly colly.master
#git clone --branch v1.2.0 https://github.com/gocolly/colly.git
cd colly
go build

cd $GOPATH/src/github.com/bronzels/gocrawler
govendor fetch github.com/qiniu/qlang@=v5.1.0
cd vendor/github.com/qiniu
mv qlang qlang.govendor
git clone --branch v5.1.0 https://github.com/qiniu/qlang.git
cd qlang
go build

cd $GOPATH/src/github.com
mkdir golang;cd golang
git clone --branch v1.3.2  https://github.com/golang/protobuf.git
cd $GOPATH/src/github.com/golang/protobuf
cd proto
go build
go install
cd $GOPATH/src/github.com/golang/protobuf
cd protoc-gen-go
go build
go install

cd $GOPATH/src/github.com/bronzels/gocrawler
govendor fetch google.golang.org/grpc
#@=v1.26.0
#govendor add +external

cd $GOPATH/src/github.com/bronzels/gocrawler
mkdir -p vendor/github.com/bronzels
make all