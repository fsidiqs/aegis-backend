package model

import "github.com/google/uuid"

type MoodTracker struct {
	ID     uuid.UUID `gorm:"primary_key; unique; type:uuid; column:id; not null; default:uuid_generate_v4()" json:"id"`
	UserID uuid.UUID `json:"user_id"`
	Mood   string    `json:"mood"`
	Notes  string    `json:"notes"`

	DefaultColumns
}

func NewMoodTracker(createdBy string) MoodTracker {
	model := MoodTracker{}
	model.NewDefaultCols(createdBy)
	return model
}

// REQUEST DATA

type QueryReqMoodTracker struct {
	QueryReqDateFilter
	QueryReqPagination
}

type ReqMoodTrackerSubmit struct {
	Mood  string `json:"mood"`
	Notes string `json:"notes"`
}

func (r ReqMoodTrackerSubmit) CreateModel(createdBy string) (*MoodTracker, error) {
	model := NewMoodTracker(createdBy)
	model.Mood = r.Mood
	model.Notes = r.Notes
	return &model, nil
}
