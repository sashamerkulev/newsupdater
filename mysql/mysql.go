package mysql

import (
	"database/sql"
	"fmt"
	"github.com/sashamerkulev/newsupdater/config"
	"github.com/sashamerkulev/newsupdater/models"
	"time"
)

var cfg config.Config

func init() {
	cfg = config.GetConfig()
	fmt.Println(cfg.Connection.Mysql)
}

func AddArticles(articles []models.Article) {
	db, err := sql.Open("mysql", cfg.Connection.Mysql)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return
	}
	insertStmt, err := db.Prepare("INSERT INTO articles(SourceName, Title, Link, Description, PubDate, Category, PictureUrl) VALUES(?,?,?,?,?,?,?)")
	if err != nil {
		tx.Rollback()
		return
	}
	defer tx.Commit()
	defer insertStmt.Close()
	for i := 0; i < len(articles); i++ {
		_, err = insertStmt.Exec(articles[i].SourceName, articles[i].Title, articles[i].Link, articles[i].Description, articles[i].PubDate, articles[i].Category, articles[i].PictureUrl)
	}
}

func WipeOldArticles(wipeTime time.Time) {
	db, err := sql.Open("mysql", cfg.Connection.Mysql)
	if err != nil {
		return
	}
	defer db.Close()
	result, err := db.Exec("DELETE FROM articles WHERE "+
		"ArticleId not in (SELECT * FROM (SELECT a1.ArticleId FROM articles a1 JOIN articleLikes ual on ual.ArticleId = a1.ArticleId "+
		" UNION "+
		" SELECT a1.ArticleId FROM articles a1 JOIN articleComments uac on uac.ArticleId = a1.ArticleId) as art) "+
		"AND PubDate <= ?", wipeTime)
	if err != nil {
		return
	}
	_, err = result.RowsAffected()
	if err != nil {
		return
	}
}

func WipeOldActivities(wipeTime time.Time) {
	db, err := sql.Open("mysql", cfg.Connection.Mysql)
	if err != nil {
		return
	}
	defer db.Close()
	result, err := db.Exec(`DELETE FROM articles WHERE ArticleId in ( 
select b.articleId
from (
select a.articleId, max(coalesce(uac.timestamp, a.pubdate)) as timestamp
from articles a
left join articleComments uac on uac.articleId = a.articleId
group by a.ArticleId
union
select a.articleId, max(coalesce(uac.timestamp, a.pubdate)) as timestamp
from articles a
left join articleLikes uac on uac.articleId = a.articleId
group by a.ArticleId
) b WHERE b.timestamp < ?)
`, wipeTime)
	if err != nil {
		return
	}
	_, err = result.RowsAffected()
	if err != nil {
		return
	}
}
