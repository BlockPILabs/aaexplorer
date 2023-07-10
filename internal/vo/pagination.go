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
	PerPage int `json:"perPage" query:"perPage" params:"perPage" validate:"required,min=1"`
	Page    int `json:"page" query:"page" params:"page" validate:"required,min=1"`  // page number
	Sort    int `json:"sort" query:"sort" params:"sort" validate:"min=0,max=50"`    // sort field idx
	Order   int `json:"order" query:"order" params:"order" validate:"min=-1,max=1"` // order sort : -1 desc   1 asc
}
