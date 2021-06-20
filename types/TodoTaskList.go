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

package types

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

type TodoTaskList []TodoTaskListItem

type TodoTaskListResponse struct {
	Value TodoTaskList `json:"value"`
}
