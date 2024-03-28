package model


type Item struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	Amount      int64   `json:"amount"`
	Created_at  string  `json:"created_at"`
	Updated_at  string  `json:"updated_at"`
}

type Status struct {
	Success bool `json:"success"`
}

type ListItems struct {
	Count    int64   `json:"count"`
	Products []*Item `json:"products"`
}

type BoughtItemsList struct {
	Products []*Item `json:"products"`
}

type BuyItemRequest struct {
	UserId    string `json:"user_id"`
	ProductId string `json:"product_id"`
	Amount    int64  `json:"amount"`
}

type BuyItemResponse struct {
	Message     string `json:"message"`
	UserId      string `json:"user_id"`
	ProductId   string `json:"product_id"`
	ProductName string `json:"product_name"`
	Amount      int64  `json:"amount"`
}

type ItemAmountRequest struct {
	ProductId int64 `json:"product_id"`
	Amount    int64 `json:"amount"`
}

