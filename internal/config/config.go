package config

import "crypto/rsa"

type envType string

const (
	EnvDEV  envType = "DEV"
	EnvPROD envType = "PROD"
)

type Config struct {
	Env  envType
	Port int
	AllowedOrigins []string
	DatabaseUser string
	DatabasePassword string
	DatabaseName string
	SessionSecret string
	JWTPublicKey *rsa.PublicKey
	JWTPrivateKey *rsa.PrivateKey
	RedisAddr string
	RedisPassword string
	RedisDatabase int
	S3Region *string
	S3Bucket *string
	S3Endpoint *string
	S3ForcePathStyle *bool
	SupervisorAddr string
	SupervisorKeyPath string
}

func (c *Config) IsDev() bool {
	return c.Env == EnvDEV
}