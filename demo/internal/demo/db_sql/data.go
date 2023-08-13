package db_sql

import "gorm.io/gorm"

type HelloData struct {
	gorm.Model
	Id      string
	Message string
}

func FirstOrCreate(db *gorm.DB, id, message string) error {
	err := db.AutoMigrate(&HelloData{})
	if err != nil {
		return err
	}
	return db.FirstOrCreate(&HelloData{Id: id, Message: message}).Error
}
