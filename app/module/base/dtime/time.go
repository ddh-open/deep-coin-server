package dtime

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"github.com/pkg/errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	// Refer: http://php.net/manual/en/function.date.php
	formats = map[byte]string{
		'd': "02",                        // Day:    Day of the month, 2 digits with leading zeros. Eg: 01 to 31.
		'D': "Mon",                       // Day:    A textual representation of a day, three letters. Eg: Mon through Sun.
		'w': "Monday",                    // Day:    Numeric representation of the day of the week. Eg: 0 (for Sunday) through 6 (for Saturday).
		'N': "Monday",                    // Day:    ISO-8601 numeric representation of the day of the week. Eg: 1 (for Monday) through 7 (for Sunday).
		'j': "=j=02",                     // Day:    Day of the month without leading zeros. Eg: 1 to 31.
		'S': "02",                        // Day:    English ordinal suffix for the day of the month, 2 characters. Eg: st, nd, rd or th. Works well with j.
		'l': "Monday",                    // Day:    A full textual representation of the day of the week. Eg: Sunday through Saturday.
		'z': "",                          // Day:    The day of the year (starting from 0). Eg: 0 through 365.
		'W': "",                          // Week:   ISO-8601 week number of year, weeks starting on Monday. Eg: 42 (the 42nd week in the year).
		'F': "January",                   // Month:  A full textual representation of a month, such as January or March. Eg: January through December.
		'm': "01",                        // Month:  Numeric representation of a month, with leading zeros. Eg: 01 through 12.
		'M': "Jan",                       // Month:  A short textual representation of a month, three letters. Eg: Jan through Dec.
		'n': "1",                         // Month:  Numeric representation of a month, without leading zeros. Eg: 1 through 12.
		't': "",                          // Month:  Number of days in the given month. Eg: 28 through 31.
		'Y': "2006",                      // Year:   A full numeric representation of a year, 4 digits. Eg: 1999 or 2003.
		'y': "06",                        // Year:   A two digit representation of a year. Eg: 99 or 03.
		'a': "pm",                        // Time:   Lowercase Ante meridiem and Post meridiem. Eg: am or pm.
		'A': "PM",                        // Time:   Uppercase Ante meridiem and Post meridiem. Eg: AM or PM.
		'g': "3",                         // Time:   12-hour format of an hour without leading zeros. Eg: 1 through 12.
		'G': "=G=15",                     // Time:   24-hour format of an hour without leading zeros. Eg: 0 through 23.
		'h': "03",                        // Time:   12-hour format of an hour with leading zeros. Eg: 01 through 12.
		'H': "15",                        // Time:   24-hour format of an hour with leading zeros. Eg: 00 through 23.
		'i': "04",                        // Time:   Minutes with leading zeros. Eg: 00 to 59.
		's': "05",                        // Time:   Seconds with leading zeros. Eg: 00 through 59.
		'u': "=u=.000",                   // Time:   Milliseconds. Eg: 234, 678.
		'U': "",                          // Time:   Seconds since the Unix Epoch (January 1 1970 00:00:00 GMT).
		'O': "-0700",                     // Zone:   Difference to Greenwich time (GMT) in hours. Eg: +0200.
		'P': "-07:00",                    // Zone:   Difference to Greenwich time (GMT) with colon between hours and minutes. Eg: +02:00.
		'T': "MST",                       // Zone:   Timezone abbreviation. Eg: UTC, EST, MDT ...
		'c': "2006-01-02T15:04:05-07:00", // Format: ISO 8601 date. Eg: 2004-02-12T15:19:21+00:00.
		'r': "Mon, 02 Jan 06 15:04 MST",  // Format: RFC 2822 formatted date. Eg: Thu, 21 Dec 2000 16:01:07 +0200.
	}

	timeRegex1, _ = regexp.Compile(timeRegexPattern1)
	timeRegex2, _ = regexp.Compile(timeRegexPattern2)
	timeRegex3, _ = regexp.Compile(timeRegexPattern3)

	// Month words to arabic numerals mapping.
	monthMap = map[string]int{
		"jan":       1,
		"feb":       2,
		"mar":       3,
		"apr":       4,
		"may":       5,
		"jun":       6,
		"jul":       7,
		"aug":       8,
		"sep":       9,
		"sept":      9,
		"oct":       10,
		"nov":       11,
		"dec":       12,
		"january":   1,
		"february":  2,
		"march":     3,
		"april":     4,
		"june":      6,
		"july":      7,
		"august":    8,
		"september": 9,
		"october":   10,
		"november":  11,
		"december":  12,
	}
)

