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
	"errors"
	"regexp"

	"github.com/dalyisaac/mstodo/api"
	"github.com/dalyisaac/mstodo/utils"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(createViewCmd())
}

type viewParams struct {
	columns []table.ColumnConfig
	filter  *regexp.Regexp
	sort    []table.SortBy
}

func createViewCmd() *cobra.Command {
	var (
		filterFlag, sortFlag, excludeFlag string
		showIdFlag                        bool
	)

	// viewCmd represents the view command
	var viewCmd = &cobra.Command{
		Use:   "view <list name>",
		Short: "View a specific list",
		Long:  `View a specific task list`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing list name")
			}

			params, err := getViewCmdParams(filterFlag, sortFlag, excludeFlag, showIdFlag)
			if err != nil {
				return err
			}

			// Get name
			name, err := utils.CleanName(args[0])
			if err != nil {
				return nil
			}

			// Get lists
			lists, err := api.GetLists()
			if err != nil {
				return err
			}

			// Get task list id
			listId, err := lists.GetListId(name)
			if err != nil {
				return err
			}

			// Get task list
			tasks, err := api.GetTasks(listId)
			if err != nil {
				return err
			}

			// Display results
			printTaskList(*tasks, params)

			return nil
		},
	}

	viewCmd.Flags().StringVarP(&filterFlag, "filter", "f", ".", "Filter the tasks which contain this regex")
	viewCmd.Flags().StringVarP(&sortFlag, "sort", "s", "none", "Sort by the fields, for example: --sort=[title:dsc,created:asc,status]")
	viewCmd.Flags().StringVarP(&excludeFlag, "exclude", "x", "", "Exclude columns")
	viewCmd.Flags().BoolVarP(&showIdFlag, "id", "i", false, "Show the task IDs")

	return viewCmd
}

var viewCmdCols = []table.ColumnConfig{
	utils.LeftColumn("Id"),
	utils.LeftColumn("Title"),
	utils.CenterColumn("Importance"),
	utils.CenterColumnTransformer("Status", utils.StatusTransformer),
	utils.CenterColumn("Reminder"),
	utils.CenterColumn("Due Date"),
	utils.CenterColumn("Completed"),
	utils.CenterColumn("Created"),
	utils.CenterColumn("Last Modified"),
}

func getViewCmdParams(filterFlag, sortFlag, excludeFlag string, showIdFlag bool) (*viewParams, error) {
	// Validate filter
	r, err := regexp.Compile(filterFlag)
	if err != nil {
		return nil, err
	}

	// Get sort
	sortBy, err := utils.GetSortByColumns(sortFlag, viewCmdCols)
	if err != nil {
		return nil, err
	}

	// Show ID flag
	if !showIdFlag {
		excludeFlag = "ID," + excludeFlag
	}

	// Construct excluded columns
	cols, err := utils.GetAllowedColumns(excludeFlag, viewCmdCols)
	if err != nil {
		return nil, err
	}

	return &viewParams{
		columns: cols,
		filter:  r,
		sort:    sortBy,
	}, nil
}

func printTaskList(taskList api.TodoTaskList, params *viewParams) {
	headerRow := table.Row{}
	columns := params.columns

	// Add the column names to the header row
	for _, c := range columns {
		headerRow = append(headerRow, c.Name)
	}

	t := utils.CreateFormattedTable(&headerRow, &columns)

	for _, todoTask := range taskList {
		if params.filter.MatchString(todoTask.Title) {
			row := table.Row{}
			row = append(row, getAllowedTodoTaskFields(todoTask, columns)...)
			t.AppendRow(row)
		}
	}

	if len(params.sort) != 0 {
		t.SortBy(params.sort)
	}

	t.Render()
}

func getAllowedTodoTaskFields(todoTask api.TodoTask, columns []table.ColumnConfig) table.Row {
	fields := table.Row{}

	for _, col := range columns {
		switch col.Name {
		case "Id":
			fields = append(fields, todoTask.Id)
		case "Title":
			fields = append(fields, todoTask.Title)
		case "Importance":
			fields = append(fields, todoTask.Importance)
		case "Status":
			fields = append(fields, todoTask.Status)
		case "Reminder":
			fields = append(fields, todoTask.ReminderDateTime)
		case "Due Date":
			fields = append(fields, todoTask.DueDateTime)
		case "Completed":
			fields = append(fields, todoTask.Completed)
		case "Created":
			fields = append(fields, todoTask.CreatedDateTime)
		case "Last Modified":
			fields = append(fields, todoTask.LastModifiedDateTime)
		}
	}

	return fields
}
