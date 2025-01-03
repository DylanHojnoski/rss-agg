package main

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

type RSSFeed struct {
    Channel struct {
        Title string `xml:"title"`
        Link string `xml:"link"`
        Description string `xml:"description"`
        Language string `xml:"language"`
        Image Image `xml:"image"`
        Categories []CategoryTag `xml:"http://www.itunes.com/dtds/podcast-1.0.dtd category"`
        Item []RSSItem `xml:"item"`
    } `xml:"channel"`
}

type Image struct {
    Url string `xml:"url"`
    Title string `xml:"title"`
    Link string `xml:"link"`
}

type CategoryTag struct {
    Text string `xml:"text,attr"`
}

type RSSItem struct {
    Title string `xml:"title"`
    Link string `xml:"link"`
    Description string `xml:"description"`
    PubDate string `xml:"pubDate"`
    Audio Enclosure `xml:"enclosure"`
    Duration string `xml:"http://www.itunes.com/dtds/podcast-1.0.dtd duration"`
}

type Enclosure struct {
    Url    string `xml:"url,attr"`
    Length int64  `xml:"length,attr"`
    Type   string `xml:"type,attr"`
}

func urlToFeed(url string) (RSSFeed, error) {
    httpClient := http.Client {
        Timeout: 10 * time.Second,
    }

    resp, err := httpClient.Get(url)
    if err != nil {
        return RSSFeed{}, err
    }

    defer resp.Body.Close()

    dat, err := io.ReadAll(resp.Body)
    if err != nil {
        return RSSFeed{}, err
    }
    rssFeed := RSSFeed{}

    xml.Unmarshal(dat, &rssFeed)
    if err != nil {
        return RSSFeed{}, err
    }

    return rssFeed, nil

}
