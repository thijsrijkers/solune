package filestore_test

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"solune/filestore"
)

func TestFileStore(t *testing.T) {
	baseDir := "../db"
	testFile := "testdata.solstr"
	fullPath := filepath.Join(baseDir, testFile)

	// Clean up all files in the db directory before and after test
	cleanupDir(t, baseDir)
	defer cleanupDir(t, baseDir)

	store, err := filestore.New("testdata")
	if err != nil {
		t.Fatalf("Failed to create FileStore: %v", err)
	}
	defer store.Close()

	// Test Write (via Update as upsert)
	if err := store.Update("foo", "bar"); err != nil {
		t.Errorf("Write failed: %v", err)
	}

	// Verify content
	content := readFile(t, fullPath)
	if !strings.Contains(content, "foo,bar") {
		t.Errorf("Expected content 'foo,bar', got: %s", content)
	}

	// Test Update existing key
	if err := store.Update("foo", "baz"); err != nil {
		t.Errorf("Update failed: %v", err)
	}
	content = readFile(t, fullPath)
	if !strings.Contains(content, "foo,baz") || strings.Contains(content, "foo,bar") {
		t.Errorf("Expected updated content 'foo,baz', got: %s", content)
	}

	// Test Update non-existing key
	if err := store.Update("newkey", "value"); err != nil {
		t.Errorf("Update new key failed: %v", err)
	}
	content = readFile(t, fullPath)
	if !strings.Contains(content, "newkey,value") {
		t.Errorf("Expected content 'newkey,value', got: %s", content)
	}

	// Test Delete
	if err := store.Delete("foo"); err != nil {
		t.Errorf("Delete failed: %v", err)
	}
	content = readFile(t, fullPath)
	if strings.Contains(content, "foo,") {
		t.Errorf("Expected 'foo' to be deleted, got: %s", content)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	var sb strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		sb.WriteString(scanner.Text() + "\n")
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("Failed to scan file: %v", err)
	}
	return sb.String()
}

func cleanupDir(t *testing.T, dir string) {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		t.Fatalf("Failed to read dir %s: %v", dir, err)
	}
	for _, entry := range entries {
		err := os.RemoveAll(filepath.Join(dir, entry.Name()))
		if err != nil {
			t.Fatalf("Failed to remove file %s: %v", entry.Name(), err)
		}
	}
}
