package main

import (
	"flag"
	"log"

	"github.com/notEpsilon/ss-cache/config"
	"github.com/notEpsilon/ss-cache/server"
)

const (
	defaultConfigPath    = "./config.json"
	defaultListenAddress = ""
)

var (
	configPath    = flag.String("config-path", defaultConfigPath, "path to configuration file. defaults to: ./config.json")
	listenAddr    = flag.String("listen", defaultListenAddress, "address to listen on in the form host:port. defaults to: config exposed address")
	shardName     = flag.String("shard-name", "", "the shard that this server instance is responsible for")
	cacheCapacity = flag.Int("cache-capacity", 1000, "capacity of cache before it needs to evict elements")
)

func parseFlags() {
	flag.Parse()

	if *shardName == "" {
		log.Fatalln("[EROR]: please provide -shard-name flag")
	}
}

func main() {
	parseFlags()

	config, err := config.ParseConfig(*configPath)
	if err != nil {
		log.Fatalln("[EROR]: unable to parse configuration, check your config.json file and the -config-path flag")
	}

	srv, err := server.New(*cacheCapacity)
	if err != nil {
		log.Fatalln("[EROR]: unable to create server instance, make sure your -cache-capacity is a positive integer")
	}

	shardIndex := -1
	for _, v := range config.Shards {
		if v.Name == *shardName {
			shardIndex = v.Idx
			break
		}
	}

	if shardIndex == -1 {
		log.Fatalf("[EROR]: no shard with name %q exists\n", *shardName)
	}

	srv.SetShard(&config.Shards[shardIndex])
	srv.SetConfig(&config)

	log.Fatalln(srv.Start(*listenAddr))
}
