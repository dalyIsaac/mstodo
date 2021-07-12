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
	"regexp"
	"strings"
	"time"
	_ "time/tzdata"
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

func insert(slice []string, value string, idx int) ([]string, int) {
	slice[idx] = value
	return slice, idx + 1
}

func generateDateTimeLayouts() []string {
	layouts := make([]string, len(dateLayouts)*len(timeLayouts)*8)

	idx := 0
	for _, d := range dateLayouts {
		for _, v := range timeLayouts {
			layouts, idx = insert(layouts, d+" at "+v, idx)
			layouts, idx = insert(layouts, d+", at "+v, idx)
			layouts, idx = insert(layouts, d+" "+v, idx)
			layouts, idx = insert(layouts, d+", "+v, idx)
			layouts, idx = insert(layouts, v+" "+d, idx)
			layouts, idx = insert(layouts, v+", "+d, idx)
			layouts, idx = insert(layouts, v+" on "+d, idx)
			layouts, idx = insert(layouts, v+", on "+d, idx)
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
		if parseType == dateParseType {
			return date, nil
		}

		// parsing datetime
		parsedTime, err := parser.getTime(input)
		if err != nil {
			return nil, err
		}

		// combine parsed time and date
		combined := parser.combineDateTime(date, parsedTime)
		return combined, nil
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

	input = strings.ToLower(input)

	if strings.Contains(input, "last") {
		relative = last
	} else if strings.Contains(input, "this") {
		relative = this
	} else if strings.Contains(input, "next") {
		relative = next
	}

	// Get day of week
	day, err := parser.getDayPart(input)
	if err != nil {
		return nil, err
	}

	// Get date based on day of week
	result := parser.getDateFromWeekday(relative, day)

	return &result, nil
}

func (parser *parserWrapper) getDayPart(input string) (time.Weekday, error) {
	input = strings.ToLower(input)

	if strings.Contains(input, "mon") {
		return time.Monday, nil
	} else if strings.Contains(input, "tue") {
		return time.Tuesday, nil
	} else if strings.Contains(input, "wed") {
		return time.Wednesday, nil
	} else if strings.Contains(input, "thu") {
		return time.Thursday, nil
	} else if strings.Contains(input, "fri") {
		return time.Friday, nil
	} else if strings.Contains(input, "sat") {
		return time.Saturday, nil
	} else if strings.Contains(input, "sun") {
		return time.Sunday, nil
	}

	return -1, errors.New("invalid day")
}

var timeRegexp = regexp.MustCompile(`[0-9]{1,2}:{0,1}[0-9]{1,2} *(am|pm){0,1}`)

func (parser *parserWrapper) getTime(input string) (*time.Time, error) {
	input = strings.ToLower(input)
	match := timeRegexp.FindStringSubmatch(input)

	if len(match) == 0 {
		return nil, errors.New("invalid time")
	}

	for _, layout := range timeLayouts {
		if res, err := time.Parse(layout, match[0]); err == nil {
			return &res, nil
		}
	}

	return nil, errors.New("invalid time")
}

func (parser *parserWrapper) combineDateTime(date *time.Time, parsedTime *time.Time) *time.Time {
	combined := time.Date(date.Year(), date.Month(), date.Day(), parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(), parsedTime.Nanosecond(), date.Location())
	return &combined
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
