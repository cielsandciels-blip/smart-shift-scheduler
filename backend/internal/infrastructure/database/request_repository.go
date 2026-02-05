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

func (r *RequestRepository) Save(req *domain.ShiftRequest) error {
	return r.db.Create(req).Error
}

func (r *RequestRepository) FindAll() ([]domain.ShiftRequest, error) {
	var reqs []domain.ShiftRequest
	if err := r.db.Find(&reqs).Error; err != nil {
		return nil, err
	}
	return reqs, nil
}

// ★修正: id uint -> id int
func (r *RequestRepository) Delete(id int) error {
	return r.db.Delete(&domain.ShiftRequest{}, id).Error
}