package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	db := initDB()
	quotas := getCoffeeQuotas()

	coffeeShop := &CoffeeShop{
		DB:     db,
		Quotas: quotas,
	}

	router := gin.Default()
	router.Use(CORSMiddleware())
	router.GET("/api/shop/quota-config/", coffeeShop.getQuotaConfigController)
	router.GET("/api/shop/orders/", coffeeShop.getOrdersController)
	router.POST("/api/shop/order/", coffeeShop.orderCoffeeController)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
	log.Println("API is running!")
}
