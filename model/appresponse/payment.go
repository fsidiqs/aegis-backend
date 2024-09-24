package appresponse

const (
	TUserHasUnfinishedInv ErrResponseType = "USER_HAS_UNFINISHED_INVOICE"
)

const (
	RespUserHasUnfinishedInv ErrMsgResp = "user has unfinished invoice"
)

const (
	ErrorPaymentStatusAmount ErrMsgResp = "status unsuccessful and incorrect amount"
	ErrorPaymentStatus       ErrMsgResp = "status unsuccessful"
	ErrorPaymentAmount       ErrMsgResp = "incorrect amount"
)
