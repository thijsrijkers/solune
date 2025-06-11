package workerlogic

import (
	"log"
	"solune/sharding"
)

func Run(port string) {
	shardManager := shard.NewShardManager(port)

	if shardManager.HasActiveShards() {
		log.Println("Shards detected. Starting them...")
		shardManager.StartAll()
	} else {
		log.Println("No active shards found, starting a new one....")
		s := shard.NewShard("9000")
		go s.Start()
	}

	select {}
}
