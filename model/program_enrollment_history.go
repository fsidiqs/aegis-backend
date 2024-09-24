package model

import (
	"time"

	"github.com/fsidiqs/aegis-backend/model/apperror"
	"github.com/google/uuid"
)

const (
	EnrInProgress TENROLLSTATUS = "IN_PROGRESS"
	EnrFinish     TENROLLSTATUS = "FINISH"
	EnrCancel     TENROLLSTATUS = "CANCEL"
)

type ProgramEnrollmentHistory struct {
	ID uuid.UUID `gorm:"primary_key; unique; type:uuid; column:id; not null; default:uuid_generate_v4()" json:"id"`

	ProgramID uuid.UUID `json:"program_id"`
	Program   *MProgram `gorm:"foreignKey:ProgramID" json:"program,omitempty"`

	UserID uuid.UUID `json:"user_id"`
	User   *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`

	Status TENROLLSTATUS `json:"status"`

	CurrentSequence int              `json:"current_sequence"`
	ProgramSession  *MProgramSession `gorm:"-" json:"program_session"`

	Rating int    `json:"rating"`
	Review string `json:"review"`

	DefaultColumns
}

func (EnrollmentReview) TableName() string {
	return "program_enrollment_history"
}

type ProgramEnrollmentHistoryUpdate struct {
	Status TENROLLSTATUS `json:"status"`

	DefaultColumns
}

func (ProgramEnrollmentHistoryUpdate) TableName() string {
	return "program_enrollment_history"
}

type ProgramEnrollmentHistoryStatus struct {
	ID uuid.UUID `gorm:"primary_key; unique; type:uuid; column:id; not null; default:uuid_generate_v4()" json:"id"`

	ProgramID uuid.UUID `json:"program_id"`

	Status TENROLLSTATUS `json:"status"`

	CurrentSequence int `json:"current_sequence"`
}

func (ProgramEnrollmentHistoryStatus) TableName() string {
	return "program_enrollment_history"
}

type SpecifiedProgramEnrollmentHistoryWithProgram struct {
	ID uuid.UUID `gorm:"column:peh.id; primary_key; unique; type:uuid; not null; default:uuid_generate_v4()" json:"id"`

	ProgramID uuid.UUID          `gorm:"column:peh.program_id" json:"program_id"`
	Program   *SpecifiedMProgram `json:"program"`
	UserID    uuid.UUID          `gorm:"column:peh.user_id" json:"user_id"`

	Status          TENROLLSTATUS `gorm:"column:peh.status" json:"status"`
	CurrentSequence int           `gorm:"column:peh.current_sequence" json:"current_sequence"`
	Rating          int           `gorm:"column:peh.rating" json:"rating"`
	Review          string        `gorm:"column:peh.review" json:"review"`
	CreatedAt       *time.Time    `gorm:"column:peh.created_at" json:"created_at"`
	UpdatedAt       *time.Time    `gorm:"column:peh.updated_at" json:"updated_at"`
	CreatedBy       string        `gorm:"column:peh.created_by" json:"-"`
	UpdatedBy       string        `gorm:"column:peh.updated_by" json:"-"`
	RecordFlag      TRecordFlag   `gorm:"column:peh.record_flag" json:"record_flag"`
}

func (SpecifiedProgramEnrollmentHistoryWithProgram) TableName() string {
	return "program_enrollment_history"
}

type TENROLLSTATUS string

// custom table naming according to db
func (ProgramEnrollmentHistory) TableName() string {
	return "program_enrollment_history"
}

type EnrollmentReview struct {
	ID        uuid.UUID `gorm:"primary_key; unique; type:uuid; column:id; not null; default:uuid_generate_v4()" json:"id"`
	ProgramID uuid.UUID `json:"program_id"`
	UserID    uuid.UUID `json:"user_id"`
	Rating    int       `json:"rating"`
	Review    string    `json:"review"`

	DefaultColumns
}

type ProgramEnrollmentHistoryArr []ProgramEnrollmentHistory

func (ProgramEnrollmentHistoryArr) TableName() string {
	return "program_enrollment_history"
}

type SpecifiedProgramEnrollmentHistoryWithProgramArr []SpecifiedProgramEnrollmentHistoryWithProgram

func (srcEnrArr ProgramEnrollmentHistoryArr) FilterRelax5() ProgramEnrollmentHistoryArr {
	enr := len(srcEnrArr)
	resultRelax5 := ProgramEnrollmentHistoryArr{}
	for i := 0; i < enr; i++ {
		// TODO add program validation whether nil or not,
		if srcEnrArr[i].Program.Category == MProgCatRelax {
			resultRelax5 = append(resultRelax5, srcEnrArr[i])
		}
	}
	return resultRelax5
}

func (srcEnrArr ProgramEnrollmentHistoryArr) FilterMeditation() ProgramEnrollmentHistoryArr {
	enr := len(srcEnrArr)
	resultMeditations := ProgramEnrollmentHistoryArr{}
	for i := 0; i < enr; i++ {
		// TODO add program validation whether nil or not,
		if srcEnrArr[i].Program.Category == MProgCatMed {
			resultMeditations = append(resultMeditations, srcEnrArr[i])
		}
	}
	return resultMeditations
}

func (srcEnrArr ProgramEnrollmentHistoryArr) FilterMorningFitness() ProgramEnrollmentHistoryArr {
	enr := len(srcEnrArr)
	resultEnrArr := ProgramEnrollmentHistoryArr{}

	for i := 0; i < enr; i++ {
		// TODO add program validation whether nil or not,
		if srcEnrArr[i].Program.Category == MProgCatMorningFitness {
			resultEnrArr = append(resultEnrArr, srcEnrArr[i])
		}
	}
	return resultEnrArr
}

func (srcEnrArr ProgramEnrollmentHistoryArr) FilterNightFitness() ProgramEnrollmentHistoryArr {
	enr := len(srcEnrArr)
	resultEnrArr := ProgramEnrollmentHistoryArr{}

	for i := 0; i < enr; i++ {
		// TODO add program validation whether nil or not,
		if srcEnrArr[i].Program.Category == MProgCatNightFitness {
			resultEnrArr = append(resultEnrArr, srcEnrArr[i])
		}
	}
	return resultEnrArr
}

func (srcEnrArr ProgramEnrollmentHistoryArr) FilterCoaching() ProgramEnrollmentHistoryArr {
	enr := len(srcEnrArr)
	resultEnrArr := ProgramEnrollmentHistoryArr{}

	for i := 0; i < enr; i++ {
		// TODO add program validation whether nil or not,
		if IsCoaching(srcEnrArr[i].Program.Category) {
			resultEnrArr = append(resultEnrArr, srcEnrArr[i])
		}
	}
	return resultEnrArr
}

func (srcEnrArr ProgramEnrollmentHistoryArr) CountFitness() (int, error) {
	enr := len(srcEnrArr)
	count := 0

	for i := 0; i < enr; i++ {
		if srcEnrArr[i].Program == nil {
			return 0, apperror.NewResourceNotFound()
		}
		if IsFitness(srcEnrArr[i].Program.Category) {
			count++
		}
	}
	return count, nil
}

func (srcEnrArr ProgramEnrollmentHistoryArr) CountCoaching() (int, error) {
	enr := len(srcEnrArr)
	count := 0

	for i := 0; i < enr; i++ {
		if srcEnrArr[i].Program == nil {
			return 0, apperror.NewResourceNotFound()
		}
		if IsCoaching(srcEnrArr[i].Program.Category) {
			count++
		}
	}
	return count, nil
}

func (srcEnrArr ProgramEnrollmentHistoryArr) IsExists(programID string) bool {
	for _, v := range srcEnrArr {
		if v.ProgramID.String() == programID {
			return true
		}
	}

	return false
}

// REQUEST HANDLER ------------------------------

type ReqProgramEnrollmentHistory struct {
	ProgramId uuid.UUID `json:"program_id" binding:"required"`
}

type ReqPEHUpdateCurrentSeq struct {
	CurrentSequence string `json:"current_sequence" binding:"required"`
}

type QueryReqEnrHist struct {
	Status string `form:"status"`
}

func (data QueryReqEnrHist) GetEnrStatus() (TENROLLSTATUS, error) {
	switch data.Status {
	case string(EnrInProgress):
		return EnrInProgress, nil
	case string(EnrFinish):
		return EnrFinish, nil
	}

	return "", &apperror.Error{Type: apperror.BadRequest, Message: apperror.BadRequestMessage}
}
