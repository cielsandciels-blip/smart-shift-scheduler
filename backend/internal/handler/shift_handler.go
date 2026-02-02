package handler

import (
	"fmt" // ★これを追加しました！
	"net/http"
	"smart-shift-scheduler/internal/domain"
	"smart-shift-scheduler/internal/usecase"

	"github.com/gin-gonic/gin"
)

type ShiftHandler struct {
	usecase *usecase.ShiftUsecase
}

func NewShiftHandler(u *usecase.ShiftUsecase) *ShiftHandler {
	return &ShiftHandler{usecase: u}
}

// Generate: POST /api/shift
func (h *ShiftHandler) Generate(c *gin.Context) {
	type Request struct {
		StartDate    string                    `json:"start_date"`
		Requests     []domain.ShiftRequest     `json:"requests"`
		Requirements []domain.DailyRequirement `json:"requirements"`
		Days         int                       `json:"days"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	input := domain.ShiftInput{
		Requests:     req.Requests,
		Requirements: req.Requirements,
		Days:         req.Days,
	}

	if err := h.usecase.GenerateAndSave(input, req.StartDate); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "シフトを作成し保存しました"})
}

// List: GET /api/shift
func (h *ShiftHandler) List(c *gin.Context) {
	shifts, err := h.usecase.GetAllShifts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, shifts)
}

// Update: PUT /api/shift/:id
func (h *ShiftHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	
	type Request struct {
		Date string `json:"date"`
	}
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	var id uint
	fmt.Sscanf(idStr, "%d", &id) // ここで fmt を使っています

	if err := h.usecase.MoveShift(id, req.Date); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "更新しました"})
}

// Delete: DELETE /api/shift/:id
func (h *ShiftHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	var id uint
	fmt.Sscanf(idStr, "%d", &id) // ここでも fmt を使っています

	if err := h.usecase.DeleteShift(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "削除しました"})
}