const (
	// Regular expression1(datetime separator supports '-', '/', '.').
	// Eg:
	// "2017-12-14 04:51:34 +0805 LMT",
	// "2017-12-14 04:51:34 +0805 LMT",
	// "2006-01-02T15:04:05Z07:00",
	// "2014-01-17T01:19:15+08:00",
	// "2018-02-09T20:46:17.897Z",
	// "2018-02-09 20:46:17.897",
	// "2018-02-09T20:46:17Z",
	// "2018-02-09 20:46:17",
	// "2018/10/31 - 16:38:46"
	// "2018-02-09",
	// "2018.02.09",
	timeRegexPattern1 = `(\d{4}[-/\.]\d{2}[-/\.]\d{2})[:\sT-]*(\d{0,2}:{0,1}\d{0,2}:{0,1}\d{0,2}){0,1}\.{0,1}(\d{0,9})([\sZ]{0,1})([\+-]{0,1})([:\d]*)`

	// Regular expression2(datetime separator supports '-', '/', '.').
	// Eg:
	// 01-Nov-2018 11:50:28
	// 01/Nov/2018 11:50:28
	// 01.Nov.2018 11:50:28
	// 01.Nov.2018:11:50:28
	timeRegexPattern2 = `(\d{1,2}[-/\.][A-Za-z]{3,}[-/\.]\d{4})[:\sT-]*(\d{0,2}:{0,1}\d{0,2}:{0,1}\d{0,2}){0,1}\.{0,1}(\d{0,9})([\sZ]{0,1})([\+-]{0,1})([:\d]*)`

	// Regular expression3(time).
	// Eg:
	// 11:50:28
	// 11:50:28.897
	timeRegexPattern3 = `(\d{2}):(\d{2}):(\d{2})\.{0,1}(\d{0,9})`
)

// wrapper is a wrapper for stdlib struct time.Time.
// It's used for overwriting some functions of time.Time, for example: String.
type wrapper struct {
	time.Time
}

// String overwrites the String function of time.Time.
func (t wrapper) String() string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

// Time is a wrapper for time.Time for additional features.
type Time struct {
	wrapper
}

// apiUnixNano is an interface definition commonly for custom time.Time wrapper.
type apiUnixNano interface {
	UnixNano() int64
}

// New creates and returns a Time object with given parameter.
// The optional parameter can be type of: time.Time/*time.Time, string or integer.
func New(param ...interface{}) *Time {
	if len(param) > 0 {
		switch r := param[0].(type) {
		case time.Time:
			return NewFromTime(r)
		case *time.Time:
			return NewFromTime(*r)
		case Time:
			return &r
		case *Time:
			return r
		case int:
			return NewFromTimeStamp(int64(r))
		case int64:
			return NewFromTimeStamp(r)
		default:
			if v, ok := r.(apiUnixNano); ok {
				return NewFromTimeStamp(v.UnixNano())
			}
		}
	}
	return &Time{
		wrapper{time.Time{}},
	}
}

// Now creates and returns a time object of now.
func Now() *Time {
	return &Time{
		wrapper{time.Now()},
	}
}

// NewFromTime creates and returns a Time object with given time.Time object.
func NewFromTime(t time.Time) *Time {
	return &Time{
		wrapper{t},
	}
}

// NewFromTimeStamp creates and returns a Time object with given timestamp,
// which can be in seconds to nanoseconds.
// Eg: 1600443866 and 1600443866199266000 are both considered as valid timestamp number.
func NewFromTimeStamp(timestamp int64) *Time {
	if timestamp == 0 {
		return &Time{}
	}
	var sec, nano int64
	if timestamp > 1e9 {
		for timestamp < 1e18 {
			timestamp *= 10
		}
		sec = timestamp / 1e9
		nano = timestamp % 1e9
	} else {
		sec = timestamp
	}
	return &Time{
		wrapper{time.Unix(sec, nano)},
	}
}

