package commands

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestConfig_CurrentDate(t *testing.T) {
	db, cleanup := setupTestDb(t)
	defer cleanup()

	cfg := &Config{Db: db}

	want := time.Now()
	got := cfg.CurrentDate()

	wantStr := fmt.Sprintf("%d-%02d-%02d", want.Year(), want.Month(), want.Day())
	if string(got) != wantStr {
		t.Errorf("CurrentDate() = %s; want %s", got, wantStr)
	}
}

func TestConfig_PreviousDate(t *testing.T) {
	db, cleanup := setupTestDb(t)
	defer cleanup()

	cfg := &Config{Db: db}

	want := time.Now().Add(-time.Hour * 24)
	got := cfg.PreviousDate()

	wantStr := fmt.Sprintf("%d-%02d-%02d", want.Year(), want.Month(), want.Day())
	if string(got) != wantStr {
		t.Errorf("PreviousDate() = %s; want %s", got, wantStr)
	}
}

func TestConfig_SaveConfig(t *testing.T) {
	db, cleanup := setupTestDb(t)
	defer cleanup()

	cfgFile = "/tmp/did-this-test-config.json"
	cfg := &Config{Db: db, DbPath: "/tmp"}

	cfg.SaveConfig()

	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		t.Errorf("SaveConfig() did not create config file: %s", cfgFile)
	}

	os.Remove(cfgFile)
}
