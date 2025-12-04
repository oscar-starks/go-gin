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

func MarkNotificationAsRead(c *gin.Context) {
	currentUserID := c.GetUint("userID")
	notificationID := c.Param("notificationID")

	var notification models.Notification
	if err := config.DB.Where("id = ? AND user_id = ?", notificationID, currentUserID).First(&notification).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Notification not found",
		})
		return
	}

	notification.Read = true
	config.DB.Save(&notification)

	c.JSON(http.StatusOK, gin.H{
		"message": "Notification marked as read",
	})

}
