package usecase

import "smart-shift-scheduler/internal/domain"

// StaffRepository: データ保存のインターフェース
type StaffRepository interface {
	Save(staff *domain.Staff) error
	FindAll() ([]domain.Staff, error)
	Delete(id uint) error // ★追加
}

type StaffUsecase struct {
	repo StaffRepository
}

func NewStaffUsecase(repo StaffRepository) *StaffUsecase {
	return &StaffUsecase{repo: repo}
}

func (u *StaffUsecase) CreateStaff(staff *domain.Staff) error {
	return u.repo.Save(staff)
}

func (u *StaffUsecase) GetAllStaff() ([]domain.Staff, error) {
	return u.repo.FindAll()
}

// DeleteStaff: スタッフ削除（★追加！）
func (u *StaffUsecase) DeleteStaff(id uint) error {
	return u.repo.Delete(id)
}