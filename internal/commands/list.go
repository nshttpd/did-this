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
	"encoding/binary"
	"fmt"
	bolt "go.etcd.io/bbolt"
	"html/template"

	"os"
	"time"

	"github.com/spf13/cobra"
)

const (
	stdTmpl   = "{{ range . }}{{ .Count }} - {{ .Did }}\n{{ end }}"
	slackTmpl = "```\n{{ range . }}{{ .Count }} - {{ .Did }}\n{{ end }}```\n"
)

var slack bool

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list out the things that you've done",
	Long: `List out those things that you have done to report into whatever
you have to. Default is to list what you did yesterday. Use the keyword 'today'
to remind yourself what you've been doing. Provide a date to get those thing
that you did in the past. Examples:

	did-this list
	did-this list today
	did-this list 2018-04-15

The date format is that of YYYY-MM-DD for getting specific dates.`,
	Run: func(cmd *cobra.Command, args []string) {
		var date []byte
		if len(args) == 0 {
			date = cfg.PreviousDate()
		} else {
			if args[0] == "today" {
				date = cfg.CurrentDate()
			} else {
				_, err := time.Parse("2006-01-02", args[0])
				if err != nil {
					fmt.Println("invalid data format.")
					os.Exit(1)
				}
				date = []byte(args[0])
			}
		}
		cfg.Db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket(date)

			type done struct {
				Count uint64
				Did   string
			}

			var items []done

			if b != nil {
				b.ForEach(func(k, v []byte) error {
					items = append(items, done{btoi(k), string(v)})
					return nil
				})
				var t *template.Template
				var err error

				if !slack {
					t, err = template.New("list").Parse(stdTmpl)
				} else {
					t, err = template.New("list").Parse(slackTmpl)
				}

				if err != nil {
					fmt.Printf("error parsing output template : %v", err)
					os.Exit(1)
				}
				err = t.Execute(os.Stdout, items)
				if err != nil {
					fmt.Printf("error processing template : %v", err)
					os.Exit(1)
				}

			}
			return nil
		})

	},
}

func btoi(v []byte) uint64 {
	return binary.BigEndian.Uint64(v)
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVar(&slack, "slack", false, "output with slack preformat wrappers")
}
