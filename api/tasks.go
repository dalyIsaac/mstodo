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

package api

import (
	"fmt"
	"time"

	"github.com/dalyisaac/mstodo/datetime"
)

type TodoTask struct {
	Id                   string              `json:"id"`
	Title                string              `json:"title"`
	Importance           string              `json:"importance"`
	IsReminderOn         bool                `json:"isReminderOn"`
	Status               GraphStatus         `json:"status"`
	ReminderDateTime     *datetime.GraphTime `json:"reminderDateTime"`
	DueDateTime          *datetime.GraphTime `json:"dueDateTime"`
	Completed            *datetime.GraphTime `json:"completedDateTime"`
	CreatedDateTime      time.Time           `json:"createdDateTime"`
	LastModifiedDateTime time.Time           `json:"lastModifiedDateTime"`
}

type TodoTaskList []TodoTask

type todoTaskListResponse struct {
	Value TodoTaskList `json:"value"`
}

func GetTasks(listId string) (*TodoTaskList, error) {
	// Create request
	req, err := CreateRequest()
	if err != nil {
		return nil, err
	}

	// Get request
	url := fmt.Sprintf("/me/todo/lists/%v/tasks", listId)
	resp, err := req.SetResult(&todoTaskListResponse{}).Get(url)
	if err != nil {
		return nil, err
	}

	tasks := resp.Result().(*todoTaskListResponse).Value
	return &tasks, nil
}

func CreateTask(listId string, task *TodoTask) error {
	// Create request
	req, err := CreateRequest()
	if err != nil {
		return err
	}

	// Post request
	url := fmt.Sprintf("/me/todo/lists/%v/tasks", listId)
	resp, err := req.SetBody(task).Post(url)
	if err != nil {
		return err
	}

	if code := resp.StatusCode(); code != 201 {
		return fmt.Errorf("http code: %v\n%v", code, string(resp.Body()))
	}

	return nil
}
