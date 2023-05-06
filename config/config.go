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
		Dnspod     []*DNSpodConfig
	}
}

type CloudflareConfig struct {
	Email   string
	Zone    string
	AuthKey string
	Host    string
}

type DNSpodConfig struct {
	SecretId  string
	SecretKey string
	Domain    string
	SubDomain string
}

func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.ddns-service/")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err.Error()))
	}

	if err := viper.Unmarshal(&Section); err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err.Error()))
	}
}
