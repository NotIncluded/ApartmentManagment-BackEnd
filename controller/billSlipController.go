package controller

import (
	"net/http"
	"strings" // ✅ Added to help us check the error message

	"github.com/PunMung-66/ApartmentSys/service"
	"github.com/gin-gonic/gin"
)

type BillSlipController struct {
	service service.BillSlipService 
}

func NewBillSlipController(service service.BillSlipService) *BillSlipController { 
	return &BillSlipController{service: service}
}

func (ctrl *BillSlipController) UploadBillSlip(c *gin.Context) {
	billID := c.PostForm("bill_id")
	roomID := c.PostForm("room_id")

	// 1. Get the file header from the form
	fileHeader, err := c.FormFile("slip")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file is received"})
		return
	}

	// 2. Open the file into memory instead of saving to disk
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file stream"})
		return
	}
	defer file.Close() // Ensure the file stream is closed when done

	contentType := fileHeader.Header.Get("Content-Type")

	// 3. Pass the opened file (io.Reader) and the Filename to the service layer.
	slipURL, err := ctrl.service.UploadSlip(c.Request.Context(), billID, roomID, file, fileHeader.Filename, contentType)
	if err != nil {
		// ✅ PRO-TIP APPLIED: If the error is because the bill is missing, return 404!
		if strings.Contains(err.Error(), "does not exist") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		// Otherwise, it's a real server/upload error, so we keep the 500
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"slip_url": slipURL})
}