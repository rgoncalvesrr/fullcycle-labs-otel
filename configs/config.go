package configs

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	WeatherApiKey              string `mapstructure:"WEATHER_API_KEY"`
	WeatherApiUrl              string `mapstructure:"WEATHER_API_URL"`
	CepApiUrl                  string `mapstructure:"CEP_API_URL"`
	InputApiHttpPort           string `mapstructure:"INPUT_API_HTTP_PORT"`
	InputApiOtelServiceName    string `mapstructure:"INPUT_API_OTEL_SERVICE_NAME"`
	OrchestratorApiPort        string `mapstructure:"ORCHESTRATOR_API_PORT"`
	OrchestratorApiHost        string `mapstructure:"ORCHESTRATOR_API_HOST"`
	OrchestratorApiServiceName string `mapstructure:"ORCHESTRATOR_API_SERVICE_NAME"`
}

var Cfg *Config

func init() {
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.AddConfigPath("../..")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	//viper.addconfSetConfigFile(cfgFile)

	e := viper.ReadInConfig()
	if e != nil {
		log.Fatal("Can't find the file app.env : ", e)
	}

	e = viper.Unmarshal(&Cfg)
	if e != nil {
		log.Fatal("Can't unmarshal the file app.env : ", e)
	}
}
