
package controller

import (
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
	"github.com/PunMung-66/ApartmentSys/service"
)

type BillSlipController struct {
	service *service.BillSlipService
}

func NewBillSlipController(service *service.BillSlipService) *BillSlipController {
	return &BillSlipController{service: service}
}

func (ctrl *BillSlipController) UploadBillSlip(c *gin.Context) {
	billID := c.PostForm("bill_id")
	roomID := c.PostForm("room_id")
	file, err := c.FormFile("slip")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file is received"})
		return
	}
	filePath := "/tmp/" + file.Filename
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	defer func() { _ = os.Remove(filePath) }()
	contentType := file.Header.Get("Content-Type")
	slipURL, err := ctrl.service.UploadSlip(c.Request.Context(), billID, roomID, filePath, contentType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"slip_url": slipURL})
}
