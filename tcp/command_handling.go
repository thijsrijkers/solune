package tcp

import (
	"fmt"
	"strconv"
)

func (s *Server) HandleGet(storeName, key string) ([]map[string]interface{}, error) {
	if storeName == "" {
		return s.handleGetAllStores()
	}
	if key == "" {
		return s.handleGetAll(storeName)
	}
	return s.handleGetSingle(storeName, key)
}

func (s *Server) handleGetAllStores() ([]map[string]interface{}, error) {
	stores := s.manager.GetStores()
	result := make([]map[string]interface{}, 0, len(stores))
	for _, name := range stores {
		result = append(result, map[string]interface{}{"store": name})
	}
	return result, nil
}

func (s *Server) handleGetSingle(storeName, key string) ([]map[string]interface{}, error) {
	store, exists := s.manager.GetStore(storeName)
	if !exists {
		return nil, fmt.Errorf("store '%s' not found", storeName)
	}

	keyInt, err := strconv.Atoi(key)
	if err != nil {
		return nil, fmt.Errorf("invalid integer key '%s': %v", key, err)
	}
	if keyInt < 0 {
		return nil, fmt.Errorf("invalid integer key '%s': must be >= 0", key)
	}

	value, err := store.Get(keyInt)
	if err != nil {
		return nil, err
	}

	return []map[string]interface{}{{"key": keyInt, "value": value}}, nil
}

func (s *Server) handleGetAll(storeName string) ([]map[string]interface{}, error) {
	store, exists := s.manager.GetStore(storeName)
	if !exists {
		return nil, fmt.Errorf("store '%s' not found", storeName)
	}

	allData := store.GetAllData()
	result := make([]map[string]interface{}, 0, len(allData))
	for k, v := range allData {
		result = append(result, map[string]interface{}{
			"key":   k,
			"value": v,
		})
	}
	return result, nil
}

func (s *Server) HandleDelete(storeName string, key string) ([]map[string]interface{}, error) {
	if storeName == "" {
		return nil, fmt.Errorf("failed to remove store: no store provided")
	}

	store, exists := s.manager.GetStore(storeName)
	if !exists {
		return nil, fmt.Errorf("store '%s' not found", storeName)
	}

	if key != "" {
		keyInt, err := strconv.Atoi(key)
		if err != nil {
			return nil, fmt.Errorf("invalid integer key '%s': %v", key, err)
		}
		if keyInt < 0 {
			return nil, fmt.Errorf("invalid integer key '%s': must be >= 0", key)
		}

		err = store.Delete(keyInt)
		if err != nil {
			return nil, fmt.Errorf("failed to delete key %d: %v", keyInt, err)
		}

		return []map[string]interface{}{{"status": 200}}, nil
	}

	if s.manager.RemoveStore(storeName) {
		return []map[string]interface{}{{"status": 200}}, nil
	}
	return nil, fmt.Errorf("failed to remove store '%s'", storeName)
}

func (s *Server) HandleSet(storeName string, key string, data string) ([]map[string]interface{}, error) {
	if storeName == "" {
		return nil, fmt.Errorf("failed to set data: no store provided")
	}

	store, exists := s.manager.GetStore(storeName)
	if !exists {
		s.manager.AddStore(storeName)
		var ok bool
		store, ok = s.manager.GetStore(storeName)
		if !ok {
			return nil, fmt.Errorf("failed to create store '%s'", storeName)
		}
		if key == "" && data == "" {
			return []map[string]interface{}{{"status": 200}}, nil
		}
	}

	if key != "" {
		keyInt, err := strconv.Atoi(key)
		if err != nil {
			return nil, fmt.Errorf("invalid integer key '%s': %v", key, err)
		}
		if keyInt < 0 {
			return nil, fmt.Errorf("invalid integer key '%s': must be >= 0", key)
		}
		if data == "" {
			return nil, fmt.Errorf("failed to update key %d: no data provided", keyInt)
		}

		err = store.Update(keyInt, data)
		if err != nil {
			return nil, fmt.Errorf("failed to update key %d: %v", keyInt, err)
		}
		return []map[string]interface{}{{"status": 200}}, nil
	} else if data != "" {
		newKey := int(store.NextKey.Load())
		err := store.Set(newKey, data)
		if err != nil {
			return nil, fmt.Errorf("failed to set data: %v", err)
		}
		return []map[string]interface{}{{"status": 200}}, nil
	}

	return nil, fmt.Errorf("failed to set data: no key or data provided")
}
