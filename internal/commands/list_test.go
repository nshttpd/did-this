package commands

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	bolt "go.etcd.io/bbolt"
)

func TestListCommand(t *testing.T) {
	db, cleanup := setupTestDb(t)
	defer cleanup()

	cfg = &Config{Db: db}

	testTask1 := "this is a test task"
	testTask2 := "this is another test task"

	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(cfg.PreviousDate())
		if err != nil {
			return err
		}
		id, _ := b.NextSequence()
		b.Put(itob(id), []byte(testTask1))
		id, _ = b.NextSequence()
		b.Put(itob(id), []byte(testTask2))
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	listCmd.Run(listCmd, []string{})

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)

	got := buf.String()

	if !strings.Contains(got, testTask1) {
		t.Errorf("list command output does not contain expected task: %s", testTask1)
	}
	if !strings.Contains(got, testTask2) {
		t.Errorf("list command output does not contain expected task: %s", testTask2)
	}
}
