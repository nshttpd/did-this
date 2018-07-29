// Copyright © 2018 Steve Brunton <sbrunton@gmail.com>
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice,
//    this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors
//    may be used to endorse or promote products derived from this software
//    without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.

package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/coreos/bbolt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

const (
	defaultConfigFile = ".didthis/config.json"
)

var (
	cfgFile  string
	logLevel string
	cfg      *Config
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "didthis",
	Short: "task completion tracking app",
	Long: `A command line task completion tracking app. Added completed
tasks will end up in a bucket based on the date they are added. All tasks
can be listed back out for daily next day reporting.

	didthis add "this is something I did today"
	didthis list
	didthis yesterday

`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		l, err := log.ParseLevel(logLevel)
		if err != nil {
			fmt.Printf("error setting log level : %v", err)
			os.Exit(1)
		}
		log.SetLevel(l)

		// load the config. this will also have a handle to the Bolt DB
		cfg = loadConfig()

		err = cfg.Db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists(cfg.CurrentDate())
			return err
		})

		if err != nil {
			log.WithField("error", err).Fatal("error creating daily DB bucket")
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		cfg.SaveConfig()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	hd, err := homedir.Dir()
	if err != nil {
		log.Fatalf("error discerning homedir : %v", err)
	}

	cfgFile = fmt.Sprintf("%s/%s", hd, defaultConfigFile)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", cfgFile, "config file")
	rootCmd.PersistentFlags().StringVar(&logLevel, "loglevel", "info", "log level")
}
