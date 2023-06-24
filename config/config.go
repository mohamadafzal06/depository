package config

import (
	"os"
)

func getEnv(key string, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

var DatabaseUser = getEnv("DEPOSITORY_DATABASE_USER", "postgres")
var DatabasePass = getEnv("DEPOSITORY_DATABASE_PASS", "postgres")
var DatabaseAddress = getEnv("DEPOSITORY_DATABASE_ADDRESS", "127.0.0.1:5432")
var DatabaseDBName = getEnv("DEPOSITORY_DATABASE_DBNAME", "depository")
