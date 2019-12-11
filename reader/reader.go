package reader

import (
	"encoding/xml"
	"github.com/sashamerkulev/newsupdater/models"
	"io/ioutil"
	"net/http"
	"time"
)

func readRssSource(url string, bytes chan []byte) {
	resp, err := http.Get(url)
	if err != nil {
		bytes <- make([]byte, 0)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		bytes <- make([]byte, 0)
		return
	}
	bytes <- body
}

func parseRssItems(bytes []byte, items chan []models.Item) {
	if len(bytes) <= 0 {
		items <- make([]models.Item, 0)
		return
	}
	r := models.Rss2{}
	err := xml.Unmarshal(bytes, &r)
	if err != nil {
		items <- make([]models.Item, 0)
		return
	}
	items <- r.ItemList
}

func saveRssItems(saver models.SaverArticlesFunc, id int, items []models.Item) {
	if len(items) <= 0 {
		return
	}
	var articles = make([]models.Article, 0)
	for i := 0; i < len(items); i++ {
		date, err := time.Parse(models.Urls[id].Layout, items[i].PubDate)
		if err == nil {
			article := models.Article{
				SourceName:  models.Urls[id].Name,
				Category:    items[i].Category,
				Description: items[i].Description,
				Link:        items[i].Link,
				PubDate:     date,
				Title:       items[i].Title,
				PictureUrl:  items[i].Enclosure.Url,
			}
			articles = append(articles, article)
		}
	}
	saver(articles)
}

func Reader(saver models.SaverArticlesFunc) {
	for i := 0; i < len(models.Urls); i++ {
		bytes := make(chan []byte)
		items := make(chan []models.Item)
		go readRssSource(models.Urls[i].Link, bytes)
		go parseRssItems(<-bytes, items)
		go saveRssItems(saver, i, <-items)
	}
}
