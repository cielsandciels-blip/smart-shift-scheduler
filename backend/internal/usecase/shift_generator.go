package usecase

import (
	"fmt"
	"smart-shift-scheduler/internal/domain"
	"smart-shift-scheduler/internal/infrastructure/engine"
	"time" // ★追加: 日付計算のために必要
)

// ... (Interface定義やStruct定義はそのまま) ...
// ShiftRepository, RequestRepository, RequirementRepository, ShiftUsecase など
// 変更がない部分は省略せずに全部書きます↓

type ShiftRepository interface {
	Save(shifts []domain.Shift) error
	FindAll() ([]domain.Shift, error)
	Update(shift *domain.Shift) error
	Delete(id int) error
	DeleteByStaffID(staffID int) error
	DeleteRange(startDate string, endDate string) error // 追加
}

type RequestRepository interface {
	Save(req *domain.ShiftRequest) error
	FindAll() ([]domain.ShiftRequest, error)
	Delete(id int) error
}

type RequirementRepository interface {
	Save(req *domain.DailyRequirement) error
	FindAll() ([]domain.DailyRequirement, error)
	Delete(id int) error // 追加
}

type ShiftUsecase struct {
	engine      *engine.ShiftEngine
	staffRepo   domain.StaffRepository
	shiftRepo   ShiftRepository
	requestRepo RequestRepository
	requireRepo RequirementRepository
}

func NewShiftUsecase(engine *engine.ShiftEngine, staffRepo domain.StaffRepository, shiftRepo ShiftRepository, requestRepo RequestRepository, requireRepo RequirementRepository) *ShiftUsecase {
	return &ShiftUsecase{
		engine:      engine,
		staffRepo:   staffRepo,
		shiftRepo:   shiftRepo,
		requestRepo: requestRepo,
		requireRepo: requireRepo,
	}
}

// GenerateAndSave: 計算して保存する
func (u *ShiftUsecase) GenerateAndSave(input domain.ShiftInput, startDateStr string) error {
	// 1. スタッフ一覧を取得
	staffList, err := u.staffRepo.FindAll()
	if err != nil {
		return err
	}
	input.StaffList = staffList

	// 2. 希望休を取得
	requests, err := u.requestRepo.FindAll()
	if err != nil {
		return err
	}
	input.Requests = requests

	// 3. 必要人数ルールを取得
	requirements, err := u.requireRepo.FindAll()
	if err == nil {
		input.Requirements = requirements
	}

	// 4. Pythonで計算
	result, err := u.engine.Generate(input)
	if err != nil {
		return err
	}
// Pythonのステータス判定
	if result.Status != "OPTIMAL" && result.Status != "FEASIBLE" && result.Status != "Optimal" && result.Status != "Feasible" {
		return fmt.Errorf("解が見つかりませんでした: %s", result.Status)
	}

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, startDateStr)
	if err != nil {
		return fmt.Errorf("日付形式エラー: %v", err)
	}

	// ★追加: 古いシフトを消す処理
	// 作成期間（デフォルト30日と仮定）の古いデータを削除
    // input.Days が 0 の場合もあるのでデフォルト30を入れておく
    days := input.Days
    if days == 0 { days = 30 }
    
	endDate := startDate.AddDate(0, 0, days-1).Format(layout)
	if err := u.shiftRepo.DeleteRange(startDateStr, endDate); err != nil {
		return fmt.Errorf("既存シフト削除失敗: %v", err)
	}

	// 5. 結果をDBに保存
	var shifts []domain.Shift
	for staffID, shiftTypes := range result.Schedule {
		for i, st := range shiftTypes {
			if st == 0 {
				continue
			}
			
			// ★修正: ちゃんとした日付を計算する
			// 開始日 + i日後 を計算して文字列に戻す
			dateStr := startDate.AddDate(0, 0, i).Format(layout)

			shifts = append(shifts, domain.Shift{
				StaffID:   staffID,
				Date:      dateStr, // "2026-02-02" のようになる
				ShiftType: st,
			})
		}
	}

	return u.shiftRepo.Save(shifts)
}

// ... (以下の ListShifts などは変更なし) ...
func (u *ShiftUsecase) ListShifts() ([]domain.Shift, error) {
	return u.shiftRepo.FindAll()
}
func (u *ShiftUsecase) UpdateShift(shift *domain.Shift) error {
	return u.shiftRepo.Update(shift)
}
func (u *ShiftUsecase) DeleteShift(id int) error {
	return u.shiftRepo.Delete(id)
}
func (u *ShiftUsecase) CreateRequest(req *domain.ShiftRequest) error {
	req.Type = "NG"
	return u.requestRepo.Save(req)
}
func (u *ShiftUsecase) ListRequests() ([]domain.ShiftRequest, error) {
	return u.requestRepo.FindAll()
}
func (u *ShiftUsecase) DeleteRequest(id int) error {
	return u.requestRepo.Delete(id)
}
func (u *ShiftUsecase) SaveRequirement(req *domain.DailyRequirement) error {
	return u.requireRepo.Save(req)
}
func (u *ShiftUsecase) GetRequirements() ([]domain.DailyRequirement, error) {
	return u.requireRepo.FindAll()
}
// ⭕️ Usecaseにはこれを貼る
func (u *ShiftUsecase) DeleteRequirement(id int) error {
	return u.requireRepo.Delete(id)
}
// ListStaff: スタッフ一覧を取得 (Handler用)
func (u *ShiftUsecase) ListStaff() ([]domain.Staff, error) {
	return u.staffRepo.FindAll()
}