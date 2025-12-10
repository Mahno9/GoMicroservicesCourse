package config

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type GrpcConfig interface {
	Port() string
}