// Value insert timestamp into mysql need this function.
func (t Time) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan value time.Time
func (t *Time) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = *NewFromTime(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

// Timestamp returns the timestamp in seconds.
func (t *Time) Timestamp() int64 {
	return t.UnixNano() / 1e9
}

// TimestampMilli returns the timestamp in milliseconds.
func (t *Time) TimestampMilli() int64 {
	return t.UnixNano() / 1e6
}

// TimestampMicro returns the timestamp in microseconds.
func (t *Time) TimestampMicro() int64 {
	return t.UnixNano() / 1e3
}

// TimestampNano returns the timestamp in nanoseconds.
func (t *Time) TimestampNano() int64 {
	return t.UnixNano()
}

// TimestampStr is a convenience method which retrieves and returns
// the timestamp in seconds as string.
func (t *Time) TimestampStr() string {
	return strconv.FormatInt(t.Timestamp(), 10)
}

// TimestampMilliStr is a convenience method which retrieves and returns
// the timestamp in milliseconds as string.
func (t *Time) TimestampMilliStr() string {
	return strconv.FormatInt(t.TimestampMilli(), 10)
}

// TimestampMicroStr is a convenience method which retrieves and returns
// the timestamp in microseconds as string.
func (t *Time) TimestampMicroStr() string {
	return strconv.FormatInt(t.TimestampMicro(), 10)
}

// TimestampNanoStr is a convenience method which retrieves and returns
// the timestamp in nanoseconds as string.
func (t *Time) TimestampNanoStr() string {
	return strconv.FormatInt(t.TimestampNano(), 10)
}

// Month returns the month of the year specified by t.
func (t *Time) Month() int {
	return int(t.Time.Month())
}

// Second returns the second offset within the minute specified by t,
// in the range [0, 59].
func (t *Time) Second() int {
	return t.Time.Second()
}

// Millisecond returns the millisecond offset within the second specified by t,
// in the range [0, 999].
func (t *Time) Millisecond() int {
	return t.Time.Nanosecond() / 1e6
}

// Microsecond returns the microsecond offset within the second specified by t,
// in the range [0, 999999].
func (t *Time) Microsecond() int {
	return t.Time.Nanosecond() / 1e3
}

// Nanosecond returns the nanosecond offset within the second specified by t,
// in the range [0, 999999999].
func (t *Time) Nanosecond() int {
	return t.Time.Nanosecond()
}

// String returns current time object as string.
func (t *Time) String() string {
	if t == nil {
		return ""
	}
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

// Clone returns a new Time object which is a clone of current time object.
func (t *Time) Clone() *Time {
	return New(t.Time)
}

// Add adds the duration to current time.
func (t *Time) Add(d time.Duration) *Time {
	newTime := t.Clone()
	newTime.Time = newTime.Time.Add(d)
	return newTime
}

// AddStr parses the given duration as string and adds it to current time.
func (t *Time) AddStr(duration string) (*Time, error) {
	if d, err := time.ParseDuration(duration); err != nil {
		return nil, err
	} else {
		return t.Add(d), nil
	}
}

// UTC converts current time to UTC timezone.
func (t *Time) UTC() *Time {
	newTime := t.Clone()
	newTime.Time = newTime.Time.UTC()
	return newTime
}

// AddDate adds year, month and day to the time.
func (t *Time) AddDate(years int, months int, days int) *Time {
	newTime := t.Clone()
	newTime.Time = newTime.Time.AddDate(years, months, days)
	return newTime
}

// Round returns the result of rounding t to the nearest multiple of d (since the zero time).
// The rounding behavior for halfway values is to round up.
// If d <= 0, Round returns t stripped of any monotonic clock reading but otherwise unchanged.
//
// Round operates on the time as an absolute duration since the
// zero time; it does not operate on the presentation form of the
// time. Thus, Round(Hour) may return a time with a non-zero
// minute, depending on the time's Location.
func (t *Time) Round(d time.Duration) *Time {
	newTime := t.Clone()
	newTime.Time = newTime.Time.Round(d)
	return newTime
}

// Truncate returns the result of rounding t down to a multiple of d (since the zero time).
// If d <= 0, Truncate returns t stripped of any monotonic clock reading but otherwise unchanged.
//
// Truncate operates on the time as an absolute duration since the
// zero time; it does not operate on the presentation form of the
// time. Thus, Truncate(Hour) may return a time with a non-zero
// minute, depending on the time's Location.
func (t *Time) Truncate(d time.Duration) *Time {
	newTime := t.Clone()
	newTime.Time = newTime.Time.Truncate(d)
	return newTime
}

// Equal reports whether t and u represent the same time instant.
// Two times can be equal even if they are in different locations.
// For example, 6:00 +0200 CEST and 4:00 UTC are Equal.
// See the documentation on the Time type for the pitfalls of using == with
// Time values; most code should use Equal instead.
func (t *Time) Equal(u *Time) bool {
	return t.Time.Equal(u.Time)
}

// Before reports whether the time instant t is before u.
func (t *Time) Before(u *Time) bool {
	return t.Time.Before(u.Time)
}

// After reports whether the time instant t is after u.
func (t *Time) After(u *Time) bool {
	return t.Time.After(u.Time)
}

// Sub returns the duration t-u. If the result exceeds the maximum (or minimum)
// value that can be stored in a Duration, the maximum (or minimum) duration
// will be returned.
// To compute t-d for a duration d, use t.Add(-d).
func (t *Time) Sub(u *Time) time.Duration {
	return t.Time.Sub(u.Time)
}

// StartOfMinute clones and returns a new time of which the seconds is set to 0.
func (t *Time) StartOfMinute() *Time {
	newTime := t.Clone()
	newTime.Time = newTime.Time.Truncate(time.Minute)
	return newTime
}

// StartOfHour clones and returns a new time of which the hour, minutes and seconds are set to 0.
func (t *Time) StartOfHour() *Time {
	y, m, d := t.Date()
	newTime := t.Clone()
	newTime.Time = time.Date(y, m, d, newTime.Time.Hour(), 0, 0, 0, newTime.Time.Location())
	return newTime
}

// StartOfDay clones and returns a new time which is the start of day, its time is set to 00:00:00.
func (t *Time) StartOfDay() *Time {
	y, m, d := t.Date()
	newTime := t.Clone()
	newTime.Time = time.Date(y, m, d, 0, 0, 0, 0, newTime.Time.Location())
	return newTime
}

// StartOfWeek clones and returns a new time which is the first day of week and its time is set to
// 00:00:00.
func (t *Time) StartOfWeek() *Time {
	weekday := int(t.Weekday())
	return t.StartOfDay().AddDate(0, 0, -weekday)
}

// StartOfMonth clones and returns a new time which is the first day of the month and its is set to
// 00:00:00
func (t *Time) StartOfMonth() *Time {
	y, m, _ := t.Date()
	newTime := t.Clone()
	newTime.Time = time.Date(y, m, 1, 0, 0, 0, 0, newTime.Time.Location())
	return newTime
}

// StartOfQuarter clones and returns a new time which is the first day of the quarter and its time is set
// to 00:00:00.
func (t *Time) StartOfQuarter() *Time {
	month := t.StartOfMonth()
	offset := (int(month.Month()) - 1) % 3
	return month.AddDate(0, -offset, 0)
}

// StartOfHalf clones and returns a new time which is the first day of the half year and its time is set
// to 00:00:00.
func (t *Time) StartOfHalf() *Time {
	month := t.StartOfMonth()
	offset := (int(month.Month()) - 1) % 6
	return month.AddDate(0, -offset, 0)
}

// StartOfYear clones and returns a new time which is the first day of the year and its time is set to
// 00:00:00.
func (t *Time) StartOfYear() *Time {
	y, _, _ := t.Date()
	newTime := t.Clone()
	newTime.Time = time.Date(y, time.January, 1, 0, 0, 0, 0, newTime.Time.Location())
	return newTime
}

// EndOfMinute clones and returns a new time of which the seconds is set to 59.
func (t *Time) EndOfMinute() *Time {
	return t.StartOfMinute().Add(time.Minute - time.Nanosecond)
}

// EndOfHour clones and returns a new time of which the minutes and seconds are both set to 59.
func (t *Time) EndOfHour() *Time {
	return t.StartOfHour().Add(time.Hour - time.Nanosecond)
}

// EndOfDay clones and returns a new time which is the end of day the and its time is set to 23:59:59.
func (t *Time) EndOfDay() *Time {
	y, m, d := t.Date()
	newTime := t.Clone()
	newTime.Time = time.Date(y, m, d, 23, 59, 59, int(time.Second-time.Nanosecond), newTime.Time.Location())
	return newTime
}

// EndOfWeek clones and returns a new time which is the end of week and its time is set to 23:59:59.
func (t *Time) EndOfWeek() *Time {
	return t.StartOfWeek().AddDate(0, 0, 7).Add(-time.Nanosecond)
}

// EndOfMonth clones and returns a new time which is the end of the month and its time is set to 23:59:59.
func (t *Time) EndOfMonth() *Time {
	return t.StartOfMonth().AddDate(0, 1, 0).Add(-time.Nanosecond)
}

// EndOfQuarter clones and returns a new time which is end of the quarter and its time is set to 23:59:59.
func (t *Time) EndOfQuarter() *Time {
	return t.StartOfQuarter().AddDate(0, 3, 0).Add(-time.Nanosecond)
}

// EndOfHalf clones and returns a new time which is the end of the half year and its time is set to 23:59:59.
func (t *Time) EndOfHalf() *Time {
	return t.StartOfHalf().AddDate(0, 6, 0).Add(-time.Nanosecond)
}

// EndOfYear clones and returns a new time which is the end of the year and its time is set to 23:59:59.
func (t *Time) EndOfYear() *Time {
	return t.StartOfYear().AddDate(1, 0, 0).Add(-time.Nanosecond)
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (t *Time) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.String() + `"`), nil
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (t *Time) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		t.Time = time.Time{}
		return nil
	}
	newTime, err := StrToTime(string(bytes.Trim(b, `"`)))
	if err != nil {
		return err
	}
	t.Time = newTime.Time
	return nil
}

