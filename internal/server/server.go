package server

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/smtp"

	"time"

	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"github.com/whitewolfaster/filin/internal/models"
	"github.com/whitewolfaster/filin/internal/repository"
	"github.com/whitewolfaster/filin/internal/repository/sqlrepository"
	"github.com/whitewolfaster/filin/internal/sessionstore"
)

const (
	COOKIE_NAME         = "filin_shop"
	COOKIE_ADMIN        = "filin_admin"
	ctxUserKey   ctxKey = 0
)

type ctxKey uint8

type Server struct {
	config       *Config
	logger       *logrus.Logger
	router       *mux.Router
	repository   repository.Repository
	sessionStore *sessionstore.SessionStore
}

func NewServer(config *Config) (*Server, error) {
	db, err := sql.Open("sqlite3", config.DatabaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	server := &Server{
		config:       config,
		logger:       logrus.New(),
		router:       mux.NewRouter(),
		repository:   sqlrepository.NewRepository(db),
		sessionStore: sessionstore.New(),
	}
	return server, nil
}

func (server *Server) Start() error {
	if err := server.configureLogger(); err != nil {
		return err
	}

	server.configureRouter()

	server.logger.Infof("server started on %s with loglevel '%s'", server.config.Port, server.config.LogLevel)

	port := os.Getenv("PORT")
	if port == "" {
		port = server.config.Port
	}

	return http.ListenAndServe(":"+port, server.router)
}

func (server *Server) configureLogger() error {
	level, err := logrus.ParseLevel(server.config.LogLevel)
	if err != nil {
		return err
	}
	server.logger.SetLevel(level)
	return nil
}

func (server *Server) configureRouter() {
	/*
		FileServer for static files(images, css, js, fonts etc.)
		Serve requests like "/static/..." - for example "/static/css/style.css"
	*/
	server.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static/"))))

	// API routes

	server.router.HandleFunc("/api/private/createAdminSession", server.APIauth(server.APICreateAdminSession())).Methods(http.MethodPost)
	server.router.HandleFunc("/api/private/deleteAdminSession", server.APIauth(server.APIDeleteAdminSession())).Methods(http.MethodGet)
	server.router.HandleFunc("/api/private/createBook", server.APIauth(server.APICreateBook())).Methods(http.MethodPost)
	server.router.HandleFunc("/api/private/updateBook", server.APIauth(server.APIUpdateBook())).Methods(http.MethodPost)
	server.router.HandleFunc("/api/private/deleteBook", server.APIauth(server.APIDeleteBook())).Methods(http.MethodPost)
	server.router.HandleFunc("/api/private/GetAllOrders", server.APIauth(server.APIGetAllOrders())).Methods(http.MethodGet)
	server.router.HandleFunc("/api/private/changeAdminPassword", server.APIauth(server.APIChangeAdminPassword())).Methods(http.MethodPost)
	server.router.HandleFunc("/api/private/createAdmin", server.APICreateAdmin()).Methods(http.MethodPost)
	server.router.HandleFunc("/api/private/deleteAdmin", server.APIDeleteAdmin()).Methods(http.MethodPost)
	server.router.HandleFunc("/api/private/GetAllAdmins", server.APIauth(server.APIGetAllAdmins())).Methods(http.MethodGet)

	server.router.HandleFunc("/api/public/getAllBooks", server.APIGetAllBooks()).Methods(http.MethodGet)
	server.router.HandleFunc("/api/public/getAllGenres", server.APIGetAllGenres()).Methods(http.MethodGet)
	server.router.HandleFunc("/api/public/isLoggedIn", server.APIIsLoggedIn()).Methods(http.MethodGet)
	server.router.HandleFunc("/api/public/createUser", server.APICreateUser()).Methods(http.MethodPost)
	server.router.HandleFunc("/api/public/createSession", server.WEBauth(server.APICreateSession())).Methods(http.MethodPost)
	server.router.HandleFunc("/api/public/deleteSession", server.WEBauth(server.APIDeleteSession())).Methods(http.MethodGet)
	server.router.HandleFunc("/api/public/getCart", server.WEBauth(server.APIGetCart())).Methods(http.MethodGet)
	server.router.HandleFunc("/api/public/addToCart", server.WEBauth(server.APIAddToCart())).Methods(http.MethodPost)
	server.router.HandleFunc("/api/public/deleteFromCart", server.WEBauth(server.APIDeleteFromCart())).Methods(http.MethodPost)
	server.router.HandleFunc("/api/public/cleanCart", server.WEBauth(server.APICleanCart())).Methods(http.MethodGet)
	server.router.HandleFunc("/api/public/createOrder", server.WEBauth(server.APICreateOrder())).Methods(http.MethodPost)
	server.router.HandleFunc("/api/public/changePassword", server.WEBauth(server.APIChangePassword())).Methods(http.MethodPost)
	server.router.HandleFunc("/api/public/activate/{uid}/{token}", server.APIActivate()).Methods(http.MethodGet)

	// Web site routes

	server.router.HandleFunc("/admin/createBook", server.APIauth(server.handleAdminCreateBook())).Methods(http.MethodGet)
	server.router.HandleFunc("/admin/changepswd", server.APIauth(server.handleAdminChangePSWD())).Methods(http.MethodGet)
	server.router.HandleFunc("/admin/login", server.APIauth(server.handleAdminLogin())).Methods(http.MethodGet)
	server.router.HandleFunc("/admin/admin_list", server.APIauth(server.handleAdminList())).Methods(http.MethodGet)
	server.router.HandleFunc("/admin/book_list", server.APIauth(server.handleAdminBookList())).Methods(http.MethodGet)
	server.router.HandleFunc("/admin/create_admin", server.APIauth(server.handleCreateAdmin())).Methods(http.MethodGet)
	server.router.HandleFunc("/admin", server.APIauth(server.handleAdminHome())).Methods(http.MethodGet)
	server.router.HandleFunc("/admin/updateBook/{id}", server.APIauth(server.handleAdminUpdateBook())).Methods(http.MethodGet)

	server.router.HandleFunc("/catalog", server.WEBauth(server.handleCatalog())).Methods(http.MethodGet)
	server.router.HandleFunc("/cart", server.WEBauth(server.handleCart())).Methods(http.MethodGet)
	server.router.HandleFunc("/order", server.WEBauth(server.handleOrder())).Methods(http.MethodGet)
	server.router.HandleFunc("/register", server.handleRegister()).Methods(http.MethodGet)
	server.router.HandleFunc("/login", server.WEBauth(server.handleLogin())).Methods(http.MethodGet)
	server.router.HandleFunc("/changepswd", server.WEBauth(server.handleChangePSWD())).Methods(http.MethodGet)
	server.router.HandleFunc("/", server.WEBauth(server.handleHome())).Methods(http.MethodGet)
}

// MiddleWare Functions

func (server *Server) SendConfirmMail(to string, message string) error {
	from := server.config.Smtp_login
	password := server.config.Smtp_password
	receiver := []string{to}
	msg := "To: " + to + "\n" + "Subject: FilinShop Activate Account\n" +
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
		message
	auth := smtp.PlainAuth("", from, password, server.config.Smtp_host)
	err := smtp.SendMail(server.config.Smtp_host+":"+server.config.Smtp_port, auth, from, receiver, []byte(msg))
	if err != nil {
		return err
	}
	return nil
}

func (server *Server) SendOrderMail(to string, message string) error {
	from := server.config.Smtp_login
	password := server.config.Smtp_password
	receiver := []string{to}
	msg := "To: " + to + "\n" + "Subject: FilinShop Order Information\n" +
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
		message
	auth := smtp.PlainAuth("", from, password, server.config.Smtp_host)
	err := smtp.SendMail(server.config.Smtp_host+":"+server.config.Smtp_port, auth, from, receiver, []byte(msg))
	if err != nil {
		return err
	}
	return nil
}

func (server *Server) WEBauth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(COOKIE_NAME)
		if err != nil {
			fmt.Println(1)
			next(w, r)
			return
		}
		session, err := server.sessionStore.Get(cookie.Value, false)
		if err != nil {
			fmt.Println(2)
			next(w, r)
			return
		}
		user, err := server.repository.Users().FindByID(session.UserID)
		if err != nil {
			fmt.Println(3)
			next(w, r)
			return
		}
		next(w, r.WithContext(context.WithValue(r.Context(), ctxUserKey, user)))
	})
}

