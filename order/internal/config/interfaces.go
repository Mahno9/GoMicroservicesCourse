package config

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type PostgresConfig interface {
	URI() string
	MigrationsDir() string
}

type HttpConfig interface {
	Host() string
	Port() string
	Address() string
}

type ClientsConfig interface {
	InventoryAddress() string
	PaymentAddress() string
}
