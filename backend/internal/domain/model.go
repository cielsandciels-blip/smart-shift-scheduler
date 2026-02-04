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

// ★これが抜けていました！★
// RoleConstraint: 役割ごとの必要人数ルール
type RoleConstraint struct {
	Role  string `json:"role"`
	Count int    `json:"count"`
}

// ShiftInput: Pythonに渡すデータ
type ShiftInput struct {
	StaffList       []Staff            `json:"staff_list"`
	Requests        []ShiftRequest     `json:"requests"`
	RoleConstraints []RoleConstraint   `json:"role_constraints"` // ★ここで使っています
	Requirements    []DailyRequirement `json:"requirements"`
	Days            int                `json:"days"`
}

// DailyRequirement: その日に必要な人数 (今は使っていないがPython互換のため残す)
type DailyRequirement struct {
	Date        string `json:"date"`
	MorningNeed int    `json:"morning_need"`
	EveningNeed int    `json:"evening_need"`
}

// ShiftResult: 計算結果
type ShiftResult struct {
	Status   string          `json:"status"`
	Schedule map[int][]int   `json:"schedule"`
}