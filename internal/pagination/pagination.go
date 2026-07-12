package pagination

type Request struct {
	Page     int
	PageSize int
}

func (r Request) Limit() int {
	return r.PageSize
}

func (r Request) Offset() int {
	return (r.Page - 1) * r.PageSize
}

type Meta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	TotalItems int64 `json:"totalItems"`
	TotalPages int   `json:"totalPages"`
}

func NewMeta(
	page int,
	pageSize int,
	totalItems int64,
) Meta {

	totalPages := 0

	if pageSize > 0 {
		totalPages = int(
			(totalItems + int64(pageSize) - 1) / int64(pageSize),
		)
	}

	return Meta{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
}
