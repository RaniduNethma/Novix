package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	App       AppConfig
	Server    ServerConfig
	MinIO     MinIOConfig
	Redis     RedisConfig
	JWT       JWTConfig
	Streaming StreamingConfig
	Consul    ConsulConfig
}

type AppConfig struct {
	Env string
}

type ServerConfig struct {
	Port string
}

type MinIOConfig struct {
	Endpoint        string
	AccessKey       string
	SecretKey       string
	UseSSL          bool
	ProcessedBucket string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret string
}

type StreamingConfig struct {
	SessionTTL          int
	PresignTTL          int
	MaxSessionsPerUser  int
}

type ConsulConfig struct {
	Host        string
	Port        int
	ServiceName string
	ServicePort int
}

func Load() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("No .env file found, reading from environment: %v", err)
	}

	return &Config{
		App: AppConfig{
			Env: viper.GetString("APP_ENV"),
		},
		Server: ServerConfig{
			Port: viper.GetString("SERVER_PORT"),
		},
		MinIO: MinIOConfig{
			Endpoint:        viper.GetString("MINIO_ENDPOINT"),
			AccessKey:       viper.GetString("MINIO_ACCESS_KEY"),
			SecretKey:       viper.GetString("MINIO_SECRET_KEY"),
			UseSSL:          viper.GetBool("MINIO_USE_SSL"),
			ProcessedBucket: viper.GetString("MINIO_PROCESSED_BUCKET"),
		},
		Redis: RedisConfig{
			Host:     viper.GetString("REDIS_HOST"),
			Port:     viper.GetString("REDIS_PORT"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
		},
		JWT: JWTConfig{
			Secret: viper.GetString("JWT_SECRET"),
		},
		Streaming: StreamingConfig{
			SessionTTL:         viper.GetInt("STREAM_SESSION_TTL"),
			PresignTTL:         viper.GetInt("STREAM_PRESIGN_TTL"),
			MaxSessionsPerUser: viper.GetInt("STREAM_MAX_SESSIONS_PER_USER"),
		},
		Consul: ConsulConfig{
			Host:        viper.GetString("CONSUL_HOST"),
			Port:        viper.GetInt("CONSUL_PORT"),
			ServiceName: viper.GetString("SERVICE_NAME"),
			ServicePort: viper.GetInt("SERVICE_PORT"),
		},
	}
}