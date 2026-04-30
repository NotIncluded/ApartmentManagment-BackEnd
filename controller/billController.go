package controller

import (
	"net/http"
	"time"

	"github.com/PunMung-66/ApartmentSys/service"
	"github.com/gin-gonic/gin"
)

type BillController struct {
	billService service.BillService
}

func NewBillController(bs service.BillService) *BillController {
	return &BillController{billService: bs}
}

// GenerateRequest represents the expected JSON payload from the frontend
type GenerateRequest struct {
	ContractID string `json:"contract_id" binding:"required"`
	RecordDate string `json:"record_date" binding:"required"` // Format: YYYY-MM-DD
	DueDate    string `json:"due_date" binding:"required"`    // Format: YYYY-MM-DD
}

// GenerateBill handles the POST /bills/generate route
func (ctrl *BillController) GenerateBill(c *gin.Context) {
	var req GenerateRequest

	// 1. Bind and validate the JSON payload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	// 2. Parse the dates
	recordDate, err := time.Parse("2006-01-02", req.RecordDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid record_date format. Use YYYY-MM-DD"})
		return
	}

	dueDate, err := time.Parse("2006-01-02", req.DueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid due_date format. Use YYYY-MM-DD"})
		return
	}

	// 3. Call the Service Coordinator
	bill, err := ctrl.billService.GenerateMonthlyBill(req.ContractID, recordDate, dueDate)
	if err != nil {
		// Return a 400 Bad Request if a business rule fails (e.g., BR-12)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 4. Return the generated bill as a success response
	c.JSON(http.StatusCreated, gin.H{
		"message": "Bill generated successfully",
		"data":    bill,
	})
}