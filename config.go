package main

import (
	"api/consts"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

type Config struct {
	GRPCSettings struct {
		Host                   string        `yaml:"Host" envconfig:"GRPC_HOST"`
		Port                   string        `yaml:"Port" envconfig:"GRPC_PORT"`
		GatewayPort            string        `yaml:"GatewayPort" envconfig:"GRPC_GATEWAY_PORT"`
		MaxConcurrentStreams   int           `yaml:"MaxConcurrentStreams" envconfig:"GRPC_MAX_CONCURRENT_STREAMS"`
		MaxGoriutinesPerStream int           `yaml:"MaxGoriutinesPerStream" envconfig:"GRPC_MAX_GOROUTINES_PER_STREAM"`
		ConnDeadlineDuration   time.Duration `yaml:"ConnDeadlineDuration" envconfig:"GRPC_CONN_DEADLINE_DURATION"`
	} `yaml:"GRPCSettings"`
	DBSettings struct {
		Login    string `yaml:"Login" envconfig:"DB_LOGIN"`
		Password string `yaml:"Password" envconfig:"DB_PASSWORD"`
		Addr     string `yaml:"Addr" envconfig:"DB_ADDR"`
		Port     string `yaml:"Port" envconfig:"DB_PORT"`
		DB       string `yaml:"DB" envconfig:"DB_DATABASE"`
		Table    string `yaml:"Table" envconfig:"DB_TABLE"`
	} `yaml:"DBSettings"`
	MetricsSettings struct {
		Port          string        `yaml:"ServerPort" envconfig:"METRICS_SERVER_PORT"`
		ServerRuntime time.Duration `yaml:"ServerRuntime" envconfig:"METRICS_SERVER_RUNTIME"`
		Path          string        `yaml:"MetricsPath" envconfig:"METRICS_PATH"`
	} `yaml:"MetricsSettings"`
	LogsSettings struct {
		Prefix    string        `yaml:"Prefix" envconfig:"LOGS_PREFIX"`
		Frequency time.Duration `yaml:"Frequency" envconfig:"LOGS_FREQUENCY_CREATING"`
		Path      string        `yaml:"LogsPath" envconfig:"LOGS_PATH"`
	} `yaml:"LogsSettings"`
}

func (c *Config) ReadConfig(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(c)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) ReadEnv() error {
	err := envconfig.Process("", c)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) SetDefaultConsts() {

	// grpc
	if c.GRPCSettings.Port == "" {
		c.GRPCSettings.Port = consts.GRPC_PORT
	}
	if c.GRPCSettings.MaxConcurrentStreams == 0 {
		c.GRPCSettings.MaxConcurrentStreams = consts.GRPC_MAX_CONCURRENT_STREAMS
	}
	if c.GRPCSettings.MaxGoriutinesPerStream == 0 {
		c.GRPCSettings.MaxGoriutinesPerStream = consts.GRPC_MAX_GOROUTINES_PER_STREAM
	}
	if c.GRPCSettings.ConnDeadlineDuration == 0 {
		c.GRPCSettings.ConnDeadlineDuration = consts.GRPC_CONN_DEADLINE_DURATION
	}

	// metrics
	if c.MetricsSettings.Port == "" {
		c.MetricsSettings.Port = consts.METRICS_PORT
	}
	if c.MetricsSettings.ServerRuntime == 0 {
		c.MetricsSettings.ServerRuntime = consts.METRICS_SERVER_TIMEOUT
	}
	if c.MetricsSettings.Path == "" {
		c.MetricsSettings.Path = consts.METRICS_PATH
	}

}
