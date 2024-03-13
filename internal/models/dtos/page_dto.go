package dtos

type PageDTO struct {
	Total      int64       `json:"total"`
	Page       int64       `json:"page"`
	Data       interface{} `json:"data"`
	TotalPages int64       `json:"totalPages"`
}

func NewPageDTO(total int64, page int64, totalPages int64, data interface{}) *PageDTO {
	return &PageDTO{
		TotalPages: totalPages,
		Total:      total,
		Page:       page,
		Data:       data,
	}
}
