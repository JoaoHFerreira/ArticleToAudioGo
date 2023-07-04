package database

import (
	"articletoaudio/models"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDBInstance() *gorm.DB {
	dsn := "host=localhost user=postgres password=postgres dbname=golang_text_to_audio port=5432"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}
	return db
}

func Migrate() {
	log.Println("Migration started")
	db := GetDBInstance()

	// explicitly create the table name in the array to automatclly create the table
	models := []map[string]interface{}{
		{"name": "Article", "model": models.Article{}},
		{"name": "Paragraph", "model": models.Paragraph{}},
	}

	for i, m := range models {
		log.Printf("Migration started for %s\n", m["name"].(string))
		log.Printf("Migrating %d of %d\n", i+1, len(models))
		db.AutoMigrate(m["model"])
	}

	log.Println("Migration ended")

}

func InsertArticle(title, url string) (uint, error) {
	db := GetDBInstance()
	article := models.Article{
		Name: title,
		URL:  url,
	}
	result := db.Create(&article)
	if result.Error != nil {
		return 0, result.Error
	}
	return article.ID, nil
}

func InsertParagraphs(paragraphs []map[string]interface{}, articleID uint) {
	db := GetDBInstance()

	for _, p := range paragraphs {
		paragraph := models.Paragraph{
			ArticleID:       articleID,
			ParagraphNumber: uint(p["index"].(int)),
			Content:         p["content"].(string),
		}
		result := db.Create(&paragraph)
		if result.Error != nil {
			fmt.Println(result.Error)
		}
	}
}