func (server *Server) APIauth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(COOKIE_ADMIN)
		if err != nil {
			fmt.Println(1)
			next(w, r)
			return
		}
		session, err := server.sessionStore.Get(cookie.Value, true)
		if err != nil {
			fmt.Println(2)
			next(w, r)
			return
		}
		admin, err := server.repository.Users().FindAdminByID(session.UserID)
		if err != nil {
			fmt.Println(3)
			next(w, r)
			return
		}
		next(w, r.WithContext(context.WithValue(r.Context(), ctxUserKey, admin)))
	})
}

// Web site Handler Functions

func (server *Server) handleAdminChangePSWD() http.HandlerFunc {
	type TemplateData struct {
		Login string
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context().Value(ctxUserKey)

		if ctx == nil {
			http.Redirect(w, r, "/admin/login", http.StatusTemporaryRedirect)
			return
		}
		admin := ctx.(*models.Admin)
		data := TemplateData{}
		data.Login = admin.Login
		tmpl, err := template.ParseFiles("./web/templates/admin/changepswd.html",
			"./web/templates/admin/header.html",
			"./web/templates/admin/auth.html",
		)
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
		err = tmpl.ExecuteTemplate(w, "changepswd", data)
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
	}
}

func (server *Server) handleAdminBookList() http.HandlerFunc {
	type TemplateData struct {
		Login  string
		Books  []models.Book
		Genres []models.Genre
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context().Value(ctxUserKey)
		if ctx == nil {
			http.Redirect(w, r, "/admin/login", http.StatusTemporaryRedirect)
			return
		}
		admin := ctx.(*models.Admin)
		data := TemplateData{}
		data.Login = admin.Login

		books, err := server.repository.Books().GetAll()
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
		genres, err := server.repository.Books().GetAllGenres()
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}

		data.Books = *books
		data.Genres = genres

		tmpl, err := template.ParseFiles("./web/templates/admin/book_list.html",
			"./web/templates/admin/header.html",
			"./web/templates/admin/auth.html",
		)
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
		err = tmpl.ExecuteTemplate(w, "book_list", data)
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
	}
}

