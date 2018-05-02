package menu

import (
	"regexp"
	"fmt"
)

const (
	Sunday Weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

type Weekday byte

var weekdayRegexps = []string{
	`ned[eě]le`,
	`pond[eě]l[ií]`,
	`[uú]ter[yý]`,
	`st[rř]eda`,
	`[cč]tvrtek`,
	`p[aá]tek`,
	`sobota`,
}

func (w Weekday) buildRegexp(regexpTemplate string) (*regexp.Regexp) {
	sRegexp := fmt.Sprintf(regexpTemplate, weekdayRegexps[w])
	return regexp.MustCompile(sRegexp)
}