package mysql

import (
	"github.com/spf13/viper"
	"go-connect/connecter/config/db"
	logger "go-connect/log/gorm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"time"
)

const configPath = "./config/database.json"

type Config struct {
	v *viper.Viper
}

func NewConfig() *Config {

	v := viper.New()

	v.SetConfigFile(configPath)

	return &Config{
		v: v,
	}
}

// Get Connect Config
func (c *Config) Get(clusterName string) (res *ConnectConfig, err error) {

	var config *db.ServerMysqlConfig

	// Load Config
	if err = c.v.ReadInConfig(); nil != err {

		return nil, err

	}

	// Unmarshal To Struct .
	if err = c.v.UnmarshalKey(clusterName, &config); nil != err {

		return nil, err

	}

	return &ConnectConfig{
		MysqlConfig: mysql.Config{
			DSN:                       config.MysqlConfig.Dsn,
			SkipInitializeWithVersion: config.MysqlConfig.SkipInitializeWithVersion,
			DefaultStringSize:         config.MysqlConfig.DefaultStringSize,
			DisableDatetimePrecision:  config.MysqlConfig.DisableDatetimePrecision,
			DontSupportRenameIndex:    config.MysqlConfig.DontSupportRenameIndex,
			DontSupportRenameColumn:   config.MysqlConfig.DontSupportRenameColumn,
		},
		GormConfig: &gorm.Config{
			SkipDefaultTransaction:   config.GormConfig.SkipDefaultTransaction,
			DisableNestedTransaction: config.GormConfig.DisableNestedTransaction,
			AllowGlobalUpdate:        config.GormConfig.AllowGlobalUpdate,
			Logger: logger.NewGORMLogger(logger.Config{
				SlowThreshold:             time.Duration(config.GormConfig.SlowThreshold) * time.Millisecond,
				IgnoreRecordNotFoundError: config.GormConfig.IgnoreRecordNotFoundError,
				LogLevel:                  glogger.LogLevel(config.GormConfig.LogLevel),
				Dsn:                       config.MysqlConfig.Dsn,
			}),
		},
		MysqlPoolConfig: db.PoolConfig{
			ConnMaxLifetime: config.PoolConfig.ConnMaxLifetime,
			MaxIdleConn:     config.PoolConfig.MaxIdleConn,
			MaxOpenConn:     config.PoolConfig.MaxOpenConn,
		},
	}, nil
}
