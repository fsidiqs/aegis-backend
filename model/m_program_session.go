package model

import (
	"encoding/json"

	"github.com/fsidiqs/aegis-backend/model/apperror"
	"github.com/google/uuid"
)

const (
	Unlocked              SessionLock = "UNLOCKED" // unlocked
	LockedNeedEnroll      SessionLock = "LOCKED_NEED_ENROLL"
	UnlockedWrongSequence SessionLock = "UNLOCKED_WRONG_SEQUENCE" // locked wrong seq
	LockedNeedSubscribe   SessionLock = "LOCKED_NEED_SUBSCRIBE"   // free account limitation
	LockedDailyLimit      SessionLock = "LOCKED_DAILY_LIMIT"      // open next day
)

type SessionLock string

type MProgramSessionWithLock struct {
	MProgramSession
	SessionLock SessionLock `json:"eligible_status"`
}
type MProgramSession struct {
	ID              uuid.UUID           `gorm:"primary_key; unique; type:uuid; column:id; not null; default:uuid_generate_v4()" json:"id"`
	SequenceNumber  int                 `json:"sequence_number"`
	Title           string              `json:"title"`
	Description     string              `json:"description"`
	Duration        int                 `json:"duration"`
	ProgramSections *MProgramSectionArr `gorm:"foreignKey:ProgramSessionID" json:"program_sections,omitempty"`
	ProgramID       uuid.UUID           `json:"program_id"`

	DefaultColumns
}

// API REQUEST

type ReqProgramSessionActivityResults struct {
	EnrollmentID    uuid.UUID             `json:"enrollment_id"`
	ActivityResults ReqProgramActivityArr `json:"activity_results"`
}

func (data ReqProgramSessionActivityResults) ToModelWithNullActivity(progActIds []uuid.UUID, createdBy string) ([]ProgramActivityHistory, error) {
	activityLen := len(progActIds)
	paHistory := NewProgramActivityHistoryArr(activityLen, createdBy)

	for idxActivity := 0; idxActivity < activityLen; idxActivity++ {
		paHistory[idxActivity].ActivityResult = nil
		paHistory[idxActivity].EnrollmentID = data.EnrollmentID
		paHistory[idxActivity].ProgramActivityID = progActIds[idxActivity]
	}

	return paHistory, nil
}

func (data ReqProgramSessionActivityResults) ToModel(createdBy string) ([]ProgramActivityHistory, error) {
	err := data.ActivityResults.Validate()
	if err != nil {
		return nil, err
	}

	activityLen := len(data.ActivityResults)
	paHistory := NewProgramActivityHistoryArr(activityLen, createdBy)

	for idxActivity := 0; idxActivity < activityLen; idxActivity++ {
		var parsedActRes map[string]interface{}
		if data.ActivityResults[idxActivity].ActivityResult != nil {

			err := json.Unmarshal([]byte(*data.ActivityResults[idxActivity].ActivityResult), &parsedActRes)
			if err != nil {
				return nil, apperror.NewBadRequest("bad request")
			}
			paHistory[idxActivity].ActivityResult = parsedActRes
		} else {
			paHistory[idxActivity].ActivityResult = nil
		}
		paHistory[idxActivity].EnrollmentID = data.EnrollmentID
		paHistory[idxActivity].ProgramActivityID = data.ActivityResults[idxActivity].ProgramActivityID
	}
	return paHistory, nil
}