func (server *Server) handleCreateAdmin() http.HandlerFunc {
	type TemplateData struct {
		Login string
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context().Value(ctxUserKey)
		if ctx == nil {
			http.Redirect(w, r, "/admin/login", http.StatusTemporaryRedirect)
			return
		}
		admin := ctx.(*models.Admin)
		data := TemplateData{
			Login: admin.Login,
		}
		tmpl, err := template.ParseFiles("./web/templates/admin/create_admin.html",
			"./web/templates/admin/header.html",
			"./web/templates/admin/auth.html",
		)
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
		err = tmpl.ExecuteTemplate(w, "create_admin", data)
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
	}
}

func (server *Server) handleAdminList() http.HandlerFunc {
	type TemplateData struct {
		Login  string
		Admins []models.Admin
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context().Value(ctxUserKey)
		if ctx == nil {
			http.Redirect(w, r, "/admin/login", http.StatusTemporaryRedirect)
			return
		}
		admin := ctx.(*models.Admin)
		data := TemplateData{
			Login: admin.Login,
		}
		admins, err := server.repository.Users().GetAllAdmins()
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
		data.Admins = *admins
		tmpl, err := template.ParseFiles("./web/templates/admin/admin_list.html",
			"./web/templates/admin/header.html",
			"./web/templates/admin/auth.html",
		)
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
		err = tmpl.ExecuteTemplate(w, "admin_list", data)
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
	}
}

func (server *Server) handleAdminHome() http.HandlerFunc {
	type TemplateData struct {
		Login string
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context().Value(ctxUserKey)
		if ctx == nil {
			http.Redirect(w, r, "/admin/login", http.StatusTemporaryRedirect)
			return
		}
		admin := ctx.(*models.Admin)
		data := TemplateData{
			Login: admin.Login,
		}
		tmpl, err := template.ParseFiles("./web/templates/admin/home.html",
			"./web/templates/admin/header.html",
			"./web/templates/admin/auth.html",
		)
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
		err = tmpl.ExecuteTemplate(w, "home", data)
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
	}
}

func (server *Server) handleAdminLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context().Value(ctxUserKey)
		if ctx == nil {
			tmpl, err := template.ParseFiles("./web/templates/admin/login.html")
			if err != nil {
				server.error(w, r, http.StatusInternalServerError, err)
				return
			}
			err = tmpl.ExecuteTemplate(w, "login", nil)
			if err != nil {
				server.error(w, r, http.StatusInternalServerError, err)
				return
			}
			return
		}
		http.Redirect(w, r, "/admin", http.StatusTemporaryRedirect)
	}
}

