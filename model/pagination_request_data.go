package model

import "strconv"

type QueryReqPagination struct {
	Limit string `form:"limit"`
	Page  string `form:"page"`
}

// ToPagination Build Pagination model with offset value
func (data QueryReqPagination) ToPagination() Pagination {
	var page, limit, offset int
	// var err error

	page, _ = strconv.Atoi(data.Page)
	if page < 1 {
		page = defaultPage
	}

	limit, _ = strconv.Atoi(data.Limit)
	if limit == 0 {
		limit = defaultLimit
	}

	offset = (limit * page) - limit
	return Pagination{
		Limit:  limit,
		Page:   page,
		Offset: int(offset),
	}
}
