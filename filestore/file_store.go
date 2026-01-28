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
	fileLock sync.Mutex
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

	return &FileStore{
		filePath: fullPath,
		file:     f,
		writer:   bufio.NewWriterSize(f, 64*1024),
	}, nil
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

	found := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, key+",") {
			line = fmt.Sprintf("%s,%s", key, value)
			found = true
		}
		if _, err := writer.WriteString(line + "\n"); err != nil {
			origFile.Close()
			tempFile.Close()
			return err
		}
	}
	origFile.Close()

	if !found {
		if _, err := writer.WriteString(fmt.Sprintf("%s,%s\n", key, value)); err != nil {
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

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, key+",") {
			if _, err := writer.WriteString(line + "\n"); err != nil {
				origFile.Close()
				tempFile.Close()
				return err
			}
		}
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
