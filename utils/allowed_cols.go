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

package utils

import (
	"fmt"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

const ignoredColChars = " [{}]"

func GetAllowedColumns(colStr string, columnConfigs []table.ColumnConfig) ([]table.ColumnConfig, error) {
	// If no columns have been ignored, return all columnConfigs
	if colStr = strings.Trim(colStr, ignoredColChars); colStr == "" {
		return columnConfigs, nil
	}

	allowed := []table.ColumnConfig{}

	// Get the list of excluded names
	excluded := []string{}
	for _, s := range strings.Split(colStr, ",") {
		s = strings.Trim(s, ignoredColChars)
		if !columnsContainName(columnConfigs, s) {
			return allowed, fmt.Errorf("column '%s' is not a valid column to exclude", s)
		}
		excluded = append(excluded, s)
	}

	// Get the list of allowed column configs
	for _, c := range columnConfigs {
		if !ContainsString(excluded, c.Name) {
			allowed = append(allowed, c)
		}
	}

	return allowed, nil
}

// Checks to see if the []ColumnConfig cols contains the name.
// c.Name and name are compared with lowercase.
func columnsContainName(cols []table.ColumnConfig, name string) bool {
	name = strings.ToLower(name)

	for _, c := range cols {
		if strings.ToLower(c.Name) == name {
			return true
		}
	}
	return false
}
