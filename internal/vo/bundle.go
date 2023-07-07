package vo

type GetBundlersRequest struct {
	PaginationRequest
	Network string `json:"network" params:"network"`
}
