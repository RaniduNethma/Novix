package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct{
	App AppConfig
	Server ServerConfig
	Database DatabaseConfig
	MinIO MinIOConfig
	Kafka KafkaConfig
	FFmpeg FFmpegConfig
	Worker WorkerConfig
	Consul ConsulConfig
	JWT JWTConfig
}

type AppConfig struct{
	Env string
}

type ServerConfig struct{
	Port string
}

type DatabaseConfig struct{
	Host string
	Port string
	User string
	Password string
	Name string
	SSLMode string
}

type MinIOConfig struct{
	Endpoint string
	AccessKey string
	SecretKey string
	UseSSL bool
	RawBucket string
	ProcessedBucket string
}

type KafkaConfig struct{
	Brokers []string
	GroupId string
	TopicUploaded string
	TopicProcessed string
	TopicFailed string
}

type FFmpegConfig struct{
	Path string
	Threads int
}

type WorkerConfig struct{
	PoolSize int
	MaxRetries int
}

type ConsulConfig struct{
	Host string
	Port int
	ServiceName string
	ServicePort int
}

type JWTConfig struct{
	Secret string
}

func Load() *Config{
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
		Database: DatabaseConfig{
			Host: viper.GetString("DB_HOST"),
			Port: viper.GetString("DB_PORT"),
			User: viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Name: viper.GetString("DB_NAME"),
			SSLMode: viper.GetString("DB_SSLMODE"),
		},
		MinIO: MinIOConfig{
			Endpoint: viper.GetString("MINIO_ENDPOINT"),
			AccessKey: viper.GetString("MINIO_ACCESS_KEY"),
			SecretKey: viper.GetString("MINIO_SECRET_KEY"),
			UseSSL: viper.GetBool("MINIO_USE_SSL"),
			RawBucket: viper.GetString("MINIO_RAW_BUCKET"),
			ProcessedBucket: viper.GetString("MINIO_PROCESSED_BUCKET"),
		},
		Kafka: KafkaConfig{
			Brokers: []string{viper.GetString("KAFKA_BROKERS")},
			GroupId: viper.GetString("KAFKA_GROUP_ID"),
			TopicUploaded: viper.GetString("KAFKA_TOPIC_VIDEO_UPLOADED"),
			TopicProcessed: viper.GetString("KAFKA_TOPIC_VIDEO_PROCESSED"),
			TopicFailed: viper.GetString("KAFKA_TOPIC_VIDEO_FAILED"),
		},
		FFmpeg: FFmpegConfig{
			Path: viper.GetString("FFMPEG_PATH"),
			Threads: viper.GetInt("FFMPEG_THREADS"),
		},
		Worker: WorkerConfig{
			PoolSize: viper.GetInt("WORKER_POOL_SIZE"),
			MaxRetries: viper.GetInt("WORKER_MAX_RETRIES"),
		},
		Consul: ConsulConfig{
			Host: viper.GetString("CONSUL_HOST"),
			Port: viper.GetInt("CONSUL_PORT"),
			ServiceName: viper.GetString("SERVICE_NAME"),
			ServicePort: viper.GetInt("SERVICE_PORT"),
		},
		JWT: JWTConfig{
			Secret: viper.GetString("JWT_SECRET"),
		},
	}
}
