package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDB() *gorm.DB {
	dsn := "host=go_coffee_db user=postgres dbname=go_coffee port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

    if err != nil {
        log.Fatalf("Error: %v", err)
    }
	db.AutoMigrate(&OrderRequest{})

	return db
}

func getCoffeeQuotas() Membership {
	return Membership{
		map[string]map[string]CoffeeQuota{
			"Basic": {
				"Cappuccino": CoffeeQuota{
					Amount:    1,
					Duration: 1440,
				},
				"Espresso": CoffeeQuota{
					Amount:    2,
					Duration: 1440,
				},
				"Americano": CoffeeQuota{
					Amount:    3,
					Duration: 1440,
				},
			},
			"Coffee Lover": {
				"Cappuccino": CoffeeQuota{
					Amount:    5,
					Duration: 1440,
				},
				"Espresso": CoffeeQuota{
					Amount:    5,
					Duration: 1440,
				},
				"Americano": CoffeeQuota{
					Amount:    5,
					Duration: 1440,
				},
			},
			"Americano Maniac": {
				"Cappuccino": CoffeeQuota{
					Amount:    1,
					Duration: 1440,
				},
				"Espresso": CoffeeQuota{
					Amount:    2,
					Duration: 1440,
				},
				"Americano": CoffeeQuota{
					Amount:    5,
					Duration: 60,
				},
			},
		},
	}
}

func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {

        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Credentials", "true")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}
