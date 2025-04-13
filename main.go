package main

import (
	"solune/sharding"
)

func main() {
	ports := []string{"9000"}

	for _, port := range ports {
		s := shard.NewShard(port)
		go s.Start()
	}

	select {} // block forever
}

