package config

import (
	"flag"
	"strings"
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
		Host       string `yaml:"host"`
		User       string `yaml:"user"`
		Pass       string `yaml:"pass"`
		Bucket     string `yaml:"bucket"`
		PathPrefix string `yaml:"path_prefix"`
	} `yaml:"datastore"`

	Server struct {
		Port         string   `yaml:"port"`
		AllowOrigins []string `yaml:"allow_origins"`

		Auth struct {
			Issuer   string `yaml:"issuer"`
			Audience string `yaml:"audience"`
		} `yaml:"auth"`
	} `yaml:"server"`

	Datasource struct {
		APIKeys APIKeys `yaml:"api_keys"`
	} `yaml:"datasource"`
}

func (c *config) ApplyDefaults() {
	c.Verbose = assertBool("verbose", &c.Verbose, false)
	c.Debug = assertBool("debug", &c.Debug, false)
	c.Database.Host = assertStr("database.host", &c.Database.Host, "mongodb://localhost:27017")
	c.Database.Name = assertStr("database.name", &c.Database.Name, "anywhere")
	// c.Database.User
	// c.Database.Pass
	c.Datastore.Host = assertStr("datastore.host", &c.Datastore.Host, "http://localhost:9000")
	// c.Datastore.User
	// c.Datastore.Pass
	c.Datastore.Bucket = assertStr("datastore.bucket", &c.Datastore.Bucket, "anywhere")
	c.Datastore.PathPrefix = assertStr("datastore.path_prefix", &c.Datastore.PathPrefix, "/")
	c.Server.Port = assertStr("server.port", &c.Server.Port, "8042")
	c.Server.AllowOrigins = assertStrs("server.allow_origins", c.Server.AllowOrigins, []string{"http://localhost:3000"})
	// c.Server.Auth.Issuer
	// c.Server.Auth.Audience
	// c.Datasource.APIKeys
}

func (c *config) ApplyFromEnvironment() {
	c.Verbose = assertBool("VERBOSE", envBool("VERBOSE"), c.Verbose)
	c.Debug = assertBool("DEBUG", envBool("DEBUG"), c.Debug)
	c.Database.Host = assertStr("DATABASE_HOST", envStr("DATABASE_HOST"), c.Database.Host)
	c.Database.Name = assertStr("DATABASE_NAME", envStr("DATABASE_NAME"), c.Database.Name)
	c.Database.User = assertStr("DATABASE_USER", envStr("DATABASE_USER"), c.Database.User)
	c.Database.Pass = assertStr("DATABASE_PASS", envStr("DATABASE_PASS"), c.Database.Pass)
	c.Datastore.Host = assertStr("DATASTORE_HOST", envStr("DATASTORE_HOST"), c.Datastore.Host)
	c.Datastore.User = assertStr("DATASTORE_USER", envStr("DATASTORE_USER"), c.Datastore.User)
	c.Datastore.Pass = assertStr("DATASTORE_PASS", envStr("DATASTORE_PASS"), c.Datastore.Pass)
	c.Datastore.Bucket = assertStr("DATASTORE_BUCKET", envStr("DATASTORE_BUCKET"), c.Datastore.Bucket)
	c.Datastore.PathPrefix = assertStr("DATASTORE_PATH_PREFIX", envStr("DATASTORE_PATH_PREFIX"), c.Datastore.PathPrefix)
	c.Server.Port = assertStr("SERVER_PORT", envStr("SERVER_PORT"), c.Server.Port)
	c.Server.AllowOrigins = assertStrs("SERVER_ALLOW_ORIGINS", envStrs("SERVER_ALLOW_ORIGINS"), c.Server.AllowOrigins)
	c.Server.Auth.Issuer = assertStr("SERVER_AUTH_ISSUER", envStr("SERVER_AUTH_ISSUER"), c.Server.Auth.Issuer)
	c.Server.Auth.Audience = assertStr("SERVER_AUTH_AUDIENCE", envStr("SERVER_AUTH_AUDIENCE"), c.Server.Auth.Audience)
	if rawapikeys := envStrs("DATASOURCE_API_KEYS"); len(rawapikeys) > 0 {
		c.Datasource.APIKeys = *new(APIKeys).MustFromStrings(rawapikeys)
	}
}

func (c *config) ApplyFromFlags() {
	flagVerbose := flag.Bool("verbose", c.Verbose, "enable verbose logging")
	flagDebug := flag.Bool("debug", c.Debug, "enable debugging mode")
	flagDatabaseHost := flag.String("database-host", c.Database.Host, "database host")
	flagDatabaseName := flag.String("database-name", c.Database.Name, "database name")
	flagDatabaseUser := flag.String("database-user", "", "database user")
	flagDatabasePass := flag.String("database-pass", "", "database pass")
	flagDatastoreHost := flag.String("datastore-host", c.Datastore.Host, "database host")
	flagDatastoreUser := flag.String("datastore-user", "", "database user")
	flagDatastorePass := flag.String("datastore-pass", "", "database pass")
	flagDatastoreBucket := flag.String("datastore-bucket", c.Datastore.Bucket, "database bucket")
	flagDatastorePathPrefix := flag.String("datastore-path-prefix", c.Datastore.PathPrefix, "database bucket path prefix")
	flagServerPort := flag.String("server-port", c.Server.Port, "server listen port")
	flagServerAllowedOrigins := flag.String("server-allowed-origins", strings.Join(c.Server.AllowOrigins, ","), "server allowed origins")
	flagServerAuthIssuer := flag.String("server-auth-issuer", c.Server.Auth.Issuer, "server auth issuer")
	flagServerAuthAudience := flag.String("server-auth-audience", c.Server.Auth.Audience, "server auth audience")
	flagDatasourceAPIKeys := flag.String("datasource-api-keys", c.Datasource.APIKeys.String(), "datasource api keys")
	flag.Parse()

	c.Verbose = *flagVerbose
	c.Debug = *flagDebug
	c.Database.Host = *flagDatabaseHost
	c.Database.Name = *flagDatabaseName
	c.Database.User = assertStr("database.user", flagDatabaseUser, c.Database.User)
	c.Database.Pass = assertStr("database.pass", flagDatabasePass, c.Database.Pass)
	c.Datastore.Host = *flagDatastoreHost
	c.Datastore.User = assertStr("datastore.user", flagDatastoreUser, c.Datastore.User)
	c.Datastore.Pass = assertStr("datastore.pass", flagDatastorePass, c.Datastore.Pass)
	c.Datastore.Bucket = *flagDatastoreBucket
	c.Datastore.PathPrefix = *flagDatastorePathPrefix
	c.Server.Port = *flagServerPort
	c.Server.AllowOrigins = strings.Split(*flagServerAllowedOrigins, ",")
	c.Server.Auth.Issuer = *flagServerAuthIssuer
	c.Server.Auth.Audience = *flagServerAuthAudience
	c.Datasource.APIKeys = *new(APIKeys).MustFromStrings(assertStrs("datasource.api_keys", strings.Split(*flagDatasourceAPIKeys, ","), strings.Split(c.Datasource.APIKeys.String(), ",")))
}

func init() {
	initConfig()
	initLogger()

	Log.Debugw("configured successfully",
		"verbose", Values.Verbose,
		"debug", Values.Debug,
		"api_keys", Values.Datasource.APIKeys.Names(),
		"allowed_origins", Values.Server.AllowOrigins,
	)
}
