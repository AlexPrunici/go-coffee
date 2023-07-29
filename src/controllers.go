package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func (coffeeShop *CoffeeShop) getQuotaConfigController(context *gin.Context) {
	context.JSON(http.StatusOK, coffeeShop.Quotas.Membership)
}

func (coffeeShop *CoffeeShop) getOrdersController(context *gin.Context) {

	var orders []OrderRequest

	coffeeShop.DB.Find(&orders)

	context.JSON(http.StatusOK, orders)
}

func (coffeeShop *CoffeeShop) orderCoffeeController(context *gin.Context) {
	var userIdBody string
	var membershipType string
	var coffeeType string
	var requestBody OrderCoffeeRequestBody

	if context.BindJSON(&requestBody) == nil {
		userIdBody = requestBody.UserID
		membershipType = requestBody.MembershipType
		coffeeType = requestBody.CoffeeType
	}
	currentTime := time.Now()

	userID, isIDValid := isValidUserID(userIdBody)
	if !isIDValid {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UserID"})
		return
	}

	if !isValidMembershipType(membershipType) {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid membership type"})
		return
	}

	if !isValidCoffeeType(coffeeType) {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid coffee type"})
		return
	}

	coffeeQuota, found := coffeeShop.getMembershipCoffeeQuota(membershipType, coffeeType)

	if !found {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Quota not found"})
	}

	ordersCount := coffeeShop.countUserOrders(userID, coffeeType, coffeeQuota.Duration)

	if int(ordersCount) >= coffeeQuota.Amount {
		var waitTime time.Duration
		if lastUsage, err := coffeeShop.getLastOrderTimestamp(userID, coffeeType); err == nil {
			waitTime = lastUsage.Add(time.Duration(coffeeQuota.Duration) * time.Minute).Sub(currentTime)
		} else {
			waitTime = 0
		}

		context.JSON(http.StatusTooManyRequests, gin.H{"error": fmt.Sprintf("Limit exceded. Wait %.0f minutes", waitTime.Minutes())})
		return
	}

	orderRequest := OrderRequest{
		UserID:     userID,
		CoffeeType: coffeeType,
		Timestamp:  currentTime,
	}

	coffeeShop.DB.Create(&orderRequest)
	context.JSON(http.StatusOK, gin.H{"success": "Order created with success"})
}

func (coffeeShop *CoffeeShop) countUserOrders(userID int, coffeeType string, quotaDuration int) int64 {
	duration := time.Now().Add(-time.Duration(quotaDuration) * time.Minute)

	var count int64
	coffeeShop.DB.Model(&OrderRequest{}).Where("user_id = ? AND coffee_type = ? AND timestamp > ?", userID, coffeeType, duration).Count(&count)

	return count
}

func (coffeeShop *CoffeeShop) getLastOrderTimestamp(userID int, coffeeType string) (time.Time, error) {
	var lastOrder OrderRequest
	if err := coffeeShop.DB.Model(&OrderRequest{}).Where("user_id = ? AND coffee_type = ?", userID, coffeeType).Order("timestamp desc").First(&lastOrder).Error; err != nil {
		return lastOrder.Timestamp, err
	}
	return lastOrder.Timestamp, nil

}

func (coffeeShop *CoffeeShop) getMembershipCoffeeQuota(membershipType string, coffeeType string) (CoffeeQuota, bool) {
	membership, found := coffeeShop.Quotas.Membership[membershipType]
	if !found {
		return CoffeeQuota{}, false
	}

	coffeeQuota, found := membership[coffeeType]
	if !found {
		return CoffeeQuota{}, false
	}

	return coffeeQuota, true
}

func isValidUserID(userID string) (int, bool) {
	if userID, err := strconv.Atoi(userID); err == nil {
		return userID, true
	}
	return 0, false
}

func isValidMembershipType(membershipType string) bool {
	return membershipType == "Basic" || membershipType == "Coffee Lover" || membershipType == "Americano Maniac"
}

func isValidCoffeeType(coffeeType string) bool {
	return coffeeType == "Espresso" || coffeeType == "Americano" || coffeeType == "Cappuccino"
}
