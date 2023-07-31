package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"strconv"
)

type DBConfig struct {
	Host    string
	User    string
	DBName  string
	Port    int
	SSLMode string
}

func readConfigFromEnv() (*DBConfig, error) {
	portStr := os.Getenv("DB_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid port number: %s", portStr)
	}

	return &DBConfig{
		Host:    os.Getenv("DB_HOST"),
		User:    os.Getenv("DB_USER"),
		DBName:  os.Getenv("DB_NAME"),
		Port:    port,
		SSLMode: os.Getenv("DB_SSLMODE"),
	}, nil
}

func initDB(config *DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s port=%d sslmode=%s",
		config.Host, config.User, config.DBName, config.Port, config.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("failed to initialize the database: %w", err)
	}

	db.AutoMigrate(&OrderRequest{})
	return db, nil
}

func getCoffeeQuotas() Membership {
	return Membership{
		map[string]map[string]CoffeeQuota{
			"Basic": {
				"Cappuccino": CoffeeQuota{
					Amount:   1,
					Duration: 1440,
				},
				"Espresso": CoffeeQuota{
					Amount:   2,
					Duration: 1440,
				},
				"Americano": CoffeeQuota{
					Amount:   3,
					Duration: 1440,
				},
			},
			"Coffee Lover": {
				"Cappuccino": CoffeeQuota{
					Amount:   5,
					Duration: 1440,
				},
				"Espresso": CoffeeQuota{
					Amount:   5,
					Duration: 1440,
				},
				"Americano": CoffeeQuota{
					Amount:   5,
					Duration: 1440,
				},
			},
			"Americano Maniac": {
				"Cappuccino": CoffeeQuota{
					Amount:   1,
					Duration: 1440,
				},
				"Espresso": CoffeeQuota{
					Amount:   2,
					Duration: 1440,
				},
				"Americano": CoffeeQuota{
					Amount:   5,
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
