package usecase

import (
	"fmt"
	"time"
	"smart-shift-scheduler/internal/domain"
	"smart-shift-scheduler/internal/infrastructure/engine"
)

// リポジトリのインターフェース
type ShiftRepository interface {
	SaveAll(shifts []domain.Shift) error
	FindAll() ([]domain.Shift, error)
	UpdateDate(id uint, newDate string) error
	Delete(id uint) error
}

// ★追加: リクエスト用インターフェース
type RequestRepository interface {
	Save(req *domain.ShiftRequest) error
	FindAll() ([]domain.ShiftRequest, error)
	Delete(id uint) error
}

type ShiftUsecase struct {
	engine      *engine.ShiftEngine
	staffRepo   domain.StaffRepository
	shiftRepo   ShiftRepository
	requestRepo RequestRepository // ★追加
}

// コンストラクタも更新
func NewShiftUsecase(engine *engine.ShiftEngine, staffRepo domain.StaffRepository, shiftRepo ShiftRepository, requestRepo RequestRepository) *ShiftUsecase {
	return &ShiftUsecase{
		engine:      engine,
		staffRepo:   staffRepo,
		shiftRepo:   shiftRepo,
		requestRepo: requestRepo,
	}
}

// GenerateAndSave: 計算して保存する
func (u *ShiftUsecase) GenerateAndSave(input domain.ShiftInput, startDateStr string) error {
	// 1. スタッフ取得
	staffList, err := u.staffRepo.FindAll()
	if err != nil {
		return err
	}
	input.StaffList = staffList

	// 2. ★ここでDBから希望休を取得してAIへの入力に追加する！
	requests, err := u.requestRepo.FindAll()
	if err != nil {
		return err
	}
	input.Requests = requests

	// 3. Pythonで計算
	// 3. Pythonで計算
	result, err := u.engine.Generate(input)
	if err != nil {
		return err
	}
    // ★全部大文字でもOKにする！
	if result.Status != "OPTIMAL" && result.Status != "FEASIBLE" && result.Status != "Optimal" && result.Status != "Feasible" {
		return fmt.Errorf("解が見つかりませんでした: %s", result.Status)
	}

	// 4. 結果を保存
	startDate, _ := time.Parse("2006-01-02", startDateStr)
	var shiftsToSave []domain.Shift

	for staffID, shiftTypes := range result.Schedule {
		for i, sType := range shiftTypes {
			if sType == 0 { continue } 
			date := startDate.AddDate(0, 0, i).Format("2006-01-02")
			shiftsToSave = append(shiftsToSave, domain.Shift{
				StaffID:   staffID,
				Date:      date,
				ShiftType: sType,
			})
		}
	}
	return u.shiftRepo.SaveAll(shiftsToSave)
}

// その他のメソッドはそのまま
func (u *ShiftUsecase) GetAllShifts() ([]domain.Shift, error) {
	return u.shiftRepo.FindAll()
}
func (u *ShiftUsecase) MoveShift(id uint, newDate string) error {
	return u.shiftRepo.UpdateDate(id, newDate)
}
func (u *ShiftUsecase) DeleteShift(id uint) error {
	return u.shiftRepo.Delete(id)
}

// ★希望休を登録・削除するメソッドも追加
func (u *ShiftUsecase) AddRequest(req *domain.ShiftRequest) error {
	return u.requestRepo.Save(req)
}
func (u *ShiftUsecase) GetAllRequests() ([]domain.ShiftRequest, error) {
	return u.requestRepo.FindAll()
}
func (u *ShiftUsecase) DeleteRequest(id uint) error {
	return u.requestRepo.Delete(id)
}