func (server *Server) handleOrder() http.HandlerFunc {
	type TemplateData struct {
		Email string
		Cart  []models.Book
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context().Value(ctxUserKey)

		if ctx == nil {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
		user := ctx.(*models.User)
		data := TemplateData{}

		cart, err := server.repository.Users().Cart(user.ID)
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
		books := []models.Book{}
		for _, value := range *cart {
			book, err := server.repository.Books().FindBookByID(value)
			if err != nil {
				server.error(w, r, http.StatusInternalServerError, err)
				return
			}
			books = append(books, *book)
		}

		data.Cart = books

		data.Email = user.Email
		tmpl, err := template.ParseFiles("./web/templates/site/order.html",
			"./web/templates/site/header.html",
			"./web/templates/site/auth.html",
			"./web/templates/site/footer.html",
		)
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
		err = tmpl.ExecuteTemplate(w, "order", data)
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
	}
}

func (server *Server) handleChangePSWD() http.HandlerFunc {
	type TemplateData struct {
		Email string
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context().Value(ctxUserKey)

		if ctx == nil {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
		user := ctx.(*models.User)
		data := TemplateData{}
		data.Email = user.Email
		tmpl, err := template.ParseFiles("./web/templates/site/changepswd.html",
			"./web/templates/site/header.html",
			"./web/templates/site/auth.html",
			"./web/templates/site/footer.html",
		)
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
		err = tmpl.ExecuteTemplate(w, "changepswd", data)
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
	}
}

func (server *Server) handleCart() http.HandlerFunc {
	type TemplateData struct {
		Email string
		Books []models.Book
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context().Value(ctxUserKey)

		if ctx == nil {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
		user := ctx.(*models.User)
		data := TemplateData{}

		cart, err := server.repository.Users().Cart(user.ID)
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
		books := []models.Book{}
		for _, value := range *cart {
			book, err := server.repository.Books().FindBookByID(value)
			if err != nil {
				server.error(w, r, http.StatusInternalServerError, err)
				return
			}
			books = append(books, *book)
		}

		data.Books = books

		data.Email = user.Email
		tmpl, err := template.ParseFiles("./web/templates/site/cart.html",
			"./web/templates/site/header.html",
			"./web/templates/site/auth.html",
			"./web/templates/site/footer.html",
		)
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
		err = tmpl.ExecuteTemplate(w, "cart", data)
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
	}
}

func (server *Server) handleCatalog() http.HandlerFunc {
	type TemplateData struct {
		Email  string
		Books  []models.Book
		Genres []models.Genre
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context().Value(ctxUserKey)
		data := TemplateData{}

		books, err := server.repository.Books().GetAll()
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
		}

		genres, err := server.repository.Books().GetAllGenres()
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
		}

		data.Books = *books
		data.Genres = genres

		if ctx == nil {
			fmt.Println("contextdsfsdfsdfsf")
			tmpl, err := template.ParseFiles("./web/templates/site/catalog.html",
				"./web/templates/site/header.html",
				"./web/templates/site/non.auth.html",
				"./web/templates/site/footer.html",
			)
			if err != nil {
				server.error(w, r, http.StatusInternalServerError, err)
				return
			}
			err = tmpl.ExecuteTemplate(w, "catalog", data)
			if err != nil {
				server.error(w, r, http.StatusInternalServerError, err)
				return
			}
			return
		}
		user := ctx.(*models.User)
		data.Email = user.Email
		tmpl, err := template.ParseFiles("./web/templates/site/catalog.html",
			"./web/templates/site/header.html",
			"./web/templates/site/auth.html",
			"./web/templates/site/footer.html",
		)
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
		err = tmpl.ExecuteTemplate(w, "catalog", data)
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
	}
}

func (server *Server) handleRegister() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./web/templates/site/register.html",
			"./web/templates/site/header.html",
			"./web/templates/site/non.auth.html",
			"./web/templates/site/footer.html",
		)
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
		err = tmpl.ExecuteTemplate(w, "reg", nil)
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
	}
}

func (server *Server) handleAdminUpdateBook() http.HandlerFunc {
	type TemplateData struct {
		Login  string
		Genres []models.Genre
		Book   models.Book
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context().Value(ctxUserKey)
		if ctx == nil {
			http.Redirect(w, r, "/admin/login", http.StatusTemporaryRedirect)
			return
		}
		admin := ctx.(*models.Admin)
		vars := mux.Vars(r)
		id := vars["id"]

		tmpl, err := template.ParseFiles("./web/templates/admin/update_book.html",
			"./web/templates/admin/header.html",
			"./web/templates/admin/auth.html",
		)
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
		book, err := server.repository.Books().FindBookByID(id)
		if err != nil {
			server.error(w, r, http.StatusBadRequest, err)
			return
		}

		genres, err := server.repository.Books().GetAllGenres()
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}

		data := TemplateData{}
		data.Login = admin.Login
		data.Genres = genres
		data.Book = *book

		err = tmpl.ExecuteTemplate(w, "update_book", data)
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
	}
}

func (server *Server) handleAdminCreateBook() http.HandlerFunc {
	type TemplateData struct {
		Login  string
		Genres []models.Genre
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context().Value(ctxUserKey)
		if ctx == nil {
			http.Redirect(w, r, "/admin/login", http.StatusTemporaryRedirect)
			return
		}
		admin := ctx.(*models.Admin)
		tmpl, err := template.ParseFiles("./web/templates/admin/add_book.html",
			"./web/templates/admin/header.html",
			"./web/templates/admin/auth.html",
		)
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
		genres, err := server.repository.Books().GetAllGenres()
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
		data := TemplateData{
			Login:  admin.Login,
			Genres: genres,
		}
		err = tmpl.ExecuteTemplate(w, "add_book", data)
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
	}
}

func (server *Server) handleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context().Value(ctxUserKey)
		if ctx == nil {
			tmpl, err := template.ParseFiles("./web/templates/site/login.html",
				"./web/templates/site/header.html",
				"./web/templates/site/non.auth.html",
				"./web/templates/site/footer.html",
			)
			if err != nil {
				server.error(w, r, http.StatusInternalServerError, err)
				return
			}
			err = tmpl.ExecuteTemplate(w, "login", nil)
			if err != nil {
				server.error(w, r, http.StatusInternalServerError, err)
				return
			}
			return
		}
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
}

