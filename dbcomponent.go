package sgulengine

import (
	"fmt"

	"github.com/itross/sgul"
	"github.com/jinzhu/gorm"
)

// DBComponent is the default Sgul Engine DB component.
// It opens a DB connection and provide it to the clients.
type DBComponent struct {
	BaseComponent
	config sgul.DB
	db     *gorm.DB
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
	connectionString := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		dbc.config.User,
		dbc.config.Password,
		dbc.config.Host,
		dbc.config.Port,
		dbc.config.Database)
	dbc.logger.Debugf("Connecting to DB at %s", connectionString)

	var err error
	dbc.db, err = gorm.Open(dbc.config.Type, connectionString)
	if err != nil {
		dbc.logger.Errorw("Unable to connect to Database server", "error", err)
		return err
	}

	dbc.db.LogMode(false)
	dbc.logger.Info("DB connection established")
	return nil
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
