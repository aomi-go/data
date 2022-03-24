package data_common

type PageRequest struct {
	Page uint
	Size uint
}

func (r PageRequest) GetPageSize() uint {
	return r.Size
}

func (r PageRequest) GetPageNumber() uint {
	return r.Page
}

func (r PageRequest) GetOffset() uint {
	return r.Page * r.Size
}
