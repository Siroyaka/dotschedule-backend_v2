package discordpost

import (
	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type DiscordPostRepository interface {
	Post(param domain.DiscordWebhookParams) utility.IError
}
