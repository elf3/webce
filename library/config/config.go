package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"strings"
)

type Config struct {
	Name string
}

func Init(cfg string) error {
	config := Config{
		Name: cfg,
	}
	if err := config.initConfig(); err != nil {
		return err
	}
	config.watchConfig()
	return nil
}

/**
初始化配置文件
*/
func (c *Config) initConfig() error {
	if c.Name != "" {
		viper.AddConfigPath("./conf")
		viper.SetConfigFile(c.Name)

	} else {
		viper.AddConfigPath("./conf")
		viper.SetConfigName("app")
	}
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("SITE_CONF")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}

// 监控配置文件变更
func (c *Config) watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Printf("Config file Changed: %s %s %s \n", in.Name, in.String(), in.Op)
		_ = recover()

	})
}
