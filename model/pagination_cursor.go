package model

const (
	defaultPage  int = 1
	defaultLimit int = 10
)

type Pagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Page   int `json:"page"`
}

type DataPaginated struct {
	Data  interface{}            `json:"items"`
	Other map[string]interface{} `json:"other,omitempty"`
	Pagination
}

func (d *DataPaginated) Paginate(m Pagination) {
	page, limit, offset := BuildPagination(m)
	d.Pagination = Pagination{
		Limit:  limit,
		Offset: offset,
		Page:   page,
	}
}

func (p *DataPaginated) UpdateData(d map[string]interface{}) {
	for k, v := range d {
		p.Other[k] = v
	}
}

// BuildPagination returns page, limit, offset,
func BuildPagination(m Pagination) (int, int, int) {
	page := (m).Page
	if page < 1 {
		page = defaultPage
	}

	limit := (m).Limit
	if limit == 0 {
		limit = defaultLimit
	}

	offset := (limit * page) - limit

	return page, limit, offset
}

func BuildPages(elems int, limit int) int {
	var pages int
	total := int(elems)
	if pages = (total / limit); (total % limit) != 0 {
		pages++
	}
	return pages
}
