package usecase

import (
	"fmt"
	"time"
	"smart-shift-scheduler/internal/domain"
	"smart-shift-scheduler/internal/infrastructure/engine"
)

// リポジトリのインターフェース定義
type ShiftRepository interface {
	SaveAll(shifts []domain.Shift) error
	FindAll() ([]domain.Shift, error)
	// ★ここが重要: UpdateDate と Delete を定義
	UpdateDate(id uint, newDate string) error
	Delete(id uint) error
}

type ShiftUsecase struct {
	engine    *engine.ShiftEngine
	staffRepo domain.StaffRepository
	shiftRepo ShiftRepository
}

func NewShiftUsecase(engine *engine.ShiftEngine, staffRepo domain.StaffRepository, shiftRepo ShiftRepository) *ShiftUsecase {
	return &ShiftUsecase{
		engine:    engine,
		staffRepo: staffRepo,
		shiftRepo: shiftRepo,
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

	// 2. Pythonで計算
	result, err := u.engine.Generate(input)
	if err != nil {
		return err
	}
	if result.Status != "Optimal" && result.Status != "Feasible" {
		return fmt.Errorf("解が見つかりませんでした: %s", result.Status)
	}

	// 3. 結果を変換
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

	// 4. 保存
	return u.shiftRepo.SaveAll(shiftsToSave)
}

// GetAllShifts: 全件取得
func (u *ShiftUsecase) GetAllShifts() ([]domain.Shift, error) {
	return u.shiftRepo.FindAll()
}

// MoveShift: シフトの日付を移動する (UpdateDateを使う)
func (u *ShiftUsecase) MoveShift(id uint, newDate string) error {
	return u.shiftRepo.UpdateDate(id, newDate)
}

// DeleteShift: シフトを削除する
func (u *ShiftUsecase) DeleteShift(id uint) error {
	return u.shiftRepo.Delete(id)
}