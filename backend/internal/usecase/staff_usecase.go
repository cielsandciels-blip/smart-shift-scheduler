package usecase

import "smart-shift-scheduler/internal/domain"

// StaffRepository: データ保存のインターフェース
type StaffRepository interface {
	Save(staff *domain.Staff) error
	FindAll() ([]domain.Staff, error)
}

type StaffUsecase struct {
	repo StaffRepository
}

func NewStaffUsecase(repo StaffRepository) *StaffUsecase {
	return &StaffUsecase{repo: repo}
}

// CreateStaff: スタッフを登録する関数 (★これが足りなかった！)
func (u *StaffUsecase) CreateStaff(staff *domain.Staff) error {
	return u.repo.Save(staff)
}

// GetAllStaff: スタッフ一覧を取得する関数
func (u *StaffUsecase) GetAllStaff() ([]domain.Staff, error) {
	return u.repo.FindAll()
}