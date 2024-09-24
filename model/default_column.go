package model

import "time"

// swagger:ignore
type DefaultColumns struct {
	CreatedAt  *time.Time  `json:"created_at"`
	UpdatedAt  *time.Time  `json:"updated_at"`
	CreatedBy  string      `json:"-"`
	UpdatedBy  string      `json:"-"`
	RecordFlag TRecordFlag `json:"record_flag"`
}

func NewDefaultColumn() DefaultColumns {
	return DefaultColumns{
		RecordFlag: RecActive,
	}
}

func NewLogDefaultCol(createdBy string) DefaultColumns {
	now := time.Now()
	return DefaultColumns{
		CreatedAt:  &now,
		UpdatedAt:  &now,
		RecordFlag: RecActive,
		CreatedBy:  createdBy,
	}
}

func (d *DefaultColumns) NewDefaultCols(createdBy string) {
	now := time.Now()
	d.CreatedAt = &now
	d.RecordFlag = RecActive
	d.CreatedBy = createdBy
}

func (d *DefaultColumns) UpdUpdatedBy(cb string) {
	d.CreatedBy = cb
}

const (
	RecActive  TRecordFlag = "ACTIVE"
	RecDeleted TRecordFlag = "DELETED"
)

const (
	Published TProgramStatus = "PUBLISHED"
)

type TRecordFlag string
type TProgramStatus string

// DefaultUpdateColumns
type DefaultUpdateColumns struct {
	UpdatedAt  time.Time   `json:"updated_at"`
	UpdatedBy  string      `json:"-"`
	RecordFlag TRecordFlag `json:"record_flag"`
}

func NewDefaultUpdateCol(updatedBy string) DefaultUpdateColumns {
	return DefaultUpdateColumns{
		UpdatedAt: time.Now(),
		UpdatedBy: updatedBy,
	}
}

func NewDefaultUpdateColWithCreatedBy(createdBy string) DefaultUpdateColumns {
	return DefaultUpdateColumns{
		RecordFlag: RecActive,
	}
}
