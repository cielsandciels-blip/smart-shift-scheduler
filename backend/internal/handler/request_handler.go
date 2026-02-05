package handler

import (
	"net/http"
	"smart-shift-scheduler/internal/domain"
	"smart-shift-scheduler/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RequestHandler struct {
	usecase *usecase.ShiftUsecase
}

func NewRequestHandler(u *usecase.ShiftUsecase) *RequestHandler {
	return &RequestHandler{usecase: u}
}

// Create: 希望休の登録
func (h *RequestHandler) Create(c *gin.Context) {
	var req domain.ShiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// ★修正: AddRequest -> CreateRequest
	if err := h.usecase.CreateRequest(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, req)
}

// List: 希望休の一覧
func (h *RequestHandler) List(c *gin.Context) {
	// ★修正: GetAllRequests -> ListRequests
	requests, err := h.usecase.ListRequests()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, requests)
}

// Delete: 希望休の削除
func (h *RequestHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	// ★修正: Atoiを使って int 型にする
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// ★修正: DeleteRequest (int型を渡す)
	if err := h.usecase.DeleteRequest(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}