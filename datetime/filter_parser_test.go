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
	"reflect"
	"testing"
	"time"
)

func Test_parserWrapper_parser(t *testing.T) {
	type fields struct {
		now func() time.Time
	}
	type args struct {
		input     string
		parseType parseType
	}

	testFields := fields{now: func() time.Time {
		// Today is Wednesday
		return date(7, 7)
	}}

	filterP := func(f DateFilters) *DateFilters {
		return &f
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *DateFilters
		wantErr bool
	}{
		{name: "start end", fields: testFields, wantErr: false, args: args{input: "start Monday; end Friday", parseType: dateParseType}, want: filterP(DateFilters{Start: p(date(5, 7)), End: p(date(9, 7))})},
		{name: "end start", fields: testFields, wantErr: false, args: args{input: "end Friday; start Monday", parseType: dateParseType}, want: filterP(DateFilters{Start: p(date(5, 7)), End: p(date(9, 7))})},
		{name: "only start", fields: testFields, wantErr: false, args: args{input: "start Monday", parseType: dateParseType}, want: filterP(DateFilters{Start: p(date(5, 7))})},
		{name: "only end", fields: testFields, wantErr: false, args: args{input: "end Friday", parseType: dateParseType}, want: filterP(DateFilters{End: p(date(9, 7))})},
		{name: "no qualifier", fields: testFields, wantErr: true, args: args{input: "monday", parseType: dateParseType}, want: nil},
		{name: "force error", fields: testFields, wantErr: true, args: args{input: "start garbage", parseType: dateParseType}, want: nil},
		{name: "too many parts", fields: testFields, wantErr: true, args: args{input: "a;b;c", parseType: dateParseType}, want: nil},
		{name: "empty filter", fields: testFields, wantErr: true, args: args{input: "a;b;c", parseType: dateParseType}, want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := &parserWrapper{
				now: tt.fields.now,
			}
			got, err := parser.filterParser(tt.args.input, tt.args.parseType)
			if (err != nil) != tt.wantErr {
				t.Errorf("parserWrapper.parser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parserWrapper.parser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_contains(t *testing.T) {
	type args struct {
		s         string
		substring string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "'start mon' contains 'start'", args: args{s: "start mon", substring: "start"}, want: 5},
		{name: "'start mon' contains 'end'", args: args{s: "start mon", substring: "end"}, want: -1},
		{name: "too short", args: args{s: "end", substring: "start"}, want: -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := contains(tt.args.s, tt.args.substring); got != tt.want {
				t.Errorf("contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parserWrapper_parseDate(t *testing.T) {
	type fields struct {
		now func() time.Time
	}
	type args struct {
		input string
		parseType parseType
	}

	testFields := fields{now: func() time.Time {
		// Today is Wednesday
		return date(7, 7)
	}}

	wantDate := p(time.Date(2021, time.January, 2, 0, 0, 0, 0, time.UTC))
	wantDatetime := p(time.Date(2021, time.January, 2, 20, 13, 0, 0, time.UTC))

	tests := []struct {
		fields  fields
		args    args
		want    *time.Time
		wantErr bool
	}{
		// dateParseType
		{args: args{input: "02/Jan/2021", parseType: dateParseType}, fields: testFields, want: wantDate, wantErr: false},
		{args: args{input: "02/JAN/2021", parseType: dateParseType}, fields: testFields, want: wantDate, wantErr: false},
		{args: args{input: "02-Jan-2021", parseType: dateParseType}, fields: testFields, want: wantDate, wantErr: false},
		{args: args{input: "02-JAN-2021", parseType: dateParseType}, fields: testFields, want: wantDate, wantErr: false},
		{args: args{input: "02-Jan-21", parseType: dateParseType}, fields: testFields, want: wantDate, wantErr: false},
		{args: args{input: "02-JAN-21", parseType: dateParseType}, fields: testFields, want: wantDate, wantErr: false},
		{args: args{input: "2/Jan/2021", parseType: dateParseType}, fields: testFields, want: wantDate, wantErr: false},
		{args: args{input: "2/JAN/2021", parseType: dateParseType}, fields: testFields, want: wantDate, wantErr: false},
		{args: args{input: "02/01/2021", parseType: dateParseType}, fields: testFields, want: wantDate, wantErr: false},
		{args: args{input: "2/01/2021", parseType: dateParseType}, fields: testFields, want: wantDate, wantErr: false},
		{args: args{input: "2/01", parseType: dateParseType}, fields: testFields, want: wantDate, wantErr: false},
		{args: args{input: "02-Jan 2021", parseType: dateParseType}, fields: testFields, want: wantDate, wantErr: false},
		{args: args{input: "02-Jan", parseType: dateParseType}, fields: testFields, want: wantDate, wantErr: false},
		{args: args{input: "Jan 02, 2021", parseType: dateParseType}, fields: testFields, want: wantDate, wantErr: false},
		{args: args{input: "Jan 02", parseType: dateParseType}, fields: testFields, want: wantDate, wantErr: false},
		{args: args{input: "Jan 2", parseType: dateParseType}, fields: testFields, want: wantDate, wantErr: false},
		{args: args{input: "January 02, 2021", parseType: dateParseType}, fields: testFields, want: wantDate, wantErr: false},
		{args: args{input: "last Mon", parseType: dateParseType}, fields: testFields, want: p(date(5, 7)), wantErr: false},
		{args: args{input: "Monday", parseType: dateParseType}, fields: testFields, want: p(date(5, 7)), wantErr: false},
		{args: args{input: "garbage", parseType: dateParseType}, fields: testFields, want: nil, wantErr: true},

		// datetimeParseType
		{args: args{input: "02/Jan/2021 at 20:13", parseType: datetimeParseType}, fields: testFields, want: wantDatetime, wantErr: false},
		{args: args{input: "02/JAN/2021, at 20:13", parseType: datetimeParseType}, fields: testFields, want: wantDatetime, wantErr: false},
		{args: args{input: "02-Jan-2021 20:13", parseType: datetimeParseType}, fields: testFields, want: wantDatetime, wantErr: false},
		{args: args{input: "02-JAN-2021, 20:13", parseType: datetimeParseType}, fields: testFields, want: wantDatetime, wantErr: false},
		{args: args{input: "20:13 02-Jan-21", parseType: datetimeParseType}, fields: testFields, want: wantDatetime, wantErr: false},
		{args: args{input: "20:13, 02-JAN-21", parseType: datetimeParseType}, fields: testFields, want: wantDatetime, wantErr: false},
		{args: args{input: "2/Jan/2021 at 08:13 PM", parseType: datetimeParseType}, fields: testFields, want: wantDatetime, wantErr: false},
		{args: args{input: "2/JAN/2021 at 08:13 pm", parseType: datetimeParseType}, fields: testFields, want: wantDatetime, wantErr: false},
		{args: args{input: "02/01/2021 at 08:13PM", parseType: datetimeParseType}, fields: testFields, want: wantDatetime, wantErr: false},
		{args: args{input: "2/01/2021 at 08:13pm", parseType: datetimeParseType}, fields: testFields, want: wantDatetime, wantErr: false},
		{args: args{input: "2/01, 08:13pm", parseType: datetimeParseType}, fields: testFields, want: wantDatetime, wantErr: false},
		{args: args{input: "8:13PM, 02-Jan 2021", parseType: datetimeParseType}, fields: testFields, want: wantDatetime, wantErr: false},
		{args: args{input: "8:13pm, 02-Jan", parseType: datetimeParseType}, fields: testFields, want: wantDatetime, wantErr: false},
		{args: args{input: "08:13PM on Jan 02, 2021", parseType: datetimeParseType}, fields: testFields, want: wantDatetime, wantErr: false},
		{args: args{input: "08:13 pm, on Jan 02", parseType: datetimeParseType}, fields: testFields, want: wantDatetime, wantErr: false},
		{args: args{input: "Jan 2 8:13PM", parseType: datetimeParseType}, fields: testFields, want: wantDatetime, wantErr: false},
		{args: args{input: "January 02, 2021, 8:13pm", parseType: datetimeParseType}, fields: testFields, want: wantDatetime, wantErr: false},
		{args: args{input: "last Mon", parseType: datetimeParseType}, fields: testFields, want: p(date(5, 7)), wantErr: false},
		{args: args{input: "Monday", parseType: datetimeParseType}, fields: testFields, want: p(date(5, 7)), wantErr: false},
		{args: args{input: "garbage", parseType: datetimeParseType}, fields: testFields, want: nil, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.args.input, func(t *testing.T) {
			parser := &parserWrapper{
				now: tt.fields.now,
			}
			got, err := parser.parse(tt.args.input, 0, tt.args.parseType)
			if (err != nil) != tt.wantErr {
				t.Errorf("parserWrapper.parseDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parserWrapper.parseDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parserWrapper_fixYear(t *testing.T) {
	type fields struct {
		now func() time.Time
	}
	type args struct {
		date time.Time
	}

	testFields := fields{now: func() time.Time {
		// Today is Wednesday
		return date(7, 7)
	}}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   *time.Time
	}{
		{name: "0 year", fields: testFields, args: args{date: time.Date(0, 7, 7, 0, 0, 0, 0, time.UTC)}, want: p(time.Date(2021, 7, 7, 0, 0, 0, 0, time.UTC))},
		{name: "2021 year", fields: testFields, args: args{date: time.Date(2021, 7, 7, 0, 0, 0, 0, time.UTC)}, want: p(time.Date(2021, 7, 7, 0, 0, 0, 0, time.UTC))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := &parserWrapper{
				now: tt.fields.now,
			}
			if got := parser.fixYear(tt.args.date); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parserWrapper.fixYear() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parserWrapper_parseDay(t *testing.T) {
	type fields struct {
		now func() time.Time
	}
	type args struct {
		input string
	}

	testFields := fields{now: func() time.Time {
		// Today is Wednesday
		return date(7, 7)
	}}

	tests := []struct {
		fields  fields
		args    args
		want    *time.Time
		wantErr bool
	}{
		// failures
		{fields: testFields, args: args{input: "mo"}, want: nil, wantErr: true},
		{fields: testFields, args: args{input: "day"}, want: nil, wantErr: true},
		{fields: testFields, args: args{input: "last"}, want: nil, wantErr: true},

		// no relative specifier
		{fields: testFields, args: args{input: "Monday"}, wantErr: false, want: p(date(5, 7))},
		{fields: testFields, args: args{input: "mon"}, wantErr: false, want: p(date(5, 7))},
		{fields: testFields, args: args{input: "Tuesday"}, wantErr: false, want: p(date(6, 7))},
		{fields: testFields, args: args{input: "tue"}, wantErr: false, want: p(date(6, 7))},
		{fields: testFields, args: args{input: "Wednesday"}, wantErr: false, want: p(date(7, 7))},
		{fields: testFields, args: args{input: "wed"}, wantErr: false, want: p(date(7, 7))},
		{fields: testFields, args: args{input: "Thursday"}, wantErr: false, want: p(date(8, 7))},
		{fields: testFields, args: args{input: "thu"}, wantErr: false, want: p(date(8, 7))},
		{fields: testFields, args: args{input: "Friday"}, wantErr: false, want: p(date(9, 7))},
		{fields: testFields, args: args{input: "fri"}, wantErr: false, want: p(date(9, 7))},
		{fields: testFields, args: args{input: "Saturday"}, wantErr: false, want: p(date(10, 7))},
		{fields: testFields, args: args{input: "sat"}, wantErr: false, want: p(date(10, 7))},
		{fields: testFields, args: args{input: "Sunday"}, wantErr: false, want: p(date(4, 7))},
		{fields: testFields, args: args{input: "sun"}, wantErr: false, want: p(date(4, 7))},

		// this relative specifier
		{fields: testFields, args: args{input: "this Monday"}, wantErr: false, want: p(date(12, 7))},
		{fields: testFields, args: args{input: "this mon"}, wantErr: false, want: p(date(12, 7))},
		{fields: testFields, args: args{input: "this Tuesday"}, wantErr: false, want: p(date(13, 7))},
		{fields: testFields, args: args{input: "this tue"}, wantErr: false, want: p(date(13, 7))},
		{fields: testFields, args: args{input: "this Wednesday"}, wantErr: false, want: p(date(7, 7))},
		{fields: testFields, args: args{input: "this wed"}, wantErr: false, want: p(date(7, 7))},
		{fields: testFields, args: args{input: "this Thursday"}, wantErr: false, want: p(date(8, 7))},
		{fields: testFields, args: args{input: "this thu"}, wantErr: false, want: p(date(8, 7))},
		{fields: testFields, args: args{input: "this Friday"}, wantErr: false, want: p(date(9, 7))},
		{fields: testFields, args: args{input: "this fri"}, wantErr: false, want: p(date(9, 7))},
		{fields: testFields, args: args{input: "this Saturday"}, wantErr: false, want: p(date(10, 7))},
		{fields: testFields, args: args{input: "this sat"}, wantErr: false, want: p(date(10, 7))},
		{fields: testFields, args: args{input: "this Sunday"}, wantErr: false, want: p(date(11, 7))},
		{fields: testFields, args: args{input: "this sun"}, wantErr: false, want: p(date(11, 7))},

		// last relative specifier
		{fields: testFields, args: args{input: "last Monday"}, wantErr: false, want: p(date(5, 7))},
		{fields: testFields, args: args{input: "last mon"}, wantErr: false, want: p(date(5, 7))},
		{fields: testFields, args: args{input: "last Tuesday"}, wantErr: false, want: p(date(6, 7))},
		{fields: testFields, args: args{input: "last tue"}, wantErr: false, want: p(date(6, 7))},
		{fields: testFields, args: args{input: "last Wednesday"}, wantErr: false, want: p(date(30, 6))},
		{fields: testFields, args: args{input: "last wed"}, wantErr: false, want: p(date(30, 6))},
		{fields: testFields, args: args{input: "last Thursday"}, wantErr: false, want: p(date(1, 7))},
		{fields: testFields, args: args{input: "last thu"}, wantErr: false, want: p(date(1, 7))},
		{fields: testFields, args: args{input: "last Friday"}, wantErr: false, want: p(date(2, 7))},
		{fields: testFields, args: args{input: "last fri"}, wantErr: false, want: p(date(2, 7))},
		{fields: testFields, args: args{input: "last Saturday"}, wantErr: false, want: p(date(3, 7))},
		{fields: testFields, args: args{input: "last sat"}, wantErr: false, want: p(date(3, 7))},
		{fields: testFields, args: args{input: "last Sunday"}, wantErr: false, want: p(date(4, 7))},
		{fields: testFields, args: args{input: "last sun"}, wantErr: false, want: p(date(4, 7))},

		// next relative specifier
		{fields: testFields, args: args{input: "next Monday"}, wantErr: false, want: p(date(12, 7))},
		{fields: testFields, args: args{input: "next mon"}, wantErr: false, want: p(date(12, 7))},
		{fields: testFields, args: args{input: "next Tuesday"}, wantErr: false, want: p(date(13, 7))},
		{fields: testFields, args: args{input: "next tue"}, wantErr: false, want: p(date(13, 7))},
		{fields: testFields, args: args{input: "next Wednesday"}, wantErr: false, want: p(date(14, 7))},
		{fields: testFields, args: args{input: "next wed"}, wantErr: false, want: p(date(14, 7))},
		{fields: testFields, args: args{input: "next Thursday"}, wantErr: false, want: p(date(15, 7))},
		{fields: testFields, args: args{input: "next thu"}, wantErr: false, want: p(date(15, 7))},
		{fields: testFields, args: args{input: "next Friday"}, wantErr: false, want: p(date(16, 7))},
		{fields: testFields, args: args{input: "next fri"}, wantErr: false, want: p(date(16, 7))},
		{fields: testFields, args: args{input: "next Saturday"}, wantErr: false, want: p(date(17, 7))},
		{fields: testFields, args: args{input: "next sat"}, wantErr: false, want: p(date(17, 7))},
		{fields: testFields, args: args{input: "next Sunday"}, wantErr: false, want: p(date(18, 7))},
		{fields: testFields, args: args{input: "next sun"}, wantErr: false, want: p(date(18, 7))},
	}
	for _, tt := range tests {
		t.Run(tt.args.input, func(t *testing.T) {
			parser := &parserWrapper{
				now: tt.fields.now,
			}
			got, err := parser.parseDay(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parserWrapper.parseDay() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parserWrapper.parseDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parserWrapper_getDayPart(t *testing.T) {
	type fields struct {
		now func() time.Time
	}
	type args struct {
		input string
		start int
	}

	testFields := fields{now: func() time.Time {
		// Today is Wednesday
		return date(7, 7)
	}}

	tests := []struct {
		fields  fields
		args    args
		want    time.Weekday
		wantErr bool
	}{
		{fields: testFields, args: args{input: "mo", start: 0}, want: -1, wantErr: true},
		{fields: testFields, args: args{input: "day", start: 0}, want: -1, wantErr: true},
		{fields: testFields, args: args{input: "Monday", start: 0}, want: time.Monday, wantErr: false},
		{fields: testFields, args: args{input: "Tuesday", start: 0}, want: time.Tuesday, wantErr: false},
		{fields: testFields, args: args{input: "Wednesday", start: 0}, want: time.Wednesday, wantErr: false},
		{fields: testFields, args: args{input: "Thursday", start: 0}, want: time.Thursday, wantErr: false},
		{fields: testFields, args: args{input: "Friday", start: 0}, want: time.Friday, wantErr: false},
		{fields: testFields, args: args{input: "Saturday", start: 0}, want: time.Saturday, wantErr: false},
		{fields: testFields, args: args{input: "Sunday", start: 0}, want: time.Sunday, wantErr: false},
		{fields: testFields, args: args{input: "mon", start: 0}, want: time.Monday, wantErr: false},
		{fields: testFields, args: args{input: "tue", start: 0}, want: time.Tuesday, wantErr: false},
		{fields: testFields, args: args{input: "wed", start: 0}, want: time.Wednesday, wantErr: false},
		{fields: testFields, args: args{input: "thu", start: 0}, want: time.Thursday, wantErr: false},
		{fields: testFields, args: args{input: "fri", start: 0}, want: time.Friday, wantErr: false},
		{fields: testFields, args: args{input: "sat", start: 0}, want: time.Saturday, wantErr: false},
		{fields: testFields, args: args{input: "sun", start: 0}, want: time.Sunday, wantErr: false},
		{fields: testFields, args: args{input: "the date is mon", start: 12}, want: time.Monday, wantErr: false},
		{fields: testFields, args: args{input: "the date is tue", start: 12}, want: time.Tuesday, wantErr: false},
		{fields: testFields, args: args{input: "the date is wed", start: 12}, want: time.Wednesday, wantErr: false},
		{fields: testFields, args: args{input: "the date is thu", start: 12}, want: time.Thursday, wantErr: false},
		{fields: testFields, args: args{input: "the date is fri", start: 12}, want: time.Friday, wantErr: false},
		{fields: testFields, args: args{input: "the date is sat", start: 12}, want: time.Saturday, wantErr: false},
		{fields: testFields, args: args{input: "the date is sun", start: 12}, want: time.Sunday, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.args.input, func(t *testing.T) {
			parser := &parserWrapper{
				now: tt.fields.now,
			}
			got, err := parser.getDayPart(tt.args.input, tt.args.start)
			if (err != nil) != tt.wantErr {
				t.Errorf("parserWrapper.getDayPart() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parserWrapper.getDayPart() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parserWrapper_getDateFromWeekday(t *testing.T) {
	type fields struct {
		now func() time.Time
	}
	type args struct {
		relative relative
		day      time.Weekday
	}

	testFields := fields{now: func() time.Time {
		// Today is Wednesday
		return date(7, 7)
	}}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   time.Time
	}{
		{name: "none wednesday", fields: testFields, args: args{relative: none, day: time.Wednesday}, want: date(7, 7)},
		{name: "this wednesday", fields: testFields, args: args{relative: this, day: time.Wednesday}, want: date(7, 7)},
		{name: "this friday", fields: testFields, args: args{relative: this, day: time.Friday}, want: date(9, 7)},
		{name: "next friday", fields: testFields, args: args{relative: next, day: time.Friday}, want: date(16, 7)},
		{name: "last friday", fields: testFields, args: args{relative: last, day: time.Friday}, want: date(2, 7)},
		{name: "last tuesday", fields: testFields, args: args{relative: last, day: time.Tuesday}, want: date(6, 7)},
		{name: "next tuesday", fields: testFields, args: args{relative: next, day: time.Tuesday}, want: date(13, 7)},
		{name: "next thursday", fields: testFields, args: args{relative: next, day: time.Thursday}, want: date(15, 7)},
		{name: "this thursday", fields: testFields, args: args{relative: this, day: time.Thursday}, want: date(8, 7)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := &parserWrapper{
				now: tt.fields.now,
			}
			if got := parser.getDateFromWeekday(tt.args.relative, tt.args.day); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parserWrapper.getDateFromWeekday() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parserWrapper_getClosestDayInstance(t *testing.T) {
	type fields struct {
		now func() time.Time
	}
	type args struct {
		day time.Weekday
	}

	testFields := fields{now: func() time.Time {
		// Today is Wednesday
		return date(7, 7)
	}}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   time.Time
	}{
		{name: "tuesday", fields: testFields, args: args{day: time.Tuesday}, want: date(6, 7)},
		{name: "wednesday", fields: testFields, args: args{day: time.Wednesday}, want: date(7, 7)},
		{name: "friday", fields: testFields, args: args{day: time.Friday}, want: date(9, 7)},
		{name: "saturday", fields: testFields, args: args{day: time.Saturday}, want: date(10, 7)},
		{name: "sunday", fields: testFields, args: args{day: time.Sunday}, want: date(4, 7)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := &parserWrapper{
				now: tt.fields.now,
			}
			if got := parser.getClosestDayInstance(tt.args.day); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parserWrapper.getClosestDayInstance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parserWrapper_getThisDayInstance(t *testing.T) {
	type fields struct {
		now func() time.Time
	}
	type args struct {
		day time.Weekday
	}

	testFields := fields{now: func() time.Time {
		// Today is Wednesday
		return date(7, 7)
	}}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   time.Time
	}{
		{name: "this Friday", fields: testFields, args: args{day: time.Friday}, want: date(9, 7)},
		{name: "this Tuesday", fields: testFields, args: args{day: time.Tuesday}, want: date(13, 7)},
		{name: "this Wednesday", fields: testFields, args: args{day: time.Wednesday}, want: date(7, 7)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := &parserWrapper{
				now: tt.fields.now,
			}
			if got := parser.getThisDayInstance(tt.args.day); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parserWrapper.getPreviousDayInstance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parserWrapper_getLastDayInstance(t *testing.T) {
	type fields struct {
		now func() time.Time
	}
	type args struct {
		day time.Weekday
	}

	testFields := fields{now: func() time.Time {
		// Today is Wednesday
		return date(7, 7)
	}}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   time.Time
	}{
		{name: "last Friday", fields: testFields, args: args{day: time.Friday}, want: date(2, 7)},
		{name: "last Tuesday", fields: testFields, args: args{day: time.Tuesday}, want: date(6, 7)},
		{name: "last Wednesday", fields: testFields, args: args{day: time.Wednesday}, want: date(30, 6)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := &parserWrapper{
				now: tt.fields.now,
			}
			if got := parser.getLastDayInstance(tt.args.day); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parserWrapper.getPreviousDayInstance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parserWrapper_getNextDayInstance(t *testing.T) {
	type fields struct {
		now func() time.Time
	}
	type args struct {
		day time.Weekday
	}

	testFields := fields{now: func() time.Time {
		// Today is Wednesday
		return date(7, 7)
	}}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   time.Time
	}{
		{name: "next Friday", fields: testFields, args: args{day: time.Friday}, want: date(16, 7)},
		{name: "next Tuesday", fields: testFields, args: args{day: time.Tuesday}, want: date(13, 7)},
		{name: "next Wednesday", fields: testFields, args: args{day: time.Wednesday}, want: date(14, 7)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := &parserWrapper{
				now: tt.fields.now,
			}
			if got := parser.getNextDayInstance(tt.args.day); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parserWrapper.getNextDayInstance() = %v, want %v", got, tt.want)
			}
		})
	}
}

// p takes a `Time` instance and returns its pointer
func p(t time.Time) *time.Time {
	return &t
}

func date(day int, month time.Month) time.Time {
	return time.Date(2021, month, day, 0, 0, 0, 0, time.UTC)
}

func Test_mod(t *testing.T) {
	type args struct {
		d int
		m int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "4 % 7", args: args{d: 4, m: 7}, want: 4},
		{name: "4 % 7", args: args{d: 4, m: 7}, want: 4},
		{name: "-3 % 7", args: args{d: -3, m: 7}, want: 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mod(tt.args.d, tt.args.m); got != tt.want {
				t.Errorf("mod() = %v, want %v", got, tt.want)
			}
		})
	}
}
