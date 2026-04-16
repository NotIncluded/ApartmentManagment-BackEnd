package controller

import (
	"net/http"

	"github.com/PunMung-66/ApartmentSys/internal/response"
	"github.com/PunMung-66/ApartmentSys/service"
	"github.com/gin-gonic/gin"
)

type ContractController struct {
	contractService *service.ContractService
}

func NewContractController(contractService *service.ContractService) *ContractController {
	return &ContractController{
		contractService: contractService,
	}
}

func (cc *ContractController) CreateContract(c *gin.Context) {
	var body struct {
		UserID    string `json:"user_id" binding:"required"`
		RoomID    string `json:"room_id" binding:"required"`
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

	contract, err := cc.contractService.CreateContract(body.UserID, body.RoomID, body.StartDate, body.EndDate, body.Status)
	if err != nil {
		switch err.Error() {
		case "room not found":
			res := response.NewAppResponse(http.StatusNotFound, "Room not found", nil)
			c.JSON(res.Status, res.Response())
		case "room is not available":
			res := response.NewAppResponse(http.StatusBadRequest, "Room is not available", nil)
			c.JSON(res.Status, res.Response())
		case "user not found":
			res := response.NewAppResponse(http.StatusNotFound, "User not found", nil)
			c.JSON(res.Status, res.Response())
		case "user already has an active contract":
			res := response.NewAppResponse(http.StatusBadRequest, "User already has an active contract", nil)
			c.JSON(res.Status, res.Response())
		case "invalid start date format (use YYYY-MM-DD)":
			res := response.NewAppResponse(http.StatusBadRequest, "Invalid start date format", nil)
			c.JSON(res.Status, res.Response())
		case "invalid end date format (use YYYY-MM-DD)":
			res := response.NewAppResponse(http.StatusBadRequest, "Invalid end date format", nil)
			c.JSON(res.Status, res.Response())
		case "end date must be after start date":
			res := response.NewAppResponse(http.StatusBadRequest, "End date must be after start date", nil)
			c.JSON(res.Status, res.Response())
		default:
			res := response.NewAppResponse(http.StatusInternalServerError, "Failed to create contract", err.Error())
			c.JSON(res.Status, res.Response())
		}
		return
	}

	res := response.NewAppResponse(http.StatusCreated, "Contract created successfully", contract)
	c.JSON(res.Status, res.Response())
}

func (cc *ContractController) GetContracts(c *gin.Context) {
	contracts, err := cc.contractService.GetContracts()
	if err != nil {
		res := response.NewAppResponse(http.StatusInternalServerError, "Failed to retrieve contracts", err.Error())
		c.JSON(res.Status, res.Response())
		return
	}

	res := response.NewAppResponse(http.StatusOK, "Contracts retrieved successfully", contracts)
	c.JSON(res.Status, res.Response())
}

func (cc *ContractController) GetContractByID(c *gin.Context) {
	contractID := c.Param("id")

	contract, err := cc.contractService.GetContractByID(contractID)
	if err != nil {
		res := response.NewAppResponse(http.StatusNotFound, "Contract not found", nil)
		c.JSON(res.Status, res.Response())
		return
	}

	res := response.NewAppResponse(http.StatusOK, "Contract retrieved successfully", contract)
	c.JSON(res.Status, res.Response())
}

func (cc *ContractController) GetContractsByUserID(c *gin.Context) {
	userID := c.Param("userID")

	contracts, err := cc.contractService.GetContractsByUserID(userID)
	if err != nil {
		switch err.Error() {
		case "user not found":
			res := response.NewAppResponse(http.StatusNotFound, "User not found", nil)
			c.JSON(res.Status, res.Response())
		default:
			res := response.NewAppResponse(http.StatusInternalServerError, "Failed to retrieve contracts", err.Error())
			c.JSON(res.Status, res.Response())
		}
		return
	}

	res := response.NewAppResponse(http.StatusOK, "Contracts retrieved successfully", contracts)
	c.JSON(res.Status, res.Response())
}

func (cc *ContractController) GetContractsByRoomID(c *gin.Context) {
	roomID := c.Param("roomID")

	contracts, err := cc.contractService.GetContractsByRoomID(roomID)
	if err != nil {
		switch err.Error() {
		case "room not found":
			res := response.NewAppResponse(http.StatusNotFound, "Room not found", nil)
			c.JSON(res.Status, res.Response())
		default:
			res := response.NewAppResponse(http.StatusInternalServerError, "Failed to retrieve contracts", err.Error())
			c.JSON(res.Status, res.Response())
		}
		return
	}

	res := response.NewAppResponse(http.StatusOK, "Contracts retrieved successfully", contracts)
	c.JSON(res.Status, res.Response())
}

func (cc *ContractController) UpdateContract(c *gin.Context) {
	contractID := c.Param("id")

	var body struct {
		UserID    string `json:"user_id"`
		RoomID    string `json:"room_id"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		Status    string `json:"status"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		res := response.NewAppResponse(http.StatusBadRequest, "Invalid request body", err.Error())
		c.JSON(res.Status, res.Response())
		return
	}

	contract, err := cc.contractService.UpdateContract(contractID, body.UserID, body.RoomID, body.StartDate, body.EndDate, body.Status)
	if err != nil {
		switch err.Error() {
		case "contract not found":
			res := response.NewAppResponse(http.StatusNotFound, "Contract not found", nil)
			c.JSON(res.Status, res.Response())
		case "user not found":
			res := response.NewAppResponse(http.StatusNotFound, "User not found", nil)
			c.JSON(res.Status, res.Response())
		case "user already has an active contract":
			res := response.NewAppResponse(http.StatusBadRequest, "User already has an active contract", nil)
			c.JSON(res.Status, res.Response())
		case "room not found":
			res := response.NewAppResponse(http.StatusNotFound, "Room not found", nil)
			c.JSON(res.Status, res.Response())
		case "target room is not available":
			res := response.NewAppResponse(http.StatusBadRequest, "Target room is not available", nil)
			c.JSON(res.Status, res.Response())
		case "invalid start date format (use YYYY-MM-DD)":
			res := response.NewAppResponse(http.StatusBadRequest, "Invalid start date format", nil)
			c.JSON(res.Status, res.Response())
		case "invalid end date format (use YYYY-MM-DD)":
			res := response.NewAppResponse(http.StatusBadRequest, "Invalid end date format", nil)
			c.JSON(res.Status, res.Response())
		case "end date must be after start date":
			res := response.NewAppResponse(http.StatusBadRequest, "End date must be after start date", nil)
			c.JSON(res.Status, res.Response())
		default:
			res := response.NewAppResponse(http.StatusInternalServerError, "Failed to update contract", err.Error())
			c.JSON(res.Status, res.Response())
		}
		return
	}

	res := response.NewAppResponse(http.StatusOK, "Contract updated successfully", contract)
	c.JSON(res.Status, res.Response())
}

func (cc *ContractController) DeleteContract(c *gin.Context) {
	contractID := c.Param("id")

	err := cc.contractService.DeleteContract(contractID)
	if err != nil {
		switch err.Error() {
		case "contract not found":
			res := response.NewAppResponse(http.StatusNotFound, "Contract not found", nil)
			c.JSON(res.Status, res.Response())
		default:
			res := response.NewAppResponse(http.StatusInternalServerError, "Failed to delete contract", err.Error())
			c.JSON(res.Status, res.Response())
		}
		return
	}

	res := response.NewAppResponse(http.StatusOK, "Contract deleted successfully", nil)
	c.JSON(res.Status, res.Response())
}
