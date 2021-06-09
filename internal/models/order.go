package models

type Product struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Price    int    `json:"price"`
}

type Order struct {
	ID         string    `json:"id"`
	UserID     string    `json:"userID"`
	First_name string    `json:"first_name"`
	Last_name  string    `json:"last_name"`
	Patronymic string    `json:"patronymic"`
	Phone      string    `json:"phone"`
	City       string    `json:"city"`
	Address    string    `json:"address"`
	Post_index string    `json:"post_index"`
	Date       string    `json:"date"`
	Products   []Product `json:"products"`
	OrderSumm  int       `json:"order_summ"`
}
