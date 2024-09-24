package appresponse

const (
	HdlMsgFreeUserMaximumEnroll    ErrMsgResp = "reached maximum enrollment"
	HdlMsgNotEnrolled              ErrMsgResp = "user not enrolled"
	HdlMsgAlrEnroll                ErrMsgResp = "already enrolled"
	HdlMsgUserMaxEnrReached        ErrMsgResp = "user has reached maximum enrollment"
	HdlMsgEnrNotExist              ErrMsgResp = "enrollment not exist"
	ENROLL_FREE_ENROLL_TO_PAID_MSG ErrMsgResp = "free user cannot enroll to paid programs"

	HdlMsgSuccessEnroll   SuccessMsgResp = "program enrolled successfully"
	HdlMsgSuccessUnenroll SuccessMsgResp = "program unenrolled successfully"
)

const (
	THdlNotEnrolled                 ErrResponseType = "NOT_ENROLLED"
	THdlAlrEnroll                   ErrResponseType = "ALREADY_ENROLLED"
	THdlMaxEnrReached               ErrResponseType = "MAX_ENROLL_REACHED"
	THdlEnrNotExist                 ErrResponseType = "ENROLLMENT_NOT_EXIST"
	THdlUnenrollFailed              ErrResponseType = "UNENROLL_FAILED"
	ENROLL_FREE_ENROLL_TO_PAID_TYPE ErrResponseType = "FREE_ENROLL_TO_PAID"

	THdlEnrSuccess       SuccessResponseType = "ENROLLED"
	THdlUnenrollSucceess SuccessResponseType = "UNENROLLED"
)
