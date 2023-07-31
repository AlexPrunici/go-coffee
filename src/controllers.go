package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func (coffeeShop *CoffeeShop) getQuotaConfigController(context *gin.Context) {
	if coffeeShop.Quotas == nil || len(coffeeShop.Quotas.Membership) == 0 {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Quota configuration not available"})
		return
	}
	context.JSON(http.StatusOK, coffeeShop.Quotas.Membership)
}

func (coffeeShop *CoffeeShop) getOrdersController(context *gin.Context) {
	var orders []OrderRequest
	if err := coffeeShop.DB.Find(&orders).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}
	context.JSON(http.StatusOK, orders)
}

func (coffeeShop *CoffeeShop) orderCoffeeController(context *gin.Context) {
	var requestBody OrderCoffeeRequestBody
	currentTime := time.Now()

	if err := context.BindJSON(&requestBody); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	userID, isIDValid := isValidUserID(requestBody.UserID)
	if !isIDValid {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UserID"})
		return
	}

	if !isValidMembershipType(requestBody.MembershipType) {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid membership type"})
		return
	}

	if !isValidCoffeeType(requestBody.CoffeeType) {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid coffee type"})
		return
	}

	coffeeQuota, err := coffeeShop.getMembershipCoffeeQuota(requestBody.MembershipType, requestBody.CoffeeType)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get quota", "details": err.Error()})
		return
	}

	ordersCount, err := coffeeShop.countUserOrders(userID, requestBody.CoffeeType, coffeeQuota.Duration)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count user orders", "details": err.Error()})
		return
	}

	if int(ordersCount) >= coffeeQuota.Amount {
		lastUsage, err := coffeeShop.getLastOrderTimestamp(userID, requestBody.CoffeeType)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count user orders", "details": err.Error()})
			return
		}
		waitTime := lastUsage.Add(time.Duration(coffeeQuota.Duration) * time.Minute).Sub(currentTime)

		context.JSON(http.StatusTooManyRequests, gin.H{"error": fmt.Sprintf("Limit exceded. Wait %.0f minutes", waitTime.Minutes())})
		return
	}

	orderRequest := OrderRequest{
		UserID:     userID,
		CoffeeType: requestBody.CoffeeType,
		Timestamp:  currentTime,
	}

	coffeeShop.DB.Create(&orderRequest)
	context.JSON(http.StatusOK, gin.H{"success": "Order created with success"})
}

func (coffeeShop *CoffeeShop) countUserOrders(userID int, coffeeType string, quotaDuration int) (int64, error) {
	duration := time.Now().Add(-time.Duration(quotaDuration) * time.Minute)

	var count int64
	if err := coffeeShop.DB.Model(&OrderRequest{}).Where("user_id = ? AND coffee_type = ? AND timestamp > ?", userID, coffeeType, duration).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (coffeeShop *CoffeeShop) getLastOrderTimestamp(userID int, coffeeType string) (time.Time, error) {
	var lastOrder OrderRequest
	if err := coffeeShop.DB.Model(&OrderRequest{}).Where("user_id = ? AND coffee_type = ?", userID, coffeeType).Order("timestamp desc").First(&lastOrder).Error; err != nil {
		return time.Time{}, err
	}
	return lastOrder.Timestamp, nil
}

func (coffeeShop *CoffeeShop) getMembershipCoffeeQuota(membershipType string, coffeeType string) (CoffeeQuota, error) {
	membership, found := coffeeShop.Quotas.Membership[membershipType]
	if !found {
		return CoffeeQuota{}, fmt.Errorf("quota not found for membership type: %s", membershipType)
	}

	coffeeQuota, found := membership[coffeeType]
	if !found {
		return CoffeeQuota{}, fmt.Errorf("quota not found for coffee type: %s", coffeeType)
	}

	return coffeeQuota, nil
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
