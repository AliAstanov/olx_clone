package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

const (
	// DebugMode indicates service mode is debug.
	DebugMode string = "debug"
	// TestMode indicates service mode is test.
	TestMode string = "test"
	// ReleaseMode indicates service mode is release.
	ReleaseMode string = "release"
)

const (
	Email         = "alibekastan1998gmail.com"
	EmailPassword = ""
	SMTPHost      = "smtp.example.com"
	SMTPPort      = "587"
)

type GeneralConfig struct {
	AppName     string
	Environment string // debug,
	Version     string
	HTTPPort    string
	HTTPScheme  string
	SignInKey   string
}

type PgConfig struct {
	Username     string
	Password     string
	Host         string
	Port         int
	DatabaseName string
}

type RedisConfig struct {
	Username     string
	Password     string
	Host         string
	Port         int
	DatabaseName int
}

type Config struct {
	GeneralConfig GeneralConfig
	PgConfig      PgConfig
	RedisConfig   RedisConfig
}

func NewConfig() Config {
	return Config{GeneralConfig: GeneralConfig{AppName: "olx_clone"}}
}

func Load() Config {

	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}

	var config = NewConfig()

	// general configs
	config.GeneralConfig.Environment = cast.ToString(getOrReturnDefaultValue("ENVIRONMENT", DebugMode))
	config.GeneralConfig.Version = cast.ToString(getOrReturnDefaultValue("VERSION", "1.0"))
	config.GeneralConfig.HTTPPort = cast.ToString(getOrReturnDefaultValue("HTTP_PORT", ":8080"))
	config.GeneralConfig.HTTPScheme = cast.ToString(getOrReturnDefaultValue("HTTP_SCHEME", "http"))
	config.GeneralConfig.SignInKey = cast.ToString(getOrReturnDefaultValue("SIGN_IN_KEY", "ASJDKLFJASasdFASE2SD2dafa"))

	// postgres config
	config.PgConfig.Username = cast.ToString(getOrReturnDefaultValue("POSTGRES_USER", "superuser"))
	config.PgConfig.Password = cast.ToString(getOrReturnDefaultValue("POSTGRES_PASSWORD", "superuser"))
	config.PgConfig.Host = cast.ToString(getOrReturnDefaultValue("POSTGRES_HOST", "localhost"))
	config.PgConfig.Port = cast.ToInt(getOrReturnDefaultValue("POSTGRES_PORT", 5432))
	config.PgConfig.DatabaseName = cast.ToString(getOrReturnDefaultValue("POSTGRES_DATABASE", config.GeneralConfig.AppName))

	// redis config
	config.RedisConfig.Host = cast.ToString(getOrReturnDefaultValue("REDIS_HOST", "localhost"))
	config.RedisConfig.Port = cast.ToInt(getOrReturnDefaultValue("REDIS_PORT", 6379))

	return config

}

// urlExample := "postgres://username:password@localhost:5432/database_name"
func getOrReturnDefaultValue(key string, defaultValue interface{}) interface{} {

	val, exists := os.LookupEnv(key)

	if exists {
		return val
	}

	return defaultValue
}
