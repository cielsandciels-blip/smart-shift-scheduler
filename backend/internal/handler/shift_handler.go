package handler

import (
	"encoding/csv"
	"net/http"
	"smart-shift-scheduler/internal/domain"
	"smart-shift-scheduler/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ShiftHandler struct {
	usecase *usecase.ShiftUsecase
}

func NewShiftHandler(u *usecase.ShiftUsecase) *ShiftHandler {
	return &ShiftHandler{usecase: u}
}

// Generate: シフト生成
func (h *ShiftHandler) Generate(c *gin.Context) {
	var input domain.ShiftInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startDate := input.StartDate
	if startDate == "" {
		// input.StartDateが無い場合のフォールバック（JSON構造によってはこっち）
		// 今回はJSONにstart_dateが含まれている前提
		startDate = "2026-02-01" 
	}

	if err := h.usecase.GenerateAndSave(input, startDate); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "シフトを作成・保存しました"})
}

// List: シフト一覧
func (h *ShiftHandler) List(c *gin.Context) {
	// ★修正: GetAllShifts -> ListShifts
	shifts, err := h.usecase.ListShifts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, shifts)
}

// Update: シフト修正（移動など）
func (h *ShiftHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// 日付の更新
	date, ok := req["date"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Date is required"})
		return
	}

	// 更新用オブジェクト作成
	// IDは uint だが DB定義に合わせる。Usecase側で int を受ける形にした方が早いが、
	// ここでは構造体を作る
	shift := &domain.Shift{
		ID:   uint(id), 
		Date: date,
		// ShiftTypeなどは必要ならフロントから送るが、
		// 現在のドラッグ&ドロップ実装だと日付変更がメイン
	}

	// ★修正: MoveShift -> UpdateShift
	// UpdateShiftは全フィールド更新の可能性があるので、
	// 本来は「既存データを取得して書き換える」のが安全だが、
	// GORMのUpdatesを使っていれば指定フィールドのみ更新される
	if err := h.usecase.UpdateShift(shift); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Updated"})
}

// Delete: シフト削除
func (h *ShiftHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	// ★修正: intに変換
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// ★修正: DeleteShift (intを渡す)
	if err := h.usecase.DeleteShift(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

// SaveRequirement: 必要人数保存 (先ほど追加したもの)
func (h *ShiftHandler) SaveRequirement(c *gin.Context) {
	var req domain.DailyRequirement
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}
	if err := h.usecase.SaveRequirement(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, req)
}

// ListRequirements: 必要人数一覧 (先ほど追加したもの)
func (h *ShiftHandler) ListRequirements(c *gin.Context) {
	list, err := h.usecase.GetRequirements()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}
// ⭕️ Handlerにはこれを貼る
func (h *ShiftHandler) DeleteRequirement(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := h.usecase.DeleteRequirement(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}



// Export: CSV出力
func (h *ShiftHandler) Export(c *gin.Context) {
	// 1. 全シフト取得
	shifts, err := h.usecase.ListShifts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 2. スタッフ名を取得
	// ★修正: StaffRepoを直接触らず、専用のメソッドを使うようにしました
	staffList, err := h.usecase.ListStaff()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// IDと名前の対照表を作る
	staffMap := make(map[int]string)
	for _, s := range staffList {
		staffMap[int(s.ID)] = s.Name
	}

	// 3. CSV作成
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment;filename=shift.csv")
	c.Writer.Write([]byte{0xEF, 0xBB, 0xBF})
	// Shift-JISにするなら変換が必要ですが、一旦UTF-8で出力します
	writer := csv.NewWriter(c.Writer)

	
	
	// ヘッダー書き込み
	writer.Write([]string{"日付", "スタッフ名", "シフト種別", "時間"})

	for _, s := range shifts {
		name := staffMap[s.StaffID]
		typeStr := ""
		timeStr := ""
		
		switch s.ShiftType {
		case 1:
			typeStr = "早番"
			timeStr = "09:00-18:00"
		case 2:
			typeStr = "遅番"
			timeStr = "18:00-23:00"
		}

		writer.Write([]string{
			s.Date,
			name,
			typeStr,
			timeStr,
		})
	}
	writer.Flush()
}