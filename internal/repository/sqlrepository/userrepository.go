package sqlrepository

import (
	"database/sql"
	//"fmt"

	"github.com/google/uuid"
	"github.com/whitewolfaster/filin/internal/models"
	"github.com/whitewolfaster/filin/internal/repository"
)

type UserRepository struct {
	repository *Repository
}

func (ur *UserRepository) Create(u *models.User) error {
	_, err := ur.repository.db.Exec(
		"INSERT INTO users (id, email, password, token, activated) VALUES ($1, $2, $3, $4, $5)",
		u.ID,
		u.Email,
		u.EncPassword,
		u.Token,
		u.Activated,
	)
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) GetAllUserID() ([]string, error) {
	rows, err := ur.repository.db.Query("select id from users")
	if err != nil {
		return nil, err
	}
	var users_id []string
	for rows.Next() {
		user_id := ""
		err = rows.Scan(&user_id)
		if err != nil {
			return nil, err
		}
		users_id = append(users_id, user_id)
	}
	return users_id, nil
}

func (ur *UserRepository) GetAllTokens() ([]string, error) {
	rows, err := ur.repository.db.Query("select token from users")
	if err != nil {
		return nil, err
	}
	var tokens []string
	for rows.Next() {
		token := ""
		err = rows.Scan(&token)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}

func (ur *UserRepository) GenerateID(idArray *[]string) string {
	uid := uuid.New().String()
	for _, id := range *idArray {
		if uid == id {
			uid = ur.GenerateID(idArray)
		}
	}
	return uid
}

func (ur *UserRepository) GenerateToken(idArray *[]string) string {
	uid := uuid.New().String()
	for _, id := range *idArray {
		if uid == id {
			uid = ur.GenerateToken(idArray)
		}
	}
	return uid
}

func (ur *UserRepository) FindByEmail(email string) (*models.User, error) {
	u := models.User{}
	if err := ur.repository.db.QueryRow(
		"SELECT id, email, password, token, activated FROM users WHERE email = $1",
		email,
	).Scan(
		&u.ID,
		&u.Email,
		&u.EncPassword,
		&u.Token,
		&u.Activated,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrRecordNotFound
		}
		return nil, err
	}

	return &u, nil
}

func (ur *UserRepository) FindByID(id string) (*models.User, error) {
	u := models.User{}
	if err := ur.repository.db.QueryRow(
		"SELECT id, email, password, token, activated FROM users WHERE id = $1",
		id,
	).Scan(
		&u.ID,
		&u.Email,
		&u.EncPassword,
		&u.Token,
		&u.Activated,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrRecordNotFound
		}
		return nil, err
	}

	return &u, nil
}

func (ur *UserRepository) Cart(user_id string) (*[]string, error) {
	rows, err := ur.repository.db.Query("SELECT book FROM user_book where user_id=$1", user_id)
	if err != nil {
		return nil, err
	}

	books := []string{}
	book_id := ""
	for rows.Next() {
		err = rows.Scan(&book_id)
		if err != nil {
			return nil, err
		}
		books = append(books, book_id)
	}
	return &books, nil
}

