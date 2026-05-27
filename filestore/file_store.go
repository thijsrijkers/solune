package filestore

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type FileStore struct {
	filePath string
	file     *os.File
	writer   *bufio.Writer
	fileLock sync.RWMutex
	index    map[string]int64
}

func New(filename string) (*FileStore, error) {
	if !strings.HasSuffix(filename, ".solstr") {
		filename += ".solstr"
	}
	baseDir := "db"
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, err
	}
	fullPath := filepath.Join(baseDir, filename)
	f, err := os.OpenFile(fullPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	store := &FileStore{
		filePath: fullPath,
		file:     f,
		writer:   bufio.NewWriterSize(f, 64*1024),
		index:    make(map[string]int64),
	}
	if err := store.buildIndex(); err != nil {
		return nil, err
	}
	return store, nil
}

func (s *FileStore) buildIndex() error {
	f, err := os.Open(s.filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	var offset int64
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if comma := strings.IndexByte(line, ','); comma > 0 {
			s.index[line[:comma]] = offset
		}
		offset += int64(len(line)) + 1
	}
	return scanner.Err()
}

func (s *FileStore) Keys() []string {
	s.fileLock.RLock()
	defer s.fileLock.RUnlock()
	keys := make([]string, 0, len(s.index))
	for k := range s.index {
		keys = append(keys, k)
	}
	return keys
}

func (s *FileStore) Get(key string) (string, error) {
	s.fileLock.RLock()
	offset, ok := s.index[key]
	s.fileLock.RUnlock()

	if !ok {
		return "", fmt.Errorf("key %q not found", key)
	}

	f, err := os.Open(s.filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := f.Seek(offset, 0); err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(f)
	if scanner.Scan() {
		line := scanner.Text()
		if comma := strings.IndexByte(line, ','); comma > 0 {
			return line[comma+1:], nil
		}
	}
	return "", fmt.Errorf("key %q not found", key)
}

func (s *FileStore) Update(key, value string) error {
	s.fileLock.Lock()
	defer s.fileLock.Unlock()

	tempPath := s.filePath + ".tmp"
	tempFile, err := os.Create(tempPath)
	if err != nil {
		return err
	}
	origFile, err := os.Open(s.filePath)
	if err != nil {
		tempFile.Close()
		return err
	}

	scanner := bufio.NewScanner(origFile)
	writer := bufio.NewWriter(tempFile)
	newIndex := make(map[string]int64, len(s.index))
	found := false
	var offset int64

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, key+",") {
			line = fmt.Sprintf("%s,%s", key, value)
			found = true
		}
		if comma := strings.IndexByte(line, ','); comma > 0 {
			newIndex[line[:comma]] = offset
		}
		n, err := writer.WriteString(line + "\n")
		if err != nil {
			origFile.Close()
			tempFile.Close()
			return err
		}
		offset += int64(n)
	}
	origFile.Close()

	if !found {
		line := fmt.Sprintf("%s,%s\n", key, value)
		newIndex[key] = offset
		if _, err := writer.WriteString(line); err != nil {
			tempFile.Close()
			return err
		}
	}
	if err := writer.Flush(); err != nil {
		tempFile.Close()
		return err
	}
	tempFile.Close()
	s.writer.Flush()
	s.file.Close()
	if err := os.Rename(tempPath, s.filePath); err != nil {
		return err
	}
	s.index = newIndex
	return s.reopenFile()
}

func (s *FileStore) Delete(key string) error {
	s.fileLock.Lock()
	defer s.fileLock.Unlock()

	tempPath := s.filePath + ".tmp"
	tempFile, err := os.Create(tempPath)
	if err != nil {
		return err
	}
	origFile, err := os.Open(s.filePath)
	if err != nil {
		tempFile.Close()
		return err
	}

	scanner := bufio.NewScanner(origFile)
	writer := bufio.NewWriter(tempFile)
	newIndex := make(map[string]int64, len(s.index))
	var offset int64

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, key+",") {
			continue
		}
		if comma := strings.IndexByte(line, ','); comma > 0 {
			newIndex[line[:comma]] = offset
		}
		n, err := writer.WriteString(line + "\n")
		if err != nil {
			origFile.Close()
			tempFile.Close()
			return err
		}
		offset += int64(n)
	}
	origFile.Close()

	if err := writer.Flush(); err != nil {
		tempFile.Close()
		return err
	}
	tempFile.Close()
	s.writer.Flush()
	s.file.Close()
	if err := os.Rename(tempPath, s.filePath); err != nil {
		return err
	}
	s.index = newIndex
	return s.reopenFile()
}

func (s *FileStore) reopenFile() error {
	f, err := os.OpenFile(s.filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	s.file = f
	s.writer = bufio.NewWriterSize(f, 64*1024)
	return nil
}

func (s *FileStore) Close() error {
	s.fileLock.Lock()
	defer s.fileLock.Unlock()
	if err := s.writer.Flush(); err != nil {
		return err
	}
	return s.file.Close()
}
