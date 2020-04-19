package models

type Vote struct {
	Nickname string `json:"nickname"`
	Voice    int64  `json:"voice"`
	ThreadID uint64 `json:"_"`
}
