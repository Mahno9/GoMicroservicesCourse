package model

import (
	"time"

	"github.com/google/uuid"
)

type Part struct {
	Uuid          uuid.UUID
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Category      int32
	Dimensions    *Dimensions
	Manufacturer  *Manufacturer
	Tags          []string
	Metadata      map[string]*any
	CreatedAt     time.Time
	UpdatedAt     *time.Time
}

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

type PartsFilter struct {
	Uuids                 []uuid.UUID
	Names                 []string
	Categories            []int32
	ManufacturerCountries []string
	Tags                  []string
}
