package apperror

const (
	EmailNotFound Type = "EMAIL_NOT_FOUND"
	OTPNotFound   Type = "OTP_NOT_FOUND"
	OTPExpired    Type = "OTP_EXPIRED"
)

const (
	MsgEmailNotFound string = "email not found"
	MsgOtpNotFound   string = "otp not found"
	OTPExpire        string = "otp has expired"
)
