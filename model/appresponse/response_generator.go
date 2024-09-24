package appresponse

import (
	"github.com/fsidiqs/aegis-backend/model/apperror"
)

func HdlRespBadRequest() ErrorResponse {
	return ErrorResponse{Message: HdlMsgBadRequest, Type: THdlBadRequest}
}

func HdlRespInternalServerError() ErrorResponse {
	return ErrorResponse{Message: HdlMsgInternal, Type: THandlerInternal}
}

func HdlRespUnautthorized() ErrorResponse {
	return ErrorResponse{Message: HdlMsgUnauthorized, Type: THdlUnauthorized}
}

// PrepareErrResp check whether the error is an apperror and then return error response with Internal Server Error as default value
func PrepareErr(err error, msg string) (errMsgInternal string, errResp ErrorResponse, apperr *apperror.Error) {
	errResp = HdlRespInternalServerError()
	apperr, ok := err.(*apperror.Error)
	if !ok {
		apperr = apperror.NewInternal()
	}
	errMsgInternal = apperror.ErrorWrapper(err, msg)
	return
}
