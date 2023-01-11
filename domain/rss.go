package domain

import "github.com/Siroyaka/dotschedule-backend_v2/utility"

type RSSMaster struct {
	ID         string
	Name       string
	Url        string
	LastUpdate utility.WrappedTime
}

func NewRSSMaster(id string, name string, url string, lastupdate utility.WrappedTime) RSSMaster {
	return RSSMaster{
		ID:         id,
		Name:       name,
		Url:        url,
		LastUpdate: lastupdate,
	}
}
