package mysql

import (
	"context"
	"database/sql"
	"github.com/sirupsen/logrus"
	"go-connect/connecter/config/db"
	"golang.org/x/sync/singleflight"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
	"time"
)

// Connect entity
var entity *MysqlConnector

// Option Customize
type Option func(*MysqlConnector)

// Configer interface
type Configer interface {
	Get(clusterName string) (*ConnectConfig, error)
}

// ConnectConfig Connect Config
type ConnectConfig struct {
	MysqlConfig     mysql.Config  `json:"mysql_config"`
	GormConfig      *gorm.Config  `json:"gorm_config"`
	MysqlPoolConfig db.PoolConfig `json:"mysql_pool_config"`
}

// Mysql Connector
type MysqlConnector struct {
	mutex       sync.RWMutex
	connections map[string]*gorm.DB
	sf          singleflight.Group
	config      Configer
}

func init() {

	entity = new(MysqlConnector)

	entity.connections = make(map[string]*gorm.DB)

}

// NewConnector MysqlConnector
func NewConnector(config Configer, options ...Option) *MysqlConnector {

	entity.config = config

	for _, fn := range options {

		fn(entity)
	}

	return entity

}

// New Db Connected .
func (m *MysqlConnector) connected(ctx context.Context, clusterName string) (*gorm.DB, error) {

	var (
		dbConn *gorm.DB
		sqlDb  *sql.DB
		err    error
		conf   *ConnectConfig
	)

	if conf, err = m.config.Get(clusterName); nil != err {

		logrus.WithFields(logrus.Fields{
			"cluster_name": clusterName,
			"conf":         conf,
			"error":        err.Error(),
		}).Error("mysql connect get config  error:")

		return nil, err

	}

	logrus.WithFields(logrus.Fields{
		"cluster_name": clusterName,
		"conf":         conf,
	}).Info("mysql connect")

	// Connect Mysql
	if dbConn, err = gorm.Open(mysql.New(conf.MysqlConfig), conf.GormConfig); nil != err {

		logrus.WithFields(logrus.Fields{
			"cluster_name": clusterName,
			"conf":         conf,
			"error":        err.Error(),
		}).Error("mysql connect error:")

		return nil, err
	}

	if sqlDb, err = dbConn.DB(); nil != err {

		logrus.WithFields(logrus.Fields{
			"cluster_name": clusterName,
			"conf":         conf,
			"error":        err.Error(),
		}).Error("fetch sqlDb error:")

		return nil, err

	}

	sqlDb.SetMaxIdleConns(conf.MysqlPoolConfig.MaxIdleConn)

	sqlDb.SetMaxOpenConns(conf.MysqlPoolConfig.MaxOpenConn)

	sqlDb.SetConnMaxLifetime(time.Duration(conf.MysqlPoolConfig.ConnMaxLifeTime) * time.Second)

	m.storage(clusterName, dbConn)

	return dbConn, err
}

// Make Db Connect Return gorm db.
func (m *MysqlConnector) Make(ctx context.Context, clusterName string) (db *gorm.DB, err error) {

	var res interface{}

	if db = m.fetch(clusterName); nil != db {

		return db, nil

	}

	// Avoid introducing concurrency
	if res, err, _ = m.sf.Do(m.connectLockFlag(clusterName), func() (res interface{}, err error) {

		return m.connected(ctx, clusterName)

	}); nil != err {

		return nil, err
	}

	return res.(*gorm.DB), nil
}

// Storage Db Connect .
func (m *MysqlConnector) storage(clusterName string, db *gorm.DB) {

	m.mutex.Lock()

	defer m.mutex.Unlock()

	m.connections[clusterName] = db

}

// Fetch Db Connect .
func (m *MysqlConnector) fetch(clusterName string) *gorm.DB {

	m.mutex.RLock()

	defer m.mutex.RUnlock()

	return m.connections[clusterName]

}

// Connect Lock Flag
func (m *MysqlConnector) connectLockFlag(clusterName string) string {

	return "connect_lock_" + clusterName
}
