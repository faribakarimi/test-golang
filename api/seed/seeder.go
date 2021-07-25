package seed

import (
	"log"

	"github.com/faribakarimi/test-golang/api/models"
	"github.com/jinzhu/gorm"
)

var items = []models.Item{
	models.Item{
		Name:  "Product Number 1",
		Price: 1000,
	},
	models.Item{
		Name:  "Product Number 2",
		Price: 2000,
	},
	models.Item{
		Name:  "Product Number 3",
		Price: 3000,
	},
	models.Item{
		Name:  "Product Number 4",
		Price: 4000,
	},
	models.Item{
		Name:  "Product Number 5",
		Price: 5000,
	},
}

func Load(db *gorm.DB) {
	var err error
	for i, _ := range items {
		err = db.Debug().Model(&models.Item{}).Create(&items[i]).Error
		if err != nil {
			log.Fatalf("Cannot seed items table: %v", err)
		}
	}

}
