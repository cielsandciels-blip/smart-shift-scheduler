package database

import (
	"smart-shift-scheduler/internal/domain"
	"gorm.io/gorm"
)

type RequestRepository struct {
	db *gorm.DB
}

func NewRequestRepository(db *gorm.DB) *RequestRepository {
	return &RequestRepository{db: db}
}

// Save: 希望を保存する（同じ日の重複は簡易的に無視せず追加）
func (r *RequestRepository) Save(req *domain.ShiftRequest) error {
	return r.db.Create(req).Error
}

// FindAll: 全員の希望を取得
func (r *RequestRepository) FindAll() ([]domain.ShiftRequest, error) {
	var requests []domain.ShiftRequest
	if err := r.db.Find(&requests).Error; err != nil {
		return nil, err
	}
	return requests, nil
}

// Delete: 取り消し用
func (r *RequestRepository) Delete(id uint) error {
	return r.db.Delete(&domain.ShiftRequest{}, id).Error
}