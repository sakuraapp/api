package config

import "crypto/rsa"

type envType string

const (
	EnvDEV envType = "DEV"
	EnvPROD envType = "PROD"
)

type Config struct {
	Env envType
	Port string
	DatabaseUser string
	DatabasePassword string
	DatabaseName string
	JWTPublicKey *rsa.PublicKey
	JWTPrivateKey *rsa.PrivateKey
	RedisAddr string
	RedisPassword string
	RedisDatabase int
}

func (c *Config) IsDev() bool {
	return c.Env == EnvDEV
}