func (ur *UserRepository) AddToCart(user_id string, book_id string) error {
	_, err := ur.repository.db.Exec("INSERT INTO user_book values($1, $2)", user_id, book_id)
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) DeleteFromCart(user_id string, book_id string) error {
	_, err := ur.repository.db.Exec("DELETE FROM user_book WHERE user_id=$1 AND book=$2", user_id, book_id)
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) CreateOrder(order *models.Order) error {
	_, err := ur.repository.db.Exec("INSERT INTO orders values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		order.ID,
		order.UserID,
		order.First_name,
		order.Last_name,
		order.Patronymic,
		order.Phone,
		order.City,
		order.Address,
		order.Post_index,
		order.Date,
		order.OrderSumm,
	)
	if err != nil {
		return err
	}
	for _, prod := range (*order).Products {
		_, err := ur.repository.db.Exec("INSERT INTO order_book values($1, $2, $3)",
			order.ID,
			prod.ID,
			prod.Quantity,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ur *UserRepository) GetAllOrderID() ([]string, error) {
	rows, err := ur.repository.db.Query("select id from orders")
	if err != nil {
		return nil, err
	}
	var orders_id []string
	for rows.Next() {
		order_id := ""
		err = rows.Scan(&order_id)
		if err != nil {
			return nil, err
		}
		orders_id = append(orders_id, order_id)
	}
	return orders_id, nil
}

func (ur *UserRepository) ChangePassword(userID string, newPassword string) error {
	encPswd := models.EncryptString(newPassword)
	_, err := ur.repository.db.Exec("UPDATE users SET password=$1 WHERE id=$2", encPswd, userID)
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) ChangeAdminPassword(adminID string, newPassword string) error {
	encPswd := models.EncryptString(newPassword)
	_, err := ur.repository.db.Exec("UPDATE admins SET password=$1 WHERE id=$2", encPswd, adminID)
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) FindAdminByLogin(login string) (*models.Admin, error) {
	admin := models.Admin{}
	if err := ur.repository.db.QueryRow(
		"SELECT id, login, password FROM admins WHERE login = $1",
		login,
	).Scan(
		&admin.ID,
		&admin.Login,
		&admin.EncPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrRecordNotFound
		}
		return nil, err
	}

	return &admin, nil
}

func (ur *UserRepository) FindAdminByID(id string) (*models.Admin, error) {
	admin := models.Admin{}
	if err := ur.repository.db.QueryRow(
		"SELECT id, login, password FROM admins WHERE id = $1",
		id,
	).Scan(
		&admin.ID,
		&admin.Login,
		&admin.EncPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrRecordNotFound
		}
		return nil, err
	}

	return &admin, nil
}

func (ur *UserRepository) CreateAdmin(a *models.Admin) error {
	_, err := ur.repository.db.Exec(
		"INSERT INTO admins (id, login, password) VALUES ($1, $2, $3)",
		a.ID,
		a.Login,
		a.EncPassword,
	)
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) DeleteAdmin(adminID string) error {
	_, err := ur.repository.db.Exec("DELETE FROM admins WHERE id=$1", adminID)
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) GetAllAdminID() ([]string, error) {
	rows, err := ur.repository.db.Query("select id from admins")
	if err != nil {
		return nil, err
	}
	var admins_id []string
	for rows.Next() {
		admin_id := ""
		err = rows.Scan(&admin_id)
		if err != nil {
			return nil, err
		}
		admins_id = append(admins_id, admin_id)
	}
	return admins_id, nil
}

func (ur *UserRepository) GetOrders() (*[]models.Order, error) {
	orders := []models.Order{}
	order := models.Order{}
	rows, err := ur.repository.db.Query("SELECT * FROM orders")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(
			&order.ID,
			&order.UserID,
			&order.First_name,
			&order.Last_name,
			&order.Patronymic,
			&order.Phone,
			&order.City,
			&order.Address,
			&order.Post_index,
			&order.Date,
			&order.OrderSumm,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	for i := 0; i < len(orders); i++ {
		rows, err := ur.repository.db.Query("SELECT book_id, quantity FROM order_book WHERE order_id=$1", orders[i].ID)
		if err != nil {
			return nil, err
		}
		product := models.Product{}
		for rows.Next() {
			err = rows.Scan(
				&product.ID,
				&product.Quantity,
			)
			if err != nil {
				return nil, err
			}
			orders[i].Products = append(orders[i].Products, product)
		}
	}
	return &orders, nil
}

func (ur *UserRepository) GetAllAdmins() (*[]models.Admin, error) {
	rows, err := ur.repository.db.Query("SELECT * FROM admins")
	if err != nil {
		return nil, err
	}
	admins := []models.Admin{}
	admin := models.Admin{}
	for rows.Next() {
		err = rows.Scan(
			&admin.ID,
			&admin.Login,
			&admin.EncPassword,
		)
		if err != nil {
			return nil, err
		}
		admins = append(admins, admin)
	}
	return &admins, nil
}

func (ur *UserRepository) Activate(userID string) error {
	_, err := ur.repository.db.Exec("UPDATE users SET activated=$1 WHERE id=$2", 1, userID)
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) ClearCart(userID string) error {
	_, err := ur.repository.db.Exec("DELETE FROM user_book WHERE user_id=$1", userID)
	if err != nil {
		return err
	}
	return nil
}
