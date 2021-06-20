/*
Copyright © 2021 Isaac Daly <isaac.daly@outlook.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package cmd

import (
	"fmt"
	"regexp"

	"github.com/dalyisaac/mstodo/api"
	"github.com/dalyisaac/mstodo/types"
	"github.com/dalyisaac/mstodo/utils"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

var (
	filter string
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Get a list of the task lists",
	Long:  `Get a list of the Microsoft To Do task lists`,
	Run: func(cmd *cobra.Command, args []string) {
		r, err := regexp.Compile(filter)
		if err != nil {
			fmt.Println("Error creating regex:", err)
			return
		}

		req, err := api.CreateRequest()
		if err != nil {
			fmt.Println("error creating request:", err)
			return
		}

		resp, err := req.SetResult(&types.TodoTaskListResponse{}).Get("/me/todo/lists")
		if err != nil {
			fmt.Println("Error getting lists:", err)
			return
		}

		result := resp.Result().(*types.TodoTaskListResponse)
		printTaskList(result.Value, r)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.PersistentFlags().StringVarP(&filter, "filter", "f", ".", "Filter the lists which contain this regex")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func printTaskList(taskList types.TodoTaskList, r *regexp.Regexp) {
	t := utils.CreateFormattedTable(
		&table.Row{"Name", "Owner", "Shared"},
		&[]table.ColumnConfig{
			{Name: "Name", Align: text.AlignLeft, Transformer: utils.Transformer},
			{Name: "Owner", Align: text.AlignCenter, Transformer: utils.Transformer},
			{Name: "Shared", Align: text.AlignCenter, Transformer: utils.Transformer},
		},
	)

	for _, taskItem := range taskList {
		if r.MatchString(taskItem.DisplayName) {
			t.AppendRow(table.Row{
				taskItem.DisplayName,
				taskItem.IsOwner,
				taskItem.IsShared,
			})
		}
	}

	t.Render()
}
