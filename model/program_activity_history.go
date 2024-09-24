package model

import (
	"log"

	"github.com/fsidiqs/aegis-backend/helper"
	"github.com/fsidiqs/aegis-backend/model/apperror"
	"github.com/google/uuid"
)

type ProgramActivityHistory struct {
	ID                uuid.UUID      `gorm:"primary_key; unique; type:uuid; column:id; not null; default:uuid_generate_v4()" json:"id"`
	EnrollmentID      uuid.UUID      `json:"enrollment_id"`
	ProgramActivityID uuid.UUID      `json:"program_activity_id"`
	ActivityResult    TActivityTempl `gorm:"column:activity_result" json:"activity_result"`

	DefaultColumns
}

func NewProgramActivityHistory(createdBy string) ProgramActivityHistory {
	model := ProgramActivityHistory{}
	model.NewDefaultCols(createdBy)
	return model
}

func NewProgramActivityHistoryArr(len int, createdBy string) []ProgramActivityHistory {
	history := make([]ProgramActivityHistory, len)
	for i := 0; i < len; i++ {
		history[i].NewDefaultCols(createdBy)
	}
	return history
}

func (ProgramActivityHistory) TableName() string {
	return "program_activity_history"
}

type ProgramActivityHistoryArr []ProgramActivityHistory

func (d *ProgramActivityHistoryArr) FindIdxByEnrollmentIDAndProgramActivityID(enrIDStr, progActIDStr string) (int, error) {
	enrID, err := uuid.Parse(enrIDStr)
	if err != nil {
		log.Printf("failed %v err:%v\n", helper.TraceCurrentFunc(), err)
		return -1, err
	}
	actID, err := uuid.Parse(progActIDStr)
	if err != nil {
		log.Printf("failed %v err:%v\n", helper.TraceCurrentFunc(), err)
		return -1, err
	}
	for i, v := range *d {
		if v.EnrollmentID == enrID && v.ProgramActivityID == actID {
			return i, nil
		}
	}
	log.Printf("failed %v err:id not found\n", helper.TraceCurrentFunc())

	return -1, apperror.NewResourceNotFound()
}

// REQUEST DATA

type ReqProgramActivityArr []ReqProgramActivityHistory

func (d *ReqProgramActivityArr) Validate() error {
	if len(*d) == 0 {
		return apperror.NewBadRequest("program activity is empty")
	}
	return nil
}

type ReqProgramActivityHistory struct {
	ProgramActivityID uuid.UUID `json:"program_activity_id"`
	ActivityResult    *string   `json:"activity_result"`
}
