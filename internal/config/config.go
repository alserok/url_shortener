package config

import (
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"os"
	"strconv"
)

type Config struct {
	Env        string
	Port       string
	ServerType uint
	DBType     uint

	DB    DB
	Cache Cache
}

type DB struct {
	postgres struct {
		port     string
		user     string
		password string
		host     string
		name     string
	}
}

func (db *DB) PostgresDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", db.postgres.user, db.postgres.password, db.postgres.host, db.postgres.port, db.postgres.name)
}

type Cache struct {
	redis struct {
		host string
		port string
	}
}

func (cache *Cache) RedisDSN() string {
	return fmt.Sprintf("%s:%s", cache.redis.host, cache.redis.port)
}

const (
	defaultServerGRPCType = 0
	defaultDBPostgresType = 0

	defaultEnv  = "PROD"
	defaultPort = "5000"
)

func MustLoad() *Config {
	cfg := Config{
		Env:        defaultEnv,
		Port:       defaultPort,
		ServerType: defaultServerGRPCType,
		DBType:     defaultDBPostgresType,
	}

	cfg.parseMainConfig()
	cfg.parseDBConfig()
	cfg.parseCacheConfig()

	return &cfg
}

func (cfg *Config) parseMainConfig() {
	if port := os.Getenv("PORT"); port != "" {
		cfg.Port = port
	}

	if env := os.Getenv("ENV"); env != "" {
		cfg.Env = env
	}

	if adapterType, err := strconv.Atoi(os.Getenv("SERVER_TYPE")); err != nil {
		cfg.ServerType = uint(adapterType)
	}

	if dbType, err := strconv.Atoi(os.Getenv("DB_TYPE")); err != nil {
		cfg.DBType = uint(dbType)
	}
}

func (cfg *Config) parseDBConfig() {
	switch cfg.DBType {
	case defaultDBPostgresType:
		cfg.DB.postgres.host = os.Getenv("DB_HOST")
		cfg.DB.postgres.port = os.Getenv("DB_PORT")
		cfg.DB.postgres.user = os.Getenv("DB_USER")
		cfg.DB.postgres.password = os.Getenv("DB_PASSWORD")
		cfg.DB.postgres.name = os.Getenv("DB_NAME")
	}
}

func (cfg *Config) parseCacheConfig() {
	cfg.Cache.redis.host = os.Getenv("REDIS_HOST")
	cfg.Cache.redis.port = os.Getenv("REDIS_PORT")
}
