package handler

import (
	"fmt" // ★fmtを忘れずに！
	"net/http"
	"smart-shift-scheduler/internal/domain"
	"smart-shift-scheduler/internal/usecase"

	"github.com/gin-gonic/gin"
)

type StaffHandler struct {
	usecase *usecase.StaffUsecase
}

func NewStaffHandler(u *usecase.StaffUsecase) *StaffHandler {
	return &StaffHandler{usecase: u}
}

// Create: スタッフ登録
func (h *StaffHandler) Create(c *gin.Context) {
	type CreateStaffRequest struct {
		Name       string `json:"name"`
		IsLeader   bool   `json:"is_leader"`
		HourlyWage int    `json:"hourly_wage"`
		Roles      string `json:"roles"` // ★受け皿を追加
	}

	var req CreateStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	staff := &domain.Staff{
		Name:       req.Name,
		IsLeader:   req.IsLeader,
		HourlyWage: req.HourlyWage,
		Roles:      req.Roles, // ★ここも追加
	}

	if err := h.usecase.CreateStaff(staff); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, staff)
}

func (h *StaffHandler) List(c *gin.Context) {
	staffList, err := h.usecase.GetAllStaff()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, staffList)
}

// Delete: スタッフ削除API（★追加！）
func (h *StaffHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	var id uint
	fmt.Sscanf(idStr, "%d", &id)

	if err := h.usecase.DeleteStaff(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "削除しました"})
}