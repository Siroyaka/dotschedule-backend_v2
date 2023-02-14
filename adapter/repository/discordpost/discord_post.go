package discordpost

import (
	"fmt"

	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

const (
	contentType = "application/json"
)

type DiscordPostRepository struct {
	url     string
	request abstruct.HTTPRequest
}

func NewDiscordPostRepository(request abstruct.HTTPRequest, url string) DiscordPostRepository {
	return DiscordPostRepository{
		request: request,
		url:     url,
	}
}

func (repos DiscordPostRepository) Post(param domain.DiscordWebhookParams) utility.IError {

	content, err := param.ToJson()
	if err != nil {
		return err.WrapError()
	}

	postParam := discordPostParam{
		url:         repos.url,
		content:     content,
		contentType: contentType,
	}

	utility.LogDebug(fmt.Sprintf("Discord Post content: %s", content))

	res, err := repos.request.Post(postParam)
	if err != nil {
		return err.WrapError()
	}

	utility.LogDebug(fmt.Sprintf("Discord Post Status: %s, Response: %s", res.Status(), res.Body()))

	return nil
}

type discordPostParam struct {
	url         string
	content     string
	contentType string
}

func (p discordPostParam) Url() string {
	return p.url
}

func (p discordPostParam) Content() string {
	return p.content
}

func (p discordPostParam) ContentType() string {
	return p.contentType
}
