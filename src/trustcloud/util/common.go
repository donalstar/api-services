package util

import (
	"regexp"
	"strconv"
	"time"
)

const layout = "2006-01-02T15:04:05.999Z"
const layout2 = "02-Jan-2006"

//(415) 566-7777
var phoneFormat1 = Regexp{regexp.MustCompile(`\((?P<0>\d+)\)\ (?P<1>\d+)-(?P<2>\d+)`)}

// 201-203-5088
var phoneFormat2 = Regexp{regexp.MustCompile(`(?P<0>\d+)\-(?P<1>\d+)-(?P<2>\d+)`)}

func ParseDate(input string) (time.Time, error) {
	return time.Parse(layout, input)
}

func FormatDate(input time.Time) string {
	return input.Format(layout)
}

func FormatDateForDisplay(input time.Time) string {
	return input.Format(layout2)
}

func FormatPhone(phone int64) string {
	p := strconv.FormatInt(phone, 10)

	result := ""

	if len(p) >= 10 {
		result = "(" + p[0:3] + ") " + p[3:6] + "-" + p[6:]
	}

	return result
}

func ParsePhoneString(phone string) *string {
	result := ParsePhoneStringWithFormat(phone, phoneFormat1)

	if result == nil {
		result = ParsePhoneStringWithFormat(phone, phoneFormat2)
	}

	return result
}

func ParsePhoneStringWithFormat(phone string, format Regexp) *string {
	m := format.FindStringSubmatchMap(phone)

	if len(m) > 0 {
		ph := m["0"] + m["1"] + m["2"]

		return &ph
	}

	return nil
}

type Regexp struct {
	*regexp.Regexp
}

func (r *Regexp) FindStringSubmatchMap(s string) map[string]string {
	captures := make(map[string]string)

	match := r.FindStringSubmatch(s)
	if match == nil {
		return captures
	}

	for i, name := range r.SubexpNames() {
		// Ignore the whole regexp match and unnamed groups
		if i == 0 || name == "" {
			continue
		}

		captures[name] = match[i]

	}
	return captures
}
