package model

import "time"

type BaseModel struct {
	ID        string    `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type User struct {
	BaseModel
	UserCd string `json:"userCd"`
}

type Order struct {
	BaseModel
	Symbol      string        `json:"symbol"`
	Name        string        `json:"name"`
	StrikePrice float64       `json:"strikePrice"`
	OrderLines  []*OrderLines `json:"orderLines"`
}

type OrderLines struct {
	BaseModel
	OrderId     string  `json:"orderId"`
	StrikePrice float64 `json:"strikePrice"`
	OptionType  string  `json:"optionType" sql:"type:ENUM('CE', 'PE')"`
	LotQty      int     `json:"lotQty"`
	BuyAt       float64 `json:"buyAt"`
	SoldAt      float64 `json:"soldAt"`
	GrossProfit float64 `json:"grossProfit"`
	TaxAmt      float64 `json:"taxAmt"`
	NetProfit   float64 `json:"netProfit"`
}

type OrderHistory struct {
	BaseModel
	StrikePrice float64 `json:"strikePrice"`
	OptionType  string  `json:"optionType" sql:"type:ENUM('CE', 'PE')"`
	LotQty      int     `json:"lotQty"`
	BuyAt       float64 `json:"buyAt"`
	SoldAt      float64 `json:"soldAt"`
	GrossProfit float64 `json:"grossProfit"`
	TaxAmt      float64 `json:"taxAmt"`
	NetProfit   float64 `json:"netProfit"`
}
