package controller

import (
	"net/http"

	"github.com/PunMung-66/ApartmentSys/internal/response"
	"github.com/PunMung-66/ApartmentSys/service"
	"github.com/gin-gonic/gin"
)

type UtilityController struct {
	utilityService *service.UtilityService
}

type UtilityControllerInterface interface {
	ConfigureRate(c *gin.Context)
	CreateRate(c *gin.Context)
	GetRateByID(c *gin.Context)
	GetAllRates(c *gin.Context)
	UpdateRate(c *gin.Context)
	DeleteRate(c *gin.Context)
	RecordUsage(c *gin.Context)
	GetUsageByID(c *gin.Context)
	GetUsagesByContract(c *gin.Context)
	UpdateUsage(c *gin.Context)
	DeleteUsage(c *gin.Context)
	GetMyUsages(c *gin.Context)
	GetMyLatestUsage(c *gin.Context)
}

func NewUtilityController(utilityService *service.UtilityService) *UtilityController {
	return &UtilityController{
		utilityService: utilityService,
	}
}

func (uc *UtilityController) ConfigureRate(c *gin.Context) {
	var req service.ConfigureRateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := response.NewAppResponse(
			http.StatusBadRequest,
			"Invalid request body",
			err.Error(),
		)
		c.JSON(appErr.Status, appErr.Response())
		return
	}

	// Basic validation (controller-level)
	if req.Period == "" {
		appErr := response.NewAppResponse(
			http.StatusBadRequest,
			"Missing required fields",
			"period is required",
		)
		c.JSON(appErr.Status, appErr.Response())
		return
	}
	if req.WaterRate < 0 || req.ElectricRate < 0 || req.CommonFee < 0 {
		appErr := response.NewAppResponse(
			http.StatusBadRequest,
			"Invalid rates",
			"rates must be >= 0",
		)
		c.JSON(appErr.Status, appErr.Response())
		return
	}

	rate, err := uc.utilityService.ConfigureRate(req)
	if err != nil {
		appErr := response.NewAppResponse(
			http.StatusBadRequest,
			"Failed to configure rate",
			err.Error(),
		)
		c.JSON(appErr.Status, appErr.Response())
		return
	}

	appRes := response.NewAppResponse(
		http.StatusCreated,
		"Utility rate configured",
		rate,
	)
	c.JSON(appRes.Status, appRes.Response())
}

func (uc *UtilityController) CreateRate(c *gin.Context) {
	var req service.ConfigureRateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := response.NewAppResponse(
			http.StatusBadRequest,
			"Invalid request body",
			err.Error(),
		)
		c.JSON(appErr.Status, appErr.Response())
		return
	}

	if req.WaterRate < 0 || req.ElectricRate < 0 || req.CommonFee < 0 {
		appErr := response.NewAppResponse(
			http.StatusBadRequest,
			"Invalid rates",
			"rates must be >= 0",
		)
		c.JSON(appErr.Status, appErr.Response())
		return
	}

	rate, err := uc.utilityService.CreateRate(req)
	if err != nil {
		appErr := response.NewAppResponse(
			http.StatusBadRequest,
			"Failed to create rate",
			err.Error(),
		)
		c.JSON(appErr.Status, appErr.Response())
		return
	}

	appRes := response.NewAppResponse(
		http.StatusCreated,
		"Utility rate created successfully",
		rate,
	)
	c.JSON(appRes.Status, appRes.Response())
}

func (uc *UtilityController) GetRateByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		appErr := response.NewAppResponse(
			http.StatusBadRequest,
			"Rate ID is required",
			nil,
		)
		c.JSON(appErr.Status, appErr.Response())
		return
	}

	rate, err := uc.utilityService.GetRateByID(id)
	if err != nil {
		appErr := response.NewAppResponse(
			http.StatusNotFound,
			"Rate not found",
			err.Error(),
		)
		c.JSON(appErr.Status, appErr.Response())
		return
	}

	appRes := response.NewAppResponse(
		http.StatusOK,
		"Rate retrieved successfully",
		rate,
	)
	c.JSON(appRes.Status, appRes.Response())
}

