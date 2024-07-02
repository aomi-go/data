package page

import (
	"github.com/aomi-go/data/common/sort"
	"math"
)

func NewPageable(page int, size int) *Pageable {
	return &Pageable{
		Page: page,
		Size: size,
	}
}

func NewPageableWithSort(page int, size int, s sort.Sort) *Pageable {
	return &Pageable{
		Page: page,
		Size: size,
		Sort: s,
	}
}

func NewDefaultPageable() *Pageable {
	return NewPageable(0, 20)
}

// Pageable 分页请求
type Pageable struct {
	sort.Sort
	Page int `form:"page" json:"page" describe:"页码从0开始"`
	Size int `form:"size" json:"size" describe:"每页的大小"`
}

func (p *Pageable) GetOffset() int64 {
	return int64(p.Page * p.Size)
}

func EmptyPage[T interface{}]() *Page[T] {
	var content []*T
	return NewPage[T](content, 0, nil)
}

func NewPage[T interface{}](content []*T, total int64, pageable *Pageable) *Page[T] {
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
	Content          []*T  `json:"content"`
	Extra            any   `json:"extra" describe:"额外信息"`
}

// Map 转换page类型
func Map[T interface{}, N interface{}](old *Page[T], mapfunc func(t *T) *N) *Page[N] {
	var newPage = Page[N]{
		Empty:            old.Empty,
		First:            old.First,
		Last:             old.Last,
		Number:           old.Number,
		NumberOfElements: old.NumberOfElements,
		Size:             old.Size,
		TotalElements:    old.TotalElements,
		TotalPages:       old.TotalPages,
	}

	var content = make([]*N, len(old.Content))
	for i, v := range old.Content {
		content[i] = mapfunc(v)
	}
	newPage.Content = content

	return &newPage
}

func calculateTotalPages(total int64, size int) int {
	if size <= 0 {
		return 0
	}
	return int(math.Ceil(float64(total) / float64(size)))
}
