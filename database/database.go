package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/my-dev-workspace/go-pkg/logger"
	"net/url"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DBConfig struct {
	Charset                  string
	Collate                  string
	Database                 string
	Host                     string
	Port                     int
	Username                 string
	Password                 string
	TablePrefix              string
	TablePrefixSqlIdentifier string
	Timeout                  int
	ReadTimeout              int
	WriteTimeout             int
	MaxIdleConns             int
	MaxOpenConns             int
	ConnMaxLifetimeSeconds   int
}

type DB struct {
	*sqlx.DB
	Config      DBConfig
	DSN         string
	ErrHandlers []func(err error)
}

type DBGroup struct {
	defaultConfigKey string
	config           map[string]DBConfig
	dbGroup          map[string]*DB
}

func NewDBGroup(defaultConfigName string) (group *DBGroup) {
	group = &DBGroup{}
	group.defaultConfigKey = defaultConfigName
	group.config = map[string]DBConfig{}
	group.dbGroup = map[string]*DB{}
	return
}
func (g *DBGroup) RegisterGroup(cfg map[string]DBConfig) (err error) {
	g.config = cfg
	for name, config := range g.config {
		g.Register(name, config)
		if err != nil {
			return
		}
	}
	return
}
func (g *DBGroup) Register(name string, cfg DBConfig) {
	var db DB
	db = NewDB(cfg)

	g.config[name] = cfg
	g.dbGroup[name] = &db
	return
}
func (g *DBGroup) DB(name ...string) (db *DB) {
	key := ""
	if len(name) == 0 {
		key = g.defaultConfigKey
	} else {
		key = name[0]
	}
	db, _ = g.dbGroup[key]
	return
}

func NewDB(config DBConfig, errHandlerFuncs ...func(err error)) DB {
	logger.Info("Initialize the connection to the database...")

	errHandlers := make([]func(err error), len(errHandlerFuncs))
	copy(errHandlers, errHandlerFuncs)

	db := DB{ErrHandlers: errHandlers}
	err := db.init(config)
	if err != nil {
		logger.Error(err.Error())
	}

	return db
}
func (db *DB) init(config DBConfig) (err error) {
	db.Config = config
	db.DSN = db.getDSN()
	db.DB, err = db.getDB()

	return
}

func (db *DB) getDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?timeout=%dms&readTimeout=%dms&writeTimeout=%dms&charset=%s&collation=%s",
		url.QueryEscape(db.Config.Username),
		db.Config.Password,
		url.QueryEscape(db.Config.Host),
		db.Config.Port,
		url.QueryEscape(db.Config.Database),
		db.Config.Timeout,
		db.Config.ReadTimeout,
		db.Config.WriteTimeout,
		url.QueryEscape(db.Config.Charset),
		url.QueryEscape(db.Config.Collate))
}
func (db *DB) getDB() (connPool *sqlx.DB, err error) {
	connPool, err = sqlx.Connect("mysql", db.getDSN())
	if err != nil {
		return
	}
	connPool.SetConnMaxLifetime(time.Second * time.Duration(db.Config.ConnMaxLifetimeSeconds))
	connPool.SetMaxOpenConns(db.Config.MaxOpenConns)
	connPool.SetMaxIdleConns(db.Config.MaxIdleConns)
	err = connPool.Ping()
	return
}

func NewDBConfigWith(host string, port int, dbName, user, pass string) (cfg DBConfig) {
	cfg = NewDBConfig()
	cfg.Host = host
	cfg.Port = port
	cfg.Username = user
	cfg.Password = pass
	cfg.Database = dbName
	return
}
func NewDBConfig() DBConfig {
	return DBConfig{
		Charset:                  "utf8",
		Collate:                  "utf8_general_ci",
		Database:                 "test",
		Host:                     "127.0.0.1",
		Port:                     3306,
		Username:                 "root",
		Password:                 "",
		TablePrefix:              "",
		TablePrefixSqlIdentifier: "",
		Timeout:                  3000,
		ReadTimeout:              5000,
		WriteTimeout:             5000,
		MaxOpenConns:             500,
		MaxIdleConns:             50,
		ConnMaxLifetimeSeconds:   1800,
	}
}
