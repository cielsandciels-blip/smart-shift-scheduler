package database

import (
	"smart-shift-scheduler/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RequirementRepository struct {
	db *gorm.DB
}

func NewRequirementRepository(db *gorm.DB) *RequirementRepository {
	return &RequirementRepository{db: db}
}

// Save: 設定を保存（同じ日のデータがあれば上書きする "Upsert" 処理）
func (r *RequirementRepository) Save(req *domain.DailyRequirement) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "date"}}, // Dateが同じなら
		DoUpdates: clause.AssignmentColumns([]string{"morning_need", "evening_need"}), // 人数だけ更新
	}).Create(req).Error
}

// FindAll: 全ての設定を取得
func (r *RequirementRepository) FindAll() ([]domain.DailyRequirement, error) {
	var reqs []domain.DailyRequirement
	if err := r.db.Find(&reqs).Error; err != nil {
		return nil, err
	}
	return reqs, nil
}

// ★追加: IDで削除
func (r *RequirementRepository) Delete(id int) error {
	return r.db.Delete(&domain.DailyRequirement{}, id).Error
}