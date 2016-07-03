package main

import (
	"encoding/xml"
	"time"
)

// <feed xmlns="http://www.w3.org/2005/Atom" xmlns:media="http://search.yahoo.com/mrss/" xml:lang="en-US">
//   <id>tag:github.com,2008:https://github.com/hashicorp/vault/releases</id>
//   <link type="text/html" rel="alternate" href="https://github.com/hashicorp/vault/releases"/>
//   <link type="application/atom+xml" rel="self" href="https://github.com/hashicorp/vault/releases.atom"/>
//   <title>Release notes from vault</title>
//   <updated>2016-06-14T15:29:30+02:00</updated>
type feed struct {
	XMLName xml.Name    `xml:"http://www.w3.org/2005/Atom feed"`
	Lang    string      `xml:"xml:lang,attr"`
	ID      string      `xml:"id"`
	Title   string      `xml:"title"`
	Updated time.Time   `xml:"updated"`
	Entries []feedEntry `xml:"entry"`
}

type feedLink struct {
	Relation string `xml:"rel,attr" json:"rel"`
	Type     string `xml:"type,attr" json:"type"`
	Href     string `xml:"href,attr" json:"href"`
}

type feedEntryThumbnail struct {
	XMLName xml.Name `xml:"http://search.yahoo.com/mrss/ thumbnail"`
	Height  int      `xml:"height,attr" json:"height"`
	Width   int      `xml:"width,attr" json:"width"`
	URL     string   `xml:"url,attr" json:"url"`
}

type feedEntryContent struct {
	Type string `xml:"type,attr,omitempty" json:"type"`
	Text string `xml:",chardata" json:"text"`
}

//  <entry>
//    <id>tag:github.com,2008:Repository/31288958/v0.6.0-rebuild</id>
//    <updated>2016-06-14T15:29:30+02:00</updated>
//    <link rel="alternate" type="text/html" href="/hashicorp/vault/releases/tag/v0.6.0-rebuild"/>
//    <title>v0.6.0-rebuild</title>
//    <content type="html">&lt;p&gt;Version 0.6.0-rebuild&lt;/p&gt;</content>
//    <author>
//      <name>jefferai</name>
//    </author>
//    <media:thumbnail height="30" width="30" url="https://avatars0.githubusercontent.com/u/28627?v=3&amp;s=60"/>
//  </entry>
type feedEntry struct {
	XMLName   xml.Name           `xml:"entry"`
	ID        string             `xml:"id" json:"id"`
	Updated   time.Time          `xml:"updated" json:"updated"`
	Link      feedLink           `xml:"link" json:"link"`
	Title     string             `xml:"title" json:"title"`
	Content   feedEntryContent   `xml:"content" json:"content"`
	Author    string             `xml:"author>name" json:"author"`
	Thumbnail feedEntryThumbnail `xml:"thumbnail" json:"thumbnail"`
}

type feedEntries []feedEntry

func (f feedEntries) Len() int           { return len(f) }
func (f feedEntries) Less(i, j int) bool { return f[i].Updated.Before(f[j].Updated) }
func (f feedEntries) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