// StrToTimeFormat parses string <str> to *Time object with given format <format>.
// The parameter <format> is like "Y-m-d H:i:s".
func StrToTimeFormat(str string, format string) (*Time, error) {
	return StrToTimeLayout(str, formatToStdLayout(format))
}

// formatToStdLayout converts custom format to stdlib layout.
func formatToStdLayout(format string) string {
	b := bytes.NewBuffer(nil)
	for i := 0; i < len(format); {
		switch format[i] {
		case '\\':
			if i < len(format)-1 {
				b.WriteByte(format[i+1])
				i += 2
				continue
			} else {
				return b.String()
			}

		default:
			if f, ok := formats[format[i]]; ok {
				// Handle particular chars.
				switch format[i] {
				case 'j':
					b.WriteString("2")
				case 'G':
					b.WriteString("15")
				case 'u':
					if i > 0 && format[i-1] == '.' {
						b.WriteString("000")
					} else {
						b.WriteString(".000")
					}

				default:
					b.WriteString(f)
				}
			} else {
				b.WriteByte(format[i])
			}
			i++
		}
	}
	return b.String()
}

// StrToTimeLayout parses string <str> to *Time object with given format <layout>.
// The parameter <layout> is in stdlib format like "2006-01-02 15:04:05".
func StrToTimeLayout(str string, layout string) (*Time, error) {
	if t, err := time.ParseInLocation(layout, str, time.Local); err == nil {
		return NewFromTime(t), nil
	} else {
		return nil, err
	}
}

