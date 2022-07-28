package mysql

import (
	"database/sql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"sync"
	"time"
)

var entity *mysqlConnector

type Option func(*mysqlConnector)

type mysqlConnector struct {
	mutex       sync.RWMutex
	connections map[string]*gorm.DB
	config      Configer
	logger      logger.Interface
	log         Log
}

func init() {
	
	entity = new(mysqlConnector)
	
	entity.connections = make(map[string]*gorm.DB)
	
}

// With Log Customize
func WithLogger(l logger.Interface) Option {
	
	return func(s *mysqlConnector) {
		
		s.logger = l
		
	}
	
}

// With Config Customize
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
func (m *mysqlConnector) NewConnected(clusterName string) (*gorm.DB, error) {
	
	var (
		dbConn *gorm.DB
		sqbDb  *sql.DB
		err    error
	)
	
	conf := m.config.Get(clusterName)
	
	if dbConn, err = gorm.Open(mysql.New(conf.MysqlConfig), conf.GormConfig); nil != err {
		
		m.log.Error()
		
		return nil, err
	}
	
	// 链接设置
	if sqbDb, err = dbConn.DB(); nil != err {
		
		m.log.Error()
		
		return nil, err
		
	}
	
	if err = sqbDb.Ping(); nil != err {
		
		m.log.Error()
		
		return nil, err
	}
	
	sqbDb.SetMaxIdleConns(conf.MysqlPoolConfig.MaxIdleConn)
	
	sqbDb.SetMaxOpenConns(conf.MysqlPoolConfig.MaxOpenConn)
	
	sqbDb.SetConnMaxLifetime(time.Duration(conf.MysqlPoolConfig.ConnMaxLifetime) * time.Second)
	
	m.storageDbConnect(clusterName, dbConn)
	
	return dbConn, err
}

// Make Db Connect Return gorm db.
func (m *mysqlConnector) MakeConnect(clusterName string) (db *gorm.DB, err error) {
	
	if db = m.dbConnect(clusterName); nil != db {
		
		return db, nil
		
	}
	
	m.mutex.Lock()
	
	db, err = m.NewConnected(clusterName)
	
	m.mutex.Unlock()
	
	return
}

// Storage Db Connect .
func (m *mysqlConnector) storageDbConnect(clusterName string, db *gorm.DB) {
	
	m.mutex.Lock()
	
	defer m.mutex.Unlock()
	
	m.connections[clusterName] = db
	
}

// Fetch Db Connect .
func (m *mysqlConnector) dbConnect(clusterName string) *gorm.DB {
	
	m.mutex.RLock()
	
	defer m.mutex.RUnlock()
	
	return m.connections[clusterName]
	
}