func (uc *UtilityController) GetAllRates(c *gin.Context) {
	rates, err := uc.utilityService.GetAllRates()
	if err != nil {
		appErr := response.NewAppResponse(
			http.StatusInternalServerError,
			"Failed to retrieve rates",
			err.Error(),
		)
		c.JSON(appErr.Status, appErr.Response())
		return
	}

	appRes := response.NewAppResponse(
		http.StatusOK,
		"Rates retrieved successfully",
		rates,
	)
	c.JSON(appRes.Status, appRes.Response())
}

func (uc *UtilityController) UpdateRate(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		appErr := response.NewAppResponse(
			http.StatusBadRequest,
			"Rate ID is required",
			nil,
		)
		c.JSON(appErr.Status, appErr.Response())
		return
	}

	var req service.ConfigureRateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := response.NewAppResponse(
			http.StatusBadRequest,
			"Invalid request body",
			err.Error(),
		)
		c.JSON(appErr.Status, appErr.Response())
		return
	}

	if req.WaterRate < 0 || req.ElectricRate < 0 || req.CommonFee < 0 {
		appErr := response.NewAppResponse(
			http.StatusBadRequest,
			"Invalid rates",
			"rates must be >= 0",
		)
		c.JSON(appErr.Status, appErr.Response())
		return
	}

	rate, err := uc.utilityService.UpdateRate(id, req)
	if err != nil {
		appErr := response.NewAppResponse(
			http.StatusBadRequest,
			"Failed to update rate",
			err.Error(),
		)
		c.JSON(appErr.Status, appErr.Response())
		return
	}

	appRes := response.NewAppResponse(
		http.StatusOK,
		"Utility rate updated successfully",
		rate,
	)
	c.JSON(appRes.Status, appRes.Response())
}

func (uc *UtilityController) DeleteRate(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		appErr := response.NewAppResponse(
			http.StatusBadRequest,
			"Rate ID is required",
			nil,
		)
		c.JSON(appErr.Status, appErr.Response())
		return
	}

	err := uc.utilityService.DeleteRate(id)
	if err != nil {
		appErr := response.NewAppResponse(
			http.StatusBadRequest,
			"Failed to delete rate",
			err.Error(),
		)
		c.JSON(appErr.Status, appErr.Response())
		return
	}

	appRes := response.NewAppResponse(
		http.StatusOK,
		"Utility rate deleted successfully",
		nil,
	)
	c.JSON(appRes.Status, appRes.Response())
}

func (uc *UtilityController) RecordUsage(c *gin.Context) {
	var req service.RecordUsageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := response.NewAppResponse(
			http.StatusBadRequest,
			"Invalid request body",
			err.Error(),
		)
		c.JSON(appErr.Status, appErr.Response())
		return
	}

	if req.ContractID == "" {
		appErr := response.NewAppResponse(
			http.StatusBadRequest,
			"Missing required fields",
			"contract_id is required",
		)
		c.JSON(appErr.Status, appErr.Response())
		return
	}

	usage, err := uc.utilityService.RecordUsage(req)
	if err != nil {
		appErr := response.NewAppResponse(
			http.StatusBadRequest,
			"Failed to record usage",
			err.Error(),
		)
		c.JSON(appErr.Status, appErr.Response())
		return
	}

	appRes := response.NewAppResponse(
		http.StatusCreated,
		"Utility usage recorded successfully",
		usage,
	)
	c.JSON(appRes.Status, appRes.Response())
}

func (uc *UtilityController) GetUsageByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		appErr := response.NewAppResponse(
			http.StatusBadRequest,
			"Usage ID is required",
			nil,
		)
		c.JSON(appErr.Status, appErr.Response())
		return
	}

	usage, err := uc.utilityService.GetUsageByID(id)
	if err != nil {
		appErr := response.NewAppResponse(
			http.StatusNotFound,
			"Usage not found",
			err.Error(),
		)
		c.JSON(appErr.Status, appErr.Response())
		return
	}

	appRes := response.NewAppResponse(
		http.StatusOK,
		"Usage retrieved successfully",
		usage,
	)
	c.JSON(appRes.Status, appRes.Response())
}

