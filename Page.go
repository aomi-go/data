package data_common

import "math"

type Page struct {
	Number        int64       `json:"number"`
	Size          int64       `json:"size"`
	Content       interface{} `json:"content"`
	First         bool        `json:"first"`
	Last          bool        `json:"last"`
	TotalPages    int64       `json:"totalPages"`
	TotalElements int64       `json:"totalElements"`
	Value         *any        `json:"value,omitempty"`
}

func Of(data interface{}, p Pageable, totals int64) Page {

	var pageable Pageable
	if nil == p {
		pageable = PageRequest{Page: 0, Size: 0}
	} else {
		pageable = p
	}

	var totalPages int64
	if pageable.GetPageSize() == 0 {
		totalPages = 1
	} else {
		tmp := float64(totals) / float64(pageable.GetPageSize())
		totalPages = int64(math.Ceil(tmp))
	}

	var tmpData = data
	if nil == data {
		d := make([]interface{}, 0)
		tmpData = &d
	}

	return Page{
		Number:        pageable.GetPageNumber(),
		Size:          pageable.GetPageSize(),
		Content:       tmpData,
		First:         pageable.GetPageNumber() == 0,
		Last:          pageable.GetPageNumber()+1 == totalPages || totals == 0,
		TotalPages:    totalPages,
		TotalElements: totals,
	}
}
