package controller

import (
	"net/http"

	"github.com/PunMung-66/ApartmentSys/internal/response"
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

// GetListRoom handles GET /rooms
// STAFF: returns all rooms
// TENANT: returns room via active contract
func (rc *RoomController) GetListRoom(c *gin.Context) {
	role, _ := c.Get("role")
	currentUserID, _ := c.Get("user_id")

	// STAFF: full access to all rooms
	if role == "STAFF" {
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
		return
	}

	// TENANT: access room through active contract only
	if role == "TENANT" {
		userID := currentUserID.(string)
		room, err := rc.roomService.GetRoomByUserID(userID)
		if err != nil {
			res := response.NewAppResponse(
				http.StatusForbidden,
				"Access denied: tenant has no active contract",
				nil,
			)
			c.JSON(res.Status, res.Response())
			return
		}

		res := response.NewAppResponse(
			http.StatusOK,
			"Room retrieved successfully",
			[]interface{}{room},
		)
		c.JSON(res.Status, res.Response())
		return
	}

	// Invalid role
	res := response.NewAppResponse(
		http.StatusForbidden,
		"Invalid role",
		nil,
	)
	c.JSON(res.Status, res.Response())
}

// GetMyRoom handles GET /me/room (TENANT only endpoint)
// Returns the tenant's room via their active contract
func (rc *RoomController) GetMyRoom(c *gin.Context) {
	role, _ := c.Get("role")
	currentUserID, _ := c.Get("user_id")

	// Only TENANT can access this endpoint
	if role != "TENANT" {
		res := response.NewAppResponse(
			http.StatusForbidden,
			"Only tenants can access this endpoint",
			nil,
		)
		c.JSON(res.Status, res.Response())
		return
	}

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
