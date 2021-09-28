// The MIT License (MIT)
//
// Copyright (c) 2021 Alessandro De Blasis <alex@deblasis.net>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
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
