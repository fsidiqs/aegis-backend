package appresponse

const (
	HdlMsgNoUser                ErrMsgResp = "user not found"
	HdlMsgOTPExpire             ErrMsgResp = "otp has expired"
	MsgOTPNotFound              ErrMsgResp = "otp not found"
	HdlMsgInvalidToken          ErrMsgResp = "provided token is invalid"
	HdlMsgNoAuthHeader          ErrMsgResp = "no authentication header provided"
	HdlMsgUserUnverified        ErrMsgResp = "user is unverified"
	HdlMsgUserNotActive         ErrMsgResp = "user is not active"
	HdlMsgUnauthorized          ErrMsgResp = "unauthorized"
	HdlMsgInvalidEmailPassCombo ErrMsgResp = "invalid email and password combination"

	HdlMsgUserExist            ErrMsgResp = "user already exist"
	HdlMsgFailedExtractUserCtx ErrMsgResp = "failed to extract user"

	MsgEmailNotFound ErrMsgResp = "email not found"
)

const (
	TEmailNotFound   ErrResponseType = "EMAIL_NOT_FOUND"
	TOTPNotFound     ErrResponseType = "OTP_NOT_FOUND"
	TOTPExpired      ErrResponseType = "OTP_EXPIRED"
	THandlerNoUser   ErrResponseType = "USER_NOT_FOUND"
	THdlNoAuthHeader ErrResponseType = "NO_AUTHENTICATION_HEADER"
	THdlUnverified   ErrResponseType = "UNVERIFIED"
	THdlUnauthorized ErrResponseType = "UNAUTHORIZED"
	THdlUserExist    ErrResponseType = "USER_ALREADY_EXIST"
)
