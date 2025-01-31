package db

import "github.com/mohammadne/ice-global/internal/entity"

func MigrateDatabase() {
	db := GetDatabase()

	// AutoMigrate will create or update the tables based on the models
	err := db.AutoMigrate(&entity.CartEntity{}, &entity.CartItem{})
	if err != nil {
		panic(err)
	}
}
