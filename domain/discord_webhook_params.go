package domain

import (
	"encoding/json"

	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type DiscordWebhookParams struct {
	Content string                `json:"content"`
	Embeds  []DiscordWebhookEmbed `json:"embeds"`
}

func (p DiscordWebhookParams) ToJson() (string, utility.IError) {
	responseJson, err := json.Marshal(p)
	if err != nil {
		return "", utility.NewError(err.Error(), utility.ERR_JSONPARSE)
	}
	return string(responseJson), nil
}

type DiscordWebhookEmbed struct {
	Author      DiscordWebhookEmbedAuthor `json:"author"`
	Title       string                    `json:"title"`
	TimeStamp   string                    `json:"timestamp"`
	Url         string                    `json:"url"`
	Description string                    `json:"description"`
}

func (e *DiscordWebhookEmbed) AddTimeStamp(t utility.WrappedTime) {
	e.TimeStamp = t.ToUTCFormatString()
}

type DiscordWebhookEmbedAuthor struct {
	Name string `json:"name"`
}
