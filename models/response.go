package models

type Response struct {
	Status  bool        `json:"status"`
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ResponsePagination struct {
	Status     bool           `json:"status"`
	Code       int            `json:"code"`
	Message    string         `json:"message"`
	Pagination MetaPagination `json:"pagination"`
	Data       interface{}    `json:"data,omitempty"`
}

type MetaPagination struct {
	CurrentPage int64 `json:"current_page"`
	NextPage    int64 `json:"next_page"`
	PrevPage    int64 `json:"prev_page"`
	TotalPages  int64 `json:"total_pages"`
	TotalCount  int64 `json:"total_count"`
}
