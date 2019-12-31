package main

import (
	"flag"
	"path/filepath"
	"strings"

	"os"

	"github.com/chrislusf/gleam/flow"
	"github.com/chrislusf/gleam/gio"
	"github.com/chrislusf/gleam/plugins/kafka"
)

var (
	brokers = flag.String("brokers", "beta-hbase02:9092,beta-hbase03:9092,beta-hbase04:9092", "a list of comma separated broker:port")
	topic   = flag.String("topic", "testgleam", "the topic name")
	group   = flag.String("group", filepath.Base(os.Args[0]), "the consumer group name")
	timeout = flag.Int("timeout", 30, "the number of seconds for timeout connections")

	Capitalize = gio.RegisterMapper(capitalize)
)

func main() {

	gio.Init()
	flag.Parse()

	brokerList := strings.Split(*brokers, ",")

	k := kafka.New(brokerList, *topic, *group)
	k.TimeoutSeconds = *timeout

	f := flow.New("kafka "+*topic).Read(k).Map("capitalize", Capitalize).Printlnf("%s")

	f.Run()

}

func capitalize(row []interface{}) error {
	line := gio.ToString(row[0])
	gio.Emit(strings.ToUpper(line))
	return nil
}
