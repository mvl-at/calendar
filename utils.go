package calendar

import (
	"bytes"
	"fmt"
	"github.com/mvl-at/model"
	"strconv"
	"strings"
)

type strBuffer struct {
	bytes.Buffer
}

func (s *strBuffer) WriteFmt(format string, a ...interface{}) {
	str := fmt.Sprintf(format+"\n", a...)
	if len(str) > maxLine {
		str = split(str)
	}
	s.WriteString(str)
}

func split(str string) string {
	result := ""
	for i, j := 0, len(str)/maxLine+1; i < j; i++ {
		lastIndex := (i + 1) * maxLine
		if i+1 >= j {
			lastIndex = len(str)
		}
		if i > 0 {
			result += "\n"
		}
		result += strings.Repeat(" ", i) + str[i*maxLine:lastIndex]
	}
	return result
}

func rangeString(events []*model.Event, note string) string {
	rangeString := "Terminplan für gute Märsche"
	if strings.ToLower(note) == strings.ToLower(conf.King) {
		rangeString = "Gute Märsche"
	}
	if len(events) > 0 {
		firstYear := events[0].Date.Year()
		lastYear := events[len(events)-1].Date.Year()
		firstMonth := events[0].Date.Month()
		lastMonth := events[len(events)-1].Date.Month()

		if firstYear == lastYear {
			if firstMonth == lastMonth {
				rangeString = months[firstMonth] + " " + strconv.Itoa(firstYear)
			} else {
				rangeString = months[firstMonth] + " bis " + months[lastMonth] + " " + strconv.Itoa(firstYear)
			}
		} else {
			rangeString = months[firstMonth] + " " + strconv.Itoa(firstYear) + " bis " + months[lastMonth] + " " + strconv.Itoa(lastYear)
		}
		rangeString = "Terminplan " + rangeString
	}
	return rangeString
}

func normalise(str string) string {
	ret := strings.ToLower(str)
	for k, v := range map[string]string{"ä": "ae", "ö": "oe", "ü": "ue", " ": "-"} {
		ret = strings.Replace(ret, k, v, -1)
	}
	return ret
}
