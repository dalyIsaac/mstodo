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

import (
	"errors"
	"strings"
	"time"
)

const parserCutset = "[] '\""

type parseType int

const (
	dateParseType parseType = iota
	dateTimeParseType
)

var wrapperInstance = (&parserWrapper{now: time.Now})

// Exported date parser
var DateStartEndParser = func(input string) (*DateFilters, error) {
	return wrapperInstance.filterParser(input, dateParseType)
}

// DateTime parser
var DateTimeParser = func(input string) (*time.Time, error) {
	return wrapperInstance.parse(input, 0, dateTimeParseType)
}

var DateParser = func(input string) (*time.Time, error) {
	return wrapperInstance.parse(input, 0, dateParseType)
}

func (parser *parserWrapper) filterParser(input string, parseType parseType) (*DateFilters, error) {
	parts := strings.Split(input, ";")
	filters := DateFilters{}

	if len(parts) == 0 {
		return nil, errors.New("empty filter")
	}

	if len(parts) > 2 {
		return nil, errors.New("too many parts")
	}

	var err error
	var start, end *time.Time
	for _, p := range parts {
		p = strings.Trim(strings.ToLower(p), parserCutset)

		if idx := contains(p, "start"); idx != -1 {
			start, err = parser.parse(p, idx, parseType)
		} else if idx := contains(p, "end"); idx != -1 {
			end, err = parser.parse(p, idx, parseType)
		} else {
			return nil, errors.New("missing qualifier")
		}

		if err != nil {
			return nil, err
		}
	}

	filters.Start = start
	filters.End = end
	return &filters, nil
}

func contains(s string, substring string) int {
	l := len(substring)
	if len(s) < l {
		return -1
	}

	if s[:l] == substring {
		return l
	}
	return -1
}

var dateLayouts = []string{
	"02/Jan/2006",
	"02/JAN/2006",
	"02-Jan-2006",
	"02-JAN-2006",
	"02-Jan-06",
	"02-JAN-06",
	"2/Jan/2006",
	"2/JAN/2006",
	"02/01/2006",
	"2/01/2006",
	"2/01",
	"02-Jan 2006",
	"02-Jan",
	"Jan 02, 2006",
	"Jan 02",
	"Jan 2",
	"January 02, 2006",
}

var timeLayouts = []string{
	"15:04",
	"3:04 PM",
	"3:04 pm",
	"3:04PM",
	"3:04pm",
	"03:04 PM",
	"03:04 pm",
	"03:04PM",
	"03:04pm",
}

func generateDateTimeLayouts() []string {
	layouts := make([]string, len(dateLayouts)*len(timeLayouts)*8)

	for _, d := range dateLayouts {
		for _, v := range timeLayouts {
			layouts = append(layouts, d+" at "+v)
			layouts = append(layouts, d+", at "+v)
			layouts = append(layouts, d+" "+v)
			layouts = append(layouts, d+", "+v)
			layouts = append(layouts, v+" "+d)
			layouts = append(layouts, v+", "+d)
			layouts = append(layouts, v+" on "+d)
			layouts = append(layouts, v+", on "+d)
		}
	}

	return layouts
}

func (parser *parserWrapper) parse(input string, startIdx int, parseType parseType) (*time.Time, error) {
	layouts := dateLayouts

	if parseType == dateTimeParseType {
		layouts = generateDateTimeLayouts()
	}

	input = input[startIdx:]

	for _, layout := range layouts {
		if date, err := time.Parse(layout, input); err == nil {
			return parser.fixYear(date), nil
		}
	}

	if date, err := parser.parseDay(input); err == nil {
		return date, nil
	}

	return nil, errors.New("invalid date")
}

func (parser *parserWrapper) fixYear(date time.Time) *time.Time {
	if date.Year() == 0 {
		date = date.AddDate(parser.now().Year(), 0, 0)
	}
	return &date
}

type relative int

const (
	none relative = iota
	last
	this
	next
)

func (parser *parserWrapper) parseDay(input string) (*time.Time, error) {
	if len(input) < 3 {
		return nil, errors.New("day too short")
	}

	// Get adjective
	relative := none
	start := 0

	if len(input) >= 8 {
		start = 4
		// last, this, next
		switch input[:4] {
		case "last":
			relative = last
		case "this":
			relative = this
		case "next":
			relative = next
		default:
			start = 0
		}
	}

	// Get day of week
	day, err := parser.getDayPart(input, start)
	if err != nil {
		return nil, err
	}

	// Get date based on day of week
	result := parser.getDateFromWeekday(relative, day)

	return &result, nil
}

func (parser *parserWrapper) getDayPart(input string, start int) (time.Weekday, error) {
	input = strings.ToLower(input)

	dayPart := strings.Trim(input[start:], parserCutset)
	if len(dayPart) < 3 {
		return -1, errors.New("day is too short")
	}
	dayPart = dayPart[:3]

	switch dayPart {
	case "mon":
		return time.Monday, nil
	case "tue":
		return time.Tuesday, nil
	case "wed":
		return time.Wednesday, nil
	case "thu":
		return time.Thursday, nil
	case "fri":
		return time.Friday, nil
	case "sat":
		return time.Saturday, nil
	case "sun":
		return time.Sunday, nil
	default:
		return -1, errors.New("invalid day")
	}
}

func (parser *parserWrapper) getDateFromWeekday(relative relative, day time.Weekday) time.Time {
	var result time.Time

	switch relative {
	case none:
		result = parser.getClosestDayInstance(day)
	case this:
		result = parser.getThisDayInstance(day)
	case last:
		result = parser.getLastDayInstance(day)
	case next:
		result = parser.getNextDayInstance(day)
	}

	return result
}

func (parser *parserWrapper) getClosestDayInstance(day time.Weekday) time.Time {
	currentDay := parser.now().Weekday()

	if day <= currentDay {
		return parser.now().AddDate(0, 0, -int(currentDay-day))
	}
	return parser.now().AddDate(0, 0, int(day-currentDay))
}

func (parser *parserWrapper) getThisDayInstance(day time.Weekday) time.Time {
	currentDay := parser.now().Weekday()

	diff := mod(int(day-currentDay), 7)

	return parser.now().AddDate(0, 0, diff)
}

func (parser *parserWrapper) getLastDayInstance(day time.Weekday) time.Time {
	currentDay := parser.now().Weekday()

	diff := mod(int(currentDay-day), 7)
	if diff == 0 {
		diff = 7
	}

	return parser.now().AddDate(0, 0, -diff)
}

func (parser *parserWrapper) getNextDayInstance(day time.Weekday) time.Time {
	currentDay := parser.now().Weekday()

	if day < currentDay && day != 0 {
		return parser.getThisDayInstance(day)
	}

	diff := mod(int(day-currentDay), 7)
	return parser.now().AddDate(0, 0, 7+diff)
}

// mod implements Python-like modulo behavior
func mod(d, m int) int {
	res := d % m
	if (res < 0 && m > 0) || (res > 0 && m < 0) {
		return res + m
	}
	return res
}
