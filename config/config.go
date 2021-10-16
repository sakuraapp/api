package config

import "crypto/rsa"

type Config struct {
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