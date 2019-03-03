package repository

import (
	"github.com/richardcase/ingest-sample/pkg/api"
)

// Repository represents a repository of people
type Repository interface {
	GetByID(id int64) (*api.Person, error)
	Store(person *api.Person) error
	Delete(id int64) error
	Check() bool
}
