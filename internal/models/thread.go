package models

import (
	"time"
)

type Thread struct {
	Author       string    `json:"author"`
	CreationDate time.Time `json:"created"`
	Forum        string    `json:"forum"`
	ID           uint64    `json:"id"`
	About        string    `json:"message"`
	Slug         string    `json:"slug"`
	Title        string    `json:"title"`
	Votes        int64     `json:"votes"`
}
