package adapters

import (
	"errors"
	"github.com/go-pg/pg/v9"
	"github.com/hashicorp/go-multierror"
	"github.com/krajeswaran/gostartup/internal/models"
	"github.com/spf13/viper"
)

//DBInterface interface for DB adapter
type DBInterface interface {
	DeepStatus() error
	FetchUser(id string) (*models.User, error)
	CreateUser(name string) (*models.User, error)
}

//DBAdapter - Struct to logically bind all the database related functions
type DBAdapter struct {
	db *pg.DB
}

// DBInit initializes DB connection
func DBInit() *DBAdapter {
	return &DBAdapter{
		db: pg.Connect(&pg.Options{
			Addr:            viper.GetString("DB_ADDR"),
			User:            viper.GetString("DB_USER"),
			Password:        viper.GetString("DB_PASSWORD"),
			Database:        viper.GetString("DB_NAME"),
			ApplicationName: viper.GetString("DB_APPLICATION_NAME"),
			DialTimeout:     viper.GetDuration("DB_DIAL_TIMEOUT"),
			ReadTimeout:     viper.GetDuration("DB_READ_TIMEOUT"),
			PoolSize:        viper.GetInt("DB_CONN_POOL_SIZE"),
			PoolTimeout:     viper.GetDuration("DB_CONN_POOL_TIMEOUT"),
		}),
	}
}

//DeepStatus checks for a DB connection
func (d *DBAdapter) DeepStatus() error {
	var n int
	_, err := d.db.QueryOne(pg.Scan(&n), "SELECT 1")
	if err != nil {
		return multierror.Append(err, errors.New("SERVICE_DB_DOWN"))
	}
	return nil
}

//FetchUser Fetches a user based on id
func (d *DBAdapter) FetchUser(id string) (*models.User, error) {
	user := new(models.User)
	if err := d.db.Model(user).Where("id = ?", id).Select(); err != nil {
		return nil, err
	}

	return user, nil
}

//CreateUser Creates user given a user name
func (d *DBAdapter) CreateUser(name string) (*models.User, error) {
	user := &models.User{
		Name: name,
	}
	_, err := d.db.Model(user).
		Column("id").
		Where("name=?name").
		OnConflict("DO NOTHING").
		Returning("id").
		SelectOrInsert()
	if err != nil {
		return nil, err
	}

	return user, nil
}
