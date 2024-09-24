package appresponse

const (
	HdlMsgSuccess SuccessMsgResp = "success"
)

const (
	THdlSuccess SuccessResponseType = "SUCCESS"
)

const (
	HdlMsgErrorCommitTrx       ErrMsgResp = "error commiting transaction"
	HdlMsgInternal             ErrMsgResp = "internal server error"
	HdlMsgBadRequest           ErrMsgResp = "invalid request parameters"
	HdlMsgResourceExist        ErrMsgResp = "resource already exist"
	HdlMsgResourceEmpty        ErrMsgResp = "resource is empty"
	HdlMsgSubscriptionNotFound ErrMsgResp = "anda belum berlangganan"
)

const (
	THdlBadRequest           ErrResponseType = "BAD_REQUEST"
	THandlerInternal         ErrResponseType = "INTERNAL_SERVER_ERROR"
	THdlResourceEmpty        ErrResponseType = "RESOURCE_IS_EMPTY"
	THdlSubscriptionNotFound ErrResponseType = "SUBSCRIPTION_NOT_FOUND"
)
