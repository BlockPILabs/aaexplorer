package vo

import "github.com/BlockPILabs/aaexplorer/config"

type Pagination struct {
	TotalCount int `json:"totalCount"`
	PerPage    int `json:"perPage"`
	Page       int `json:"page"`
	//Records    Any   `json:"records"` // define in used struct
}

type PaginationAny struct {
	Pagination
	Records any `json:"records"`
}

type PaginationRequest struct {
	TotalCount int      `json:"totalCount" query:"totalCount" params:"totalCount" `
	PerPage    int      `json:"perPage" query:"perPage" params:"perPage" validate:"required,min=1"`
	Page       int      `json:"page" query:"page" params:"page" validate:"required,min=1,max=1000"` // page number
	Sort       int      `json:"sort" query:"sort" params:"sort" validate:"min=0,max=50"`            // sort field idx
	Order      int      `json:"order" query:"order" params:"order" validate:"min=-1,max=1"`         // order sort : -1 desc   1 asc
	Select     []string `json:"-" query:"-" params:"-" `
}

func NewDefaultPaginationRequest() PaginationRequest {
	return PaginationRequest{
		PerPage: config.DefaultPerPage,
		Page:    config.MinPage,
	}
}

func (r *PaginationRequest) GetOffset() int {
	return (r.GetPage() - 1) * r.GetPerPage()
}

func (r *PaginationRequest) GetPerPage() int {
	if r.PerPage > config.MaxPerPage {
		r.PerPage = config.MaxPerPage
		return config.MaxPerPage
	} else if r.PerPage < config.MinPerPage {
		r.PerPage = config.MinPerPage
		return config.MinPerPage
	}
	return r.PerPage
}

func (r *PaginationRequest) GetPage() int {
	if r.Page > config.MaxPage {
		r.Page = config.MaxPage
		return config.MaxPage
	} else if r.Page < config.MinPage {
		r.Page = config.MinPage
		return config.MinPage
	}
	return r.Page
}
