// Package model models that use in this project
package model

// Person : struct for user
type Person struct {
	ID           string `bson,json:"id"`
	Teams        string `bson,json:"teams"`
	xG           string   `bson,json:"xG"`
	Score        string  `bson,json:"score"`
	MVP          string `bson,json:"mvp"`
}

// Config struct create config
type Config struct {
	CurrentDB     string `env:"CURRENT_DB" envDefault:"postgres"`
	PostgresDBURL string `env:"POSTGRES_DB_URL"`
	MongoDBURL    string `env:"MONGO_DB_URL"`
	JwtKey        []byte `env:"JWT-KEY" `
}