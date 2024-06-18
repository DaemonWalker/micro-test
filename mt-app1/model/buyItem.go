package model

type BuyItemModel struct {
	User   string `json:"user" validate:"required"`
	ItemId int    `json:"itemId" validate:"required"`
	Count  int    `json:"count" validate:"required"`
}
