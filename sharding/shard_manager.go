package shard

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
)

type ShardManager struct {
	shards []*Shard
	active bool
}

// NewShardManager scans the db directory and initializes shards based on folder names
func NewShardManager(baseDir string) *ShardManager {
	manager := &ShardManager{}

	entries, err := ioutil.ReadDir(baseDir)
	if err != nil {
		log.Fatalf("Failed to read base directory %s: %v", baseDir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			port := entry.Name()
			if _, err := strconv.Atoi(port); err != nil {
				log.Printf("Skipping invalid port folder: %s", port)
				continue
			}

			shard := NewShard(port)
			manager.shards = append(manager.shards, shard)
		}
	}

	manager.active = len(manager.shards) > 0
	return manager
}

// StartAll starts all the shards managed
func (sm *ShardManager) StartAll() {
	for _, shard := range sm.shards {
		go shard.Start()
	}
	fmt.Printf("Started %d shard(s).\n", len(sm.shards))
}

// HasActiveShards returns true if any shard is active (i.e., loaded from folder)
func (sm *ShardManager) HasActiveShards() bool {
	return sm.active
}
