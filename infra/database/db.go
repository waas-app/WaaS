package database

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"github.com/waas-app/WaaS/config"
	"github.com/waas-app/WaaS/util"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"moul.io/zapgorm2"
)

var instance *gorm.DB

func Connect(logger *zap.Logger) (*gorm.DB, error) {
	u, err := url.Parse(config.Spec.Storage)
	if err != nil {
		return nil, err
	}
	lg := zapgorm2.New(logger)
	if config.Spec.Environment == config.Development {
		lg.LogMode(gormlogger.Info)
	} else {
		lg.LogMode(gormlogger.Error)
	}
	lg.SetAsDefault()

	var connectionString string
	var db *gorm.DB
	switch u.Scheme {
	case "postgresql":
		// handle `postgresql` as the scheme to be compatible with
		// standar uri style postgresql connection strings (i.e. like psql)
		u.Scheme = "postgres"
		fallthrough
	case "postgres":
		connectionString = pgconn(u)
		logger.Info("connecting to postgres", zap.String("connectionString", connectionString))
		db, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{
			Logger: lg,
		})
	case "mysql":
		connectionString = mysqlconn(u)
		db, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{
			Logger: lg,
		})
	default:
		// unreachable because our storage backend factory
		// function (contracts.go) already checks the url scheme.
		logger.Panic("unknown sql storage backend", zap.String("scheme", u.Scheme))
	}

	if err != nil {
		return db, err
	}

	if err := db.Use(otelgorm.NewPlugin()); err != nil {
		logger.Error("error adding opentelemetry to database", zap.Error(err))
		return db, err
	}

	instance = db
	AutoMigrate(db)
	return db, nil
}

func pgconn(u *url.URL) string {
	password, _ := u.User.Password()
	decodedQuery, err := url.QueryUnescape(u.RawQuery)
	if err != nil {
		util.Logger(context.Background()).Error("Failed to decode query", zap.Error(err))
		decodedQuery = ""
	}
	return fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=%s dbname=%s %s",
		u.Hostname(),
		u.Port(),
		u.User.Username(),
		password,
		"disable",
		strings.TrimLeft(u.Path, "/"),
		decodedQuery,
	)
}

func mysqlconn(u *url.URL) string {
	password, _ := u.User.Password()
	return fmt.Sprintf(
		"%s:%s@%s/%s?%s",
		u.User.Username(),
		password,
		u.Host,
		strings.TrimLeft(u.Path, "/"),
		u.RawQuery,
	)
}

func Instance(ctx context.Context) *gorm.DB {
	nCtx := ctx

	if instance == nil {
		var err error
		instance, err = Connect(util.Logger(nCtx).ZapLogger())
		if err != nil {
			util.Logger(nCtx).Fatal("Unable to connect to DB, Please check the Database configuration")
			return nil
		}
	}
	return instance.Session(&gorm.Session{NewDB: true, Context: nCtx})
}
