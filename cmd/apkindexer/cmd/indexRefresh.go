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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wenerme/tools/pkg/apki"
)

var indexRefreshOpts = struct {
	coordinate apki.IndexCoordinate
}{}

// indexRefreshCmd represents the indexRefresh command
var indexRefreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh index",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.Set("Database.AutoMigrate", true)
		idx, err := buildIndexer()
		if err != nil {
			return err
		}
		coord := indexRefreshOpts.coordinate
		if coord.Arch != "" && coord.Branch != "" && coord.Repo != "" {
			return idx.RefreshIndex(coord)
		}
		return idx.RefreshAllIndex()
	},
}

func init() {
	indexCmd.AddCommand(indexRefreshCmd)
	indexRefreshCmd.Flags().StringVar(&indexRefreshOpts.coordinate.Branch, "branch", "", "Branch")
	indexRefreshCmd.Flags().StringVar(&indexRefreshOpts.coordinate.Arch, "arch", "", "Arch")
	indexRefreshCmd.Flags().StringVar(&indexRefreshOpts.coordinate.Repo, "repo", "", "Repository")
}
