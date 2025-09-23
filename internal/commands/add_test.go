package commands

import (
	"testing"

	bolt "go.etcd.io/bbolt"
)

func TestAddCommand(t *testing.T) {
	db, cleanup := setupTestDb(t)
	defer cleanup()

	cfg = &Config{Db: db}

	testTask := "this is a test task"

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(cfg.CurrentDate())
		return err
	})

	addCmd.Run(addCmd, []string{testTask})

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(cfg.CurrentDate())
		if b == nil {
			t.Fatal("bucket not created")
		}

		c := b.Cursor()
		k, v := c.First()
		if k == nil {
			t.Fatal("no data saved in bucket")
		}

		if string(v) != testTask {
			t.Errorf("add command failed, got: %s, want: %s", string(v), testTask)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
