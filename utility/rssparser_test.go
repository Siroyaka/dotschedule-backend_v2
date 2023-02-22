package utility_test

import (
	"fmt"
	"testing"

	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

const (
	RSSData = `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns:yt="http://www.youtube.com/xml/schemas/2015" xmlns:media="http://search.yahoo.com/mrss/" xmlns="http://www.w3.org/2005/Atom">
 <link rel="self" href="http://example.com/feedurl"/>
 <id>yt:channel:yt_channelId</id>
 <yt:channelId>yt_channelId</yt:channelId>
 <title>feed title</title>
 <link rel="alternate" href="https://example.com/feedlink"/>
 <author>
  <name>feed author name</name>
  <uri>https://example.com/feedauthor</uri>
 </author>
 <published>2018-04-27T05:01:07+00:00</published>
 <entry>
  <id>yt:video:yt_videoId1</id>
  <yt:videoId>yt_videoId1</yt:videoId>
  <yt:channelId>yt_channelId1</yt:channelId>
  <title>entry title1</title>
  <link rel="alternate" href="http://example.com/link1"/>
  <author>
   <name>name name</name>
   <uri>http://example.com/uri1</uri>
  </author>
  <published>2022-12-03T15:24:04+00:00</published>
  <updated>2022-12-03T15:30:45+00:00</updated>
  <media:group>
   <media:title>media title1</media:title>
   <media:content url="http://example.com/content1" type="application/x-shockwave-flash" width="640" height="390"/>
   <media:thumbnail url="http://example.com/thumbnail1" width="480" height="360"/>
   <media:description>description
description  description1
description</media:description>
   <media:community>
    <media:starRating count="100" average="5.00" min="1" max="5"/>
    <media:statistics views="1000"/>
   </media:community>
  </media:group>
 </entry>
 <entry>
  <id>yt:video:yt_videoId2</id>
  <yt:videoId>yt_videoId2</yt:videoId>
  <yt:channelId>yt_channelId2</yt:channelId>
  <title>entry title2</title>
  <link rel="alternate" href="http://example.com/link2"/>
  <author>
   <name>name name</name>
   <uri>http://example.com/uri2</uri>
  </author>
  <published>2022-12-03T15:24:04+00:00</published>
  <updated>2022-12-03T15:30:45+00:00</updated>
  <media:group>
   <media:title>media title2</media:title>
   <media:content url="http://example.com/content2" type="application/x-shockwave-flash" width="640" height="390"/>
   <media:thumbnail url="http://example.com/thumbnail2" width="480" height="360"/>
   <media:description>description
description  description1
description</media:description>
   <media:community>
    <media:starRating count="200" average="5.00" min="1" max="5"/>
    <media:statistics views="2000"/>
   </media:community>
  </media:group>
 </entry>
</feed>`
)

func TestRSSParse(t *testing.T) {
	parser := utility.NewRSSParser("")
	_, err := parser.Parse(RSSData)
	if err != nil {
		t.Error("RSS Parse Error")
	}
}

func TestRSSItemLen(t *testing.T) {
	parser := utility.NewRSSParser("")
	feed, err := parser.Parse(RSSData)
	if err != nil {
		t.Error("RSS Parse Error")
	}
	t.Log(feed.GetItemLength())
	for i := 0; i < feed.GetItemLength(); i++ {
		item := feed.GetItem(i)
		fmt.Println(item.GetTitle())
		title, value, _ := item.GetExtension("yt", "videoId")
		t.Logf("title:%s value:%s", title, value)
		fmt.Println(item.GetUpdateAt())
	}
}
