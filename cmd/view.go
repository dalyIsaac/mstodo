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
	"time"

	"github.com/dalyisaac/mstodo/api"
	"github.com/dalyisaac/mstodo/datetime"
	"github.com/dalyisaac/mstodo/utils"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(createViewCmd())
}

type viewParamsFlags struct {
	// filter flags
	title, status, reminder, dueDate, completed, created, lastModified string
	sort, exclude                                                      string
	absoluteTime, showId                                               bool
}

type viewParams struct {
	columns            []table.ColumnConfig
	sort               []table.SortBy
	titleFilter        *regexp.Regexp
	statusFilter       *regexp.Regexp
	reminderFilter     *datetime.DateFilters
	dueDateFilter      *datetime.DateFilters
	completedFilter    *datetime.DateFilters
	createdFilter      *datetime.DateFilters
	lastModifiedFilter *datetime.DateFilters
}

const matchAll = "."

func createViewCmd() *cobra.Command {
	flags := viewParamsFlags{}

	// viewCmd represents the view command
	var viewCmd = &cobra.Command{
		Use:   "view <list name>",
		Short: "View a specific list",
		Long: `View a specific task list.

Dates can be filtered using by specifying the start and/or end date you're interested in. For example:

--reminder="start Monday; end fri"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing list name")
			}

			params, err := getViewCmdParams(flags)
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
			params.printTaskList(*tasks)

			return nil
		},
	}

	viewCmd.Flags().StringVarP(&flags.title, "title", "l", matchAll, "Filter the task names which contain this regex")
	viewCmd.Flags().StringVarP(&flags.status, "status", "u", matchAll, "Filter the status (use 'completed' for ✅)")
	viewCmd.Flags().StringVarP(&flags.reminder, "reminder", "r", matchAll, "Filter by reminder using the date syntax")
	viewCmd.Flags().StringVarP(&flags.dueDate, "due", "d", matchAll, "Filter by due using the date syntax")
	viewCmd.Flags().StringVarP(&flags.completed, "completed", "o", matchAll, "Filter by completed using the date syntax")
	viewCmd.Flags().StringVarP(&flags.created, "created", "c", matchAll, "Filter by created using the date syntax")
	viewCmd.Flags().StringVarP(&flags.lastModified, "last-modified", "m", matchAll, "Filter by last-modified using the date syntax")

	viewCmd.Flags().StringVarP(&flags.sort, "sort", "s", "none", "Sort by the fields, for example: --sort=\"[title:dsc,created:asc,status]\"")
	viewCmd.Flags().StringVarP(&flags.exclude, "exclude", "x", "", "Exclude columns")
	viewCmd.Flags().BoolVarP(&flags.absoluteTime, "absolute", "a", false, "Show absolute datetime")
	viewCmd.Flags().BoolVarP(&flags.showId, "id", "i", false, "Show the task IDs")

	return viewCmd
}

func getViewCmdParams(flags viewParamsFlags) (*viewParams, error) {
	params := viewParams{}

	timeTransformer := utils.Transformer
	if flags.absoluteTime {
		timeTransformer = utils.AbsoluteTimeTransformer
	}

	// Validate filter
	getFilters(&params, flags)

	var viewCmdCols = []table.ColumnConfig{
		utils.LeftColumn("Id"),
		utils.LeftColumn("Title"),
		utils.CenterColumn("Importance"),
		utils.CenterColumnTransformer("Status", utils.StatusTransformer),
		utils.CenterColumnTransformer("Reminder", timeTransformer),
		utils.CenterColumnTransformer("Due Date", timeTransformer),
		utils.CenterColumnTransformer("Completed", timeTransformer),
		utils.CenterColumnTransformer("Created", timeTransformer),
		utils.CenterColumnTransformer("Last Modified", timeTransformer),
	}

	// Get sort
	sortBy, err := utils.GetSortByColumns(flags.sort, viewCmdCols)
	if err != nil {
		return nil, err
	}
	params.sort = sortBy

	// Show ID flag
	if !flags.showId {
		flags.exclude = "ID," + flags.exclude
	}

	// Construct excluded columns
	cols, err := utils.GetAllowedColumns(flags.exclude, viewCmdCols)
	if err != nil {
		return nil, err
	}
	params.columns = cols

	return &params, nil
}

type createDateFilter struct {
	filter **datetime.DateFilters
	flag   string
}

func getFilters(params *viewParams, flags viewParamsFlags) error {
	if flags.title != matchAll {
		if err := setStringFilter(&params.titleFilter, flags.title); err != nil {
			return err
		}
	}

	if flags.status != matchAll {
		if err := setStringFilter(&params.statusFilter, flags.status); err != nil {
			return err
		}
	}

	filters := []createDateFilter{
		{filter: &params.reminderFilter, flag: flags.reminder},
		{filter: &params.dueDateFilter, flag: flags.dueDate},
		{filter: &params.completedFilter, flag: flags.completed},
		{filter: &params.createdFilter, flag: flags.created},
		{filter: &params.lastModifiedFilter, flag: flags.lastModified},
	}

	for _, f := range filters {
		if f.flag == matchAll {
			continue
		}

		if err := setDateFilter(f.filter, f.flag); err != nil {
			return err
		}
	}

	return nil
}

func setStringFilter(filter **regexp.Regexp, flag string) error {
	r, err := regexp.Compile(flag)
	if err != nil {
		return err
	}

	*filter = r
	return nil
}

func setDateFilter(filter **datetime.DateFilters, flag string) error {
	res, err := datetime.DateStartEndParser(flag)
	if err != nil {
		return err
	}

	*filter = res
	return nil
}

func (params *viewParams) printTaskList(taskList api.TodoTaskList) {
	headerRow := table.Row{}
	columns := params.columns

	// Add the column names to the header row
	for _, c := range columns {
		headerRow = append(headerRow, c.Name)
	}

	t := utils.CreateFormattedTable(&headerRow, &columns)

	for _, todoTask := range taskList {
		if params.canAdd(todoTask) {
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

type viewCanAddRegexp struct {
	filter *regexp.Regexp
	field  *string
}

type viewCanAddDate struct {
	filter *datetime.DateFilters
	field  *datetime.GraphTime
}

func (params *viewParams) canAdd(task api.TodoTask) bool {
	checkStrings := []viewCanAddRegexp{
		{filter: params.titleFilter, field: &task.Title},
		{filter: params.statusFilter, field: status(task.Status)},
	}

	for _, f := range checkStrings {
		if f.filter != nil && !f.filter.MatchString(*f.field) {
			return false
		}
	}

	checkDates := []viewCanAddDate{
		{filter: params.reminderFilter, field: task.ReminderDateTime},
		{filter: params.dueDateFilter, field: task.DueDateTime},
		{filter: params.completedFilter, field: task.Completed},
		{filter: params.createdFilter, field: graphtime(task.CreatedDateTime)},
		{filter: params.lastModifiedFilter, field: graphtime(task.LastModifiedDateTime)},
	}

	for _, f := range checkDates {
		if !f.filter.Contains(f.field) {
			return false
		}
	}

	return true
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

// graphtime converts `time.Time` to a `datetime.GraphTime` pointer
func graphtime(t time.Time) *datetime.GraphTime {
	g := datetime.GraphTime(t)
	return &g
}

// status converts the `GraphStatus` g into a string pointer
func status(g api.GraphStatus) *string {
	s := string(g)
	return &s
}
