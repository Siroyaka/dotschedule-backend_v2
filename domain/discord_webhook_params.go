package domain

import (
	"encoding/json"

	"github.com/Siroyaka/dotschedule-backend_v2/utility/utilerror"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/wrappedbasics"
)

type DiscordWebhookParams struct {
	Content string                `json:"content"`
	Embeds  []DiscordWebhookEmbed `json:"embeds"`
}

func (p DiscordWebhookParams) ToJson() (string, utilerror.IError) {
	responseJson, err := json.Marshal(p)
	if err != nil {
		return "", utilerror.New(err.Error(), utilerror.ERR_JSONPARSE)
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

func (e *DiscordWebhookEmbed) AddTimeStamp(t wrappedbasics.IWrappedTime) {
	e.TimeStamp = t.ToUTCFormatString(wrappedbasics.WrappedTimeProps.DateTimeFormat())
}

type DiscordWebhookEmbedAuthor struct {
	Name string `json:"name"`
}
