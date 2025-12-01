package handlers

import (
	"net/http"

	"gin-project/utils"

	"github.com/gin-gonic/gin"
)

func HasRoom(c *gin.Context) {
	currentUserID := c.GetUint("userID")
	targetUserID := c.Query("userID")

	if targetUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Target userID is required",
		})
		return
	}

	roomResult := utils.CheckRoomExists(currentUserID, targetUserID)

	c.JSON(http.StatusOK, gin.H{
		"exists":  roomResult.Exists,
		"room_id": roomResult.RoomID,
		"message": "Room existence checked successfully",
	})
}

// CreateRoom creates a new room between two users
func CreateRoom(c *gin.Context) {
	currentUserID := c.GetUint("userID")

	// Get target user ID from URL parameter
	targetUserID := c.Param("userID")
	if targetUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Target userID is required",
		})
		return
	}

	// Use utility function to create room
	result := utils.CreateRoomBetweenUsers(currentUserID, targetUserID)

	if !result.Success {
		statusCode := http.StatusInternalServerError
		if result.Error == "Room already exists" {
			statusCode = http.StatusConflict
		} else if result.Error == "Current user not found" || result.Error == "Target user not found" {
			statusCode = http.StatusNotFound
		}

		c.JSON(statusCode, gin.H{
			"error": result.Error,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Room created successfully",
		"room":    result.Room,
	})
}
