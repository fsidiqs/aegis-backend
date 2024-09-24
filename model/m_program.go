package model

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
)

type MProgramCat string
type MProgramSubCat string
type MProgramContentType string

const (
	MProgCatCareer          MProgramCat = "CAREER"
	MProgCatLove            MProgramCat = "LOVE"
	MProgCatFamily          MProgramCat = "FAMILY"
	MProgCatPersonalMastery MProgramCat = "PERSONAL_MASTERY"
	MProgCatMorningFitness  MProgramCat = "MORNING_FITNESS"
	MProgCatNightFitness    MProgramCat = "NIGHT_FITNESS"
	MProgCatRelax           MProgramCat = "RELAXATION"
	MProgCatMed             MProgramCat = "MEDITATION"
)

func IsCoaching(cat MProgramCat) bool {
	switch cat {
	case MProgCatFamily:
		return true
	case MProgCatLove:
		return true
	case MProgCatCareer:
		return true
	case MProgCatPersonalMastery:
		return true
	}
	return false
}

func IsFitness(cat MProgramCat) bool {
	switch cat {
	case MProgCatMorningFitness:
		return true
	case MProgCatNightFitness:
		return true
	}
	return false
}

const (
	MProgSubStarting   MProgramSubCat = "STARTING"
	MProgSubDeveloping MProgramSubCat = "DEVELOPING"
	MProgSubEmpowering MProgramSubCat = "EMPOWERING"
	MProgSubBasicEdu   MProgramSubCat = "BASIC_EDUCATION"
)

const (
	Video MProgramContentType = "VIDEO"
	Audio MProgramContentType = "AUDIO"
	Read  MProgramContentType = "READ"
)

type MProgram struct {
	ID          uuid.UUID           `gorm:"primary_key; unique; type:uuid; column:id; not null; default:uuid_generate_v4()" json:"id"`
	Name        string              `json:"name"`
	ContentType MProgramContentType `json:"content_type"`
	Category    MProgramCat         `json:"category"`
	SubCategory MProgramSubCat      `json:"sub_category"`
	Description string              `json:"description"`
	Objective   *PGTextArr          `json:"objective"`
	// in the actual db the object filename is 'thumbnail_url'.
	// since we fetch the actual object from cloud storage, we return the 'ThumbnailSource' (fetched from storage) as 'thumbnail_url' instead.
	ThumbnailUrl    string `json:"-"`
	ThumbnailSource string `gorm:"-" json:"thumbnail_url"`

	CoachID uuid.UUID `json:"coach_id"`
	Coach   *MCoach   `gorm:"foreignKey:CoachID" json:"coach"`

	// TODO! DELETE
	ModeratorID uuid.UUID `json:"-"`
	Moderator   *MCoach   `gorm:"foreignKey:ModeratorID" json:"-"`

	IntroVideoID uuid.UUID       `json:"intro_video_id"`
	IntroVideo   *MMediaAssetArr `gorm:"foreignKey:MediaCollectionID;references:IntroVideoID" json:"intro_video"`

	Rating           int                         `json:"rating"`
	Duration         int                         `json:"duration"`
	SubscriptionType TSubscription               `json:"subscription_type"`
	Status           TProgramStatus              `json:"program_status"`
	Enrollments      ProgramEnrollmentHistoryArr `gorm:"foreignKey:ProgramID" json:"enrollments_history,omitempty"`
	DefaultColumns
}