func (uc *UtilityController) GetUsagesByContract(c *gin.Context) {
	contractID := c.Param("contractID")
	if contractID == "" {
		appErr := response.NewAppResponse(
			http.StatusBadRequest,
			"Contract ID is required",
			nil,
		)
		c.JSON(appErr.Status, appErr.Response())
		return
	}

	usages, err := uc.utilityService.GetUsagesByContract(contractID)
	if err != nil {
		appErr := response.NewAppResponse(
			http.StatusNotFound,
			"Failed to retrieve usages",
			err.Error(),
		)
		c.JSON(appErr.Status, appErr.Response())
		return
	}

	appRes := response.NewAppResponse(
		http.StatusOK,
		"Usages retrieved successfully",
		usages,
	)
	c.JSON(appRes.Status, appRes.Response())
}

func (uc *UtilityController) UpdateUsage(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		appErr := response.NewAppResponse(
			http.StatusBadRequest,
			"Usage ID is required",
			nil,
		)
		c.JSON(appErr.Status, appErr.Response())
		return
	}

	var req service.RecordUsageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := response.NewAppResponse(
			http.StatusBadRequest,
			"Invalid request body",
			err.Error(),
		)
		c.JSON(appErr.Status, appErr.Response())
		return
	}

	usage, err := uc.utilityService.UpdateUsage(id, req)
	if err != nil {
		appErr := response.NewAppResponse(
			http.StatusBadRequest,
			"Failed to update usage",
			err.Error(),
		)
		c.JSON(appErr.Status, appErr.Response())
		return
	}

	appRes := response.NewAppResponse(
		http.StatusOK,
		"Utility usage updated successfully",
		usage,
	)
	c.JSON(appRes.Status, appRes.Response())
}

func (uc *UtilityController) DeleteUsage(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		appErr := response.NewAppResponse(
			http.StatusBadRequest,
			"Usage ID is required",
			nil,
		)
		c.JSON(appErr.Status, appErr.Response())
		return
	}

	err := uc.utilityService.DeleteUsage(id)
	if err != nil {
		appErr := response.NewAppResponse(
			http.StatusBadRequest,
			"Failed to delete usage",
			err.Error(),
		)
		c.JSON(appErr.Status, appErr.Response())
		return
	}

	appRes := response.NewAppResponse(
		http.StatusOK,
		"Utility usage deleted successfully",
		nil,
	)
	c.JSON(appRes.Status, appRes.Response())
}

func (uc *UtilityController) GetMyUsages(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		appErr := response.NewAppResponse(
			http.StatusUnauthorized,
			"Unauthorized",
			"User ID not found in token",
		)
		c.JSON(appErr.Status, appErr.Response())
		return
	}

	usages, err := uc.utilityService.GetMyUsages(userID.(string))
	if err != nil {
		appErr := response.NewAppResponse(
			http.StatusNotFound,
			"Failed to retrieve usages",
			err.Error(),
		)
		c.JSON(appErr.Status, appErr.Response())
		return
	}

	appRes := response.NewAppResponse(
		http.StatusOK,
		"Usages retrieved successfully",
		usages,
	)
	c.JSON(appRes.Status, appRes.Response())
}

func (uc *UtilityController) GetMyLatestUsage(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		appErr := response.NewAppResponse(
			http.StatusUnauthorized,
			"Unauthorized",
			"User ID not found in token",
		)
		c.JSON(appErr.Status, appErr.Response())
		return
	}

	usage, err := uc.utilityService.GetMyLatestUsage(userID.(string))
	if err != nil {
		appErr := response.NewAppResponse(
			http.StatusNotFound,
			"Failed to retrieve latest usage",
			err.Error(),
		)
		c.JSON(appErr.Status, appErr.Response())
		return
	}

	appRes := response.NewAppResponse(
		http.StatusOK,
		"Latest usage retrieved successfully",
		usage,
	)
	c.JSON(appRes.Status, appRes.Response())
}
