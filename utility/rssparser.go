package utility

import (
	"github.com/mmcdole/gofeed"
	ext "github.com/mmcdole/gofeed/extensions"
)

type RSSParser struct {
	outputTimeFormat string
}

type IFeed interface {
	GetExtension(string, ...string) (string, string, map[string]string)
	GetItemLength() int
	GetItem(int) IFeedItem
}

type IFeedItem interface {
	GetExtension(string, ...string) (string, string, map[string]string)
	GetTitle() string
	GetLink() string
	GetUpdateAt() string
}

type Feed struct {
	feed       gofeed.Feed
	timeFormat string
}

func NewRSSParser(outputTimeFormat string) RSSParser {
	if outputTimeFormat == "" {
		outputTimeFormat = DEFAULT_TIME_FORMAT
	}
	return RSSParser{
		outputTimeFormat: outputTimeFormat,
	}
}

func (rp RSSParser) Parse(text string) (IFeed, IError) {
	parser := gofeed.NewParser()
	feed, err := parser.ParseString(text)
	if err != nil {
		return nil, NewError("", ERR_RSS_PARSE)
	}
	return Feed{
		feed:       *feed,
		timeFormat: rp.outputTimeFormat,
	}, nil
}

func getExtensionValue(extensions ext.Extensions, element string, children ...string) (string, string, map[string]string) {
	extMap := extensions[element]

	if extMap == nil || len(extMap) == 0 {
		return "", "", nil
	}
	if len(children) == 0 {
		return "", "", nil
	}

	var ext ext.Extension
	for _, child := range children {
		extList := extMap[child]
		if extList == nil || len(extList) == 0 {
			return "", "", nil
		}
		ext = extList[0]
		extMap = ext.Children
	}

	return ext.Name, ext.Value, ext.Attrs
}

func (f Feed) GetExtension(element string, children ...string) (string, string, map[string]string) {
	return getExtensionValue(f.feed.Extensions, element, children...)
}

func (f Feed) GetItemLength() int {
	return len(f.feed.Items)
}

func (f Feed) GetItem(index int) IFeedItem {
	if index >= len(f.feed.Items) {
		return nil
	}
	return FeedItem{
		item:       *f.feed.Items[index],
		timeFormat: f.timeFormat,
	}
}

type FeedItem struct {
	item       gofeed.Item
	timeFormat string
}

func (fi FeedItem) GetExtension(element string, children ...string) (string, string, map[string]string) {
	return getExtensionValue(fi.item.Extensions, element, children...)
}

func (fi FeedItem) GetTitle() string {
	return fi.item.Title
}

func (fi FeedItem) GetLink() string {
	return fi.item.Link
}

func (fi FeedItem) GetUpdateAt() string {
	return fi.item.PublishedParsed.UTC().Format(fi.timeFormat)
}
