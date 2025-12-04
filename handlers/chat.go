package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"gin-project/config"
	"gin-project/models"
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

func RequestChat(c *gin.Context) {
	currentUserID := c.GetUint("userID")
	targetUserIDStr := c.Param("userID")

	// Convert string to uint
	targetUserID, err := strconv.ParseUint(targetUserIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Target userID must be a valid number",
		})
		return
	}

	targetUserIDUint := uint(targetUserID)
	if currentUserID == targetUserIDUint {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cannot request chat with yourself",
		})
		return
	}

	roomResult := utils.CheckRoomExists(currentUserID, targetUserID)
	if roomResult.Exists {
		c.JSON(http.StatusConflict, gin.H{
			"message": "Room already exists",
			"room_id": roomResult.RoomID,
		})
		return
	}

	var targetUser models.User
	if err := config.DB.First(&targetUser, targetUserIDUint).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Target user not found",
		})
	}

	requestExists := config.DB.Where("sender_id = ? AND receiver_id = ? AND status = ? OR status = ?", currentUserID, targetUserIDUint, "pending", "accepted").First(&models.ChatRequest{})
	if requestExists.Error == nil {
		c.JSON(http.StatusConflict, gin.H{
			"message": "Chat request already sent to this user",
		})
		return
	}

	chat_request := models.ChatRequest{
		SenderID:   currentUserID,
		ReceiverID: targetUserIDUint,
	}

	if err := config.DB.Create(&chat_request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create chat request",
		})
		return
	}

	// Preload the sender and receiver relationships
	if err := config.DB.Preload("Sender").Preload("Receiver").Where("id = ?", chat_request.ID).First(&chat_request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Chat request created but failed to load relationships",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "Chat request created successfully",
		"chat_request": chat_request,
	})
}

func ListReceivedChatRequests(c *gin.Context) {
	currentUserID := c.GetUint("userID")

	paginationParams := utils.GetPaginationParams(c)

	query := config.DB.Model(&models.ChatRequest{}).Preload("Sender").Preload("Receiver").Where("receiver_id = ?", currentUserID)
	paginatedQuery, paginationResult := utils.Paginate(query, paginationParams)

	var chatRequests []models.ChatRequest
	if err := paginatedQuery.Find(&chatRequests).Error; err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch chat requests",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Chat requests fetched successfully",
		"data":       chatRequests,
		"pagination": paginationResult,
	})
}

func ListSentChatRequests(c *gin.Context) {
	currentUserID := c.GetUint("userID")

	paginationParams := utils.GetPaginationParams(c)

	query := config.DB.Model(&models.ChatRequest{}).Preload("Sender").Preload("Receiver").Where("sender_id = ?", currentUserID)
	paginatedQuery, paginationResult := utils.Paginate(query, paginationParams)

	var chatRequests []models.ChatRequest
	if err := paginatedQuery.Find(&chatRequests).Error; err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch chat requests",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Sent chat requests fetched successfully",
		"data":       chatRequests,
		"pagination": paginationResult,
	})
}

func RespondToChatRequest(c *gin.Context) {
	currentUserID := c.GetUint("userID")
	chatRequestID := c.Param("requestID")

	var req models.AcceptOrRejectChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationResponse := utils.FormatValidationErrors(err)
		c.JSON(http.StatusBadRequest, validationResponse)
		return
	}

	var chatRequest models.ChatRequest
	if err := config.DB.Preload("Sender").Preload("Receiver").First(&chatRequest, "id = ? AND receiver_id = ?", chatRequestID, currentUserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Chat request not found",
		})
		return
	}

	if chatRequest.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Chat request is not pending",
		})
		return
	}

	log.Printf("Responding to chat request: %s with accept = %t", chatRequest.ID, req.Accept)

	if !req.Accept {
		log.Println("Rejecting chat request")
		chatRequest.Status = "rejected"
		chatRequest.UpdatedAt = utils.GetCurrentTimestamp()
		config.DB.Save(&chatRequest)
	} else {
		log.Println("Creating room between users")
		chatRequest.Status = "accepted"
		chatRequest.UpdatedAt = utils.GetCurrentTimestamp()
		config.DB.Save(&chatRequest)

		// Create room between users
		roomResult := utils.CreateRoomBetweenUsers(chatRequest.SenderID, chatRequest.ReceiverID)
		if !roomResult.Success || roomResult.Room == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to create room after accepting chat request",
				"error":   roomResult.Error,
			})
			return
		}

		// Create notification for the receiver
		notification := models.Notification{
			UserID:   chatRequest.SenderID,
			Message:  fmt.Sprintf("Your chat request to %s has been accepted.", chatRequest.Receiver.Name),
			Metadata: json.RawMessage(fmt.Sprintf(`{"room_id": "%s"}`, roomResult.Room.ID)),
		}
		config.DB.Create(&notification)
		log.Printf("Notification created successfully for user: %d", notification.UserID)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Chat request %s successfully", chatRequest.Status),
		"data":    chatRequest,
	})

}
