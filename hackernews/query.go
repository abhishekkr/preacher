package hackernews

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/abhishekkr/gol/golhttpclient"
)

type DataItem struct {
	By          string `json:"by"`
	Descendants int    `json:"descendants"` //only-story
	ID          int    `json:"id"`
	Parent      int    `json:"parent"` //only-comment
	Text        string `json:"text"`   //only-comment
	Score       int    `json:"score"`  //only-story
	Time        int    `json:"time"`
	Title       string `json:"title"` //only-story
	Type        string `json:"type"`
	URL         string `json:"url"` //only-story
}

func (item *DataItem) Get(itemID int) {
	uri := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", itemID)
	responseBody, err := golhttpclient.HttpGet(uri, map[string]string{}, map[string]string{})
	if err != nil {
		log.Printf("[error] %s", err.Error())
	}

	if err = json.Unmarshal([]byte(responseBody), &item); err != nil {
		log.Printf("[error] %s", err.Error())
	}
	return
}

func (item DataItem) IsComment() bool {
	if item.Type == "comment" {
		return true
	}
	return false
}

func (item DataItem) IsStory() bool {
	if item.Type == "story" {
		return true
	}
	return false
}
