package JsonParsers

type GamesOffersJson struct {
	Id       int    `json:"id"`
	ParentId int    `json:"parent section id"`
	Name     string `json:"offer name"`
	Price    string `json:"price"`
	Currency string `json:"currency"`
	Discount string `json:"discount"`
	Gift     string `json:"gift"`
	IdSeller string `json:"id_seller"`
}