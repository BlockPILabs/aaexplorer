package vo

type GetBundlersRequest struct {
	PaginationRequest
	Network string `json:"network" params:"network" validate:"required,min=3"`
}