// StrToTime converts string to *Time object. It also supports timestamp string.
// The parameter <format> is unnecessary, which specifies the format for converting like "Y-m-d H:i:s".
// If <format> is given, it acts as same as function StrToTimeFormat.
// If <format> is not given, it converts string as a "standard" datetime string.
// Note that, it fails and returns error if there's no date string in <str>.
func StrToTime(str string, format ...string) (*Time, error) {
	if len(format) > 0 {
		return StrToTimeFormat(str, format[0])
	}
	if isTimestampStr(str) {
		timestamp, _ := strconv.ParseInt(str, 10, 64)
		return NewFromTimeStamp(timestamp), nil
	}
	var (
		year, month, day     int
		hour, min, sec, nsec int
		match                []string
		local                = time.Local
	)
	if match = timeRegex1.FindStringSubmatch(str); len(match) > 0 && match[1] != "" {
		//for k, v := range match {
		//	match[k] = strings.TrimSpace(v)
		//}
		year, month, day = parseDateStr(match[1])
	} else if match = timeRegex2.FindStringSubmatch(str); len(match) > 0 && match[1] != "" {
		//for k, v := range match {
		//	match[k] = strings.TrimSpace(v)
		//}
		year, month, day = parseDateStr(match[1])
	} else if match = timeRegex3.FindStringSubmatch(str); len(match) > 0 && match[1] != "" {
		//for k, v := range match {
		//	match[k] = strings.TrimSpace(v)
		//}
		s := strings.Replace(match[2], ":", "", -1)
		if len(s) < 6 {
			s += strings.Repeat("0", 6-len(s))
		}
		hour, _ = strconv.Atoi(match[1])
		min, _ = strconv.Atoi(match[2])
		sec, _ = strconv.Atoi(match[3])
		nsec, _ = strconv.Atoi(match[4])
		for i := 0; i < 9-len(match[4]); i++ {
			nsec *= 10
		}
		return NewFromTime(time.Date(0, time.Month(1), 1, hour, min, sec, nsec, local)), nil
	} else {
		return nil, errors.New("unsupported time format")
	}

	// Time
	if len(match[2]) > 0 {
		s := strings.Replace(match[2], ":", "", -1)
		if len(s) < 6 {
			s += strings.Repeat("0", 6-len(s))
		}
		hour, _ = strconv.Atoi(s[0:2])
		min, _ = strconv.Atoi(s[2:4])
		sec, _ = strconv.Atoi(s[4:6])
	}
	// Nanoseconds, check and perform bit filling
	if len(match[3]) > 0 {
		nsec, _ = strconv.Atoi(match[3])
		for i := 0; i < 9-len(match[3]); i++ {
			nsec *= 10
		}
	}
	// If there's zone information in the string,
	// it then performs time zone conversion, which converts the time zone to UTC.
	if match[4] != "" && match[6] == "" {
		match[6] = "000000"
	}
	// If there's offset in the string, it then firstly processes the offset.
	if match[6] != "" {
		zone := strings.Replace(match[6], ":", "", -1)
		zone = strings.TrimLeft(zone, "+-")
		if len(zone) <= 6 {
			zone += strings.Repeat("0", 6-len(zone))
			h, _ := strconv.Atoi(zone[0:2])
			m, _ := strconv.Atoi(zone[2:4])
			s, _ := strconv.Atoi(zone[4:6])
			if h > 24 || m > 59 || s > 59 {
				return nil, errors.Errorf("invalid zone string: %s", match[6])
			}
			// Comparing the given time zone whether equals to current time zone,
			// it converts it to UTC if they does not equal.
			_, localOffset := time.Now().Zone()
			// Comparing in seconds.
			if (h*3600 + m*60 + s) != localOffset {
				local = time.UTC
				// UTC conversion.
				operation := match[5]
				if operation != "+" && operation != "-" {
					operation = "-"
				}
				switch operation {
				case "+":
					if h > 0 {
						hour -= h
					}
					if m > 0 {
						min -= m
					}
					if s > 0 {
						sec -= s
					}
				case "-":
					if h > 0 {
						hour += h
					}
					if m > 0 {
						min += m
					}
					if s > 0 {
						sec += s
					}
				}
			}
		}
	}
	if month <= 0 || day <= 0 {
		return nil, errors.New("invalid time string:" + str)
	}
	return NewFromTime(time.Date(year, time.Month(month), day, hour, min, sec, nsec, local)), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// Note that it overwrites the same implementer of `time.Time`.
func (t *Time) UnmarshalText(data []byte) error {
	vTime := New(data)
	if vTime != nil {
		*t = *vTime
		return nil
	}
	return errors.Errorf(`invalid time value: %s`, data)
}

// isTimestampStr checks and returns whether given string a timestamp string.
func isTimestampStr(s string) bool {
	length := len(s)
	if length == 0 {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

// parseDateStr parses the string to year, month and day numbers.
func parseDateStr(s string) (year, month, day int) {
	array := strings.Split(s, "-")
	if len(array) < 3 {
		array = strings.Split(s, "/")
	}
	if len(array) < 3 {
		array = strings.Split(s, ".")
	}
	// Parsing failed.
	if len(array) < 3 {
		return
	}
	// Checking the year in head or tail.
	if true {
		year, _ = strconv.Atoi(array[0])
		month, _ = strconv.Atoi(array[1])
		day, _ = strconv.Atoi(array[2])
	} else {
		if v, ok := monthMap[strings.ToLower(array[1])]; ok {
			month = v
		} else {
			return
		}
		year, _ = strconv.Atoi(array[2])
		day, _ = strconv.Atoi(array[0])
	}
	return
}
