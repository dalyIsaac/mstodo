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
	"encoding/json"

	"github.com/iancoleman/strcase"
)

// Used for unmarshalling camelcase into normal text
type GraphStatus string

func (status *GraphStatus) UnmarshalJSON(b []byte) (err error) {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}

	*status = GraphStatus(strcase.ToDelimited(str, ' '))
	return nil
}

// Marshal is called by TodoTask.MarshalJSON
func (status *GraphStatus) Marshal() (string) {
	return strcase.ToLowerCamel(string(*status))
}

var GraphStatusOptions = []string{"not started", "in progress", "completed", "waiting on others", "deferred"}
