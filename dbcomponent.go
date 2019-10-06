// Copyright 2019 Luca Stasio <joshuagame@gmail.com>
// Copyright 2019 IT Resources s.r.l.
//
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package sgulengine defines the Sgul Engine structure and functionalities.
package sgulengine

import (
	"fmt"
	"strings"

	"github.com/itross/sgul"
	"github.com/jinzhu/gorm"
)

const mysqlConnectionStringFormat = "%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local"

// DBProvider is the DB Component Provider contract.
type DBProvider interface {
	DB() *gorm.DB
}

// DBComponent is the default Sgul Engine DB component.
// It opens a DB connection and provide it to the clients.
type DBComponent struct {
	BaseComponent
	config       sgul.DB
	db           *gorm.DB
	repositories []sgul.GormRepository
}

// NewDBComponent returns a new DB Compoennts instance.
func NewDBComponent() *DBComponent {
	return &DBComponent{
		BaseComponent: BaseComponent{
			uniqueName: "db",
			logger:     sgul.GetLogger(),
		},
	}
}

// Configure willl configure the db component.
func (dbc *DBComponent) Configure(conf interface{}) error {
	dbc.config = conf.(sgul.DB)
	return nil
}

// Start willl start the DB Server after initialization.
// TODO: customize connectionString and gorm.Open with configured db.type.
func (dbc *DBComponent) Start(e *Engine) error {
	connectionString := dbc.connectionString()
	dbc.logger.Debugf("Connecting to DB at %s", connectionString)

	var err error
	dbc.db, err = gorm.Open(dbc.config.Type, connectionString)
	if err != nil {
		dbc.logger.Errorw("Unable to connect to Database server", "error", err)
		return err
	}

	dbc.db.LogMode(false)
	dbc.logger.Info("DB connection established")

	dbc.injectDB()
	dbc.logger.Info("DB instance injected into repositories")

	return nil
}

func (dbc *DBComponent) injectDB() {
	if len(dbc.repositories) > 0 {
		for _, repository := range dbc.repositories {
			repository.SetDB(dbc.db)
		}
	}
}

// Shutdown will stop serving the API.
func (dbc *DBComponent) Shutdown() {
	if err := dbc.db.Close(); err != nil {
		dbc.logger.Errorw("error shutting down DB Component", "error", err)
	}
}

// DB returns the DB Component db reference.
func (dbc *DBComponent) DB() *gorm.DB {
	return dbc.db
}

// AddRepository adds a sgul Repository to the managed repositories.
func (dbc *DBComponent) AddRepository(repository sgul.GormRepository) {
	dbc.repositories = append(dbc.repositories, repository)
}

// returns the righr db connection string according to the configured db type.
// Note that actually we support only MySQL.
// Default connection string is for MySQL.
func (dbc *DBComponent) connectionString() string {
	switch strings.ToLower(dbc.config.Type) {
	case "mysql":
		return fmt.Sprintf(
			mysqlConnectionStringFormat,
			dbc.config.User,
			dbc.config.Password,
			dbc.config.Host,
			dbc.config.Port,
			dbc.config.Database)
	default:
		return fmt.Sprintf(
			mysqlConnectionStringFormat,
			dbc.config.User,
			dbc.config.Password,
			dbc.config.Host,
			dbc.config.Port,
			dbc.config.Database)
	}
}
