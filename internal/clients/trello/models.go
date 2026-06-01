package trello

type member struct {
	FullName string `json:"fullName"`
}

type board struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
