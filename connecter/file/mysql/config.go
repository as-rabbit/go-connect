package mysql

import (
	"github.com/sirupsen/logrus"
	"go-connect/connecter/config/db"
	logger "go-connect/log/gorm"
	"go-micro.dev/v4/config"
	"go-micro.dev/v4/config/encoder/json"
	"go-micro.dev/v4/config/source"
	"go-micro.dev/v4/config/source/file"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"time"
)

const configPath = "./config/database.json"

type Config struct {
	c config.Config
}

func NewConfig() (c *Config, err error) {

	var conf config.Config

	fileSource := file.NewSource(
		file.WithPath(configPath),
		source.WithEncoder(json.NewEncoder()),
	)

	if conf, err = config.NewConfig(); nil != err {
		logrus.Error(
			logrus.Fields{
				"err":  err.Error(),
				"path": configPath,
			},
		)
		return nil, err
	}

	if err = conf.Load(fileSource); nil != err {

		logrus.Error(
			logrus.Fields{
				"err":  err.Error(),
				"path": configPath,
			},
		)

		return nil, err
	}

	return &Config{
		c: conf,
	}, nil
}

// Get Connect Config
func (c *Config) Get(clusterName string) (res *ConnectConfig, err error) {

	var mysqlConf *db.ServerMysqlConfig

	// Load Config
	if err = c.c.Scan(mysqlConf); nil != err {

		return nil, err

	}

	return &ConnectConfig{
		MysqlConfig: mysql.Config{
			DSN:                       mysqlConf.MysqlConfig.Dsn,
			SkipInitializeWithVersion: mysqlConf.MysqlConfig.SkipInitializeWithVersion,
			DefaultStringSize:         mysqlConf.MysqlConfig.DefaultStringSize,
			DisableDatetimePrecision:  mysqlConf.MysqlConfig.DisableDatetimePrecision,
			DontSupportRenameIndex:    mysqlConf.MysqlConfig.DontSupportRenameIndex,
			DontSupportRenameColumn:   mysqlConf.MysqlConfig.DontSupportRenameColumn,
		},
		GormConfig: &gorm.Config{
			SkipDefaultTransaction:   mysqlConf.GormConfig.SkipDefaultTransaction,
			DisableNestedTransaction: mysqlConf.GormConfig.DisableNestedTransaction,
			AllowGlobalUpdate:        mysqlConf.GormConfig.AllowGlobalUpdate,
			Logger: logger.NewGORMLogger(logger.Config{
				SlowThreshold:             time.Duration(mysqlConf.GormConfig.SlowThreshold) * time.Millisecond,
				IgnoreRecordNotFoundError: mysqlConf.GormConfig.IgnoreRecordNotFoundError,
				LogLevel:                  glogger.LogLevel(mysqlConf.GormConfig.LogLevel),
				Dsn:                       mysqlConf.MysqlConfig.Dsn,
			}),
		},
		MysqlPoolConfig: db.PoolConfig{
			ConnMaxLifetime: mysqlConf.PoolConfig.ConnMaxLifetime,
			MaxIdleConn:     mysqlConf.PoolConfig.MaxIdleConn,
			MaxOpenConn:     mysqlConf.PoolConfig.MaxOpenConn,
		},
	}, nil
}
