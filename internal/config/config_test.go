package config

import (
	"testing"

	"github.com/joho/godotenv"
)

func TestConfig(t *testing.T) {
	err := godotenv.Load(".env")
	if err != nil {
		t.Fatalf("Error loading .env file")
	}
	want := Config {
		DBName: "postgres",
		DBUser: "postgres",
		DBPassword: "postgres",
		DBPort: "5432",
		DBHost: "testhost",
	}
	conf := NewConfig()
	if want.DBHost != conf.DBHost && want.DBPort != conf.DBPort && 
	want.DBPassword != conf.DBPassword && want.DBUser != conf.DBUser && want.DBName != conf.DBName {
		t.Fatalf("Values doesn't match got\nPassword: %s\nUser: %s\nName: %s\nHost: %s\nPort: %s", conf.DBPassword, conf.DBUser,
	conf.DBName, conf.DBHost, conf.DBPort)
	}
}
