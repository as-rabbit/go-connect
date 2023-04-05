package mysql

import (
	"fmt"
	"github.com/as-rabbit/go-connect/connecter/config/db"
	"github.com/sirupsen/logrus"
	"go-micro.dev/v4/config"
	"go-micro.dev/v4/config/encoder/json"
	"go-micro.dev/v4/config/source"
	"go-micro.dev/v4/config/source/file"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gLogger "gorm.io/gorm/logger"
	"log"
	"os"
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
	if err = c.c.Get(clusterName).Scan(&mysqlConf); nil != err {

		return nil, err

	}

	fmt.Println("测试", c.c.Get(clusterName).Bytes())

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
			//Logger: logger.NewGORMLogger(logger.Config{
			//	SlowThreshold:             time.Duration(mysqlConf.GormConfig.SlowThreshold) * time.Millisecond,
			//	IgnoreRecordNotFoundError: mysqlConf.GormConfig.IgnoreRecordNotFoundError,
			//	LogLevel:                  glogger.LogLevel(mysqlConf.GormConfig.LogLevel),
			//	Dsn:                       mysqlConf.MysqlConfig.Dsn,
			//}),
			Logger: gLogger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), gLogger.Config{
				SlowThreshold: 100 * time.Millisecond, // 慢查询
				LogLevel:      4,
			}),
		},
		MysqlPoolConfig: db.PoolConfig{
			ConnMaxLifeTime: mysqlConf.PoolConfig.ConnMaxLifeTime,
			MaxIdleConn:     mysqlConf.PoolConfig.MaxIdleConn,
			MaxOpenConn:     mysqlConf.PoolConfig.MaxOpenConn,
		},
	}, nil
}
