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
	"errors"
	"fmt"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

// Sort flag mode which allows no sorting (unlike tables.SortMode)
const (
	// Asc sorts the column in Ascending order alphabetically.
	asc = "asc"
	// AscNumeric sorts the column in Ascending order numerically.
	ascNumeric = "asc-numeric"
	// Dsc sorts the column in Descending order alphabetically.
	dsc = "dsc"
	// DscNumeric sorts the column in Descending order numerically.
	dscNumeric = "dsc-numeric"
	// Don't sort
	none   = "none"
	NoSort = -1
)

func GetSortOptions() string {
	options := []string{asc, dsc, none}
	return fmt.Sprintf("[%s]", strings.Join(options, ", "))
}

func GetNumericSortOptions() string {
	options := []string{ascNumeric, dscNumeric, none}
	return fmt.Sprintf("[%s]", strings.Join(options, ", "))
}

func GetSortMode(flag string) (table.SortMode, error) {
	return getSortMode(flag, false)
}

func GetNumericSortMode(flag string) (table.SortMode, error) {
	return getSortMode(flag, true)
}

func getSortMode(flag string, numericSort bool) (table.SortMode, error) {
	switch flag {
	case asc:
		return table.Asc, nil
	case ascNumeric:
		return table.AscNumeric, nil
	case dsc:
		return table.Dsc, nil
	case dscNumeric:
		return table.DscNumeric, nil
	case none:
		return -1, nil
	}

	// an error has occurred
	var options string
	if numericSort {
		options = GetNumericSortOptions()
	} else {
		options = GetSortOptions()
	}

	return -1, errors.New("invalid sort option - valid options: " + options)
}
