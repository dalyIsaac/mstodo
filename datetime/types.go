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
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Custom time unmarshalling

type dateTimeTimeZoneTime time.Time

const msgraphTimeLayout = "2006-01-02T15:04:05.0000000"

// UnmarshalJSON parses the json string into the Microsoft Graph time format,
// as per https://docs.microsoft.com/en-us/graph/api/resources/datetimetimezone?view=graph-rest-1.0
func (ct *dateTimeTimeZoneTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	nt, err := time.Parse(msgraphTimeLayout, s)
	*ct = dateTimeTimeZoneTime(nt)
	return err
}

// Custom location time unmarshalling

type dateTimeTimeZoneLocation time.Location

func (ct *dateTimeTimeZoneLocation) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	loc, err := time.LoadLocation(s)
	*ct = dateTimeTimeZoneLocation(*loc)
	return err
}

// DateTimeTimeZone based on https://docs.microsoft.com/en-us/graph/api/resources/datetimetimezone?view=graph-rest-1.0
type dateTimeTimeZone struct {
	DateTime dateTimeTimeZoneTime     `json:"dateTime"`
	TimeZone dateTimeTimeZoneLocation `json:"timeZone"`
}

// GraphTime
type GraphTime time.Time

func (t *GraphTime) UnmarshalJSON(b []byte) (err error) {
	var dt dateTimeTimeZone
	if err := json.Unmarshal(b, &dt); err != nil {
		return err
	}

	date := time.Time(dt.DateTime)
	year, month, day := date.Date()
	hour, min, sec := date.Clock()
	nsec := date.Nanosecond()
	loc := time.Location(dt.TimeZone)

	*t = GraphTime(time.Date(year, month, day, hour, min, sec, nsec, &loc))
	return nil
}

type GraphTimeMarshal struct {
	DateTime string `json:"dateTime"`
	TimeZone string `json:"timeZone"`
}

// Marshal is called by TodoTask.MarshalJSON
func (t *GraphTime) Marshal() *GraphTimeMarshal {
	ct := time.Time(*t)

	// Get dateTime
	year, month, day := ct.Date()
	hour, minute, second := ct.Clock()
	nanosecond := ct.Nanosecond()

	dateStr := fmt.Sprintf("%04d-%02d-%02d", year, month, day)
	timeStr := fmt.Sprintf("%04d:%02d:%02d.%07d", hour, minute, second, nanosecond)
	dateTime := fmt.Sprintf("%vT%v", dateStr, timeStr)

	// Get timeZone
	timeZone := ct.Location().String()

	// Marshal
	result := GraphTimeMarshal{DateTime: dateTime, TimeZone: timeZone}
	return &result
}
