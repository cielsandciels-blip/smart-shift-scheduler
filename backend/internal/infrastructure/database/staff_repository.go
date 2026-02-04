package database

import (
	"smart-shift-scheduler/internal/domain"
	"gorm.io/gorm"
)

type StaffRepository struct {
	db *gorm.DB
}

func NewStaffRepository(db *gorm.DB) *StaffRepository {
	return &StaffRepository{db: db}
}

// Save: 保存
func (r *StaffRepository) Save(staff *domain.Staff) error {
	return r.db.Create(staff).Error
}

// FindAll: 一覧取得
func (r *StaffRepository) FindAll() ([]domain.Staff, error) {
	var staffList []domain.Staff
	if err := r.db.Find(&staffList).Error; err != nil {
		return nil, err
	}
	return staffList, nil
}

// Delete: スタッフとその人のシフトを削除（★ここを修正！）
func (r *StaffRepository) Delete(id uint) error {
	// 1. まず、そのスタッフIDに紐づくシフトを全削除する
	// ("shifts" テーブルから staff_id = id のデータを消す)
	if err := r.db.Where("staff_id = ?", id).Delete(&domain.Shift{}).Error; err != nil {
		return err // シフト削除に失敗したらエラーを返す
	}

	// 2. シフトが消えたら、スタッフ本人を削除する
	return r.db.Delete(&domain.Staff{}, id).Error
}