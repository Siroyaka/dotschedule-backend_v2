package apireference

import "github.com/Siroyaka/dotschedule-backend_v2/utility/wrappedbasics"

type StreamingSearchValues struct {
	member    []string
	from      wrappedbasics.IWrappedTime
	to        wrappedbasics.IWrappedTime
	tags      []string
	page      int
	title     string
	maxResult int
	sortOrder string
}

func NewStreamingSearchValues(
	member []string,
	from wrappedbasics.IWrappedTime,
	to wrappedbasics.IWrappedTime,
	tags []string,
	page int,
	title string,
	maxResult int,
	sortOrder string,
) StreamingSearchValues {
	return StreamingSearchValues{
		member:    member,
		from:      from,
		to:        to,
		tags:      tags,
		page:      page,
		title:     title,
		maxResult: maxResult,
		sortOrder: sortOrder,
	}
}

func EmptyStreamingSearchValues() StreamingSearchValues {
	return StreamingSearchValues{}
}

func (ssv StreamingSearchValues) Extract() ([]string, wrappedbasics.IWrappedTime, wrappedbasics.IWrappedTime, []string, int, string, int, string) {
	return ssv.member, ssv.from, ssv.to, ssv.tags, ssv.page, ssv.title, ssv.maxResult, ssv.sortOrder
}
