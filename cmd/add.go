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
	"fmt"
	"strings"

	"github.com/dalyisaac/mstodo/api"
	"github.com/dalyisaac/mstodo/datetime"
	"github.com/dalyisaac/mstodo/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(createAddCmd())
}

type addParamsFlags struct {
	list       string
	importance string
	reminder   string
	dueDate    string
	status     string
}

const emptyString = ""
const addCutset = "'\" "

var importanceChoices = []string{"low", "normal", "high"}

func createAddCmd() *cobra.Command {
	flags := addParamsFlags{}

	addCmd := &cobra.Command{
		Use:   "add <task title>",
		Short: "Add a task",
		Long:  `Add a task`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing task name")
			}

			// Construct task
			task, err := constructTaskPayload(flags, args[0])
			if err != nil {
				return err
			}

			// Get lists
			lists, err := api.GetLists()
			if err != nil {
				return err
			}

			// Get task list id
			listId, err := lists.GetListId(flags.list)
			if err != nil {
				return err
			}

			if err := api.CreateTask(listId, task); err != nil {
				return err
			}

			return nil
		},
	}

	addCmd.Flags().StringVarP(&flags.list, "list", "l", "tasks", "Add the task to the specified list")
	addCmd.Flags().StringVarP(&flags.reminder, "reminder", "r", emptyString, "Task reminder (date time). For example, --reminder=\"Next Friday at 15:00\"")
	addCmd.Flags().StringVarP(&flags.dueDate, "due-date", "d", emptyString, "Task due date (date). For example, --due-date=\"next friday\"")
	addCmd.Flags().StringVarP(&flags.importance, "importance", "i", "normal", fmt.Sprintf("Task importance - choices: [%v]", strings.Join(importanceChoices, ", ")))
	addCmd.Flags().StringVarP(&flags.status, "status", "s", "not started", fmt.Sprintf("Task status - choices: [%v]", strings.Join(api.GraphStatusOptions, ", ")))

	return addCmd
}

func constructTaskPayload(flags addParamsFlags, title string) (*api.TodoTask, error) {
	task := api.TodoTask{}

	// title
	title = strings.Trim(title, addCutset)
	if len(title) == 0 {
		return nil, errors.New("title is empty")
	}
	task.Title = title

	// reminder
	task.IsReminderOn = false
	reminder := strings.Trim(flags.reminder, addCutset)
	if reminder != emptyString {
		if reminder, err := datetime.DateTimeParser(reminder); err != nil {
			return nil, err
		} else {
			task.IsReminderOn = true
			task.ReminderDateTime = (*datetime.GraphTime)(reminder)
		}
	}

	// due date
	dueDate := strings.Trim(flags.dueDate, addCutset)
	if dueDate != emptyString {
		if dueDate, err := datetime.DateParser(dueDate); err != nil {
			return nil, err
		} else {
			task.DueDateTime = (*datetime.GraphTime)(dueDate)
		}
	}

	// status
	status := strings.Trim(flags.status, addCutset)
	if !utils.ContainsString(api.GraphStatusOptions, status) {
		return nil, fmt.Errorf("'%v' is not a valid value for status", status)
	}
	task.Status = api.GraphStatus(status)

	// importance
	importance := strings.Trim(flags.importance, addCutset)
	if !utils.ContainsString(importanceChoices, importance) {
		return nil, fmt.Errorf("'%v' is not a valid value for importance", importance)
	}
	task.Importance = importance

	return &task, nil
}
