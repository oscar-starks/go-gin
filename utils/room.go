package utils

import (
	"fmt"

	"gin-project/config"
	"gin-project/models"
)

func CheckRoomExists(userID1, userID2 any) models.RoomCheckResult {
	roomID1 := fmt.Sprintf("%v_%v", userID1, userID2)
	roomID2 := fmt.Sprintf("%v_%v", userID2, userID1)

	var room models.Room
	result := config.DB.Where("id = ? OR id = ?", roomID1, roomID2).First(&room)

	if result.Error != nil {
		return models.RoomCheckResult{
			Exists: false,
			RoomID: nil,
		}
	}

	return models.RoomCheckResult{
		Exists: true,
		RoomID: room.ID,
	}
}

func CreateRoomBetweenUsers(currentUserID any, targetUserID any) models.RoomCreateResult {
	roomID := fmt.Sprintf("%v_%v", currentUserID, targetUserID)

	existingRoom := CheckRoomExists(currentUserID, targetUserID)
	if existingRoom.Exists {
		return models.RoomCreateResult{
			Success: false,
			Error:   "Room already exists",
			Room:    nil,
		}
	}

	// Get both users
	var currentUser, targetUser models.User
	if err := config.DB.First(&currentUser, currentUserID).Error; err != nil {
		return models.RoomCreateResult{
			Success: false,
			Error:   "Current user not found",
			Room:    nil,
		}
	}

	if err := config.DB.First(&targetUser, targetUserID).Error; err != nil {
		return models.RoomCreateResult{
			Success: false,
			Error:   "Target user not found",
			Room:    nil,
		}
	}

	// Create new room
	room := models.Room{
		ID:      roomID,
		Members: []models.User{currentUser, targetUser},
	}

	if err := config.DB.Create(&room).Error; err != nil {
		return models.RoomCreateResult{
			Success: false,
			Error:   "Failed to create room",
			Room:    nil,
		}
	}

	return models.RoomCreateResult{
		Success: true,
		Error:   "",
		Room:    &room,
	}
}
