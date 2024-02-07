package models

import "gorm.io/gorm"

type Books struct {
	ID        uint    `gorm:"primary key;autoIncrement" json:"id"`
	Author    *string `json:"author"`
	Title     *string `json:"title"`
	Publisher *string `json:"publisher"`
}

type URLS struct {
	ID       uint    `gorm:"primary key;autoIncrement" json:"id"`
	LongUrl  *string `json:"longUrl"`
	ShortUrl *string `json:"shortUrl"`
}

func MigrateBooks(db *gorm.DB) error {
	err := db.AutoMigrate(&Books{})
	return err
}
func MigrateURL(db *gorm.DB) error {
	err := db.AutoMigrate(&URLS{})
	return err
}
