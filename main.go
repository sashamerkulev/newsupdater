package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/sashamerkulev/newsupdater/mysql"
	"github.com/sashamerkulev/newsupdater/reader"
	"time"
)

func readArticles() {
	ticker := time.NewTicker(time.Minute * 15)
	for range ticker.C {
		reader.Reader(mysql.AddArticles)
	}
}

func wipeArticles() {
	ticker := time.NewTicker(time.Hour * 24)
	for range ticker.C {
		wipeTime := time.Now()
		wipeTime = wipeTime.Add(-24 * 7 * time.Hour)
		mysql.WipeOldArticles(wipeTime)
	}
}

func wipeActivities() {
	ticker := time.NewTicker(time.Hour * 24)
	for range ticker.C {
		wipeTime := time.Now()
		wipeTime = wipeTime.Add(-24 * 30 * time.Hour)
		mysql.WipeOldActivities(wipeTime)
	}
}

func main() {
	go wipeArticles()
	go wipeActivities()

	readArticles()
}
