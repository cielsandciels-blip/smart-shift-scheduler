package domain

// StaffRepository: データベース操作のメニュー表
// 「保存(Save)」と「全取得(FindAll)」ができると定義
type StaffRepository interface {
	Save(staff *Staff) error
	FindAll() ([]Staff, error)
}