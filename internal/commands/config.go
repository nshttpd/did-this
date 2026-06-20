package commands

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"fmt"

	"time"

	log "github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
)

const (
	latestVersion = 1
	defaultDbName = "did-this.db"
)

type Config struct {
	Version int      `json:"version,omitempty"`
	DbPath  string   `json:"dbPath"`
	Db      *bolt.DB `json:"-"`
}

// validateConfigPath ensures the config file path is safe and within expected boundaries
func validateConfigPath(path string) error {
	// Get absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("invalid config path: %w", err)
	}

	// Clean the path to remove any .. or other traversal attempts
	cleanPath := filepath.Clean(absPath)

	// Ensure the path is in the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("cannot determine home directory: %w", err)
	}

	// Check if the clean path starts with the home directory
	if !strings.HasPrefix(cleanPath, homeDir) {
		return fmt.Errorf("config file must be within home directory")
	}

	// Ensure it's not a system-critical file
	systemPaths := []string{"/etc/", "/bin/", "/sbin/", "/usr/bin/", "/usr/sbin/", "/var/"}
	for _, sysPath := range systemPaths {
		if strings.HasPrefix(cleanPath, sysPath) {
			return fmt.Errorf("config file cannot be in system directory")
		}
	}

	return nil
}

func loadConfig() *Config {
	// Validate config file path before using it
	if err := validateConfigPath(cfgFile); err != nil {
		log.WithFields(log.Fields{"cfgFile": cfgFile, "error": err}).Fatal("invalid config file path")
	}

	c := &Config{}
	if cf, err := ioutil.ReadFile(cfgFile); err != nil {
		if os.IsNotExist(err) {
			p := strings.Split(cfgFile, string(os.PathSeparator))
			c.DbPath = strings.Join(p[0:len(p)-1], string(os.PathSeparator))
			// make sure the directory exists and if not create it
			if _, err := ioutil.ReadDir(c.DbPath); err != nil {
				err := os.MkdirAll(c.DbPath, 0700)
				if err != nil {
					log.WithFields(log.Fields{"dbPath": c.DbPath, "error": err}).Fatal("error creating db dir")
				}
			}
			c.Version = latestVersion
		} else {
			log.Error(err)
		}
	} else {
		err = json.Unmarshal(cf, c)
		if err != nil {
			log.WithFields(log.Fields{"cfgFile": cfgFile, "error": err}).Fatal("error reading config file")
		}
	}

	dbFile := fmt.Sprintf("%s/%s", c.DbPath, defaultDbName)

	var err error
	c.Db, err = bolt.Open(dbFile, 0600, &bolt.Options{Timeout: 1 * time.Second})

	if err != nil {
		log.WithFields(log.Fields{"dbFile": dbFile, "error": err}).Fatal("error opening DB file")
	}

	return c

}

func (c *Config) SaveConfig() {
	err := c.Db.Close()
	if err != nil {
		log.WithField("error", err).Fatal("error closing DB file")
	}
	sl := log.WithField("cfgfile", cfgFile)
	if cf, err := json.Marshal(c); err != nil {
		sl.WithField("error", err).Error("error marshalling config file")
	} else {
		if err := ioutil.WriteFile(cfgFile, cf, 0644); err != nil {
			sl.WithField("error", err).Error("error saving config file")
		} else {
			sl.Debug("saved config file")
		}
	}
}

func (c *Config) CurrentDate() []byte {
	t := time.Now()
	return []byte(fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day()))
}

func (c *Config) PreviousDate() []byte {
	t := time.Now().Add(-time.Hour * 24)
	return []byte(fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day()))
}
