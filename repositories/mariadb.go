package repositories

import (
	"time"

	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/mysql"
	"go.uber.org/zap"

	"github.com/github.com/steevehook/account-api/logging"
)

const (
	accountsTableName = "accounts"
)

// MariaDBSettings represents the settings for MariaDB
type MariaDBSettings struct {
	URL                string
	MaxOpenConnections int
	MaxIdleConnections int
	ConnMaxLifetime    time.Duration
}

// NewMariaDBDriver creates a new instance of MariaDB database
func NewMariaDBDriver(settings MariaDBSettings) (db.Session, error) {
	conn, err := mysql.ParseURL(settings.URL)
	if err != nil {
		logging.Logger.Error("could not parse mariadb connection url", zap.Error(err))
		return nil, err
	}
	session, err := mysql.Open(conn)
	if err != nil {
		logging.Logger.Error("could not open mariadb database", zap.Error(err))
		return nil, err
	}
	session.SetConnMaxLifetime(settings.ConnMaxLifetime)
	session.SetMaxOpenConns(settings.MaxOpenConnections)
	session.SetMaxIdleConns(settings.MaxIdleConnections)

	return session, nil
}
