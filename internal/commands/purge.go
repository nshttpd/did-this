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

package commands

import (
	"fmt"
	bolt "go.etcd.io/bbolt"
	"time"

	"os"

	"github.com/spf13/cobra"
)

const RETENTION = 30

var purgeRetentionVar int

// purgeCmd represents the purge command
var purgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "purge out old completed stored tasks",
	Long: `Loop through the database and purge out old daily buckets
of data. Default retention period is 30 days. Can be overridden on the
command line.

	did-this purge

this is destructive and no backup of data is made.`,
	Run: func(cmd *cobra.Command, args []string) {

		var purge [][]byte
		now := time.Now()

		err := cfg.Db.View(func(tx *bolt.Tx) error {
			return tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
				d, err := time.Parse("2006-01-02", string(name))
				if err != nil {
					return fmt.Errorf("error parsing time from bucket : %s", string(name))
				}

				if d.Before(now.AddDate(0, 0, -purgeRetentionVar)) {
					purge = append(purge, name)
				}
				return nil
			})
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if len(purge) > 0 {
			err := cfg.Db.Update(func(tx *bolt.Tx) error {
				for _, b := range purge {
					if e := tx.DeleteBucket(b); e != nil {
						return e
					}
					fmt.Printf("purged : %s\n", string(b))
				}
				return nil
			})
			if err != nil {
				fmt.Println(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(purgeCmd)
	purgeCmd.Flags().IntVar(&purgeRetentionVar, "retention", RETENTION,
		"retention days for completed tasks")
}
