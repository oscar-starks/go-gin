package handlers

import (
	"log"
	"net/http"

	"gin-project/config"
	"gin-project/models"
	"gin-project/utils"

	"github.com/gin-gonic/gin"
)

func ListNotifications(c *gin.Context) {
	currentUserID := c.GetUint("userID")
	paginationParams := utils.GetPaginationParams(c)

	log.Println(currentUserID)

	query := config.DB.Model(&models.Notification{}).Where("user_id = ?", currentUserID).Order("created_at DESC")
	paginatedQuery, paginationResult := utils.Paginate(query, paginationParams)

	var notifications []models.Notification
	if err := paginatedQuery.Find(&notifications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch notifications",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Notifications fetched successfully",
		"data":       notifications,
		"pagination": paginationResult,
	})
}

func MarkNotificationAsRead() {
	// Implementation for marking a notification as read
}
