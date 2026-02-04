package handler

import (
	"fmt"
	"net/http"
	"smart-shift-scheduler/internal/domain"
	"smart-shift-scheduler/internal/usecase"

	"github.com/gin-gonic/gin"
)

type RequestHandler struct {
	usecase *usecase.ShiftUsecase // ShiftUsecaseに機能をまとめたのでこれを借用
}

func NewRequestHandler(u *usecase.ShiftUsecase) *RequestHandler {
	return &RequestHandler{usecase: u}
}

// Create: 希望休の登録
func (h *RequestHandler) Create(c *gin.Context) {
	type RequestBody struct {
		StaffID int    `json:"staff_id"`
		Date    string `json:"date"` // "2026-02-10"
	}
	var req RequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	newReq := &domain.ShiftRequest{
		StaffID: req.StaffID,
		Date:    req.Date,
		Type:    "NG",
	}

	if err := h.usecase.AddRequest(newReq); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, newReq)
}

// List: 一覧取得
func (h *RequestHandler) List(c *gin.Context) {
	list, err := h.usecase.GetAllRequests()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// Delete: 削除
func (h *RequestHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	var id uint
	fmt.Sscanf(idStr, "%d", &id)
	if err := h.usecase.DeleteRequest(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}