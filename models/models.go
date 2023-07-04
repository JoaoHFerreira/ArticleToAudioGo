package models

type Article struct {
	ID   uint   `gorm:"primaryKey;index"`
	Name string `gorm:"not null"`
	URL  string `gorm:"not null;index"`
}

type Paragraph struct {
	ID              uint   `gorm:"primaryKey;index"`
	ArticleID       uint   `gorm:"not null;index"`
	ParagraphNumber uint   `gorm:"not null;index"`
	Content         string `gorm:"not null"`
}

type AudioParagraph struct {
	ID           uint   `gorm:"primaryKey;index"`
	ArticleID    uint   `gorm:"index;not null"`
	ParagraphID  uint   `gorm:"not null"`
	AudioContent []byte `gorm:"not null"`
}
