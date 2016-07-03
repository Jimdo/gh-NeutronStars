package main

import (
	"encoding/xml"
	"reflect"
	"testing"
	"time"
)

func TestUnmarshal(t *testing.T) {
	raw := `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom" xmlns:media="http://search.yahoo.com/mrss/" xml:lang="en-US">
  <id>tag:github.com,2008:https://github.com/hashicorp/vault/releases</id>
  <link type="text/html" rel="alternate" href="https://github.com/hashicorp/vault/releases"/>
  <link type="application/atom+xml" rel="self" href="https://github.com/hashicorp/vault/releases.atom"/>
  <title>Release notes from vault</title>
  <updated>2016-06-14T15:29:30+02:00</updated>
  <entry>
    <id>tag:github.com,2008:Repository/31288958/v0.6.0-rebuild</id>
    <updated>2016-06-14T15:29:30+02:00</updated>
    <link rel="alternate" type="text/html" href="/hashicorp/vault/releases/tag/v0.6.0-rebuild"/>
    <title>v0.6.0-rebuild</title>
    <content type="html">&lt;p&gt;Version 0.6.0-rebuild&lt;/p&gt;</content>
    <author>
      <name>jefferai</name>
    </author>
    <media:thumbnail height="30" width="30" url="https://avatars0.githubusercontent.com/u/28627?v=3&amp;s=60"/>
  </entry>
</feed>
`

	res := feed{}
	if err := xml.Unmarshal([]byte(raw), &res); err != nil {
		t.Errorf("Unmarshal failed: %s", err)
	}

	if len(res.Entries) != 1 {
		t.Errorf("Found %d entries, expected 1", len(res.Entries))
	}

	etime, _ := time.Parse("2006-01-02 15:04:05 -07:00", "2016-06-14 15:29:30 +02:00")
	expected := feedEntry{
		XMLName: xml.Name{Space: "http://www.w3.org/2005/Atom", Local: "entry"},
		ID:      "tag:github.com,2008:Repository/31288958/v0.6.0-rebuild",
		Updated: etime,
		Link:    feedLink{Relation: "alternate", Type: "text/html", Href: "/hashicorp/vault/releases/tag/v0.6.0-rebuild"},
		Title:   "v0.6.0-rebuild",
		Content: feedEntryContent{Type: "html", Text: "<p>Version 0.6.0-rebuild</p>"},
		Author:  "jefferai",
		Thumbnail: feedEntryThumbnail{
			XMLName: xml.Name{Space: "http://search.yahoo.com/mrss/", Local: "thumbnail"},
			Height:  30, Width: 30, URL: "https://avatars0.githubusercontent.com/u/28627?v=3&s=60",
		},
	}

	if !reflect.DeepEqual(expected, res.Entries[0]) {
		t.Errorf("Entry did not match expectation: RES=%#v EXP=%#v", res.Entries[0], expected)
	}
}

func TestMarshal(t *testing.T) {
	raw := `<feed xmlns="http://www.w3.org/2005/Atom" xml:lang="en-US"><id>gh-NeutronStars:Luzifer</id><title>Release summary for stared repos of GitHub user Luzifer</title><updated>2016-06-14T15:29:30+02:00</updated><entry><id>tag:github.com,2008:Repository/31288958/v0.6.0-rebuild</id><updated>2016-06-14T15:29:30+02:00</updated><link rel="alternate" type="text/html" href="/hashicorp/vault/releases/tag/v0.6.0-rebuild"></link><title>v0.6.0-rebuild</title><content type="html">&lt;p&gt;Version 0.6.0-rebuild&lt;/p&gt;</content><author><name>jefferai</name></author><thumbnail xmlns="http://search.yahoo.com/mrss/" height="30" width="30" url="https://avatars0.githubusercontent.com/u/28627?v=3&amp;s=60"></thumbnail></entry></feed>`

	etime, _ := time.Parse("2006-01-02 15:04:05 -07:00", "2016-06-14 15:29:30 +02:00")
	obj := feed{
		Lang:    "en-US",
		ID:      "gh-NeutronStars:Luzifer",
		Title:   "Release summary for stared repos of GitHub user Luzifer",
		Updated: etime,
		Entries: []feedEntry{{
			XMLName: xml.Name{Space: "http://www.w3.org/2005/Atom", Local: "entry"},
			ID:      "tag:github.com,2008:Repository/31288958/v0.6.0-rebuild",
			Updated: etime,
			Link:    feedLink{Relation: "alternate", Type: "text/html", Href: "/hashicorp/vault/releases/tag/v0.6.0-rebuild"},
			Title:   "v0.6.0-rebuild",
			Content: feedEntryContent{Type: "html", Text: "<p>Version 0.6.0-rebuild</p>"},
			Author:  "jefferai",
			Thumbnail: feedEntryThumbnail{
				XMLName: xml.Name{Space: "http://search.yahoo.com/mrss/", Local: "thumbnail"},
				Height:  30, Width: 30, URL: "https://avatars0.githubusercontent.com/u/28627?v=3&s=60",
			}}},
	}

	data, err := xml.Marshal(obj)

	if err != nil {
		t.Errorf("Marshal failed: %s", err)
	}

	if raw != string(data) {
		t.Errorf("Marshalling result did not match expectation:\nRES='%s'\nEXP='%s'", string(data), raw)
	}

}
