package env

import (
	"github.com/caarlos0/env/v11"
)

type clientsEnvConfig struct {
	InventoryAddress string `env:"INVENTORY_SERVICE_ADDRESS,required"`
	PaymentAddress   string `env:"PAYMENT_SERVICE_ADDRESS,required"`
}

type clientsConfig struct {
	inventoryAddress string
	paymentAddress   string
}

func NewClientsConfig() (*clientsConfig, error) {
	var raw clientsEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &clientsConfig{
		inventoryAddress: raw.InventoryAddress,
		paymentAddress:   raw.PaymentAddress,
	}, nil
}

func (c *clientsConfig) InventoryAddress() string {
	return c.inventoryAddress
}

func (c *clientsConfig) PaymentAddress() string {
	return c.paymentAddress
}
