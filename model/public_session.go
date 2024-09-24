package model

// REQUEST FROM CLIENT

type ReqPublicSession struct {
	DeviceID          string `json:"device_id" binding:"required"`
	NotificationToken string `json:"notification_token"`
}
