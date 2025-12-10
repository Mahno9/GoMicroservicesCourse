package config

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type MongoConfig interface {
	URI() string
	DatabaseName() string
}

type GrpcConfig interface {
	Host() string
	Port() string
}
