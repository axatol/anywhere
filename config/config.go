package config

import (
	_ "embed"
	"flag"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
)

//go:embed config.yml
var rawConfig []byte

type config struct {
	Verbose bool `yaml:"verbose" json:"verbose"`
	Debug   bool `yaml:"debug" json:"debug"`

	Database struct {
		Host string `yaml:"host" json:"host"`
		Port int    `yaml:"port" json:"port"`
		Name string `yaml:"name" json:"name"`
		User string `yaml:"user" json:"user"`
		Pass string `yaml:"pass" json:"pass"`
	} `yaml:"database" json:"database"`

	Server struct {
		Port         int      `yaml:"port" json:"port"`
		AllowOrigins []string `yaml:"allow_origins" json:"allow_origins"`
	} `yaml:"server" json:"server"`
}

var Config config
var Logger *zap.SugaredLogger

type flagList []string

func (l *flagList) String() string {
	return strings.Join(*l, ", ")
}

func (l *flagList) Set(value string) error {
	*l = append(*l, value)
	return nil
}

func init() {
	initConfig()
	initLogger()
}

func initConfig() {
	var defaultConfig config
	if err := yaml.Unmarshal(rawConfig, &defaultConfig); err != nil {
		panic("failed to parse config file")
	}

	verbose := flag.Bool("verbose", defaultConfig.Verbose, "enable verbose logging")
	debug := flag.Bool("debug", defaultConfig.Debug, "enable debugging mode")
	dbHost := flag.String("dbhost", defaultConfig.Database.Host, "database host")
	dbPort := flag.Int("dbport", defaultConfig.Database.Port, "database port")
	dbName := flag.String("dbname", defaultConfig.Database.Name, "database name")
	dbUser := flag.String("dbuser", defaultConfig.Database.User, "database user")
	dbPass := flag.String("dbpass", defaultConfig.Database.Pass, "database pass")
	srvPort := flag.Int("srvport", defaultConfig.Server.Port, "server listen port")
	var srvOrigins flagList
	flag.Var(&srvOrigins, "alloworigin", "allowed origins")
	flag.Parse()

	Config.Verbose = *verbose
	Config.Debug = *debug
	Config.Database.Host = *dbHost
	Config.Database.Port = *dbPort
	Config.Database.Name = *dbName
	Config.Database.User = *dbUser
	Config.Database.Pass = *dbPass
	Config.Server.Port = *srvPort
	if len(srvOrigins) > 1 {
		Config.Server.AllowOrigins = srvOrigins
	} else {
		Config.Server.AllowOrigins = defaultConfig.Server.AllowOrigins
	}
}

func initLogger() {
	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: true,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "name",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	if Config.Debug {
		config.Development = true
		config.Sampling = nil
		config.Encoding = "console"
	}

	if Config.Verbose {
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	Logger = logger.Sugar()
}
