package rest

type Response[T any] struct {
	Success    bool              `json:"success"`
	Message    string            `json:"message"`
	Data       T                 `json:"data,omitempty"`
	Errors     []ValidationError `json:"errors,omitempty"`
	Pagination *Pagination       `json:"pagination,omitempty"`
}

type ValidationError struct {
	Field string   `json:"field"`
	Tags  []string `json:"tags"`
}

type Pagination struct {
	Count     int `json:"count"`
	Page      int `json:"page"`
	TotalPage int `json:"totalPage"`
	Size      int `json:"size"`
}
