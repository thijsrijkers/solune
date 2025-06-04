package workerlogic

import (
	"fmt"
	"solune/sharding"
)

func Run(port string) {
	shardManager := shard.NewShardManager(port)

	if shardManager.HasActiveShards() {
		fmt.Println("Shards detected. Starting them...")
		shardManager.StartAll()
	} else {
		fmt.Println("No active shards found, starting a new one....")
		s := shard.NewShard("9000")
		go s.Start()
	}

	select {}
}
