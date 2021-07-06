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
	"errors"
	"strings"
)

type TodoTaskListItem struct {
	// The name of the task list.
	DisplayName string `json:"displayName"`

	// The identifier of the task list, unique in the user's mailbox
	// Read-only.
	// Inherited from entity
	Id string `json:"id"`

	// True if the user is owner of the given task list.
	IsOwner bool `json:"isOwner"`

	// True if the task list is shared with other users
	IsShared bool `json:"isShared"`

	// Property indicating the list name if the given list is a well-known list.
	// Possible values are: none, defaultList, flaggedEmails, unknownFutureValue.
	WellknownListName string `json:"wellknownListName"`
}

// List of TodoTaskList items
type TodoTaskListList []TodoTaskListItem

type todoTaskListListResponse struct {
	Value TodoTaskListList `json:"value"`
}

func (l *TodoTaskListList) GetListId(name string) (string, error) {
	name = strings.ToLower(name)
	for _, item := range *l {
		if strings.ToLower(item.DisplayName) == name {
			return item.Id, nil
		}
	}

	return "", errors.New("could not find name '" + name + "'")
}

func GetLists() (*TodoTaskListList, error) {
	// Create request
	req, err := CreateRequest()
	if err != nil {
		return nil, err
	}

	// Get request
	resp, err := req.SetResult(&todoTaskListListResponse{}).Get("/me/todo/lists")
	if err != nil {
		return nil, err
	}

	lists := resp.Result().(*todoTaskListListResponse).Value
	return &lists, nil
}
