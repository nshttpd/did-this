package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	bolt "go.etcd.io/bbolt"
)

const (
	testDbName = "test.db"
)

func setupTestDb(t *testing.T) (*bolt.DB, func()) {
	t.Helper()

	tmpDir, err := ioutil.TempDir("", "did-this-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	db, err := bolt.Open(fmt.Sprintf("%s/%s", tmpDir, testDbName), 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}

	return db, func() {
		db.Close()
		os.RemoveAll(tmpDir)
	}
}
