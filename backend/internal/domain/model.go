package domain

// Staff: スタッフ情報
type Staff struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Name       string `json:"name"`
	IsLeader   bool   `json:"is_leader"`
	HourlyWage int    `json:"hourly_wage"`
	Roles      string `json:"roles"` // "Kitchen,Leader"
}

// Shift: 確定したシフト
type Shift struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	StaffID   int    `json:"staff_id"`
	Date      string `json:"date"`
	ShiftType int    `json:"shift_type"`
}

// ShiftRequest: 希望休
type ShiftRequest struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	StaffID int    `json:"staff_id"`
	Date    string `json:"date"`
	Type    string `json:"type"`
}

// RoleConstraint: 役割ごとの必要人数ルール
type RoleConstraint struct {
	Role  string `json:"role"`
	Count int    `json:"count"`
}

// DailyRequirement: その日の必要人数設定
type DailyRequirement struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Date        string `json:"date" gorm:"unique"`
	MorningNeed int    `json:"morning_need"`
	EveningNeed int    `json:"evening_need"`
}

// ShiftInput: Pythonに渡すデータ
type ShiftInput struct {
	StaffList       []Staff            `json:"staff_list"`
	Requests        []ShiftRequest     `json:"requests"`
	RoleConstraints []RoleConstraint   `json:"role_constraints"`
	Requirements    []DailyRequirement `json:"requirements"`
	Days            int                `json:"days"`
	StartDate       string             `json:"start_date"` // ★これを追加しました！
}

// ShiftResult: 計算結果
type ShiftResult struct {
	Status   string          `json:"status"`
	Schedule map[int][]int   `json:"schedule"`
}