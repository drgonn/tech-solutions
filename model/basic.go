package model

import "time"

// Outlet represents the different outlets of TechSolutions with specific information.
type Outlet struct {
	BaseModel
	Name     string `json:"name,omitempty" gorm:"type:varchar(255)"`
	Location string `json:"location,omitempty" gorm:"type:varchar(255)"`
}

// Product represents electronic devices, accessories, and store-specific gadgets for sale.
type Product struct {
	BaseModel
	Name     string  `json:"name,omitempty" gorm:"type:varchar(255)"`
	Price    float64 `json:"price,omitempty"`
	Stock    int     `json:"stock,omitempty"`
	OutletID uint    `json:"outlet_id,omitempty"`
}

// Customer represents customers with purchase history and GadgetPoints loyalty program.
type Customer struct {
	BaseModel
	Name               string    `json:"name,omitempty" gorm:"type:varchar(255)"`
	Email              string    `json:"email,omitempty" gorm:"type:varchar(255)"`
	GadgetPoints       int       `json:"gadget_points,omitempty"`
	SubscriptionExpiry time.Time `json:"subscription_expiry"`
}

// Purchase represents customer purchase orders.
type Purchase struct {
	BaseModel
	CustomerID   uint      `json:"customer_id,omitempty"`
	ProductID    uint      `json:"product_id,omitempty"`
	PurchaseDate time.Time `json:"purchase_date,omitempty"`
	Quantity     int       `json:"quantity,omitempty"`
}

// GadgetPoints represents the GadgetPoints loyalty program.
type GadgetPoints struct {
	BaseModel
	CustomerID   uint `json:"customer_id,omitempty"`
	PointsEarned int  `json:"points_earned,omitempty"`
}

type PointRedemption struct {
	BaseModel
	CustomerID     uint      `json:"-"`
	OutletID       uint      `json:"-"`
	ProductID      uint      `json:"-"`
	RedemptionDate time.Time `json:"redemption_date"`

	// When the company replenishes the Freebie for the Outlet, it is recorded as 'true.
	Replenished bool `json:"replenished" gorm:"default:false"`
}

type MonthlySubscription struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	CustomerID    uint      `json:"-"`
	OnlineStoreID uint      `json:"-"`
	StartAt       time.Time `json:"start_at"`
	EndAt         time.Time `json:"end_at"`
	NumOfMonths   int       `json:"num_of_months"`
	Paid          bool      `json:"paid" gorm:"default:false"`
	Actived       bool      `json:"actived" gorm:"default:false"`
}

type Discount struct {
	ID             uint    `json:"id" gorm:"primaryKey"`
	Type           string  `json:"type" gorm:"type:varchar(255)"`
	Description    string  `json:"description" gorm:"type:text"`
	DiscountAmount float64 `json:"discount_amount"`
	ProductID      uint    `json:"-"`
}
