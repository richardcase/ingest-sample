package repository

import (
	"github.com/richardcase/ingest-sample/pkg/api"
)

type Repository interface {
	GetByID(pid int64) (*api.Person, error)
	Store(person *api.Person) error
	Delete(id int64) error
	Check() bool
}
