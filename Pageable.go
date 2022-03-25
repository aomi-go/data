package data_common

type Pageable interface {
	GetPageSize() int64
	GetPageNumber() int64
	GetOffset() int64
}
