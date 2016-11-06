package rss

import "encoding/xml"

type rss struct {
	Channel channel `xml:"channel"`
}

type channel struct {
	Items []item `xml:"item"`
}

type item struct {
	Link      string    `xml:"link"`
	Enclosure enclosure `xml:"enclosure"`
}

type enclosure struct {
	Url string `xml:"url,attr"`
}

func NewRss(data []byte) (r *rss, err error) {
	var Rss rss
	if err := xml.Unmarshal(data, &Rss); err != nil {
		return nil, err
	}
	return &Rss, nil
}

func (r *rss) GetLinks() []string {
	var links []string
	for _, item := range r.Channel.Items {
		link := ""
		if item.Enclosure.Url != "" {
			link = item.Enclosure.Url
		} else if item.Link != "" {
			link = item.Link
		}
		if link != "" {
			links = append(links, link)
		}
	}
	return links
}
