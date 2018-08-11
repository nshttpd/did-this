package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"fmt"

	"time"

	"github.com/coreos/bbolt"
	log "github.com/sirupsen/logrus"
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

func loadConfig() *Config {
	c := &Config{}
	if cf, err := ioutil.ReadFile(cfgFile); err != nil {
		if os.IsNotExist(err) {
			p := strings.Split(cfgFile, string(os.PathSeparator))
			c.DbPath = strings.Join(p[0:len(p)-1], string(os.PathSeparator))
			// make sure the directory exists and if not create it
			if _, err := ioutil.ReadDir(c.DbPath); err != nil {
				err := os.MkdirAll(c.DbPath, 0755)
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
