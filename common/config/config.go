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
	ListenAddr     string
	HttpServerPort string
	GrpcServerPort string

	BindOnLocalhost bool

	Db         DbConfig         `mapstructure:"db"`
	Consul     ConsulConfig     `mapstructure:"consul"`
	Zipkin     ZipkinConfig     `mapstructure:"zipkin"`
	SSL        SSLConfig        `mapstructure:"ssl"`
	APIGateway APIGatewayConfig `mapstructure:"apigateway"`
	JWT        JWTConfig        `mapstructure:"jwt"`
	Kong       KongConfig       `mapstructure:"kong"`

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

type KongConfig struct {
	BaseUrl string
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
}

type APIGatewayConfig struct {
	RetryMax       int
	RetryTimeoutMs int
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
	logger := getLogger(loglevel)

	cfg.Logger = logger
	return cfg, nil
}

func getLogger(loglevel string) log.Logger {

	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

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

	logger = level.NewFilter(logger, logFilter)
	logger = log.With(logger, "caller", log.DefaultCaller)

	return logger
}
