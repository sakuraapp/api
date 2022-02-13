package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/discord"
	"github.com/sakuraapp/api/config"
	"github.com/sakuraapp/api/server"
	"github.com/sakuraapp/shared/pkg/crypto"
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

	allowedOrigins := strings.Split(strings.ToLower(os.Getenv("ALLOWED_ORIGINS")), ", ")

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

	jwtPrivateKey, err := crypto.LoadRSAPrivateKey(jwtPrivatePath, jwtPassphrase)

	if err != nil {
		log.WithError(err).Fatal("Failed to load private key")
	}

	jwtPublicKey, err := crypto.LoadRSAPublicKey(jwtPublicPath)

	if err != nil {
		log.WithError(err).Fatal("Failed to load public key")
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDatabase := os.Getenv("REDIS_DATABASE")
	redisDb, _ := strconv.Atoi(redisDatabase)

	s3Region := os.Getenv("S3_REGION")
	s3Bucket := os.Getenv("S3_BUCKET")
	s3Endpoint := os.Getenv("S3_ENDPOINT")
	s3ForcePathStyleStr := os.Getenv("S3_FORCE_PATH_STYLE")
	s3ForcePathStyle := false

	if s3ForcePathStyleStr == "1" {
		s3ForcePathStyle = true
	}

	server.Create(&config.Config{
		Env: envType,
		Port: port,
		AllowedOrigins: allowedOrigins,
		JWTPrivateKey: jwtPrivateKey,
		JWTPublicKey: jwtPublicKey,
		DatabaseUser: os.Getenv("DB_USER"),
		DatabasePassword: os.Getenv("DB_PASSWORD"),
		DatabaseName: os.Getenv("DB_DATABASE"),
		SessionSecret: os.Getenv("SESSION_SECRET"),
		RedisAddr: redisAddr,
		SupervisorAddr: os.Getenv("SUPERVISOR_ADDR"),
		SupervisorKeyPath: os.Getenv("SUPERVISOR_KEY"),
		RedisPassword: redisPassword,
		RedisDatabase: redisDb,
		S3Bucket: aws.String(s3Bucket),
		S3Region: aws.String(s3Region),
		S3Endpoint: aws.String(s3Endpoint),
		S3ForcePathStyle: &s3ForcePathStyle,
	})
}

func GetScopes(key string) []string {
	return strings.Split(os.Getenv(key), ", ")
}