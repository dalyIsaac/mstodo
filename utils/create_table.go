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
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/viper"
)

func CreateBasicTable(header *table.Row) table.Writer {
	return createTable(header, true, nil)
}

func CreateFormattedTable(header *table.Row, columnConfigs *[]table.ColumnConfig) table.Writer {
	return createTable(header, true, columnConfigs)
}

func createTable(header *table.Row, showHeader bool, columnConfigs *[]table.ColumnConfig) table.Writer {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	t.SetStyle(*matchTableStyle(viper.GetString("table-style")))

	if header != nil && showHeader {
		t.AppendHeader(*header)
	}

	if columnConfigs != nil && len(*columnConfigs) != 0 {
		t.SetColumnConfigs(*columnConfigs)
	}

	return t
}

func IsTableStyleValid(style string) bool {
	return matchTableStyle(style) != nil
}

func matchTableStyle(style string) *table.Style {
	switch style {
	case "Default":
		return &table.StyleDefault
	case "Bold":
		return &table.StyleBold
	case "ColoredBright":
		return &table.StyleColoredBright
	case "ColoredDark":
		return &table.StyleColoredDark
	case "ColoredBlackOnBlueWhite":
		return &table.StyleColoredBlackOnBlueWhite
	case "ColoredBlackOnCyanWhite":
		return &table.StyleColoredBlackOnCyanWhite
	case "ColoredBlackOnGreenWhite":
		return &table.StyleColoredBlackOnGreenWhite
	case "ColoredBlackOnMagentaWhite":
		return &table.StyleColoredBlackOnMagentaWhite
	case "ColoredBlackOnYellowWhite":
		return &table.StyleColoredBlackOnYellowWhite
	case "ColoredBlackOnRedWhite":
		return &table.StyleColoredBlackOnRedWhite
	case "ColoredBlueWhiteOnBlack":
		return &table.StyleColoredBlueWhiteOnBlack
	case "ColoredCyanWhiteOnBlack":
		return &table.StyleColoredCyanWhiteOnBlack
	case "ColoredGreenWhiteOnBlack":
		return &table.StyleColoredGreenWhiteOnBlack
	case "ColoredMagentaWhiteOnBlack":
		return &table.StyleColoredMagentaWhiteOnBlack
	case "ColoredRedWhiteOnBlack":
		return &table.StyleColoredRedWhiteOnBlack
	case "ColoredYellowWhiteOnBlack":
		return &table.StyleColoredYellowWhiteOnBlack
	case "Double":
		return &table.StyleDouble
	case "Light":
		return &table.StyleLight
	case "Rounded":
		return &table.StyleRounded
	default:
		return nil
	}
}
