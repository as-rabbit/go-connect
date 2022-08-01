package file

import (
	"context"
	"database/sql"
	"golang.org/x/sync/singleflight"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"sync"
	"time"
)

// Connect entity
var entity *mysqlConnector

// Option Customize
type Option func(*mysqlConnector)

// Configer interface
type Configer interface {
	Get(clusterName string) (*ConnectConfig, error)
}

// ConnectConfig Connect Config
type ConnectConfig struct {
	MysqlConfig     mysql.Config `json:"mysql_config"`
	GormConfig      *gorm.Config `json:"gorm_config"`
	MysqlPoolConfig PoolConfig   `json:"mysql_pool_config"`
}

// ConnectConfig Connect Config
type mysqlConnector struct {
	mutex       sync.RWMutex
	connections map[string]*gorm.DB
	sf          singleflight.Group
	config      Configer
	logger      logger.Interface
}

func init() {
	
	entity = new(mysqlConnector)
	
	entity.connections = make(map[string]*gorm.DB)
	
	entity.config = NewConfig()
	
}

// WithLogger Customize
func WithLogger(l logger.Interface) Option {
	
	return func(s *mysqlConnector) {
		
		s.logger = l
		
	}
	
}

// WithConfig Customize
func WithConfig(c Configer) Option {
	
	return func(s *mysqlConnector) {
		
		s.config = c
		
	}
	
}

// New MysqlConnector
func (m *mysqlConnector) New(options ...Option) *mysqlConnector {
	
	for _, fn := range options {
		
		fn(entity)
	}
	
	return entity
	
}

// New Db Connected .
func (m *mysqlConnector) connected(ctx context.Context, clusterName string) (*gorm.DB, error) {
	
	var (
		dbConn *gorm.DB
		sqlDb  *sql.DB
		err    error
		conf   *ConnectConfig
	)
	
	if conf, err = m.config.Get(clusterName); nil != err {
		
		return nil, err
		
	}
	
	if dbConn, err = gorm.Open(mysql.New(conf.MysqlConfig), conf.GormConfig); nil != err {
		
		m.logger.Error(ctx, "mysql connect error:", conf)
		
		return nil, err
	}
	
	if sqlDb, err = dbConn.DB(); nil != err {
		
		m.logger.Error(ctx, "fetch sqlDb error:", conf)
		
		return nil, err
		
	}
	
	sqlDb.SetMaxIdleConns(conf.MysqlPoolConfig.MaxIdleConn)
	
	sqlDb.SetMaxOpenConns(conf.MysqlPoolConfig.MaxOpenConn)
	
	sqlDb.SetConnMaxLifetime(time.Duration(conf.MysqlPoolConfig.ConnMaxLifetime) * time.Second)
	
	m.storage(clusterName, dbConn)
	
	return dbConn, err
}

// Make Db Connect Return gorm db.
func (m *mysqlConnector) Make(ctx context.Context, clusterName string) (db *gorm.DB, err error) {
	
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
func (m *mysqlConnector) storage(clusterName string, db *gorm.DB) {
	
	m.mutex.Lock()
	
	defer m.mutex.Unlock()
	
	m.connections[clusterName] = db
	
}

// Fetch Db Connect .
func (m *mysqlConnector) fetch(clusterName string) *gorm.DB {
	
	m.mutex.RLock()
	
	defer m.mutex.RUnlock()
	
	return m.connections[clusterName]
	
}

// connect_lock
func (m *mysqlConnector) connectLockFlag(clusterName string) string {
	
	return "connect_lock_" + clusterName
}
