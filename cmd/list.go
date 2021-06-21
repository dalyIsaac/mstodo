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
	"regexp"

	"github.com/dalyisaac/mstodo/api"
	"github.com/dalyisaac/mstodo/types"
	"github.com/dalyisaac/mstodo/utils"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

var (
	filterFlag, sortFlag, excludeFlag string
	showIdFlag                        bool
)

var listColumns = []table.ColumnConfig{
	{Name: "Name", Align: text.AlignLeft, Transformer: utils.Transformer},
	{Name: "Owner", Align: text.AlignCenter, Transformer: utils.Transformer},
	{Name: "Shared", Align: text.AlignCenter, Transformer: utils.Transformer},
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Get a list of the task lists",
	Long:  `Get a list of the Microsoft To Do task lists`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validate filter
		r, err := regexp.Compile(filterFlag)
		if err != nil {
			return err
		}

		// Validate sort
		sortMode, err := utils.GetSortMode(sortFlag)
		if err != nil {
			return err
		}

		// Construct excluded columns
		cols, err := utils.GetAllowedColumns(excludeFlag, listColumns)
		if err != nil {
			return err
		}

		// Create request
		req, err := api.CreateRequest()
		if err != nil {
			return err
		}

		// Get request
		resp, err := req.SetResult(&types.TodoTaskListResponse{}).Get("/me/todo/lists")
		if err != nil {
			return err
		}

		// Display results
		result := resp.Result().(*types.TodoTaskListResponse)
		printTaskList(result.Value, cols, r, sortMode)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&filterFlag, "filter", "f", ".", "Filter the lists which contain this regex")
	listCmd.Flags().StringVarP(&sortFlag, "sort", "s", "none", "Sort by the name - choices: "+utils.GetSortOptions())
	listCmd.Flags().StringVarP(&excludeFlag, "exclude", "x", "", "Exclude columns: "+utils.GetSortOptions())
	listCmd.Flags().BoolVarP(&showIdFlag, "id", "i", false, "Show the list IDs")
}

func printTaskList(taskList types.TodoTaskList, columns []table.ColumnConfig, r *regexp.Regexp, sortMode table.SortMode) {
	rows := table.Row{}

	if showIdFlag {
		columns = append([]table.ColumnConfig{{Name: "ID", Align: text.AlignLeft, Transformer: utils.Transformer}}, columns...)
	}

	for _, c := range columns {
		rows = append(rows, c.Name)
	}

	t := utils.CreateFormattedTable(&rows, &columns)

	for _, taskItem := range taskList {
		if r.MatchString(taskItem.DisplayName) {
			row := table.Row{}
			if showIdFlag {
				row = append(row, taskItem.Id)
			}

			row = append(row, getAllowedTaskItemFields(taskItem, columns)...)
			t.AppendRow(row)
		}
	}

	if sortMode != utils.NoSort {
		t.SortBy([]table.SortBy{
			{Name: "Name", Mode: sortMode},
		})
	}

	t.Render()
}

func getAllowedTaskItemFields(taskItem types.TodoTaskListItem, columns []table.ColumnConfig) table.Row {
	fields := table.Row{}

	for _, col := range columns {
		switch col.Name {
		case "Name":
			fields = append(fields, taskItem.DisplayName)
		case "Owner":
			fields = append(fields, taskItem.IsOwner)
		case "Shared":
			fields = append(fields, taskItem.IsShared)
		}
	}

	return fields
}
