package model

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/google/uuid"
)

type MProgActType string

const (
	MPActTypeAskSingleChoice   MProgActType = "ASK_SINGLE_CHOICE"
	MPActTypeAskMultipleChoice MProgActType = "ASK_MULTIPLE_CHOICE"
	MPActTypeAskEssay          MProgActType = "ASK_ESSAY"
	MPActTypeAskTakePicture    MProgActType = "ASK_TAKE_PICTURE"
	MPActTypeAskInstruction    MProgActType = "ASK_INSTRUCTION"
	MPActTypeMediaVideo        MProgActType = "MEDIA_VIDEO"
	MPActTypeMediaAudio        MProgActType = "MEDIA_AUDIO"
	MPActTypeMediaIllustration MProgActType = "MEDIA_ILLUSTRATION"
	MPActTypeInfoText          MProgActType = "INFO_TEXT"
	MPActTypeInfoTakeway       MProgActType = "INFO_TAKEAWAY"
	MPActTypeInfoTextEmoji     MProgActType = "INFO_TEXT_EMOJI"
	MPActTypeAskLimitedAnswer  MProgActType = "ASK_LIMITED_ANSWER"
)

type MProgramActivity struct {
	ID uuid.UUID `gorm:"primary_key; unique; type:uuid; column:id; not null; default:uuid_generate_v4()" json:"id"`

	ActivityType     MProgActType   `json:"activity_type"`
	SequenceNumber   int            `json:"sequence_number"`
	ActivityTemplate TActivityTempl `json:"activity_template"`

	ProgramSectionID uuid.UUID `json:"program_section_id"`

	MediaCollectionID uuid.UUID       `json:"media_collection_id"`
	MediaCollections  *MMediaAssetArr `gorm:"foreignKey:MediaCollectionID;references:MediaCollectionID" json:"media_collections"`

	DefaultColumns
}

type MProgramActivityArr []MProgramActivity

type MediaVideoTemplate struct {
	Type         MProgActType                   `json:"type"`
	ThumbnailUrl string                         `json:"thumbnail_url"`
	AudioURL     string                         `json:"audio_url"`
	Collections  []MediaVideoActivityCollection `json:"collections"`
}
type MediaVideoActivityCollection struct {
	Res int    `json:"res"`
	URL string `json:"url"`
}
type TActivityTempl map[string]interface{}

func (j TActivityTempl) GormDataType() string {
	return "text"
}

func (j TActivityTempl) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}

	return json.Marshal(j)
}

func (j *TActivityTempl) Scan(value interface{}) error {
	if err := json.Unmarshal(value.([]byte), &j); err != nil {
		return err
	}
	return nil
}

type MediaVideoTempl struct {
	Type         MProgActType                   `json:"type"`
	ThumbnailURL string                         `json:"thumbnail_url"`
	AudioURL     string                         `json:"audio_url"`
	Collections  []MediaVideoActivityCollection `json:"collections"`
}

type MediaAudioTempl struct {
	Type     MProgActType `json:"type"`
	AudioURL string       `json:"audio_url"`
}
