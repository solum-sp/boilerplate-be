package utils


import (
	"log"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config interface {
	Load() error
}

type ConfigManager struct {
	cfgs []Config
}

func NewConfigManager(cfgs []Config) *ConfigManager {
	return &ConfigManager{
		cfgs: cfgs,
	}
}

func (c *ConfigManager) LoadAllConfig(envFolderPath string) error {
	loadEnv(envFolderPath)
	for _, cfg := range c.cfgs {
		if err := cfg.Load(); err != nil {
			return err
		}
	}
	return nil
}
func loadEnv(envFolderPath string) error {
	env := os.Getenv("APP_ENV")
	envFileName := ""

	switch env {
	case "development":
		envFileName = ".env.dev"
	case "test":
		envFileName = ".env.test"
	case "production":
		envFileName = ".env.production"
	}
	if envFileName != "" {
		envFilePath := envFolderPath + "/" + envFileName
		log.Printf("Loading config from file:%s\n", envFilePath)
		_ = godotenv.Load(envFilePath)
	}
	log.Printf("Loading config from environment\n")
	_ = godotenv.Load()
	return nil
}

func ParseConfig(c Config) error {
	return env.Parse(c)
}
