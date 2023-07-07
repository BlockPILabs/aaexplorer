package vo

type Pagination struct {
	TotalCount int64 `json:"totalCount"`
	PerPage    int   `json:"perPage"`
	Page       int   `json:"page"`
	//Records    Any   `json:"records"` // define in used struct
}

type PaginationAny struct {
	Pagination
	Records any `json:"records"`
}
