package data_common

import "math"

type Page struct {
	Number        int64          `json:"number,omitempty"`
	Size          int64          `json:"size,omitempty"`
	Content       *[]interface{} `json:"content,omitempty"`
	First         bool           `json:"first,omitempty"`
	Last          bool           `json:"last,omitempty"`
	TotalPages    int64          `json:"totalPages,omitempty"`
	TotalElements int64          `json:"totalElements,omitempty"`
	Value         any            `json:"value,omitempty"`
}

func Of(data *[]interface{}, p Pageable, totals int64) *Page {

	var pageable Pageable
	if nil == p {
		pageable = PageRequest{Page: 0, Size: 0}
	}

	var totalPages int64
	if pageable.GetPageSize() == 0 {
		totalPages = 1
	} else {
		tmp := totals / pageable.GetPageSize()
		totalPages = int64(math.Ceil(float64(tmp)))
	}

	return &Page{
		Number:        pageable.GetPageNumber(),
		Size:          pageable.GetPageSize(),
		Content:       data,
		First:         pageable.GetPageNumber() == 0,
		Last:          pageable.GetPageNumber()+1 == totalPages,
		TotalPages:    totalPages,
		TotalElements: totals,
	}
}
