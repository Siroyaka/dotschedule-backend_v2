package domain

type ScheduleIDWithIsComplete struct {
	data map[string]bool
}

func NewScheduleIDWithIsComplete() ScheduleIDWithIsComplete {
	return ScheduleIDWithIsComplete{
		data: make(map[string]bool),
	}
}

func (sc ScheduleIDWithIsComplete) Add(value string, isComplete bool) ScheduleIDWithIsComplete {
	sc.data[value] = isComplete
	return sc
}

func (sc ScheduleIDWithIsComplete) Has(value string) bool {
	_, ok := sc.data[value]
	return ok
}

func (sc ScheduleIDWithIsComplete) IsComplete(value string) bool {
	v, ok := sc.data[value]
	if !ok {
		return false
	}
	return v
}

func (sc ScheduleIDWithIsComplete) List() []string {
	var list []string
	for k := range sc.data {
		list = append(list, k)
	}
	return list
}
