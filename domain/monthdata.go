package domain

type MonthData struct {
	Date       string
	MemberData []MonthMemberData
	Icons      []string
}

func NewMonthData(date string, memberList []MonthMemberData, iconList []string) MonthData {
	return MonthData{
		Date:       date,
		MemberData: memberList,
		Icons:      iconList,
	}
}

type MonthMemberData struct {
	Id   string
	Name string
}

func NewMonthMemberData(id, name string) MonthMemberData {
	return MonthMemberData{
		Id:   id,
		Name: name,
	}
}
