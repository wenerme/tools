/*
Copyright © 2020 wener <wenermail@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/wenerme/tools/pkg/apki"

	"github.com/spf13/cobra"
)

// dbMirrorLoadCmd represents the mirrorLoad command
var dbMirrorLoadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load mirror",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		idx, err := buildIndexer()
		if err != nil {
			return err
		}

		r, err := http.Get("https://mirrors.alpinelinux.org/mirrors.json")
		if err != nil {
			return err
		}
		var all []mirrorRecord
		err = jsoniter.NewDecoder(r.Body).Decode(&all)
		if err != nil {
			return err
		}

		all = append(all, mirrorRecord{
			Name:     "mirrors.aliyun.com",
			Location: "China",
			URLs: []string{
				"http://mirrors.aliyun.com/alpine/",
				"https://mirrors.aliyun.com/alpine/",
			},
		})

		_ = idx
		for _, v := range all {
			m := apki.Mirror{
				Name:      v.Name,
				Location:  v.Location,
				Bandwidth: v.Bandwidth,
			}
			_ = m.URLs.Set(v.URLs)

			if err := idx.DB.FirstOrCreate(&m, "name = ?", v.Name).Error; err != nil {
				return errors.Wrapf(err, "create %q", v.Name)
			}
		}
		return nil
	},
}

type mirrorRecord struct {
	Name      string
	Location  string
	URLs      []string
	Bandwidth string
}

func init() {
	dbMirrorCmd.AddCommand(dbMirrorLoadCmd)
}
