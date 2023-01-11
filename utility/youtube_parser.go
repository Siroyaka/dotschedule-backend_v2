package utility

import (
	"strconv"
	"strings"
)

type YoutubeDurationParser struct {
	hourSplitter, miniteSplitter, secondSplitter string
	hour, minite, second                         int
	err                                          IError
}

func NewYoutubeDurationParser(hourSplitter, miniteSplitter, secondSplitter string) YoutubeDurationParser {
	return YoutubeDurationParser{
		hourSplitter:   hourSplitter,
		miniteSplitter: miniteSplitter,
		secondSplitter: secondSplitter,
		hour:           0,
		minite:         0,
		second:         0,
		err:            nil,
	}
}

func (parser YoutubeDurationParser) GetHour() int {
	return parser.hour
}

func (parser YoutubeDurationParser) GetMinite() int {
	return parser.minite
}

func (parser YoutubeDurationParser) GetSecond() int {
	return parser.second
}

func (parser YoutubeDurationParser) GetTotalSeconds() int {
	return (parser.hour*60+parser.minite)*60 + parser.second
}

func (parser YoutubeDurationParser) Err() IError {
	return parser.err
}

func (parser YoutubeDurationParser) Set(value string) YoutubeDurationParser {
	hold := 0
	for _, s := range strings.Split(value, "") {
		switch s {
		case "1", "2", "3", "4", "5", "6", "7", "8", "9", "0":
			hold *= 10
			a, _ := strconv.Atoi(s)
			hold += a
		case parser.hourSplitter:
			parser.hour = hold
			hold = 0
		case parser.miniteSplitter:
			parser.minite = hold
			hold = 0
		case parser.secondSplitter:
			parser.second = hold
			hold = 0
		default:
		}
	}
	return parser
}
