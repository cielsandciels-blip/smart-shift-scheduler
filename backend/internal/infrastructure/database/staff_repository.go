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

// Save: スタッフを1人保存する
func (r *StaffRepository) Save(staff *domain.Staff) error {
	return r.db.Create(staff).Error
}

// FindAll: 全員のリストを取得する
func (r *StaffRepository) FindAll() ([]domain.Staff, error) {
	var staffList []domain.Staff
	if err := r.db.Find(&staffList).Error; err != nil {
		return nil, err
	}
	return staffList, nil
}