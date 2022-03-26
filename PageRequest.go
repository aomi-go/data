package data_common

type PageRequest struct {
	Page int64
	Size int64
}

func (r PageRequest) GetPageSize() int64 {
	return r.Size
}

func (r PageRequest) GetPageNumber() int64 {
	return r.Page
}

func (r PageRequest) GetOffset() int64 {
	return r.Page * r.Size
}
