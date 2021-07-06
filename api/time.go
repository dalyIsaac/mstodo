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
	"fmt"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

// Custom time unmarshalling

type msgraphTime time.Time

const msgraphTimeLayout = "2006-01-02T15:04:05.0000000"

// UnmarshalJSON parses the json string into the Microsoft Graph time format,
// as per https://docs.microsoft.com/en-us/graph/api/resources/datetimetimezone?view=graph-rest-1.0
func (ct *msgraphTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	nt, err := time.Parse(msgraphTimeLayout, s)
	*ct = msgraphTime(nt)
	return err
}

// Custom location time unmarshalling

type msgraphLocation time.Location

func (ct *msgraphLocation) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	loc, err := time.LoadLocation(s)
	*ct = msgraphLocation(*loc)
	return err
}

// DateTimeTimeZone based on https://docs.microsoft.com/en-us/graph/api/resources/datetimetimezone?view=graph-rest-1.0
type DateTimeTimeZone struct {
	DateTime msgraphTime     `json:"dateTime"`
	TimeZone msgraphLocation `json:"timeZone"`
}

func (dt *DateTimeTimeZone) String() string {
	t := time.Time(dt.DateTime)
	l := time.Location(dt.TimeZone)
	return fmt.Sprintf("%v (%v)", humanize.Time(t), l.String())
}
