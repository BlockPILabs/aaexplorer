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

type PaginationRequest struct {
	PerPage int `json:"perPage" query:"perPage" params:"perPage"`
	Page    int `json:"page" query:"page" params:"page"`    // page number
	Sort    int `json:"sort" query:"sort" params:"sort"`    // sort field idx
	Order   int `json:"order" query:"order" params:"order"` // order sort : -1 desc   1 asc
}
