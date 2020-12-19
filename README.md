# Sakura API

## Installation
Install the dependencies
```
npm install
```
Copy the .env.sample file as .env
```
cp .env.sample .env
```
Create a discord application on Discord's [Developer Portal](https://discordapp.com/developers/applications), then fill in your information in the .env file
```
DISCORD_CLIENT_ID="application id"
DISCORD_CLIENT_SECRET="application secret"
DISCORD_SCOPES="identify"
DISCORD_REDIRECT_URI="http://your api host:8081/auth/login"
```

## Usage
To run in a development environment:
```
npm run dev
```
Note that the API requires the master server to be running