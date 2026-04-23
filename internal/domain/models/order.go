package models

type Order struct {
	ID       string `json:"ID"`
	Item     string `json:"Item"`
	Quantity int32  `json:"Quantity"`
}
