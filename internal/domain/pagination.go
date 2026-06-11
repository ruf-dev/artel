package domain

type Paging struct {
	Limit  uint32
	Offset uint32
}

type ListUsersReq struct {
	Paging Paging
	Search string
}
