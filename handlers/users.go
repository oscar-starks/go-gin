package handlers

import (
	"net/http"
	"strconv"

	"gin-project/config"
	"gin-project/models"

	"github.com/gin-gonic/gin"
)

func GetProfile(c *gin.Context) {
	userIDParam := c.Query("userID")

	var targetUserID uint
	var err error
	var user models.User

	if userIDParam != "" {
		userID64, parseErr := strconv.ParseUint(userIDParam, 10, 32)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid userID parameter",
			})
			return
		}
		targetUserID = uint(userID64)

	} else {
		loggedInUserID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Something went wrong",
			})
			return
		}
		targetUserID = loggedInUserID.(uint)
	}

	// Fetch user from database
	if err = config.DB.First(&user, targetUserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	user.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"data":    user,
		"message": "Profile fetched successfully",
	})
}
