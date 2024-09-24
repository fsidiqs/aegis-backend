package model

import (
	"strings"
	"time"

	"github.com/fsidiqs/aegis-backend/helper/datehelper"
)

type QueryReqDateFilter struct {
	From string `form:"from"`
	To   string `form:"to"`
}

func (d QueryReqDateFilter) ToDateFilter() *DateFilter {
	var from, to time.Time
	var err error
	trimFrom := strings.TrimSpace(d.From)
	trimTo := strings.TrimSpace(d.To)

	if len(trimFrom) == 0 && len(trimTo) == 0 {
		return nil
	}

	if len(trimFrom) > 0 {
		from, err = datehelper.Date(trimFrom)
		if err != nil {
			return nil
		}

	}
	// set initial to to now if no query param set
	to = time.Now()
	if len(trimTo) > 0 {
		to, err = datehelper.Date(trimTo)
		to = to.AddDate(0, 0, 1)
		// if date is invalid, make it to now
		if err != nil {
			return nil
		}
	}

	return &DateFilter{
		From: from,
		To:   to,
	}
}
