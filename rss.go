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
        ItunesSummary string `xml:"http://www.itunes.com/dtds/podcast-1.0.dtd summary"`
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

type ItunesImage struct {
    Href string `xml:"href,attr"`
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
    ItunesImage ItunesImage `xml:"http://www.itunes.com/dtds/podcast-1.0.dtd image"`
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

    if (rssFeed.Channel.Description == "" && rssFeed.Channel.ItunesSummary != "") {
        rssFeed.Channel.Description = rssFeed.Channel.ItunesSummary
    }

    if rssFeed.Channel.Image.Url == "" && len(rssFeed.Channel.Item) > 0 && rssFeed.Channel.Item[0].ItunesImage.Href != "" {
        rssFeed.Channel.Image.Url = rssFeed.Channel.Item[0].ItunesImage.Href
    }

    return rssFeed, nil
}
