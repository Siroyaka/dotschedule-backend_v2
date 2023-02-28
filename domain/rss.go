package domain

import (
	"github.com/Siroyaka/dotschedule-backend_v2/utility/wrappedbasics"
)

type RSSMaster struct {
	ID         string
	Name       string
	Url        string
	LastUpdate wrappedbasics.IWrappedTime
}

func NewRSSMaster(id string, name string, url string, lastupdate wrappedbasics.IWrappedTime) RSSMaster {
	return RSSMaster{
		ID:         id,
		Name:       name,
		Url:        url,
		LastUpdate: lastupdate,
	}
}
