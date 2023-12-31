package main

import (
	"time"

	"gorm.io/gorm"
)

type CoffeeQuota struct {
	Amount   int
	Duration int
}

type Membership struct {
	Membership map[string]map[string]CoffeeQuota
}

type CoffeeShop struct {
	DB     *gorm.DB
	Quotas *Membership
}

type OrderCoffeeRequestBody struct {
	CoffeeType     string `json:"coffee"`
	UserID         string `json:"userId"`
	MembershipType string `json:"membershipType"`
}

type OrderRequest struct {
	ID         uint `gorm:"primary_key"`
	UserID     int
	CoffeeType string
	Timestamp  time.Time
}
