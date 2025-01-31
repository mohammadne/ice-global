package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mohammadne/ice-global/internal/controllers"
	"github.com/mohammadne/ice-global/internal/db"
)

func main() {
	db.MigrateDatabase()

	ginEngine := gin.Default()

	var taxController controllers.TaxController
	ginEngine.GET("/", taxController.ShowAddItemForm)
	ginEngine.POST("/add-item", taxController.AddItem)
	ginEngine.GET("/remove-cart-item", taxController.DeleteCartItem)
	srv := &http.Server{
		Addr:    ":8088",
		Handler: ginEngine,
	}

	srv.ListenAndServe()
}
