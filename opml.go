package main

import "encoding/xml"


type OPML struct {
    XMLName   xml.Name `xml:"opml"`
    Version string `xml:"version,attr"`
    Head Head`xml:"head"`
    Body Body `xml:"body"`
}

type Head struct {
    Title string `xml:"title"`
}

type Body struct {
    OutlineContainer OutlineContainer `xml:"outline"`
}

type Outline struct {
    Text string `xml:"text,attr"`
    Url string `xml:"xmlUrl,attr"`
    Type string `xml:"rss,attr"`
}

type OutlineContainer struct {
    Text string `xml:"text,attr"`
    Outlines []Outline `xml:"outline"`
}


func feedsToOPML(feeds []Feed) OPML {
    var outlines []Outline;

    for _, feed := range feeds {
        outlines = append(outlines, Outline{
            Url: feed.Url,
            Text: feed.Name,
            Type: "rss",
        })
    } 

    return OPML{
        Version: "1.0",
        Head: Head {
            Title: "Podcast OPML",
        },
        Body: Body {
            OutlineContainer: OutlineContainer{ 
                Text: "feeds",
                Outlines: outlines,
            },
        },
    }
}
