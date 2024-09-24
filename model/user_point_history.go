package model

import (
	"github.com/google/uuid"
)

const (
	Plus  FlowType = 1
	Minus FlowType = -1
)

type FlowType int

type UserPointHistory struct {
	ID             uuid.UUID `gorm:"primary_key; unique; type:uuid; column:id; not null; default:uuid_generate_v4()" json:"id"`
	UserID         uuid.UUID `json:"user_id"`
	Flow           FlowType  `json:"flow"`
	Amount         int       `json:"amount"`
	CurrentBalance int       `json:"current_balance"`
	HistoryType    string    `json:"history_type"`
	Notes          string    `json:"note"`
	DefaultColumns
}

func (UserPointHistory) TableName() string {
	return "user_point_history"
}
