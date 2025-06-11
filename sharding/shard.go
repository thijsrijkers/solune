package shard

import (
	"log"
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
		manager: store.NewDataStoreManager(string(port)),
	}
}

func (s *Shard) Start() {
	log.Printf("Starting shard on port %s\n", s.Port)
	tcp.StartServer(s.Port, s.manager)
}