type SpecifiedMProgram struct {
	ID          uuid.UUID      `gorm:"column:mp.id; primary_key; unique; type:uuid; not null; default:uuid_generate_v4()" json:"id"`
	Name        string         `gorm:"column:mp.name" json:"name"`
	Category    MProgramCat    `gorm:"column:mp.category" json:"category"`
	SubCategory MProgramSubCat `gorm:"column:mp.sub_category" json:"sub_category"`
	Description string         `gorm:"column:mp.description" json:"description"`
	Objective   *PGTextArr     `gorm:"column:mp.objective" json:"objective"`
	// in the actual db the object filename is 'thumbnail_url'.
	// since we fetch the actual object from cloud storage, we return the 'ThumbnailSource' (fetched from storage) as 'thumbnail_url' instead.
	ThumbnailUrl    string `gorm:"column:mp.thumbnail_url" json:"-"`
	ThumbnailSource string `gorm:"-" json:"thumbnail_url"`

	CoachID uuid.UUID `gorm:"column:mp.coach_id" json:"coach_id"`
	Coach   *MCoach   `gorm:"foreignKey:CoachID" json:"coach"`

	IntroVideoID uuid.UUID       `gorm:"column:mp.intro_video_id" json:"intro_video_id"`
	IntroVideo   *MMediaAssetArr `gorm:"foreignKey:MediaCollectionID;references:IntroVideoID" json:"intro_video"`

	Rating           int                         `gorm:"column:mp.rating" json:"rating"`
	Duration         int                         `gorm:"column:mp.duration" json:"duration"`
	SubscriptionType TSubscription               `gorm:"column:mp.subscription_type" json:"subscription_type"`
	Status           TProgramStatus              `gorm:"column:mp.program_status" json:"program_status"`
	Enrollments      ProgramEnrollmentHistoryArr `gorm:"foreignKey:ProgramID" json:"enrollments_history,omitempty"`

	CreatedAt  *time.Time  `gorm:"column:mp.created_at" json:"created_at"`
	UpdatedAt  *time.Time  `gorm:"column:mp.updated_at" json:"updated_at"`
	CreatedBy  string      `gorm:"column:mp.created_by" json:"-"`
	UpdatedBy  string      `gorm:"column:mp.updated_by" json:"-"`
	RecordFlag TRecordFlag `gorm:"column:mp.record_flag" json:"record_flag"`
}

func (SpecifiedMProgram) TableName() string {
	return "m_programs"
}

func (d *SpecifiedMProgram) Value() (driver.Value, error) {
	valueString, err := json.Marshal(d)
	return string(valueString), err
}

func (j *SpecifiedMProgram) Scan(value interface{}) error {

	if err := json.Unmarshal(value.([]byte), &j); err != nil {
		return err
	}
	return nil
}

type PGTextArr []string

func (d *PGTextArr) Value() (driver.Value, error) {
	valueString, err := json.Marshal(d)
	return string(valueString), err
}

func (d *PGTextArr) Scan(value interface{}) error {
	// assert to string
	fromDBStr := value.(string)
	// split by ";"
	*d = strings.Split(fromDBStr, ";")

	return nil
}

const (
	Free TSubscription = "FREE"
	Paid TSubscription = "PAID"
)

type TSubscription string

type MProgArr []MProgram

func (mProgArr MProgArr) FilterMorningFitness() MProgArr {
	mProgLen := len(mProgArr)
	fitness := MProgArr{}

	for i := 0; i < mProgLen; i++ {
		if mProgArr[i].Category == MProgCatMorningFitness {
			fitness = append(fitness, mProgArr[i])
		}
	}
	return fitness
}

func (mProgArr MProgArr) FilterNightFitness() MProgArr {
	mProgLen := len(mProgArr)
	fitness := MProgArr{}

	for i := 0; i < mProgLen; i++ {
		if mProgArr[i].Category == MProgCatNightFitness {
			fitness = append(fitness, mProgArr[i])
		}
	}
	return fitness
}

func (mProgArr MProgArr) FilterCoaching() MProgArr {
	mProgLen := len(mProgArr)
	fitness := MProgArr{}

	for i := 0; i < mProgLen; i++ {
		if IsCoaching(mProgArr[i].Category) {
			fitness = append(fitness, mProgArr[i])
		}
	}
	return fitness
}

func MProgSliceToMProgArr(src []MProgram) MProgArr {
	progLen := len(src)
	mProgArr := make([]MProgram, progLen)
	for i := 0; i < progLen; i++ {
		mProgArr[i] = src[i]
	}
	return mProgArr
}
