package config

import (
	"fmt"
	"github.com/li1553770945/personal-file-service/biz/constant"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type ServerConfig struct {
	ServiceName   string `yaml:"service-name"`
	ListenAddress string `yaml:"listen-address"`
}

type OpenTelemetryConfig struct {
	Endpoint string `yaml:"endpoint"`
}

type CosConfig struct {
	Ak       string `yaml:"ak"`
	Sk       string `yaml:"sk"`
	Endpoint string `yaml:"endpoint"`
}

type EtcdConfig struct {
	Endpoint []string `yaml:"endpoint"`
}

type DatabaseConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Address  string `yaml:"address"`
	Port     int32  `yaml:"port"`
	UseTLS   bool   `yaml:"use-tls"`
}

type Config struct {
	Env                 string
	ServerConfig        ServerConfig        `yaml:"server"`
	OpenTelemetryConfig OpenTelemetryConfig `yaml:"open-telemetry"`
	DatabaseConfig      DatabaseConfig      `yaml:"database"`
	EtcdConfig          EtcdConfig          `yaml:"etcd"`
	CosConfig           CosConfig           `yaml:"cos"`
}

func GetConfig(env string) *Config {
	if env != constant.EnvProduction && env != constant.EnvDevelopment {
		panic(fmt.Sprintf("环境必须是%s或者%s之一", constant.EnvProduction, constant.EnvDevelopment))
	}
	conf := &Config{}
	path := filepath.Join("conf", fmt.Sprintf("%s.yml", env))
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	err = yaml.NewDecoder(f).Decode(conf)
	conf.Env = env
	if err != nil {
		panic(err)
	}

	return conf
}
