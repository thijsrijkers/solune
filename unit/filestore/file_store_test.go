package filestore_test

import (
	"bufio"
	"os"
	"testing"
	"solune/filestore"
	"strings"
	"time"
)

func TestFileStore(t *testing.T) {
	// Clean up after the test runs
	defer func() {
		err := os.RemoveAll("db")
		if err != nil {
			t.Fatalf("Error cleaning up db folder: %v", err)
		}
	}()

	// Test 1: Create a new FileStore
	store, err := filestore.New("testfile")
	if err != nil {
		t.Fatalf("Failed to create FileStore: %v", err)
	}

	// Check if the file exists
	if _, err := os.Stat("db/testfile.solstr"); os.IsNotExist(err) {
		t.Fatalf("Expected file to be created, but it doesn't exist")
	}

	// Test 2: Update a key-value pair
	key := "name"
	value := "John Doe"
	err = store.Update(key, value)
	if err != nil {
		t.Fatalf("Failed to update: %v", err)
	}

	// Open the file and verify if the content is correct
	file, err := os.Open("db/testfile.solstr")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close() // Ensure the file is closed before renaming

	// Read the file line by line
	var found bool
	var line string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line = scanner.Text()
		if strings.HasPrefix(line, key+",") {
			found = true
			if !strings.Contains(line, value) {
				t.Fatalf("Expected value %s, but got %s", value, line)
			}
		}
	}

	if !found {
		t.Fatalf("Key %s not found in file", key)
	}

	// Check that the key-value pair is no longer in the file
	found = false
	for scanner.Scan() {
		line = scanner.Text()
		if strings.HasPrefix(line, key+",") {
			found = true
			break
		}
	}

	if found {
		t.Fatalf("Expected key %s to be deleted, but it still exists", key)
	}

	// Close the file store last
	// Add a small delay before closing files, to allow the system to unlock the files
	time.Sleep(100 * time.Millisecond)

	// Ensure there is no temporary file left around
	if _, err := os.Stat("db/testfile.solstr.tmp"); err == nil {
		err := os.Remove("db/testfile.solstr.tmp")
		if err != nil {
			t.Fatalf("Failed to remove temp file: %v", err)
		}
	}

	// Final cleanup: Ensure the file store is closed properly
	err = store.Close()
	if err != nil {
		t.Fatalf("Failed to close FileStore after final cleanup: %v", err)
	}
}
