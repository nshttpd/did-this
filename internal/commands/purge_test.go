package commands

import (
	"fmt"
	"testing"
	"time"

	bolt "go.etcd.io/bbolt"
)

func TestPurgeCommand(t *testing.T) {
	db, cleanup := setupTestDb(t)
	defer cleanup()

	cfg = &Config{Db: db}

	oldDate := time.Now().AddDate(0, 0, -31)
	oldDateStr := fmt.Sprintf("%d-%02d-%02d", oldDate.Year(), oldDate.Month(), oldDate.Day())

	recentDate := time.Now()
	recentDateStr := fmt.Sprintf("%d-%02d-%02d", recentDate.Year(), recentDate.Month(), recentDate.Day())

	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(oldDateStr))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(recentDateStr))
		return err
	})
	if err != nil {
		t.Fatal(err)
	}

	purgeCmd.Run(purgeCmd, []string{})

	err = db.View(func(tx *bolt.Tx) error {
		if tx.Bucket([]byte(oldDateStr)) != nil {
			t.Errorf("purge command failed, old bucket not removed: %s", oldDateStr)
		}
		if tx.Bucket([]byte(recentDateStr)) == nil {
			t.Errorf("purge command failed, recent bucket removed: %s", recentDateStr)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
