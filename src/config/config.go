package config

import "github.com/spf13/viper"

type Conf struct {
	Server struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`
	
	Mysql struct {
		Dsn string `mapstructure:"dsn"`
	} `mapstructure:"mysql"`
	Redis struct {
		Addr     string `mapstructure:"addr"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
	} `mapstructure:"redis"`
	JWT struct {
		SecretKey     string `mapstructure:"secret_key"`
		ExpireDays    int    `mapstructure:"expire_days"`
	} `mapstructure:"jwt"`
	
}

var Global Conf 

func LoadConfig() {
	// 配置加载逻辑
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.SetConfigName("conf")



	viper.SetDefault("server.port", "8080")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}


	err = viper.Unmarshal(&Global)
	if err != nil {
		panic(err)
	}
}
