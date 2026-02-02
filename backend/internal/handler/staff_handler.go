package handler

import (
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

// Create: スタッフ登録API
func (h *StaffHandler) Create(c *gin.Context) {
	type CreateStaffRequest struct {
		Name       string `json:"name"`
		IsLeader   bool   `json:"is_leader"`
		HourlyWage int    `json:"hourly_wage"`
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
	}

	// ★修正: Usecaseの CreateStaff を呼ぶように統一
	if err := h.usecase.CreateStaff(staff); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, staff)
}

// List: スタッフ一覧API
func (h *StaffHandler) List(c *gin.Context) {
	staffList, err := h.usecase.GetAllStaff()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, staffList)
}