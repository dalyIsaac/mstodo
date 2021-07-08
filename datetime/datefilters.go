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

package datetime

import "time"

type DateFilters struct {
	Start *time.Time
	End   *time.Time
}

func (filters *DateFilters) Contains(g *GraphTime) bool {
	if filters == nil {
		return true
	}

	if g == nil {
		return false
	}

	t := time.Time(*g)

	if filters.Start != nil {
		if t.Before(*filters.Start) {
			return false
		}
	}

	if filters.End != nil {
		if t.After(*filters.End) {
			return false
		}
	}

	return true
}
