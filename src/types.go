package main

import (
	"time"

	"gorm.io/gorm"
)

// Type for the coffee quota information
type CoffeeQuota struct {
	Amount   int
	Duration int
}

type Membership struct {
	Membership map[string]map[string]CoffeeQuota
}

// Type for defining the database connection and hardcoded quotas
// for handling coffee-related operations and requests.
type CoffeeShop struct {
    DB      *gorm.DB
    Quotas  Membership
}
	
// Type for request body on handling buy-coddee request
type OrderCoffeeRequestBody struct {
	CoffeeType string `json:"coffee"`
	UserID string `json:"userId"`
	MembershipType string `json:"membershipType"`
}

type CoffeeData struct {
	Memberships map[string]map[string]map[string]int `json:"memberships"`
}

type OrderRequest struct {
	ID         uint `gorm:"primary_key"`
	UserID     int
	CoffeeType string
	Timestamp  time.Time
}

