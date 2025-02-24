package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func main() {
	configName := os.Getenv("CONFIG_NAME")
	if configName == "" {
		configName = "props"
	}
	viper.SetConfigName(configName)

	viper.SetConfigType("yml")
	viper.AddConfigPath("configs/")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("failed to read properties config: %w", err))
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}
