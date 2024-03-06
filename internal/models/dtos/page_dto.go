package dtos

type PageDTO struct {
	Total int64       `json:"total"`
	Page  int64       `json:"page"`
	Data  interface{} `json:"data"`
}

func NewPageDTO(total int64, page int64, data interface{}) *PageDTO {
	return &PageDTO{
		Total: total,
		Page:  page,
		Data:  data,
	}
}
