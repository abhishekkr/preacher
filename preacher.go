package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/abhishekkr/gol/golhttpclient"
	"github.com/abhishekkr/preacher/hackernews"
)

var (
	htmlType = regexp.MustCompile(`^text/html`)
	xmlType  = regexp.MustCompile(`^text/xml`)
	pdfType  = regexp.MustCompile(`^application/pdf`)
	jsonType = regexp.MustCompile(`^application/json`)
)

func htmlBodySummary(responseBody io.Reader) (summary string) {
	doc, err := goquery.NewDocumentFromReader(responseBody)
	if err != nil {
		panic(err)
	}

	doc.Find(".body").Each(func(i int, s *goquery.Selection) {
		summary = s.Text()[0:100]
	})
	return
}

func summaryByType(contentType string, responseBody io.Reader) (summary string) {
	var err error
	var summaryByte []byte

	switch {
	case htmlType.MatchString(contentType):
		summary = htmlBodySummary(responseBody)

	case xmlType.MatchString(contentType):
		summaryByte, err = ioutil.ReadAll(responseBody)
		summary = string(summaryByte)

	case pdfType.MatchString(contentType):
		summary = "it's a PDF file"

	case jsonType.MatchString(contentType):
		summaryByte, err = ioutil.ReadAll(responseBody)
		summary = string(summaryByte)

	default:
		fmt.Println(contentType, "is not recognized")
	}

	if err != nil {
		log.Fatalf(err.Error())
	}
	return
}

func urlSummary(uri string) {
	response, err := golhttpclient.Http("GET", uri, map[string]string{}, map[string]string{})
	if err != nil {
		log.Printf("[error] %s", err.Error())
	}

	fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
	summaryByType(response.Header["Content-Type"][0], response.Body)
}

func hn() {
	for _, storyID := range hackernews.NewStoryIDs() {
		var hackernewsItem hackernews.DataItem
		hackernewsItem.Get(storyID)
		fmt.Printf("[#%d :: %s](%s)\n", hackernewsItem.ID, hackernewsItem.Title, hackernewsItem.URL)
		urlSummary(hackernewsItem.URL)
		break
	}
}

func main() {
	hn()
}
