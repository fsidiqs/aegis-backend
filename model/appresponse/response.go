package appresponse

import "github.com/fsidiqs/aegis-backend/model/apperror"

type (
	SuccessResponseType string
	ErrResponseType     string
	SuccessMsgResp      string
	ErrMsgResp          string
)

type ResponseWrapper struct {
	StatusCode int
	Response   interface{}
}

type SuccessResponse struct {
	Type    SuccessResponseType `json:"type"`
	Message SuccessMsgResp      `json:"message"`
	Data    interface{}         `json:"data"`
}

// A Response Error
//
// swagger:response ErrorResponse
type ErrorResponse struct {
	Type    ErrResponseType `json:"type"`
	Message ErrMsgResp      `json:"message"`
}

type ErrorResponseMessageArr struct {
	Type     ErrResponseType          `json:"type"`
	Messages []map[string]interface{} `json:"messages"`
}

type ErrRespWithInvalidArgs struct {
	Type        ErrResponseType            `json:"type"`
	Message     ErrMsgResp                 `json:"message"`
	InvalidArgs []apperror.InvalidArgument `json:"invalid_args"`
}
