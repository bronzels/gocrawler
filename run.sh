#clone to $GOPATH/src/github.com/bronzels
cd gocrawler
rm -rf vendor/github.com/bronzels
mkdir -p vendor/github.com/bronzels
cp -R ../gocommon vendor/github.com/bronzels
make all
docker build -t bronzels/gocrawler:latest .
apply.sh
#helloworldclient
crawlnewsclient
