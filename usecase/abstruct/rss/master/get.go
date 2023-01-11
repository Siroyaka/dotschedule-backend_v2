package rssmaster

import (
	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type GetRepository interface {
	Get(MasterDataAdapter) ([]domain.RSSMaster, utility.IError)
}

type MasterDataAdapter func(id, name, url, date string) (domain.RSSMaster, utility.IError)
