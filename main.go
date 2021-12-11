package main

import (
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/discord"
	"github.com/sakuraapp/api/config"
	"github.com/sakuraapp/api/server"
	shared "github.com/sakuraapp/shared/pkg"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "4000"
	}

	env := os.Getenv("APP_ENV")
	envType := config.EnvDEV

	if env == string(config.EnvPROD) {
		envType = config.EnvPROD
	}

	goth.UseProviders(
		discord.New(
			os.Getenv("DISCORD_KEY"),
			os.Getenv("DISCORD_SECRET"),
			os.Getenv("DISCORD_REDIRECT"),
			GetScopes("DISCORD_SCOPES")...
		),
	)

	jwtPublicPath := os.Getenv("JWT_PUBLIC_KEY")
	jwtPrivatePath := os.Getenv("JWT_PRIVATE_KEY")
	jwtPassphrase := os.Getenv("JWT_PASSPHRASE")

	jwtPrivateKey, err := shared.LoadRSAPrivateKey(jwtPrivatePath, jwtPassphrase)

	if err != nil {
		log.WithError(err).Fatal("Failed to load private key")
	}

	jwtPublicKey, err := shared.LoadRSAPublicKey(jwtPublicPath)

	if err != nil {
		log.WithError(err).Fatal("Failed to load public key")
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDatabase := os.Getenv("REDIS_DATABASE")
	redisDb, err := strconv.Atoi(redisDatabase)

	server.Create(config.Config{
		Env: envType,
		Port: port,
		JWTPrivateKey: jwtPrivateKey,
		JWTPublicKey: jwtPublicKey,
		DatabaseUser: os.Getenv("DB_USER"),
		DatabasePassword: os.Getenv("DB_PASSWORD"),
		DatabaseName: os.Getenv("DB_DATABASE"),
		RedisAddr: redisAddr,
		RedisPassword: redisPassword,
		RedisDatabase: redisDb,
	})
}

func GetScopes(key string) []string {
	return strings.Split(os.Getenv(key), ", ")
}