package model

import (
	"time"
)

type Category string

const (
	CategoryUnknown  Category = "UNKNOWN"
	CategoryEngine   Category = "ENGINE"
	CategoryFuel     Category = "FUEL"
	CategoryPorthole Category = "PORTHOLE"
	CategoryWing     Category = "WING"
)

type Dimensions struct {
	Length float64
	Width  float64
	Height float64
	Weight float64
}

type Manufacturer struct {
	Name    string
	Country string
	Website string
}

type Part struct {
	Uuid          string         `bson:"uuid"`
	Name          string         `bson:"name"`
	Description   string         `bson:"description"`
	Price         float64        `bson:"price"`
	StockQuantity int64          `bson:"stock_quantity"`
	Category      Category       `bson:"category,omitempty"`
	Dimensions    *Dimensions    `bson:"dimensions,omitempty"`
	Manufacturer  *Manufacturer  `bson:"manufacturer,omitempty"`
	Tags          []string       `bson:"tags,omitempty"`
	Metadata      map[string]any `bson:"metadata,omitempty"`
	CreatedAt     time.Time      `bson:"created_at"`
	UpdatedAt     *time.Time     `bson:"updated_at,omitempty"`
}
