package domain

// Staff: スタッフ情報 (DBのテーブルになります)
type Staff struct {
	ID       uint   `json:"id" gorm:"primaryKey"` // uintにしてprimaryKeyを指定
	Name     string `json:"name"`
	IsLeader bool   `json:"is_leader"`
	HourlyWage int    `json:"hourly_wage"` 
}

// 他のstructはDB保存しないのでそのままでOK
type ShiftRequest struct {
	StaffID  int `json:"staff_id"`
	Day      int `json:"day"`
	Priority int `json:"priority"`
}

type DailyRequirement struct {
	Day             int `json:"day"`
	RequiredCount   int `json:"required_count"`
	RequiredLeaders int `json:"required_leaders"`
}

type ShiftInput struct {
	StaffList    []Staff            `json:"staff_list"`
	Requests     []ShiftRequest     `json:"requests"`
	Requirements []DailyRequirement `json:"requirements"`
	Days         int                `json:"days"`
}

type ShiftResult struct {
	Status   string         `json:"status"`
	Schedule map[int][]int  `json:"schedule"`
}

// Shift: 確定したシフト（DB保存用）
type Shift struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	StaffID   int    `json:"staff_id"`
	Date      string `json:"date"`       // "2026-02-01" のような文字列で保存
	ShiftType int    `json:"shift_type"` // 1:早番, 2:遅番
}