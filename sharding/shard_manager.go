package shard

import (
	"log"
	"strconv"
)

type ShardManager struct {
	shards []*Shard
	active bool
}

// NewShardManager initializes the ShardManager with a single shard based on the given port
func NewShardManager(port string) *ShardManager {
	manager := &ShardManager{}

	// Validate port
	if _, err := strconv.Atoi(port); err != nil {
		log.Fatalf("Invalid port number: %s", port)
	}

	// Create and add shard
	shard := NewShard(port)
	manager.shards = append(manager.shards, shard)

	// Mark as active
	manager.active = true
	return manager
}

// StartAll starts all the shards managed
func (sm *ShardManager) StartAll() {
	for _, shard := range sm.shards {
		go shard.Start()
	}
}

// HasActiveShards returns true if any shard is active (i.e., loaded from folder)
func (sm *ShardManager) HasActiveShards() bool {
	return sm.active
}
