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
	"github.com/dalyisaac/mstodo/utils"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(createListsCmd())
}

type listsParams struct {
	columns  []table.ColumnConfig
	filter   *regexp.Regexp
	sortMode table.SortMode
}

func createListsCmd() *cobra.Command {
	var (
		filterFlag, sortFlag, excludeFlag string
		showIdFlag                        bool
	)

	// listsCmd represents the list command
	var listsCmd = &cobra.Command{
		Use:   "lists",
		Short: "Get a list of the task lists",
		Long:  `Get a list of the Microsoft To Do task lists`,
		RunE: func(cmd *cobra.Command, args []string) error {
			params, err := getListsCmdParams(filterFlag, sortFlag, excludeFlag, showIdFlag)
			if err != nil {
				return err
			}

			// Get lists
			lists, err := api.GetLists()
			if err != nil {
				return err
			}

			// Display results
			printTaskListList(*lists, params)
			return nil
		},
	}

	listsCmd.Flags().StringVarP(&filterFlag, "filter", "f", ".", "Filter the lists which contain this regex")
	listsCmd.Flags().StringVarP(&sortFlag, "sort", "s", "none", "Sort by the name - choices: "+utils.GetSortOptions())
	listsCmd.Flags().StringVarP(&excludeFlag, "exclude", "x", "", "Exclude columns")
	listsCmd.Flags().BoolVarP(&showIdFlag, "id", "i", false, "Show the list IDs")

	return listsCmd
}

var listsCmdCols = []table.ColumnConfig{
	{Name: "ID", Align: text.AlignLeft, Transformer: utils.Transformer},
	{Name: "Name", Align: text.AlignLeft, Transformer: utils.Transformer},
	{Name: "Owner", Align: text.AlignCenter, Transformer: utils.Transformer},
	{Name: "Shared", Align: text.AlignCenter, Transformer: utils.Transformer},
}

func getListsCmdParams(filterFlag, sortFlag, excludeFlag string, showIdFlag bool) (*listsParams, error) {
	// Validate filter
	r, err := regexp.Compile(filterFlag)
	if err != nil {
		return nil, err
	}

	// Validate sort
	sortMode, err := utils.GetSortMode(sortFlag)
	if err != nil {
		return nil, err
	}

	// Show ID flag
	if !showIdFlag {
		excludeFlag = "ID," + excludeFlag
	}

	// Construct excluded columns
	cols, err := utils.GetAllowedColumns(excludeFlag, listsCmdCols)
	if err != nil {
		return nil, err
	}

	return &listsParams{
		columns:  cols,
		filter:   r,
		sortMode: sortMode,
	}, nil
}

func printTaskListList(taskListList api.TodoTaskListList, params *listsParams) {
	headerRow := table.Row{}
	columns := params.columns

	// Add the column names to the header row
	for _, c := range columns {
		headerRow = append(headerRow, c.Name)
	}

	t := utils.CreateFormattedTable(&headerRow, &columns)

	for _, taskList := range taskListList {
		if params.filter.MatchString(taskList.DisplayName) {
			row := table.Row{}

			row = append(row, getAllowedTaskListItemFields(taskList, columns)...)
			t.AppendRow(row)
		}
	}

	if params.sortMode != utils.NoSort {
		t.SortBy([]table.SortBy{
			{Name: "Name", Mode: params.sortMode},
		})
	}

	t.Render()
}

func getAllowedTaskListItemFields(taskItem api.TodoTaskListItem, columns []table.ColumnConfig) table.Row {
	fields := table.Row{}

	for _, col := range columns {
		switch col.Name {
		case "ID":
			fields = append(fields, taskItem.Id)
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
