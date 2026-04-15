package controller

import (
	"net/http"

	"github.com/PunMung-66/ApartmentSys/internal/response"
	"github.com/PunMung-66/ApartmentSys/model"
	"github.com/PunMung-66/ApartmentSys/service"
	"github.com/gin-gonic/gin"
)

type RoomController struct {
	roomService *service.RoomService
}

func NewRoomController(roomService *service.RoomService) *RoomController {
	return &RoomController{
		roomService: roomService,
	}
}

// CreateRoom handles POST /rooms (STAFF only)
// Creates a new room
func (rc *RoomController) CreateRoom(c *gin.Context) {
	var room model.Room
	if err := c.ShouldBindJSON(&room); err != nil {
		res := response.NewAppResponse(
			http.StatusBadRequest,
			"Invalid request body",
			err.Error(),
		)
		c.JSON(res.Status, res.Response())
		return
	}

	createdRoom, err := rc.roomService.CreateRoom(&room)
	if err != nil {
		res := response.NewAppResponse(
			http.StatusInternalServerError,
			"Failed to create room",
			err.Error(),
		)
		c.JSON(res.Status, res.Response())
		return
	}

	res := response.NewAppResponse(
		http.StatusCreated,
		"Room created successfully",
		createdRoom,
	)
	c.JSON(res.Status, res.Response())
}

// GetListRoom handles GET /rooms (STAFF only)
// Returns all rooms
func (rc *RoomController) GetListRoom(c *gin.Context) {
	rooms, err := rc.roomService.GetListRoom()
	if err != nil {
		res := response.NewAppResponse(
			http.StatusInternalServerError,
			"Failed to retrieve rooms",
			err.Error(),
		)
		c.JSON(res.Status, res.Response())
		return
	}

	res := response.NewAppResponse(
		http.StatusOK,
		"Rooms retrieved successfully",
		rooms,
	)
	c.JSON(res.Status, res.Response())
}

// GetRoomByID handles GET /rooms/{id} (STAFF only)
// Returns room detail
func (rc *RoomController) GetRoomByID(c *gin.Context) {
	roomID := c.Param("id")

	room, err := rc.roomService.GetRoomByID(roomID)
	if err != nil {
		res := response.NewAppResponse(
			http.StatusNotFound,
			"Room not found",
			err.Error(),
		)
		c.JSON(res.Status, res.Response())
		return
	}

	res := response.NewAppResponse(
		http.StatusOK,
		"Room retrieved successfully",
		room,
	)
	c.JSON(res.Status, res.Response())
}

// UpdateRoom handles PUT /rooms/{id} (STAFF only)
// Updates room information
func (rc *RoomController) UpdateRoom(c *gin.Context) {
	roomID := c.Param("id")

	var room model.Room
	if err := c.ShouldBindJSON(&room); err != nil {
		res := response.NewAppResponse(
			http.StatusBadRequest,
			"Invalid request body",
			err.Error(),
		)
		c.JSON(res.Status, res.Response())
		return
	}

	// Preserve the ID from URL
	room.ID = roomID

	updatedRoom, err := rc.roomService.UpdateRoom(&room)
	if err != nil {
		res := response.NewAppResponse(
			http.StatusInternalServerError,
			"Failed to update room",
			err.Error(),
		)
		c.JSON(res.Status, res.Response())
		return
	}

	res := response.NewAppResponse(
		http.StatusOK,
		"Room updated successfully",
		updatedRoom,
	)
	c.JSON(res.Status, res.Response())
}

// DeleteRoom handles DELETE /rooms/{id} (STAFF only)
// Deletes room - cannot delete if room has any contract
func (rc *RoomController) DeleteRoom(c *gin.Context) {
	roomID := c.Param("id")

	err := rc.roomService.DeleteRoom(roomID)
	if err != nil {
		// Check if error is due to existing contract
		if err.Error() == "cannot delete room with active contract" {
			res := response.NewAppResponse(
				http.StatusConflict,
				"Cannot delete room with existing contract",
				nil,
			)
			c.JSON(res.Status, res.Response())
			return
		}

		res := response.NewAppResponse(
			http.StatusInternalServerError,
			"Failed to delete room",
			err.Error(),
		)
		c.JSON(res.Status, res.Response())
		return
	}

	res := response.NewAppResponse(
		http.StatusOK,
		"Room deleted successfully",
		nil,
	)
	c.JSON(res.Status, res.Response())
}

