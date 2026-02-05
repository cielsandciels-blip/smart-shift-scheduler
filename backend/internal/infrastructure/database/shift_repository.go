package database

import (
	"smart-shift-scheduler/internal/domain"

	"gorm.io/gorm"
)

type ShiftRepository struct {
	db *gorm.DB
}

func NewShiftRepository(db *gorm.DB) *ShiftRepository {
	return &ShiftRepository{db: db}
}

// Save: シフトを保存（本来は期間指定で消してから保存などが望ましいが、今回は追記/Upsert）
func (r *ShiftRepository) Save(shifts []domain.Shift) error {
    // 既存の同日・同スタッフのシフトがあれば削除してから保存するなどのロジックが必要だが
    // 簡易的にそのままCreate（重複回避はUniqueキー等で制御推奨）
	return r.db.Create(&shifts).Error
}

func (r *ShiftRepository) FindAll() ([]domain.Shift, error) {
	var shifts []domain.Shift
	if err := r.db.Find(&shifts).Error; err != nil {
		return nil, err
	}
	return shifts, nil
}

func (r *ShiftRepository) Update(shift *domain.Shift) error {
	// 指定したフィールドのみ更新（Dateなど）
	return r.db.Model(shift).Updates(shift).Error
}

// ★修正: id uint -> id int
func (r *ShiftRepository) Delete(id int) error {
	return r.db.Delete(&domain.Shift{}, id).Error
}

// ★修正: staffID int (これは元々intの可能性が高いが念のため)
func (r *ShiftRepository) DeleteByStaffID(staffID int) error {
	return r.db.Where("staff_id = ?", staffID).Delete(&domain.Shift{}).Error
}
func (r *ShiftRepository) DeleteRange(startDate string, endDate string) error {
	return r.db.Where("date >= ? AND date <= ?", startDate, endDate).Delete(&domain.Shift{}).Error
}