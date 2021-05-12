package config

import (
	"fmt"
	"github.com/spf13/viper"
)

var Section *ConfigSection

type ConfigSection struct {
	IPGetter []string
	DnsApi   struct {
		Cloudflare []*CloudflareConfig
	}
}

type CloudflareConfig struct {
	Email   string
	Zone    string
	AuthKey string
	Host    string
}

type DnspodConfig struct {
}

func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	err = viper.Unmarshal(&Section)
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
