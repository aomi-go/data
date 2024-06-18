package page

import (
	"github.com/aomi-go/data/common/sort"
	"math"
)

func NewDefaultPageable() *Pageable {
	return &Pageable{
		Page: 0,
		Size: 20,
	}
}

// Pageable 分页请求
type Pageable struct {
	sort.Sort
	Page int `json:"page" describe:"页码从0开始"`
	Size int `json:"size" describe:"每页的大小"`
}

func NewPage[T interface{}](content []T, total int64, pageable *Pageable) *Page[T] {
	if nil == pageable {
		pageable = NewDefaultPageable()
	}

	totalPages := calculateTotalPages(total, pageable.Size)
	hasPrevious := pageable.Page > 0
	hasNext := pageable.Page+1 < totalPages

	return &Page[T]{
		Empty:            len(content) == 0,
		First:            !hasPrevious,
		Last:             !hasNext,
		Number:           pageable.Page,
		NumberOfElements: len(content),
		Size:             pageable.Size,
		TotalElements:    total,
		TotalPages:       totalPages,
		Content:          content,
	}
}

type Page[T interface{}] struct {
	Empty            bool  `json:"empty"`
	First            bool  `json:"first"`
	Last             bool  `json:"last"`
	Number           int   `json:"number"`
	NumberOfElements int   `json:"numberOfElements"`
	Size             int   `json:"size"`
	TotalElements    int64 `json:"totalElements"`
	TotalPages       int   `json:"totalPages"`
	Content          []T   `json:"content"`
}

func calculateTotalPages(total int64, size int) int {
	if size <= 0 {
		return 0
	}
	return int(math.Ceil(float64(total) / float64(size)))
}
