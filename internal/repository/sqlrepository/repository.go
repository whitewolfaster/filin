package sqlrepository

import (
	"database/sql"

	"github.com/whitewolfaster/filin/internal/repository"
)

type Repository struct {
	db             *sql.DB
	bookRepository *BookRepository
	userRepository *UserRepository
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (repository *Repository) Books() repository.BookRepository {
	if repository.bookRepository != nil {
		return repository.bookRepository
	}

	repository.bookRepository = &BookRepository{
		repository: repository,
	}
	return repository.bookRepository
}

func (repository *Repository) Users() repository.UserRepository {
	if repository.userRepository != nil {
		return repository.userRepository
	}

	repository.userRepository = &UserRepository{
		repository: repository,
	}
	return repository.userRepository
}
