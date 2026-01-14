package common

type PaginationResponse struct {
	HasNextPage bool `json:"has_next_page"`
	Page        int  `json:"page"`
	PageSize    int  `json:"page_size"`
}
