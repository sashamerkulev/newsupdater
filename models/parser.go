package models

import "encoding/xml"

type Rss2 struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	// Required
	Title       string `xml:"channel>title"`
	Link        string `xml:"channel>link"`
	Description string `xml:"channel>description"`
	// Optional
	PubDate  string `xml:"channel>pubDate"`
	ItemList []Item `xml:"channel>item"`
}

type Enclosure struct {
	Url    string `xml:"url,attr"`
	Length int64  `xml:"length,attr"`
	Type   string `xml:"type,attr"`
}

type Item struct {
	// Required
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	// Optional
	Content   string    `xml:"encoded"`
	PubDate   string    `xml:"pubDate"`
	Comments  string    `xml:"comments"`
	Enclosure Enclosure `xml:"enclosure"`
	Category  string    `xml:"category"`
	Guid      string    `xml:"guid"`
}
