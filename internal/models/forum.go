package models

type Forum struct {
	PostsCount    uint64 `json:"posts"`
	Slug          string `json:"slug"`
	ThreadsCount  uint64 `json:"threads"`
	Title         string `json:"title"`
	AdminNickname string `json:"user"`
}
