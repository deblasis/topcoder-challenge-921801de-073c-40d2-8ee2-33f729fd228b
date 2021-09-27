//
// Copyright 2021 Alessandro De Blasis <alex@deblasis.net>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package config

import (
	"os"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const configFileName = "app"

// Config declare the application configuration variables
type Config struct {
	ListenAddr        string
	HttpServerPort    string
	GrpcServerPort    string
	AuxGrpcServerPort string

	BindOnLocalhost bool

	Db         DbConfig         `mapstructure:"db"`
	Consul     ConsulConfig     `mapstructure:"consul"`
	Zipkin     ZipkinConfig     `mapstructure:"zipkin"`
	SSL        SSLConfig        `mapstructure:"ssl"`
	APIGateway APIGatewayConfig `mapstructure:"apigateway"`
	JWT        JWTConfig        `mapstructure:"jwt"`

	ShippingStation ShippingStationConfig `mapstructure:"shippingstation"`
	Clessidra       ClessidraConfig       `mapstructure:"clessidra"`

	Logger log.Logger
}

type ConsulConfig struct {
	Host string
	Port string
}

type SSLConfig struct {
	ServerCert string
	ServerKey  string
}

type ZipkinConfig struct {
	V2Url          string
	UseBridge      bool
	LightstepToken string
}

// DbConfig declare database variables
type DbConfig struct {
	Address    string
	Username   string
	Password   string
	Database   string
	Sslmode    string
	Drivername string
}

type JWTConfig struct {
	Secret        string
	TokenDuration int

	PrivKeyPath string //= "../../certs/jwt.pem.key"
	PubKeyPath  string //= "../../certs/jwt.pem.pub"
}

type APIGatewayConfig struct {
	RetryMax                          int
	RetryTimeoutMs                    int
	AuthServiceGRPCEndpoint           string `mapstructure:"authservice_grpcendpoint"`
	CentralCommandServiceGRPCEndpoint string `mapstructure:"centralcommandservice_grpcendpoint"`
	ShippingStationGRPCEndpoint       string `mapstructure:"shippingstationservice_grpcendpoint"`
}

type ShippingStationConfig struct {
	DockHoldingPeriod int
}
type ClessidraConfig struct {
	PollingInterval int64
}

// LoadConfig load config from file
func LoadConfig() (Config, error) {
	v := viper.New()

	v.SetConfigName(configFileName)
	v.SetEnvPrefix("deblasis")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("../..")
	v.AddConfigPath("./config")

	replacer := strings.NewReplacer(".", "_")
	v.SetEnvKeyReplacer(replacer)
	v.AutomaticEnv()

	var cfg Config
	if err := v.ReadInConfig(); err != nil {
		return Config{}, errors.Wrap(err, "Failed to read config")
	}

	err := v.Unmarshal(&cfg)
	if err != nil {
		return Config{}, errors.Wrap(err, "Unable to decode into struct")
	}

	loglevel := v.GetString("loglevel")

	cfg.Logger = getLogger(loglevel)
	return cfg, nil
}

func getLogger(loglevel string) log.Logger {

	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	var logFilter level.Option
	switch loglevel {
	case "debug":
		logFilter = level.AllowDebug()
	case "info":
		logFilter = level.AllowInfo()
	case "warn":
		logFilter = level.AllowWarn()
	case "error":
		logFilter = level.AllowError()
	default:
		logFilter = level.AllowAll()
	}

	return level.NewFilter(logger, logFilter)
}
