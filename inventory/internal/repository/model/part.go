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
	Uuid          string         `bson:"uuid,omitempty"`
	Name          string         `bson:"name"`
	Description   string         `bson:"description"`
	Price         float64        `bson:"price"`
	StockQuantity int64          `bson:"stock_quantity"`
	Category      Category       `bson:"category"`
	Dimensions    *Dimensions    `bson:"dimensions"`
	Manufacturer  *Manufacturer  `bson:"manufacturer"`
	Tags          []string       `bson:"tags"`
	Metadata      map[string]any `bson:"metadata"`
	CreatedAt     time.Time      `bson:"created_at,omitempty"`
	UpdatedAt     *time.Time     `bson:"updated_at"`
}

type PartsFilter struct {
	Uuids                 []string
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
	Tags                  []string
}
