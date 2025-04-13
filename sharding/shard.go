package shard

import (
	"fmt"
	"solune/store"
	"solune/tcp"
)

type Shard struct {
	Port    string
	manager *store.DataStoreManager
}

func NewShard(port string) *Shard {
	return &Shard{
		Port:    port,
		manager: store.NewDataStoreManager(),
	}
}

func (s *Shard) Start() {
	fmt.Printf("Starting shard on port %s\n", s.Port)
	tcp.StartServer(s.Port, s.manager)
}