func (server *Server) handleHome() http.HandlerFunc {
	type TemplateData struct {
		Email        string
		BooksOfMonth []models.Book
		BookLiders   []models.Book
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context().Value(ctxUserKey)
		data := TemplateData{}

		booksOfMonthID, err := server.repository.Books().BooksOfMonth()
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
		bookLidersID, err := server.repository.Books().BookLiders()
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}

		booksOfMonth := []models.Book{}
		bookLiders := []models.Book{}

		for _, book := range *booksOfMonthID {
			book, err := server.repository.Books().FindBookByID(book)
			if err != nil {
				server.error(w, r, http.StatusInternalServerError, err)
				return
			}
			booksOfMonth = append(booksOfMonth, *book)
		}

		for _, book := range *bookLidersID {
			book, err := server.repository.Books().FindBookByID(book)
			if err != nil {
				server.error(w, r, http.StatusInternalServerError, err)
				return
			}
			bookLiders = append(bookLiders, *book)
		}

		data.BooksOfMonth = booksOfMonth
		data.BookLiders = bookLiders

		if ctx == nil {
			tmpl, err := template.ParseFiles("./web/templates/site/home.html",
				"./web/templates/site/header.html",
				"./web/templates/site/non.auth.html",
				"./web/templates/site/footer.html",
			)
			if err != nil {
				server.error(w, r, http.StatusInternalServerError, err)
				return
			}
			err = tmpl.ExecuteTemplate(w, "home", data)
			if err != nil {
				server.error(w, r, http.StatusInternalServerError, err)
				return
			}
			return
		}
		user := ctx.(*models.User)
		data.Email = user.Email
		tmpl, err := template.ParseFiles("./web/templates/site/home.html",
			"./web/templates/site/header.html",
			"./web/templates/site/auth.html",
			"./web/templates/site/footer.html",
		)
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
		err = tmpl.ExecuteTemplate(w, "home", data)
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
	}
}

// API Handler Functions

func (server *Server) APICleanCart() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context().Value(ctxUserKey)
		if ctx == nil {
			server.error(w, r, http.StatusUnauthorized, ErrNullContext)
			return
		}
		user := ctx.(*models.User)

		err := server.repository.Users().ClearCart(user.ID)
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
		server.respond(w, r, http.StatusOK, nil)
	}
}

func (server *Server) APIIsLoggedIn() http.HandlerFunc {
	type data struct {
		IsLoggedIn string `json:"isLoggedIn"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		data := data{}
		cookie, err := r.Cookie(COOKIE_NAME)
		if err != nil {
			data.IsLoggedIn = "false"
			server.respond(w, r, http.StatusOK, &data)
			return
		}
		session, err := server.sessionStore.Get(cookie.Value, false)
		if err != nil {
			data.IsLoggedIn = "false"
			server.respond(w, r, http.StatusOK, &data)
			return
		}
		_, err = server.repository.Users().FindByID(session.UserID)
		if err != nil {
			data.IsLoggedIn = "false"
			server.respond(w, r, http.StatusOK, &data)
			return
		}
		data.IsLoggedIn = "true"
		server.respond(w, r, http.StatusOK, &data)
	}
}

func (server *Server) APIActivate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["uid"]
		token := vars["token"]

		user, err := server.repository.Users().FindByID(id)
		if err != nil {
			server.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		if user.Activated == 1 {
			server.error(w, r, http.StatusBadRequest, ErrActivated)
			return
		}

		if user.Token == token {
			err = server.repository.Users().Activate(user.ID)
			if err != nil {
				server.error(w, r, http.StatusInternalServerError, err)
				return
			}
			server.respond(w, r, http.StatusOK, nil)
			return
		}
		server.error(w, r, http.StatusUnprocessableEntity, ErrInvalidToken)
	}
}

func (server *Server) APIGetAllAdmins() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context().Value(ctxUserKey)
		if ctx == nil {
			server.error(w, r, http.StatusUnauthorized, ErrNullContext)
			return
		}
		admins, err := server.repository.Users().GetAllAdmins()
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
		server.respond(w, r, http.StatusOK, admins)
	}
}

func (server *Server) APIDeleteAdmin() http.HandlerFunc {
	type request struct {
		PrimaryKey string `json:"primary_key"`
		AdminID    string `json:"admin_id"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := request{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			server.error(w, r, http.StatusBadRequest, err)
			return
		}

		if req.PrimaryKey != server.config.PrimayKey {
			server.error(w, r, http.StatusBadRequest, ErrWrongPrimaryKey)
			return
		}

		admin, err := server.repository.Users().FindAdminByID(req.AdminID)
		if err != nil {
			server.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		err = server.repository.Users().DeleteAdmin(admin.ID)
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}

		server.respond(w, r, http.StatusOK, nil)
	}
}

func (server *Server) APICreateAdmin() http.HandlerFunc {
	type request struct {
		PrimaryKey string `json:"primary_key"`
		Login      string `json:"login"`
		Password   string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := request{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			server.error(w, r, http.StatusBadRequest, err)
			return
		}

		if req.PrimaryKey != server.config.PrimayKey {
			server.error(w, r, http.StatusBadRequest, ErrWrongPrimaryKey)
			return
		}

		a := &models.Admin{
			Login:    req.Login,
			Password: req.Password,
		}

		if err := a.BeforeCreate(); err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}

		admin_ids, err := server.repository.Users().GetAllAdminID()
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}

		id := server.repository.Users().GenerateID(&admin_ids)

		a.ID = id

		err = server.repository.Users().CreateAdmin(a)
		if err != nil {
			server.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		server.respond(w, r, http.StatusCreated, a)
	}
}

