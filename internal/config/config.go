package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	"github.com/mohammadne/shopping-cart-manager/internal"
	"github.com/mohammadne/shopping-cart-manager/pkg/mysql"
	"github.com/mohammadne/shopping-cart-manager/pkg/redis"
)

type Config struct {
	Mysql *mysql.Config `required:"true"`
	Redis *redis.Config `required:"true"`
}

func Load(print bool) (config Config, err error) {
	prefix := strings.ToUpper(internal.System)

	if err = envconfig.Process(prefix, &config); err != nil {
		return Config{}, fmt.Errorf("error processing config via envconfig: %v", err)
	}

	if print {
		fmt.Println("================ Loaded Configuration ================")
		object, _ := json.MarshalIndent(config, "", "  ")
		fmt.Println(string(object))
		fmt.Println("======================================================")
	}

	return config, nil
}

func LoadDefaults(print bool, relativePath string) (config Config, err error) {
	currentWorkingDirectory, err := os.Getwd()
	if err != nil {
		return Config{}, fmt.Errorf("got err on getting current directory: %v", err)
	}

	envFilePath := currentWorkingDirectory + relativePath + "/internal/config/defaults.env"
	fmt.Println(envFilePath)
	err = godotenv.Load(envFilePath)
	if err != nil {
		return Config{}, fmt.Errorf("got err on loading env file: %v", err)
	}

	return Load(print)
}
