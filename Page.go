package data_common

type Page[T interface{}] struct {
	Number        uint `json:"number,omitempty"`
	Size          uint `json:"size,omitempty"`
	Content       []T  `json:"content,omitempty"`
	First         bool `json:"first,omitempty"`
	Last          bool `json:"last,omitempty"`
	TotalPages    uint `json:"totalPages,omitempty"`
	TotalElements uint `json:"totalElements,omitempty"`
}
