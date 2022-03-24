package data_common

type Pageable interface {
	GetPageSize() uint
	GetPageNumber() uint
	GetOffset() uint
}
