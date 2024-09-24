package model

import "github.com/google/uuid"

type MMediaAsset struct {
	ID                uuid.UUID `gorm:"primary_key; unique; type:uuid; column:id; not null; default:uuid_generate_v4()" json:"id"`
	MediaCollectionID uuid.UUID `json:"media_collection_id"`

	MediaType MMediaAssetMediaType `json:"media_type"`
	URL       string               `json:"-"`
	AssetSrc  string               `gorm:"-" json:"asset_url"`

	ContentType TContentType `json:"content_type"`
	Resolution  int          `json:"resolution"`

	DefaultColumns
}

type MMediaAssetArr []MMediaAsset

func (d *MMediaAssetArr) GetAudioVersion() (string, bool) {
	for _, mediaAsset := range *d {
		if mediaAsset.MediaType == MediaTypeAudio {
			return mediaAsset.AssetSrc, true
		}
	}
	return "", false
}

func (d *MMediaAssetArr) GetVideoVersion() ([]MediaVideoActivityCollection, bool) {
	tempMediaVideo := []MediaVideoActivityCollection{}
	for _, mediaAsset := range *d {
		if mediaAsset.MediaType == MediaTypeVideo {
			tempMediaVideo = append(tempMediaVideo,
				MediaVideoActivityCollection{
					Res: mediaAsset.Resolution,
					URL: mediaAsset.AssetSrc,
				},
			)
		}
	}
	return tempMediaVideo, true
}

const (
	Program      TContentType = "PROGRAM"
	IntroProgram TContentType = "INTRO_PROGRAM"
	IntroMood    TContentType = "INTRO_MOOD"
	Promo        TContentType = "PROMO"
)

const (
	MediaTypeVideo MMediaAssetMediaType = "VIDEO"
	MediaTypeAudio MMediaAssetMediaType = "AUDIO"
)

type TContentType string
type MMediaAssetMediaType string
