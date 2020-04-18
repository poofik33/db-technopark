package service

type Repository interface {
	GetCountForum() (uint64, error)
	GetCountPost() (uint64, error)
	GetCountThread() (uint64, error)
	GetCountUser() (uint64, error)

	DeleteAllForum() error
	DeleteAllPost() error
	DeleteAllThread() error
	DeleteAllUser() error
}
