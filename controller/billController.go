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

// GenerateRequest represents the expected JSON payload
type GenerateRequest struct {
	RoomID     string `json:"room_id" binding:"required"`
	ContractID string `json:"contract_id" binding:"required"` // ✅ Added Contract ID
	RecordDate string `json:"record_date" binding:"required"` 
	DueDate    string `json:"due_date" binding:"required"`    
}

// GenerateBill handles the POST /bills/generate route
func (ctrl *BillController) GenerateBill(c *gin.Context) {
	var req GenerateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

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

	// ✅ Pass both the RoomID (for BR-02) and ContractID (for usage data)
	bill, err := ctrl.billService.GenerateMonthlyBill(req.RoomID, req.ContractID, recordDate, dueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Bill generated successfully",
		"data":    bill,
	})
}