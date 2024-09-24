package model

import (
	"github.com/fsidiqs/aegis-backend/helper/datehelper"
	"github.com/google/uuid"
)

type JourneytodayContainer struct {
	Relax5 []ProgJourneyable
	JourneytodayEnrolled
	Meditation []ProgJourneyable
}

type JourneytodayEnrolled struct {
	MorningFits []ProgJourneyable
	Coachings   []ProgJourneyable
	NightFits   []ProgJourneyable
}

func (data JourneytodayContainer) CreateProgramJourneyArr() []ProgJourneyable {
	var progJourArr []ProgJourneyable

	if data.Relax5 != nil {
		progJourArr = append(progJourArr, data.Relax5...)
	}

	if data.MorningFits != nil {
		progJourArr = append(progJourArr, data.MorningFits...)
	}

	if data.Coachings != nil {
		progJourArr = append(progJourArr, data.Coachings...)
	}

	if data.NightFits != nil {
		progJourArr = append(progJourArr, data.NightFits...)
	}

	if data.Meditation != nil {
		progJourArr = append(progJourArr, data.Meditation...)
	}
	return progJourArr
}

type ProgJourneyable struct {
	*MProgram
	CurrentSequence        *int             `json:"current_sequence,omitempty"`
	CurrentSequenceSession *MProgramSession `json:"session_by_current_sequence"`
	DaySinceEnroll         *int             `json:"day_since_enrollment,omitempty"`
}

func MProgArrToProgJourneyable(src ProgramEnrollmentHistoryArr) []ProgJourneyable {
	progLen := len(src)
	mProgJourneyArr := make([]ProgJourneyable, progLen)
	for i := 0; i < progLen; i++ {
		// truncated to 00:00
		intDaySince := datehelper.GetDateSince(*src[i].CreatedAt)
		mProgJourneyArr[i].MProgram = src[i].Program
		mProgJourneyArr[i].CurrentSequence = (*int)(&src[i].CurrentSequence)
		mProgJourneyArr[i].DaySinceEnroll = &intDaySince
	}
	return mProgJourneyArr
}

type ProgramIDSequenceNum struct {
	ProgramID       uuid.UUID
	CurrentSequence int
}