func (server *Server) APIGetAllOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context().Value(ctxUserKey)
		if ctx == nil {
			server.error(w, r, http.StatusUnauthorized, ErrNullContext)
			return
		}
		orders, err := server.repository.Users().GetOrders()
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
		for i := 0; i < len(*orders); i++ {
			for j := 0; j < len((*orders)[i].Products); j++ {
				book, err := server.repository.Books().FindBookByID((*orders)[i].Products[j].ID)
				if err != nil {
					server.error(w, r, http.StatusInternalServerError, err)
					return
				}
				(*orders)[i].Products[j].Name = book.Name
				(*orders)[i].Products[j].Price = book.Price
			}
		}

		server.respond(w, r, http.StatusOK, *orders)
	}
}

func (server *Server) APIChangeAdminPassword() http.HandlerFunc {
	type request struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context().Value(ctxUserKey)
		if ctx == nil {
			server.error(w, r, http.StatusUnauthorized, ErrNullContext)
			return
		}
		admin := ctx.(*models.Admin)
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			server.error(w, r, http.StatusBadRequest, err)
			return
		}
		if !admin.ComparePassword(req.OldPassword) {
			server.error(w, r, http.StatusBadRequest, ErrWrongPassword)
			return
		}

		err := server.repository.Users().ChangeAdminPassword(admin.ID, req.NewPassword)
		if err != nil {
			server.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		server.respond(w, r, http.StatusOK, nil)
	}
}

func (server *Server) APIChangePassword() http.HandlerFunc {
	type request struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context().Value(ctxUserKey)
		if ctx == nil {
			server.error(w, r, http.StatusUnauthorized, ErrNullContext)
			return
		}
		user := ctx.(*models.User)
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			server.error(w, r, http.StatusBadRequest, err)
			return
		}
		if !user.ComparePassword(req.OldPassword) {
			server.error(w, r, http.StatusBadRequest, ErrWrongPassword)
			return
		}

		err := server.repository.Users().ChangePassword(user.ID, req.NewPassword)
		if err != nil {
			server.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		server.respond(w, r, http.StatusOK, nil)
	}
}

func (server *Server) APICreateOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context().Value(ctxUserKey)
		if ctx == nil {
			server.error(w, r, http.StatusUnauthorized, ErrNullContext)
			return
		}
		user := ctx.(*models.User)

		if user.Activated == 0 {
			server.error(w, r, http.StatusBadRequest, ErrNotActivated)
			return
		}

		req := &models.Order{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			server.error(w, r, http.StatusBadRequest, err)
			return
		}

		if len(req.Products) == 0 {
			server.error(w, r, http.StatusBadRequest, ErrNullProducts)
			return
		}

		var orderSumm int = 0

		for i := 0; i < len((*req).Products); i++ {
			book, err := server.repository.Books().FindBookByID((*req).Products[i].ID)
			if err != nil {
				server.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}
			(*req).Products[i].Name = book.Name
			(*req).Products[i].Price = book.Price

			orderSumm += book.Price * (*req).Products[i].Quantity
		}

		t := time.Now()
		time_str := fmt.Sprintf("%02d-%02d-%d %02d:%02d:%02d",
			t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute(), t.Second(),
		)

		id_array, err := server.repository.Users().GetAllOrderID()
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
		order_id := server.repository.Users().GenerateID(&id_array)

		req.ID = order_id
		req.UserID = user.ID
		req.Date = time_str
		req.OrderSumm = orderSumm

		fmt.Printf("%+v\n", *req)

		err = server.repository.Users().CreateOrder(req)
		if err != nil {
			server.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		err = server.repository.Users().ClearCart(user.ID)
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}

		tmpl, err := template.ParseFiles("./web/templates/site/order_email.html")
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
		buf := new(bytes.Buffer)
		if err = tmpl.Execute(buf, *req); err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}

		err = server.SendOrderMail(user.Email, buf.String())
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
		server.respond(w, r, http.StatusOK, nil)
	}
}

func (server *Server) APIDeleteFromCart() http.HandlerFunc {
	type request struct {
		BookID string `json:"book_id"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context().Value(ctxUserKey)
		if ctx == nil {
			server.error(w, r, http.StatusUnauthorized, ErrNullContext)
			return
		}
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			server.error(w, r, http.StatusBadRequest, err)
			return
		}
		user := ctx.(*models.User)
		_, err := server.repository.Books().FindBookByID(req.BookID)
		if err != nil {
			server.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		err = server.repository.Users().DeleteFromCart(user.ID, req.BookID)
		if err != nil {
			server.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		server.respond(w, r, http.StatusOK, nil)
	}
}

func (server *Server) APIAddToCart() http.HandlerFunc {
	type request struct {
		BookID string `json:"book_id"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context().Value(ctxUserKey)
		if ctx == nil {
			server.error(w, r, http.StatusUnauthorized, ErrNullContext)
			return
		}
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			server.error(w, r, http.StatusBadRequest, err)
			return
		}
		user := ctx.(*models.User)
		_, err := server.repository.Books().FindBookByID(req.BookID)
		if err != nil {
			server.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		cart, err := server.repository.Users().Cart(user.ID)
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}

		for _, value := range *cart {
			if value == req.BookID {
				server.error(w, r, http.StatusUnprocessableEntity, ErrExistInCart)
				return
			}
		}

		err = server.repository.Users().AddToCart(user.ID, req.BookID)
		if err != nil {
			server.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		server.respond(w, r, http.StatusCreated, nil)
	}
}

