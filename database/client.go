package database

import (
	"fmt"
	"strings"

	"github.com/tunes-anywhere/anywhere/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var logger = config.Logger.Named("database")

func Init() {
	dsn := []string{
		fmt.Sprintf("host=%s", config.Config.Database.Host),
		fmt.Sprintf("user=%s", config.Config.Database.User),
		fmt.Sprintf("password=%s", config.Config.Database.Pass),
		fmt.Sprintf("port=%d", config.Config.Database.Port),
		fmt.Sprintf("dbname=%s", config.Config.Database.Name),
	}

	config := postgres.Config{
		DSN: strings.Join(dsn, " "),
	}

	conn, err := gorm.Open(postgres.New(config))
	if err != nil {
		logger.Fatalln(err)
	}

	DB = conn
}
