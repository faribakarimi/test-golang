package database

import (
	"log"

	"github.com/faribakarimi/test-golang/api/models"
	"github.com/jinzhu/gorm"
)

var Connector *gorm.DB

func Connect(connectionString string) error {
	var err error
	Connector, err = gorm.Open("mysql", connectionString)
	if err != nil {
		return err
	}
	log.Println("Connection was successful!")
	return nil
}

func MigrateUser(table *models.User) {
	Connector.AutoMigrate(&table)
	log.Println("Table Users migrated.")
}