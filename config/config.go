package config

import (
	_ "embed"
	"flag"
	"os"

	"gopkg.in/yaml.v3"
)

type config struct {
	Verbose bool `yaml:"verbose"`
	Debug   bool `yaml:"debug"`

	Database struct {
		Host string `yaml:"host"`
		Name string `yaml:"name"`
		User string `yaml:"user"`
		Pass string `yaml:"pass"`
	} `yaml:"database"`

	Datastore struct {
		Host string `yaml:"host"`
		Name string `yaml:"name"`
		User string `yaml:"user"`
		Pass string `yaml:"pass"`
	} `yaml:"datastore"`

	Server struct {
		Port         string   `yaml:"port"`
		AllowOrigins []string `yaml:"allow_origins"`
	} `yaml:"server"`
}

var (
	Config config
)

func init() {
	initConfig()
	initLogger()
}

func resolveRawConfig() []byte {
	var locationsToTry []string
	if loc, ok := os.LookupEnv("ANYWHERE_CONFIG_FILE"); ok {
		locationsToTry = []string{loc}
	} else {
		locationsToTry = []string{
			"./config.yml",
			"./example.config.yml",
		}
	}

	for _, loc := range locationsToTry {
		if _, err := os.Stat(loc); err != nil {
			continue
		}

		rawConfig, err := os.ReadFile(loc)
		if err != nil {
			panic(err)
		}

		return rawConfig
	}

	return nil
}

func initConfig() {
	if rawConfig := resolveRawConfig(); rawConfig != nil {
		if err := yaml.Unmarshal(rawConfig, &Config); err != nil {
			panic("failed to parse config file")
		}
	}

	// applying config from defaults
	Config.Verbose = assertBool("Verbose", &Config.Verbose, false)
	Config.Debug = assertBool("Debug", &Config.Debug, false)
	Config.Database.Host = assertStr("Database.Host", &Config.Database.Host, "mongodb://localhost:27017")
	Config.Database.Name = assertStr("Database.Name", &Config.Database.Name, "anywhere")
	// Config.Database.User
	// Config.Database.Pass
	Config.Datastore.Host = assertStr("Datastore.Host", &Config.Datastore.Host, "http://localhost:9000")
	Config.Datastore.Name = assertStr("Datastore.Name", &Config.Datastore.Name, "anywhere")
	// Config.Datastore.User
	// Config.Datastore.Pass
	Config.Server.Port = assertStr("Server.Port", &Config.Server.Port, "8042")
	Config.Server.AllowOrigins = assertStrs("Server.Port", Config.Server.AllowOrigins, []string{"http://localhost:3000"})

	// applying config from environment variables
	Config.Verbose = assertBool("Verbose", envBool("VERBOSE"), Config.Verbose)
	Config.Debug = assertBool("Debug", envBool("DEBUG"), Config.Debug)
	Config.Database.Host = assertStr("Database.Host", envStr("DATABASE_HOST"), Config.Database.Host)
	Config.Database.Name = assertStr("Database.Name", envStr("DATABASE_NAME"), Config.Database.Name)
	Config.Database.User = assertStr("Database.User", envStr("DATABASE_USER"), Config.Database.User)
	Config.Database.Pass = assertStr("Database.Pass", envStr("DATABASE_PASS"), Config.Database.Pass)
	Config.Datastore.Host = assertStr("Datastore.Host", envStr("DATASTORE_HOST"), Config.Datastore.Host)
	Config.Datastore.Name = assertStr("Datastore.Name", envStr("DATASTORE_NAME"), Config.Datastore.Name)
	Config.Datastore.User = assertStr("Datastore.User", envStr("DATASTORE_USER"), Config.Datastore.User)
	Config.Datastore.Pass = assertStr("Datastore.Pass", envStr("DATASTORE_PASS"), Config.Datastore.Pass)
	Config.Server.Port = assertStr("Server.Port", envStr("SERVER_PORT"), Config.Server.Port)
	Config.Server.AllowOrigins = assertStrs("Server.AllowOrigins", envStrs("SERVER_ALLOW_ORIGINS"), Config.Server.AllowOrigins)

	// applying config from flags
	verbose := flag.Bool("verbose", Config.Verbose, "enable verbose logging")
	debug := flag.Bool("debug", Config.Debug, "enable debugging mode")
	dbHost := flag.String("dbhost", Config.Database.Host, "database host")
	dbName := flag.String("dbname", Config.Database.Name, "database name")
	dbUser := flag.String("dbuser", Config.Database.User, "database user")
	dbPass := flag.String("dbpass", Config.Database.Pass, "database pass")
	dsHost := flag.String("dshost", Config.Datastore.Host, "database host")
	dsName := flag.String("dsname", Config.Datastore.Name, "database name")
	dsUser := flag.String("dsuser", Config.Datastore.User, "database user")
	dsPass := flag.String("dspass", Config.Datastore.Pass, "database pass")
	srvPort := flag.String("srvport", Config.Server.Port, "server listen port")
	var srvOrigins flagList
	flag.Var(&srvOrigins, "alloworigin", "allowed origins")
	flag.Parse()

	Config.Verbose = *verbose
	Config.Debug = *debug
	Config.Database.Host = *dbHost
	Config.Database.Name = *dbName
	Config.Database.User = *dbUser
	Config.Database.Pass = *dbPass
	Config.Datastore.Host = *dsHost
	Config.Datastore.Name = *dsName
	Config.Datastore.User = *dsUser
	Config.Datastore.Pass = *dsPass
	Config.Server.Port = *srvPort
	Config.Server.AllowOrigins = assertStrs("Server.AllowOrigins", srvOrigins, Config.Server.AllowOrigins)
}
