package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/nixpig/syringe.sh/cmd"
	"github.com/spf13/viper"
)

func main() {
	v := viper.New()
	if err := initialiseConfig(v); err != nil {
		log.Fatal(err)
	}

	if err := cmd.New(v).ExecuteContext(context.Background()); err != nil {
		log.Error("execute root command", "err", err)

	}
}

func initialiseConfig(v *viper.Viper) error {
	configPath := os.Getenv("SYRINGE_CONFIG_PATH")
	if configPath == "" {
		userConfigDir, err := os.UserConfigDir()
		if err != nil {
			return fmt.Errorf("get user config dir: %w", err)
		}

		configPath = filepath.Join(userConfigDir, "syringe")
	}

	if err := os.MkdirAll(configPath, os.ModePerm); err != nil {
		return err
	}

	log.SetLevel(log.DebugLevel)

	configFile, err := os.OpenFile(
		filepath.Join(configPath, "settings"),
		os.O_RDWR|os.O_CREATE,
		0666,
	)
	if err != nil {
		return fmt.Errorf("open config file (%s): %w", configPath, err)
	}
	configFile.Close()

	v.SetConfigFile(configFile.Name())
	v.SetConfigType("env")
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("read in config: %w", err)
	}

	return nil
}
