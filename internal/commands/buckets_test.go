package commands

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	bolt "go.etcd.io/bbolt"
)

func TestBucketsCommand(t *testing.T) {
	db, cleanup := setupTestDb(t)
	defer cleanup()

	cfg = &Config{Db: db}

	bucket1 := "2025-09-20"
	bucket2 := "2025-09-21"

	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket1))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(bucket2))
		return err
	})
	if err != nil {
		t.Fatal(err)
	}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	bucketsCmd.Run(bucketsCmd, []string{})

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)

	got := buf.String()

	if !strings.Contains(got, bucket1) {
		t.Errorf("buckets command output does not contain expected bucket: %s", bucket1)
	}
	if !strings.Contains(got, bucket2) {
		t.Errorf("buckets command output does not contain expected bucket: %s", bucket2)
	}
}