func (server *Server) APIGetCart() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context().Value(ctxUserKey)
		if ctx == nil {
			server.error(w, r, http.StatusUnauthorized, ErrNullContext)
			return
		}
		user := ctx.(*models.User)
		cart, err := server.repository.Users().Cart(user.ID)
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
		books := []models.Book{}
		for i := 0; i < len(*cart); i++ {
			book_id := (*cart)[i]
			book, err := server.repository.Books().FindBookByID(book_id)
			if err != nil {
				server.error(w, r, http.StatusInternalServerError, err)
				return
			}
			books = append(books, *book)
		}

		server.respond(w, r, http.StatusOK, &books)
	}
}

func (server *Server) APIGetAllBooks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		books, err := server.repository.Books().GetAll()
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}

		server.respond(w, r, http.StatusOK, books)
	}
}

func (server *Server) APIDeleteBook() http.HandlerFunc {
	type request struct {
		BookID string `json:"book_id"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context().Value(ctxUserKey)
		if ctx == nil {
			server.error(w, r, http.StatusUnauthorized, ErrNullContext)
			return
		}
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			server.error(w, r, http.StatusBadRequest, err)
			return
		}
		book, err := server.repository.Books().FindBookByID(req.BookID)
		if err != nil {
			server.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		err = os.Remove("./web" + book.CoverPath)
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
		err = server.repository.Books().Delete(req.BookID)
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
		server.respond(w, r, http.StatusOK, nil)
	}
}

func (server *Server) APICreateBook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context().Value(ctxUserKey)
		if ctx == nil {
			server.error(w, r, http.StatusUnauthorized, ErrNullContext)
			return
		}
		err := r.ParseMultipartForm(30 * 1024 * 1024)
		if err != nil {
			server.error(w, r, http.StatusBadRequest, err)
			return
		}

		book := models.Book{}
		book_json := r.FormValue("json")
		err = json.Unmarshal([]byte(book_json), &book)
		if err != nil {
			server.error(w, r, http.StatusBadRequest, err)
			return
		}

		err = book.BeforeCreate()
		if err != nil {
			server.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		id_array, err := server.repository.Books().GetAllBookID()
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
		id := server.repository.Books().GenerateBookID(&id_array)
		book.ID = id

		cover, file_header, err := r.FormFile("cover")
		if err != nil {
			book.CoverPath = ""
			err = server.repository.Books().Create(&book)
			if err != nil {
				server.error(w, r, http.StatusInternalServerError, err)
				return
			}
			server.respond(w, r, http.StatusCreated, nil)
			return
		}
		defer cover.Close()
		tempfile, err := os.Create("./web/static/img/covers/" + id + filepath.Ext(file_header.Filename))
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
		defer tempfile.Close()
		coverBytes, err := ioutil.ReadAll(cover)
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
		_, err = tempfile.Write(coverBytes)
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
		web_path := server.config.WebCoverPath + id + filepath.Ext(file_header.Filename)
		book.CoverPath = web_path

		err = server.repository.Books().Create(&book)
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
		server.respond(w, r, http.StatusCreated, nil)
	}
}

func (server *Server) APIUpdateBook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context().Value(ctxUserKey)
		if ctx == nil {
			server.error(w, r, http.StatusUnauthorized, ErrNullContext)
			return
		}
		err := r.ParseMultipartForm(30 * 1024 * 1024)
		if err != nil {
			server.error(w, r, http.StatusBadRequest, err)
			return
		}

		book := models.Book{}
		book_json := r.FormValue("json")
		err = json.Unmarshal([]byte(book_json), &book)
		if err != nil {
			server.error(w, r, http.StatusBadRequest, err)
			return
		}

		err = book.BeforeUpdate()
		if err != nil {
			server.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		old_book, err := server.repository.Books().FindBookByID(book.ID)
		if err != nil {
			server.error(w, r, http.StatusBadRequest, err)
			return
		}

		old_book.Name = book.Name
		old_book.Author = book.Author
		old_book.Year = book.Year
		old_book.Genre = book.Genre
		old_book.PubHouse = book.PubHouse
		old_book.Description = book.Description
		old_book.Price = book.Price

		cover, file_header, err := r.FormFile("cover")
		if err != nil {
			err = server.repository.Books().Update(old_book)
			if err != nil {
				server.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}
			server.respond(w, r, http.StatusOK, nil)
			return
		}
		defer cover.Close()

		err = os.Remove("./web" + old_book.CoverPath)
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}

		tempfile, err := os.Create("./web/static/img/covers/" + old_book.ID + filepath.Ext(file_header.Filename))
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
		defer tempfile.Close()
		coverBytes, err := ioutil.ReadAll(cover)
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
		_, err = tempfile.Write(coverBytes)
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
		web_path := server.config.WebCoverPath + old_book.ID + filepath.Ext(file_header.Filename)
		old_book.CoverPath = web_path

		err = server.repository.Books().Update(old_book)
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
		server.respond(w, r, http.StatusOK, nil)
	}
}

func (server *Server) APIGetAllGenres() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		genres, err := server.repository.Books().GetAllGenres()
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		server.respond(w, r, http.StatusOK, genres)
	}
}

func (server *Server) APICreateUser() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			server.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := &models.User{
			Email:     req.Email,
			Password:  req.Password,
			Activated: 0,
		}

		if err := u.BeforeCreate(); err != nil {
			server.error(w, r, http.StatusBadRequest, err)
			return
		}

		user_ids, err := server.repository.Users().GetAllUserID()
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
		tokens, err := server.repository.Users().GetAllTokens()
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}

		u.ID = server.repository.Users().GenerateID(&user_ids)
		u.Token = server.repository.Users().GenerateToken(&tokens)

		if err := server.repository.Users().Create(u); err != nil {
			server.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		confirmURL := "https://filin-shop.herokuapp.com" + "/api/public/activate/" + u.ID + "/" + u.Token

		tmpl, err := template.ParseFiles("./web/templates/site/confirm_email.html")
		if err != nil {
			server.error(w, r, http.StatusCreated, nil)
			return
		}
		buf := new(bytes.Buffer)
		if err = tmpl.Execute(buf, confirmURL); err != nil {
			server.error(w, r, http.StatusCreated, nil)
			return
		}

		err = server.SendConfirmMail(u.Email, buf.String())
		if err != nil {
			server.error(w, r, http.StatusCreated, nil)
			return
		}

		server.respond(w, r, http.StatusCreated, nil)
	}
}

func (server *Server) APICreateAdminSession() http.HandlerFunc {
	type request struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context().Value(ctxUserKey)
		if ctx == nil {
			req := &request{}
			if err := json.NewDecoder(r.Body).Decode(req); err != nil {
				server.error(w, r, http.StatusBadRequest, err)
				return
			}

			admin, err := server.repository.Users().FindAdminByLogin(req.Login)
			if err != nil || !admin.ComparePassword(req.Password) {
				server.error(w, r, http.StatusUnauthorized, ErrIncorrectEmailOrPassword)
				return
			}

			sessionID := server.sessionStore.Create(admin.ID, true)

			cookie := &http.Cookie{
				Name:     COOKIE_ADMIN,
				Value:    sessionID,
				Path:     "/",
				Expires:  time.Now().Add(time.Hour * 48),
				Secure:   false,
				SameSite: 0,
			}
			http.SetCookie(w, cookie)
			return
		}
		server.error(w, r, http.StatusOK, ErrAuthenticated)
	}
}

func (server *Server) APICreateSession() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context().Value(ctxUserKey)
		if ctx == nil {
			req := &request{}
			if err := json.NewDecoder(r.Body).Decode(req); err != nil {
				server.error(w, r, http.StatusBadRequest, err)
				return
			}

			u, err := server.repository.Users().FindByEmail(req.Email)
			if err != nil || !u.ComparePassword(req.Password) {
				server.error(w, r, http.StatusUnauthorized, ErrIncorrectEmailOrPassword)
				return
			}

			sessionID := server.sessionStore.Create(u.ID, false)

			cookie := &http.Cookie{
				Name:     COOKIE_NAME,
				Value:    sessionID,
				Path:     "/",
				Expires:  time.Now().Add(time.Hour * 48),
				Secure:   false,
				SameSite: 0,
			}
			http.SetCookie(w, cookie)
			return
		}
		server.error(w, r, http.StatusOK, ErrAuthenticated)
	}
}

func (server *Server) APIDeleteAdminSession() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context().Value(ctxUserKey)
		if ctx == nil {
			server.error(w, r, http.StatusUnauthorized, ErrNullContext)
			return
		}
		cookie, err := r.Cookie(COOKIE_ADMIN)
		if err != nil {
			server.error(w, r, http.StatusUnauthorized, err)
			return
		}
		server.sessionStore.Delete(cookie.Value, true)
		cookie.Expires = time.Now()
		http.SetCookie(w, cookie)
		server.respond(w, r, http.StatusOK, nil)
	}
}

func (server *Server) APIDeleteSession() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context().Value(ctxUserKey)
		if ctx == nil {
			server.error(w, r, http.StatusUnauthorized, ErrNullContext)
			return
		}
		cookie, err := r.Cookie(COOKIE_NAME)
		if err != nil {
			server.error(w, r, http.StatusUnauthorized, err)
			return
		}
		server.sessionStore.Delete(cookie.Value, false)
		cookie.Expires = time.Now()
		http.SetCookie(w, cookie)
		server.respond(w, r, http.StatusOK, nil)
	}
}

func (server *Server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	fmt.Println(err.Error())
	server.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (server *Server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
