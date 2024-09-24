package apperror

import (
	"github.com/xendit/xendit-go"
)

const (
	TInvoiceNotFound Type = "INVOICE_NOT_FOUND"
)

const (
	MsgInvoiceNotfound string = "invoice not found"
)

func NewXenditError(err *xendit.Error) *Error {
	if err.ErrorCode == "INVOICE_NOT_FOUND_ERROR" {
		return &Error{
			Type:    TInvoiceNotFound,
			Message: MsgInvoiceNotfound,
		}
	}
	return NewInternal()
}
