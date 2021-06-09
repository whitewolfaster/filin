package sqlrepository

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/whitewolfaster/filin/internal/models"
	"github.com/whitewolfaster/filin/internal/repository"
)

type BookRepository struct {
	repository *Repository
}

func (br *BookRepository) Create(book *models.Book) error {
	_, err := br.repository.db.Exec(
		"INSERT INTO books (id, name, author, year, genre, pubhouse, coverpath, description, price) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
		book.ID,
		book.Name,
		book.Author,
		book.Year,
		book.Genre,
		book.PubHouse,
		book.CoverPath,
		book.Description,
		book.Price,
	)
	if err != nil {
		return err
	}
	return nil
}

func (br *BookRepository) Update(book *models.Book) error {
	_, err := br.repository.db.Exec("UPDATE books SET name=$1, author=$2, year=$3, genre=$4, pubhouse=$5, coverpath=$6, description=$7, price=$8 WHERE id=$9",
		book.Name,
		book.Author,
		book.Year,
		book.Genre,
		book.PubHouse,
		book.CoverPath,
		book.Description,
		book.Price,
		book.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (br *BookRepository) GetAll() (*[]models.Book, error) {
	rows, err := br.repository.db.Query("SELECT * FROM books")
	if err != nil {
		return nil, err
	}
	books := []models.Book{}
	book := models.Book{}
	for rows.Next() {
		err = rows.Scan(
			&book.ID,
			&book.Name,
			&book.Author,
			&book.Year,
			&book.Genre,
			&book.PubHouse,
			&book.CoverPath,
			&book.Description,
			&book.Price,
		)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return &books, nil
}

func (br *BookRepository) GetAllBookID() ([]string, error) {
	rows, err := br.repository.db.Query("select id from books")
	if err != nil {
		return nil, err
	}
	var books_id []string
	for rows.Next() {
		book_id := ""
		err = rows.Scan(&book_id)
		if err != nil {
			return nil, err
		}
		books_id = append(books_id, book_id)
	}
	return books_id, nil
}

func (br *BookRepository) GenerateBookID(ids *[]string) string {
	uid := uuid.New().String()
	for _, id := range *ids {
		if uid == id {
			uid = br.GenerateBookID(ids)
		}
	}
	return uid
}

func (br *BookRepository) GetAllGenres() ([]models.Genre, error) {
	rows, err := br.repository.db.Query("SELECT name FROM genres")
	if err != nil {
		return nil, err
	}
	genres := []models.Genre{}
	genre := models.Genre{}
	for rows.Next() {
		err = rows.Scan(
			&genre.Name,
		)
		if err != nil {
			return nil, err
		}
		genres = append(genres, genre)
	}
	return genres, nil
}

func (br *BookRepository) FindBookByID(id string) (*models.Book, error) {
	book := models.Book{}
	if err := br.repository.db.QueryRow(
		"SELECT * FROM books WHERE id = $1",
		id,
	).Scan(
		&book.ID,
		&book.Name,
		&book.Author,
		&book.Year,
		&book.Genre,
		&book.PubHouse,
		&book.CoverPath,
		&book.Description,
		&book.Price,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrRecordNotFound
		}
		return nil, err
	}

	return &book, nil
}

func (br *BookRepository) Delete(bookID string) error {
	_, err := br.repository.db.Exec("DELETE FROM books WHERE id=$1", bookID)
	if err != nil {
		return err
	}
	return nil
}

func (br *BookRepository) BooksOfMonth() (*[]string, error) {
	rows, err := br.repository.db.Query("SELECT * FROM books_of_month")
	if err != nil {
		return nil, err
	}
	books := []string{}
	book := ""
	for rows.Next() {
		err = rows.Scan(
			&book,
		)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return &books, nil
}

func (br *BookRepository) BookLiders() (*[]string, error) {
	rows, err := br.repository.db.Query("SELECT * FROM book_liders")
	if err != nil {
		return nil, err
	}
	books := []string{}
	book := ""
	for rows.Next() {
		err = rows.Scan(
			&book,
		)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return &books, nil
}

func (br *BookRepository) CreateGenre(genre *models.Genre) error {
	_, err := br.repository.db.Exec(
		"INSERT INTO genres (name) VALUES ($1)",
		genre.Name,
	)
	if err != nil {
		return err
	}
	return nil
}

func (br *BookRepository) DeleteGenre(genre *models.Genre) error {
	_, err := br.repository.db.Exec(
		"DELETE FROM genres WHERE name=$1",
		genre.Name,
	)
	if err != nil {
		return err
	}
	return nil
}
