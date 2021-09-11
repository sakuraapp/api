# Sakura API

## Installation
Install the dependencies
```
go get
```
Copy the .env.sample file as .env
```
cp .env.sample .env
```
Create a discord application on Discord's [Developer Portal](https://discordapp.com/developers/applications), then fill in your information in the .env file
```
DISCORD_KEY="application id"
DISCORD_SECRET="application secret"
DISCORD_SCOPES="identify"
DISCORD_REDIRECT="http://website url/auth/discord/callback"
```

## Usage
To run in a development environment:
```
go run main.gos
```