// GetRoomActiveContract handles GET /rooms/{id}/contract (STAFF only)
// Returns the active contract of this room
func (rc *RoomController) GetRoomActiveContract(c *gin.Context) {
	roomID := c.Param("id")

	contract, err := rc.roomService.GetRoomActiveContract(roomID)
	if err != nil {
		res := response.NewAppResponse(
			http.StatusNotFound,
			"No active contract found for this room",
			nil,
		)
		c.JSON(res.Status, res.Response())
		return
	}

	res := response.NewAppResponse(
		http.StatusOK,
		"Active contract retrieved successfully",
		contract,
	)
	c.JSON(res.Status, res.Response())
}

// GetRoomContractHistory handles GET /rooms/{id}/contracts (STAFF only)
// Returns contract history for this room
func (rc *RoomController) GetRoomContractHistory(c *gin.Context) {
	roomID := c.Param("id")

	contracts, err := rc.roomService.GetRoomContractHistory(roomID)
	if err != nil {
		res := response.NewAppResponse(
			http.StatusInternalServerError,
			"Failed to retrieve contract history",
			err.Error(),
		)
		c.JSON(res.Status, res.Response())
		return
	}

	res := response.NewAppResponse(
		http.StatusOK,
		"Contract history retrieved successfully",
		contracts,
	)
	c.JSON(res.Status, res.Response())
}

// GetRoomTenant handles GET /rooms/{id}/tenant (STAFF only)
// Returns current tenant via active contract
func (rc *RoomController) GetRoomTenant(c *gin.Context) {
	roomID := c.Param("id")

	user, err := rc.roomService.GetRoomTenant(roomID)
	if err != nil {
		res := response.NewAppResponse(
			http.StatusNotFound,
			"No tenant assigned to this room",
			nil,
		)
		c.JSON(res.Status, res.Response())
		return
	}

	res := response.NewAppResponse(
		http.StatusOK,
		"Tenant retrieved successfully",
		user,
	)
	c.JSON(res.Status, res.Response())
}

// AssignRoom handles POST /rooms/{id}/assign (STAFF only)
// Assigns tenant to room by creating a contract
// Request body should include: userID, startDate, endDate, status
func (rc *RoomController) AssignRoom(c *gin.Context) {
	roomID := c.Param("id")

	var body struct {
		UserID    string `json:"user_id" binding:"required"`
		StartDate string `json:"start_date" binding:"required"`
		EndDate   string `json:"end_date"`
		Status    string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		res := response.NewAppResponse(
			http.StatusBadRequest,
			"Invalid request body",
			err.Error(),
		)
		c.JSON(res.Status, res.Response())
		return
	}

	contract, err := rc.roomService.AssignRoom(roomID, body.UserID, body.StartDate, body.EndDate, body.Status)
	if err != nil {
		switch err.Error() {
		case "room not found":
			res := response.NewAppResponse(
				http.StatusNotFound,
				"Room not found",
				nil,
			)
			c.JSON(res.Status, res.Response())
		case "user not found":
			res := response.NewAppResponse(
				http.StatusNotFound,
				"User not found",
				nil,
			)
			c.JSON(res.Status, res.Response())
		case "room is occupied":
			fallthrough
		case "room is in maintenance":
			res := response.NewAppResponse(
				http.StatusConflict,
				err.Error(),
				nil,
			)
			c.JSON(res.Status, res.Response())
		default:
			res := response.NewAppResponse(
				http.StatusInternalServerError,
				"Failed to assign room",
				err.Error(),
			)
			c.JSON(res.Status, res.Response())
		}
		return
	}

	res := response.NewAppResponse(
		http.StatusCreated,
		"Room assigned to tenant successfully",
		contract,
	)
	c.JSON(res.Status, res.Response())
}

// GetMyRoom handles GET /me/room (TENANT only endpoint)
// Returns the tenant's room via their active contract
func (rc *RoomController) GetMyRoom(c *gin.Context) {
	currentUserID, _ := c.Get("user_id")
	userID := currentUserID.(string)

	room, err := rc.roomService.GetRoomByUserID(userID)
	if err != nil {
		res := response.NewAppResponse(
			http.StatusForbidden,
			"Access denied: you have no active contract",
			nil,
		)
		c.JSON(res.Status, res.Response())
		return
	}

	res := response.NewAppResponse(
		http.StatusOK,
		"Your room retrieved successfully",
		room,
	)
	c.JSON(res.Status, res.Response())
}
