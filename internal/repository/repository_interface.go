package repository

import (
	"github.com/whitewolfaster/filin/internal/models"
)

type BookRepository interface {
	Create(*models.Book) error
	Update(*models.Book) error
	Delete(string) error
	GenerateBookID(*[]string) string
	GetAll() (*[]models.Book, error)
	GetAllGenres() ([]models.Genre, error)
	GetAllBookID() ([]string, error)
	FindBookByID(string) (*models.Book, error)
	BookLiders() (*[]string, error)
	BooksOfMonth() (*[]string, error)
	CreateGenre(*models.Genre) error
	DeleteGenre(*models.Genre) error
}

type UserRepository interface {
	Create(*models.User) error
	GetAllUserID() ([]string, error)
	GetAllTokens() ([]string, error)
	GenerateID(*[]string) string
	GenerateToken(*[]string) string
	FindByEmail(string) (*models.User, error)
	FindByID(string) (*models.User, error)
	Cart(string) (*[]string, error)
	AddToCart(string, string) error
	DeleteFromCart(string, string) error
	GetAllOrderID() ([]string, error)
	CreateOrder(*models.Order) error
	ChangePassword(string, string) error
	FindAdminByLogin(string) (*models.Admin, error)
	FindAdminByID(string) (*models.Admin, error)
	CreateAdmin(*models.Admin) error
	GetAllAdminID() ([]string, error)
	GetOrders() (*[]models.Order, error)
	ChangeAdminPassword(string, string) error
	DeleteAdmin(string) error
	GetAllAdmins() (*[]models.Admin, error)
	Activate(string) error
	ClearCart(string) error
}

type Repository interface {
	Books() BookRepository
	Users() UserRepository
}
