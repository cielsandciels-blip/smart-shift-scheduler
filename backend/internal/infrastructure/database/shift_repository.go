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

func (r *ShiftRepository) SaveAll(shifts []domain.Shift) error {
	// 簡易実装: 全消しして保存
	r.db.Exec("DELETE FROM shifts")
	return r.db.Create(&shifts).Error
}

func (r *ShiftRepository) FindAll() ([]domain.Shift, error) {
	var shifts []domain.Shift
	if err := r.db.Find(&shifts).Error; err != nil {
		return nil, err
	}
	return shifts, nil
}

// Update: シフトの内容（日付やタイプ）を更新する
func (r *ShiftRepository) Update(shift *domain.Shift) error {
	return r.db.Save(shift).Error
}

// Delete: シフトを削除する
func (r *ShiftRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Shift{}, id).Error
}

// UpdateDate: 指定したIDの日付だけを変更する
func (r *ShiftRepository) UpdateDate(id uint, newDate string) error {
	return r.db.Model(&domain.Shift{}).Where("id = ?", id).Update("date", newDate).Error
}