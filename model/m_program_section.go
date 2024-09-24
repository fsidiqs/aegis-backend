package model

import "github.com/google/uuid"

type MProgSection string

const (
	MPSectionSelfTrivia   MProgSection = "SELF_TRIVIA"
	MPSectionSelfPractice MProgSection = "SELF_PRACTICE"
	MPSectionNewMindset   MProgSection = "NEW_MINDSET"
)

type MProgramSection struct {
	ID                uuid.UUID            `gorm:"primary_key; unique; type:uuid; column:id; not null; default:uuid_generate_v4()" json:"id"`
	Title             string               `json:"title"`
	Description       string               `json:"description"`
	SectionType       MProgSection         `json:"section_type"`
	SequenceNumber    int                  `json:"sequence_number"`
	ProgramSessionID  uuid.UUID            `json:"program_session_id"`
	ProgramActivities *MProgramActivityArr `gorm:"foreignKey:ProgramSectionID" json:"program_activities,omitempty"`

	DefaultColumns
}

type MProgramSectionArr []MProgramSection
