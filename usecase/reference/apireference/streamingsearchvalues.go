package apireference

import "github.com/Siroyaka/dotschedule-backend_v2/utility/wrappedbasics"

type StreamingSearchValues struct {
	member []string
	from   wrappedbasics.IWrappedTime
	to     wrappedbasics.IWrappedTime
	tags   []string
	page   int
	title  string
}

func NewStreamingSearchValues(
	member []string,
	from wrappedbasics.IWrappedTime,
	to wrappedbasics.IWrappedTime,
	tags []string,
	page int,
	title string,
) StreamingSearchValues {
	return StreamingSearchValues{
		member: member,
		from:   from,
		to:     to,
		tags:   tags,
		page:   page,
		title:  title,
	}
}

func EmptyStreamingSearchValues() StreamingSearchValues {
	return StreamingSearchValues{}
}

func (ssv StreamingSearchValues) Extract() ([]string, wrappedbasics.IWrappedTime, wrappedbasics.IWrappedTime, []string, int, string) {
	return ssv.member, ssv.from, ssv.to, ssv.tags, ssv.page, ssv.title
}
