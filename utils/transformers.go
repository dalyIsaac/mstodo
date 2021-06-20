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

import "github.com/jedib0t/go-pretty/v6/text"

func boolToEmoji(v bool) string {
	if v {
		return "✅"
	}
	return "❌"
}

// Cell data transformer for data based on their type
var Transformer = text.Transformer(func(val interface{}) string {
	switch val := val.(type) {
	case bool:
		return boolToEmoji(val)
	case string:
		return val
	default:
		return "unknown type"
	}
})